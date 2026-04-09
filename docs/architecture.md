# 기술 아키텍처 설계 (architecture.md)

> 작성: Architect 에이전트  
> 최종 갱신: 2026-04-09 (rev.2 -- review-notes 피드백 전체 반영)  
> 상태: **확정**

---

## 1. 디렉토리 구조

```
new-lb/
├── main.go                  # 엔트리포인트: CLI 플래그 파싱, 수집/분류/서버 모드 분기
├── go.mod
├── go.sum
├── collector/
│   ├── collector.go         # 수집 오케스트레이터 (ALB+Route53 병렬 → 매칭 → TG+헬스체크 워커풀)
│   ├── alb.go               # ALB 목록 수집 (ELBv2 DescribeLoadBalancers)
│   ├── route53.go           # Route53 HostedZone/레코드 전체 조회 → 메모리 적재
│   ├── targetgroup.go       # Target Group 목록 + 타겟 헬스 수집
│   └── healthcheck.go       # Route53 도메인 HTTPS(443) 헬스체크
├── classifier/
│   ├── classifier.go        # 자동 분류 (이름 파싱 + 키워드 매칭 + 조치상태 추론)
│   └── tikitaka.go          # 미분류 항목 사용자 질의 인터페이스 (CLI)
├── model/
│   └── entry.go             # 데이터 모델 (Entry, Record, TGInfo, MergePlan)
├── store/
│   └── store.go             # JSON 파일 읽기/쓰기 + sync.Mutex 동시성 제어
├── server/
│   ├── server.go            # HTTP 서버 (localhost only) + API 라우팅
│   └── handler.go           # API 핸들러 구현
├── web/
│   └── index.html           # 단일 HTML/JS 파일 (FR-17)
├── data/
│   └── entries.json         # 수집/분류 결과 저장 (런타임 생성)
└── docs/
    ├── requirements.md
    ├── plan.md
    ├── architecture.md      # (본 문서)
    ├── aws-report.md
    └── review-notes.md
```

### 패키지 역할 요약

| 패키지 | 역할 | 의존 |
|--------|------|------|
| `main` | CLI 플래그 파싱, 모드 분기 (collect/classify/serve) | collector, classifier, server, store |
| `collector` | AWS API 호출로 ALB/Route53/TG/헬스체크 데이터 수집 | model, AWS SDK v2 |
| `classifier` | ALB 이름 파싱, 솔루션/환경 매핑, 조치상태 자동 추론, 티키타카 | model |
| `model` | 데이터 구조체 정의 (Entry, Record, TGInfo, MergePlan) | (의존 없음) |
| `store` | JSON 파일 CRUD + sync.Mutex 동시성 보호 | model |
| `server` | HTTP API + 정적 파일 서빙 (localhost:8080) | model, store, collector, classifier |
| `web` | 단일 HTML/JS 파일 (테이블 뷰, 필터, 드래그앤드롭) | (서버에서 서빙) |

---

## 2. 데이터 모델

### 2.1 Entry 구조체 (핵심)

```go
package model

type Entry struct {
    ALBName     string    `json:"albName"`
    ALBArn      string    `json:"albArn"`
    ALBDNS      string    `json:"albDns"`
    Solution    string    `json:"solution"`    // signage|cms|nserise|bss|ncount|srt|aiagent|wine|ws2025|dooh|unknown
    Environment string    `json:"environment"` // dev|stg|prd|unknown
    Status      string    `json:"status"`      // healthy|unhealthy|no_target|no_record|error|unknown
    Action      string    `json:"action"`      // 유지|합치기|삭제|미정
    Records     []Record  `json:"records"`     // Route53 레코드 (구조체 배열)
    TGs         []TGInfo  `json:"targetGroups"` // Target Group 상세 (구조체 배열)
    MergeTarget string    `json:"mergeTarget,omitempty"` // 합칠 대상 ALB 이름
    Note        string    `json:"note,omitempty"`
}
```

> **review-notes AI-01 반영**: records와 targetGroups를 string 배열이 아닌 구조체 배열로 확장. 헬스체크 결과와 TG 상세 정보를 함께 저장한다.

### 2.2 Record 구조체

```go
type Record struct {
    Name       string `json:"name"`       // Route53 레코드 이름 (예: dev-scms.nextpay.co.kr)
    ZoneID     string `json:"zoneId"`     // Hosted Zone ID
    ZoneName   string `json:"zoneName"`   // Hosted Zone 이름 (예: nextpay.co.kr)
    Type       string `json:"type"`       // ALIAS 또는 CNAME
    HealthCode int    `json:"healthCode"` // HTTP 응답코드 (200, 403, 502 등), -1=인증서에러, 0=타임아웃
}
```

### 2.3 TGInfo 구조체

```go
type TGInfo struct {
    Name          string `json:"name"`          // Target Group 이름
    ARN           string `json:"arn"`
    HealthyCount  int    `json:"healthyCount"`
    UnhealthyCount int   `json:"unhealthyCount"`
    UnusedCount   int    `json:"unusedCount"`
    TargetCount   int    `json:"targetCount"`   // 등록된 타겟 총 수
}
```

### 2.4 MergePlan 구조체

```go
type MergePlan struct {
    TargetALB   string   `json:"targetAlb"`   // 합칠 대상 ALB 이름
    SourceALBs  []string `json:"sourceAlbs"`  // 합쳐질 ALB 이름들
    Records     []Record `json:"records"`     // 이동해야 할 Route53 레코드
    TGs         []TGInfo `json:"targetGroups"` // 이동해야 할 Target Group
    Note        string   `json:"note,omitempty"`
}
```

### 2.5 Status 필드 정의

| 값 | 의미 | 판정 기준 |
|----|------|-----------|
| `healthy` | 모든 TG의 타겟이 healthy | TG 존재 + 모든 타겟 healthy |
| `unhealthy` | 하나 이상의 TG에서 unhealthy 타겟 존재 | unhealthy 타겟 > 0 |
| `no_target` | TG는 있으나 등록된 타겟이 없음 | TG 존재 + targetCount == 0 |
| `no_record` | Route53 레코드에 연결되지 않음 | records 배열 길이 == 0 |
| `error` | HTTPS 헬스체크에서 모든 경로 에러 | 500+ 응답 또는 인증서 에러 또는 타임아웃 |
| `unknown` | TG 조회 실패 등 판정 불가 | API 호출 실패 시 (review-notes AI-03 반영) |

> **Status 판정 우선순위**: no_target > no_record > error > unhealthy > healthy. unknown은 API 실패 시에만 사용.

### 2.6 Action(조치상태) 정의

| 값 | 의미 |
|----|------|
| `유지` | 정상 운영, 통합 불필요 |
| `합치기` | 동일 환경-솔루션의 대표 ALB에 통합 |
| `삭제` | 타겟 없음, 레코드 없음, 에러 등으로 삭제 대상 |
| `미정` | 자동 판단 불가, 사용자 입력 대기 |

---

## 3. 수집 전략

### 3.1 전체 수집 파이프라인

```
Phase A (병렬)                    Phase B              Phase C (워커풀)
┌─────────────┐                 ┌──────────┐          ┌──────────────┐
│ ALB 목록 수집 │────┐           │  매칭    │          │ TG 상태 수집   │
│ (ELBv2 API)  │    ├──────────▶│ ALB DNS ↔│──────────▶│ (워커 10개)   │
│              │    │           │ Route53  │          │              │
└─────────────┘    │           └──────────┘          │ HTTPS 헬스체크 │
┌─────────────┐    │                                 │ (워커 5개)    │
│ Route53 전체  │────┘                                 └──────────────┘
│ 조회 + 메모리 │
│ 적재          │
└─────────────┘
```

### 3.2 핵심 설계 결정

| 항목 | 결정 | 근거 |
|------|------|------|
| Route53 조회 방식 | 전체 조회 → 메모리 적재 → ALB DNS 매칭 | FR-03: 개별 조회 금지. ALB 48개 x 레코드별 조회 시 API 과다 호출 |
| ALB-레코드 매칭 | ALIAS의 DNSName 또는 CNAME 값이 ALB DNS와 일치하면 매칭 | Route53 레코드 타입에 따라 분기 |
| 병렬 전략 | errgroup + 세마포어 패턴 | Go 표준 병렬 패턴, AWS SDK v2와 호환 |
| k8s ALB 필터링 | ALB 이름이 "k8s"로 시작하면 제외 | FR-02: strings.HasPrefix(name, "k8s") |

### 3.3 Route53 메모리 매칭 상세

```go
// 1단계: Route53 전체 조회
func collectAllRecords(ctx context.Context) (map[string][]Record, error) {
    // ListHostedZones → 각 Zone에 대해 ListResourceRecordSets
    // 반환: map[albDNS][]Record  (ALB DNS를 키로, 연결된 레코드 목록을 값으로)
    //
    // ALIAS 레코드: record.AliasTarget.DNSName → ALB DNS와 비교
    // CNAME 레코드: record.ResourceRecords[0].Value → ALB DNS와 비교
    // DNS 비교 시 trailing dot(.) 제거 후 대소문자 무시 비교
}
```

> **실데이터 참고**: aws-report.md 기준 Hosted Zone 2개(nextpay.co.kr, nextaistore.co.kr/com), 레코드 약 80개. 메모리 부담 없음.

### 3.4 동시성 상수

```go
const (
    maxTGWorkers     = 10  // TG 상태 조회 동시성 (ALB 48개 → 10개 워커로 분산)
    maxHealthWorkers = 5   // HTTPS 헬스체크 동시성 (외부 서버 부하 고려, 보수적 설정)
    httpTimeout      = 10 * time.Second  // HTTPS 요청 타임아웃
)
```

### 3.5 Rate Limit 대응

- AWS SDK v2 내장 retry (기본 3회, 지수 백오프)를 활용
- 세마포어(golang.org/x/sync/semaphore 또는 채널)로 동시 API 호출 수 제한
- TG 조회 부분 실패 시 해당 ALB만 `status="unknown"`으로 마킹하고 나머지 계속 수집 (review-notes A3)

---

## 4. HTTPS 헬스체크 모듈

### 4.1 설계 개요

Route53 레코드에 연결된 도메인으로 HTTPS(443) 요청을 보내 서비스 가용성을 판정한다. ALB DNS가 아니라 **실제 도메인**으로 요청하여 사용자 관점의 접근성을 확인한다.

### 4.2 요청 흐름

```
도메인 (예: dev-scms.nextpay.co.kr)
  │
  ├─ GET https://도메인/          → 응답코드 확인
  │   └─ 200~499 → "정상" (즉시 종료)
  │   └─ 500+ → 다음 경로 시도
  │
  ├─ GET https://도메인/health    → 응답코드 확인
  │   └─ 200~499 → "정상" (즉시 종료)
  │   └─ 500+ → 다음 경로 시도
  │
  └─ GET https://도메인/healthz   → 응답코드 확인
      └─ 200~499 → "정상" (즉시 종료)
      └─ 500+ → "에러"
```

### 4.3 응답 판정 기준

| 응답 | healthCode | 판정 | 근거 |
|------|-----------|------|------|
| 200~499 | 실제 코드 | 정상 | 400대는 서버가 응답 중 (review-notes A2, 3인 합의 완료) |
| 500~599 | 실제 코드 | 에러 | 서버 에러 |
| TLS/인증서 에러 | -1 | 에러 | 인증서 만료/불일치 등 (requirements Q4) |
| 타임아웃 (10초) | 0 | 에러 | 서버 무응답 |
| DNS 해석 실패 | -2 | 에러 | 도메인 자체가 해석 불가 |

### 4.4 구현 핵심

```go
func checkHealth(domain string) int {
    // TLS 검증 비활성화하지 않음 → 인증서 에러도 감지
    // 단, 인증서 에러 시 InsecureSkipVerify=true로 재요청하여 서버 자체는 살아있는지 확인
    //
    // 경로 순차 시도: "/", "/health", "/healthz"
    // 하나라도 200~499 응답 시 해당 코드 반환 (정상)
    // 모든 경로에서 500+ 또는 에러 → 마지막 에러 코드 반환
    //
    // http.Client{
    //     Timeout: httpTimeout,
    //     CheckRedirect: 최대 3회 리다이렉트 허용,
    // }
}
```

### 4.5 ALB의 최종 status 결정

하나의 ALB에 여러 레코드가 연결될 수 있다 (예: webService는 12개 레코드).

- **모든 레코드가 정상(200~499)** → `healthy`
- **하나라도 에러(500+, -1, 0, -2)** → `error` (보수적 판정)
- **레코드가 없음** → `no_record`

> TG 상태(no_target, unhealthy)가 헬스체크보다 우선. Status 판정 우선순위는 2.5 참조.

---

## 5. ALB 분류 로직

### 5.1 실데이터 기반 현실 분석

aws-report.md 분석 결과:

- **전체 ALB**: 48개 (k8s 제외)
- **네이밍 패턴 준수**: 3개 (`alb-dev-signage`, `lb-dev-scms`, `lb-dev-swaiting`)
- **패턴 미준수**: 45개
- **솔루션 추정 가능**: 46개 (키워드 기반, rev.3 보강 후)
- **솔루션 불명(unknown)**: 2개 (apiService, next-office-dev)

> 이름 기반 정규식 파싱(`{리소스}-{환경}-{솔루션}-{서비스}`)은 48개 중 3개만 매칭되므로 **보조 수단**으로만 사용한다. 주력은 키워드 매칭이다.

### 5.2 분류 전략 (3단계)

```
단계 1: 정규식 패턴 파싱
  └─ ^(alb|lb|elb)-?(dev|stg|prd)-(.+)$ → 환경+솔루션 추출
  └─ 매칭되면 솔루션 키워드와 교차 검증

단계 2: 키워드 매칭 (주력)
  └─ ALB 이름에 솔루션 키워드가 포함되면 자동 분류
  └─ 환경은 이름에 "dev", "staging", "stg" 포함 여부로 판정, 없으면 "prd"

단계 3: 티키타카 Fallback
  └─ 단계 1, 2 모두 매칭 실패 시 사용자에게 질의
```

### 5.3 솔루션 키워드 매핑

| 솔루션 | 키워드 (ALB 이름에서 탐색) | 매칭 ALB 예시 |
|--------|--------------------------|-------------|
| signage | `signage`, `scms`, `swaiting`, `tizenweb`, `oss-cms`, `oss-waiting`, `product-8080`, `waiting-808` | BSSTizenWebSignage, lb-scms, lb-swaiting, dev-product-8080, waiting-8081 |
| cms | `cms-elb`, `dev-cms` | cms-elb, dev-cms-elb |
| aiagent | `aiagent`, `knowledge-graph` | aiagent-staging, elb-knowledge-graph |
| nserise | `nseries`, `nextpay-kiosk`, `nextpay-nkds`, `nextpay-norder`, `nextpay-npos`, `nextpay-order`, `nextpay-shop`, `socket-dev`, `webService` | nseries-elb, nextpay-kiosk-dev, socket-dev, webService |
| bss | `bss-` | bss-cms, bss-kiosk, bss-order |
| ncount | `ncount` | NCount-Live-ELB, ncount-dev |
| srt | `srt` | elb-srt-device-api, srt-store-api |
| wine | `wine` | wine-curation-elb |
| ws2025 | `ws2025` | ws2025-lb |
| dooh | `dooh` | dooh-dev-elb |

> **rev.3 반영**: aws-report.md 분류 결과에 따라 키워드 보강. signage에 `product-8080`, `waiting-808` 추가 (dev-product-8080, waiting-8081 등 매칭). nserise에 `socket-dev`, `webService` 추가 (이름에 nserise 키워드가 없으나 실데이터 분류 결과 반영). 이 두 항목은 일반적인 키워드 매칭이 아닌 **ALB 이름 완전 일치** 또는 **고유 키워드**이므로, 구현 시 키워드 매칭 순서에서 범용 키워드보다 우선 적용해야 한다.

### 5.4 환경 판정 규칙

```go
func detectEnvironment(albName string) string {
    lower := strings.ToLower(albName)
    switch {
    case strings.Contains(lower, "-dev") || strings.HasPrefix(lower, "dev-"):
        return "dev"
    case strings.Contains(lower, "staging") || strings.Contains(lower, "-stg"):
        return "stg"
    default:
        return "prd"  // 환경 키워드 없으면 운영으로 간주
    }
}
```

### 5.5 조치상태 자동 추론 규칙

```go
func inferAction(entry *Entry) string {
    // 1. 타겟 없음 → 삭제
    if entry.Status == "no_target" {
        return "삭제"
    }
    // 2. 레코드 없음 → 삭제
    if entry.Status == "no_record" {
        return "삭제"
    }
    // 3. HTTPS 에러 → 삭제 (공지 후)
    if entry.Status == "error" {
        return "삭제"
    }
    // 4. 중지 프로젝트 → 미정 (컨펌 대기)
    if isSuspendedProject(entry.Solution) {
        return "미정"
    }
    // 5. 정상 → 유지
    return "유지"
}
```

### 5.6 티키타카 질의 (미분류 2개 ALB)

aws-report.md 분류 결과 반영 후, 키워드 매칭으로도 분류 불가능한 ALB는 2개:

| ALB 이름 | 환경 | 상태 | 레코드 | 질의 내용 |
|---------|------|------|--------|----------|
| apiService | prd | unused | api.nextpay.co.kr | 솔루션 + 조치상태 선택 필요 |
| next-office-dev | dev | unused | dev-nextoffice.nextpay.co.kr | 솔루션 + 조치상태 선택 필요 |

> **rev.3 반영**: 기존 unknown 10개 중 8개가 aws-report.md 분석으로 분류됨.
> - signage: dev-product-8080, dev-waiting-8081, dev-waiting-8082, product-8080, waiting-8081, waiting-8082
> - nserise: socket-dev, webService
> - 이들의 키워드는 5.3 매핑 테이블에 반영 완료.

**CLI 티키타카 형식**:
```
[1/2] ALB: apiService (prd, unused, 레코드: api.nextpay.co.kr)
  솔루션을 선택해주세요:
  1. signage  2. cms  3. nserise  4. bss  5. ncount
  6. srt  7. aiagent  8. wine  9. ws2025  10. dooh
  11. unknown (모름)
  선택: _

  조치상태를 선택해주세요:
  1. 유지  2. 합치기  3. 삭제
  선택: _
```

---

## 6. API 설계

### 6.1 엔드포인트 목록

| 메서드 | 경로 | 설명 | 요청 | 응답 |
|--------|------|------|------|------|
| `GET` | `/api/entries` | 전체 ALB 목록 조회 | query: ?solution=&status=&action= | `[]Entry` |
| `GET` | `/api/entries/:name` | 단일 ALB 조회 | path param | `Entry` |
| `PATCH` | `/api/entries/:name` | 솔루션/조치상태 수정 (인라인 편집) | `{"solution":"...", "action":"...", "note":"..."}` | `Entry` |
| `POST` | `/api/collect` | 수집 실행 트리거 | (없음) | `{"status":"started", "count":48}` |
| `POST` | `/api/classify` | 자동 분류 실행 (review-notes AI-07 반영) | (없음) | `{"classified":46, "unknown":2}` |
| `GET` | `/api/report` | 현황 리포트 (솔루션별/상태별 요약) | (없음) | `Report` |
| `GET` | `/` | 웹 UI (index.html 서빙) | (없음) | HTML |

> **review-notes AI-07 반영**: `POST /api/classify` 엔드포인트를 추가하여 plan.md Phase 8.2와의 불일치를 해소함. 자동 분류 가능 항목만 분류하고, unknown은 PATCH로 개별 수정한다.

### 6.2 Report 응답 구조

```go
type Report struct {
    Total         int                    `json:"total"`
    BySolution    map[string]int         `json:"bySolution"`
    ByStatus      map[string]int         `json:"byStatus"`
    ByAction      map[string]int         `json:"byAction"`
    ByEnvironment map[string]int         `json:"byEnvironment"`
    DeleteCandidates []string            `json:"deleteCandidates"`
    MergeCandidates  []MergePlan         `json:"mergeCandidates"`
}
```

### 6.3 API 설계 원칙

- **net/http 표준 라이브러리만 사용** (review-notes A5: gin/echo 등 외부 프레임워크 불필요)
- 라우팅: `http.ServeMux` 사용, 패턴 매칭으로 처리
- JSON 응답: `encoding/json`
- 에러 응답: `{"error": "message"}` + 적절한 HTTP 상태 코드
- CORS 불필요 (localhost 동일 오리진)

---

## 7. 웹 UI 구조

### 7.1 단일 HTML/JS 파일 (FR-17)

`web/index.html` 하나로 구현. 외부 CDN 최소화.

- **CSS**: 인라인 `<style>` 또는 최소한의 Pico CSS CDN
- **JS**: 인라인 `<script>`, fetch API로 서버와 통신
- **프레임워크 없음**: 바닐라 JS

### 7.2 화면 레이아웃

```
┌─────────────────────────────────────────────────────────┐
│  ALB 리소스 정리 도구                    [수집] [분류]     │
├─────────────────────────────────────────────────────────┤
│  필터: [솔루션 ▼] [상태 ▼] [조치상태 ▼] [환경 ▼]  [검색]   │
├──┬──────────┬────────┬─────┬────────┬──────┬───────────┤
│# │ ALB 이름  │ 솔루션  │ 환경 │ 상태   │조치상태│ 레코드    │
├──┼──────────┼────────┼─────┼────────┼──────┼───────────┤
│1 │ alb-xxx  │signage │ dev │healthy │ 유지  │ 3개       │
│2 │ bss-cms  │ bss    │ prd │healthy │ 미정  │ 1개       │
│..│  ...     │  ...   │ ... │  ...   │ ...  │  ...      │
├──┴──────────┴────────┴─────┴────────┴──────┴───────────┤
│  통계: 전체 48 | 유지 15 | 합치기 8 | 삭제 15 | 미정 10    │
└─────────────────────────────────────────────────────────┘
```

### 7.3 기능

- **테이블 뷰**: ALB 현황 표시 (P1)
- **필터링**: 솔루션별, 상태별, 조치상태별, 환경별 필터 (P1)
- **인라인 편집**: 솔루션/조치상태 셀 클릭 시 드롭다운으로 수정 → PATCH API 호출 (P1, review-notes P4 반영)
- **드래그앤드롭 합치기**: 합칠 대상 ALB에 드래그 시 관련 정보 표시 (P2)
- **리포트 보기**: GET /api/report 호출하여 솔루션별/상태별 요약 표시 (P1)

### 7.4 서버 바인딩

```go
http.ListenAndServe("127.0.0.1:8080", mux)  // localhost only (FR-19)
```

---

## 8. 파일 저장

### 8.1 저장 경로

```
data/entries.json     # ALB 엔트리 배열
```

### 8.2 저장 형식

```json
[
  {
    "albName": "alb-dev-signage",
    "albArn": "arn:aws:elasticloadbalancing:...",
    "albDns": "alb-dev-signage-299073066.ap-northeast-2.elb.amazonaws.com",
    "solution": "signage",
    "environment": "dev",
    "status": "healthy",
    "action": "유지",
    "records": [
      {
        "name": "dev-scms.nextpay.co.kr",
        "zoneId": "Z1234567890",
        "zoneName": "nextpay.co.kr",
        "type": "ALIAS",
        "healthCode": 200
      }
    ],
    "targetGroups": [
      {
        "name": "tg-dev-signage",
        "arn": "arn:aws:elasticloadbalancing:...",
        "healthyCount": 2,
        "unhealthyCount": 1,
        "unusedCount": 0,
        "targetCount": 3
      }
    ],
    "mergeTarget": "",
    "note": ""
  }
]
```

### 8.3 동시성 제어 (review-notes AI-06 반영)

```go
type Store struct {
    mu       sync.Mutex
    filePath string
}

func (s *Store) Load() ([]Entry, error) {
    s.mu.Lock()
    defer s.mu.Unlock()
    // 파일 읽기 → JSON 디코드
}

func (s *Store) Save(entries []Entry) error {
    s.mu.Lock()
    defer s.mu.Unlock()
    // 임시 파일에 쓰기 → os.Rename (원자적 교체)
}

func (s *Store) UpdateEntry(name string, updates map[string]interface{}) (*Entry, error) {
    s.mu.Lock()
    defer s.mu.Unlock()
    // Load → 수정 → 임시 파일 쓰기 → Rename
}
```

> **원자적 저장**: 임시 파일(`data/entries.json.tmp`)에 먼저 쓰고 `os.Rename`으로 교체하여 부분 쓰기 방지.
> **sync.Mutex**: 웹 UI PATCH와 POST /api/collect 동시 실행 시 경합 방지 (review-notes Planner 지적 반영).

---

## 9. 시퀀스 다이어그램

### 9.1 수집 흐름 (POST /api/collect)

```
Client          Server          Collector        AWS API          Store
  │               │               │                │               │
  │──POST /collect──▶│             │                │               │
  │               │──Run()────────▶│               │               │
  │               │               │                │               │
  │               │               │──[Phase A: 병렬]──▶            │
  │               │               │  ListLoadBalancers ──▶         │
  │               │               │  ListHostedZones ────▶         │
  │               │               │  ListResourceRecordSets ──▶    │
  │               │               │◀── ALB목록 + 레코드맵 ──────────│
  │               │               │                │               │
  │               │               │──[Phase B: 매칭]               │
  │               │               │  ALB DNS ↔ 레코드맵             │
  │               │               │  → Entry.Records 채움           │
  │               │               │  → k8s ALB 필터링               │
  │               │               │                │               │
  │               │               │──[Phase C: 워커풀]──▶           │
  │               │               │  TG 상태 조회 (10 workers) ──▶  │
  │               │               │  HTTPS 헬스체크 (5 workers)     │
  │               │               │◀── TG상태 + 헬스코드 ──────────  │
  │               │               │                │               │
  │               │               │──Status 판정──────────────────  │
  │               │               │──────────────────────────────▶ Save()
  │◀──200 OK──────│◀──완료────────│               │               │
```

### 9.2 분류 흐름 (POST /api/classify)

```
Client          Server          Classifier       Store
  │               │               │               │
  │──POST /classify──▶│           │               │
  │               │──Classify()──▶│              │
  │               │               │◀──Load()────│
  │               │               │              │
  │               │               │──[Step 1: 패턴 파싱]
  │               │               │  ALB 이름 → 정규식 매칭
  │               │               │
  │               │               │──[Step 2: 키워드 매칭]
  │               │               │  ALB 이름에 솔루션 키워드 탐색
  │               │               │
  │               │               │──[Step 3: 조치상태 추론]
  │               │               │  Status 기반 규칙 적용
  │               │               │
  │               │               │──Save()────▶│
  │◀──200 {classified:38}──│◀────│              │
```

### 9.3 합치기 계획 흐름

```
Client(웹UI)     Server          Store
  │               │               │
  │──드래그앤드롭──▶│              │
  │  (sourceALB → targetALB)     │
  │               │◀──Load()────│
  │               │               │
  │               │──MergePlan 생성:
  │               │  - 이동할 Route53 레코드
  │               │  - 이동할 TG 목록
  │               │  - 리스너 룰 정보
  │               │               │
  │◀──MergePlan JSON──│          │
  │               │               │
  │  (사용자가 계획 확인)          │
  │  ※ 실제 AWS 변경은 별도 절차   │
```

---

## 10. 기술 스택

| 항목 | 선택 | 비고 |
|------|------|------|
| 언어 | Go 1.22+ | TC-01 |
| AWS SDK | aws-sdk-go-v2 | TC-03 |
| 병렬 처리 | golang.org/x/sync/errgroup | 세마포어 패턴 |
| HTTP 서버 | net/http 표준 라이브러리 | A5: 외부 프레임워크 불필요 |
| JSON | encoding/json | TC-05 |
| 웹 UI | 바닐라 HTML/JS | FR-17: 단일 파일 |
| 외부 의존 | AWS SDK v2 + x/sync만 | NFR-02: 과설계 방지 |

---

## 11. 에러 처리 전략

| 상황 | 대응 | 결과 |
|------|------|------|
| ALB 목록 조회 실패 | 전체 수집 중단, 에러 반환 | 재시도 필요 |
| Route53 조회 실패 | 전체 수집 중단, 에러 반환 | 재시도 필요 |
| 특정 ALB의 TG 조회 실패 | 해당 ALB만 status="unknown", 나머지 계속 | A3 반영 |
| 특정 도메인 헬스체크 실패 | healthCode=0(타임아웃) 기록, 나머지 계속 | 부분 실패 허용 |
| entries.json 파일 없음 | 빈 배열로 초기화 | 최초 실행 |
| entries.json 파싱 에러 | 에러 로그 + 빈 배열 반환 | 재수집으로 복구 |

---

## 12. review-notes 피드백 반영 요약

| 액션 ID | 내용 | 반영 위치 |
|---------|------|-----------|
| AI-01 | records/targetGroups를 구조체 배열로 확장 | 2.2, 2.3 |
| AI-02 | 400대 응답 = 정상 판정 기준 명시 | 4.3 |
| AI-03 | status 필드에 "unknown" 추가 | 2.5 |
| AI-06 | store 패키지에 sync.Mutex 설계 추가 | 8.3 |
| AI-07 | POST /api/classify 엔드포인트 추가 | 6.1 |
| A2 | 200~499=정상, 500+=에러 (3인 합의) | 4.3 |
| A3 | TG 조회 부분 실패 허용 | 3.5, 11 |
| A5 | 외부 프레임워크 불사용 | 10 |
| P4 | 웹 UI 인라인 편집 + PATCH API | 6.1, 7.3 |

---

## 변경 이력

| 날짜 | 변경 내용 | 변경자 |
|------|-----------|--------|
| 2026-04-09 | 초안 작성 | Architect |
| 2026-04-09 | rev.2 -- review-notes 피드백 전체 반영 (AI-01~07, A2/A3/A5, P4), 실데이터 기반 분류 전략 보강 | Architect |
| 2026-04-09 | rev.3 -- unknown ALB 10->2 반영: 5.3 키워드 매핑 보강(signage: product-8080/waiting-808, nserise: socket-dev/webService), 5.6 티키타카 목록 축소, 5.1 수치 갱신, API classify 응답 예시 갱신 | Architect |
