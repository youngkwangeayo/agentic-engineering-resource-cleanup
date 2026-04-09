# ALB Resource Cleanup Tool - 기능정의서

> 버전: 1.0 | 작성일: 2026-04-09

---

## 1. 개요

### 1.1 이 도구는 무엇인가?
AWS ALB(Application Load Balancer) 리소스를 **조사하고 정리 계획을 세우는** 도구이다.
약 60개의 ALB와 연결된 Route53 레코드, Target Group을 자동으로 수집하고,
어떤 프로젝트(솔루션)에 속하는지 분류한 뒤, 유지/합치기/삭제를 판단한다.

> **중요:** 이 도구는 계획 수립까지만 수행한다. 실제 AWS 리소스를 삭제하거나 변경하지 않는다.

### 1.2 누가 사용하는가?
- 인프라 운영자
- 개발팀 리더 (ALB 정리 의사결정)

### 1.3 왜 필요한가?
- ALB가 60개 이상 누적되어 어떤 것이 사용 중이고 어떤 것이 폐기 가능한지 파악이 어려움
- 수동으로 하나씩 확인하면 시간이 많이 소요됨
- 자동 수집 + 분류로 빠르게 현황을 파악하고 정리 계획을 세울 수 있음

---

## 2. 사용 방법

### 2.1 실행 모드

```bash
# 방법 1: 한번에 전부 (수집 → 분류 → 웹UI)
go run . -mode=all

# 방법 2: 단계별 실행
go run . -mode=collect    # AWS 데이터 수집 (1~2분 소요)
go run . -mode=classify   # 자동 분류
go run . -mode=serve      # 웹 UI 시작 (기본값)
```

| 모드 | 설명 | 소요시간 |
|------|------|----------|
| `collect` | AWS에서 ALB, Route53, Target Group, 헬스체크 데이터 수집 | 1~2분 |
| `classify` | 수집된 ALB를 솔루션/환경/조치상태로 자동 분류 | 즉시 |
| `serve` | 웹 UI 시작 (기본 모드) | 즉시 |
| `all` | 위 3개를 순서대로 실행 | 1~2분 |

### 2.2 옵션

| 플래그 | 기본값 | 설명 |
|--------|--------|------|
| `-mode` | `serve` | 실행 모드 선택 |
| `-data` | `data` | 데이터 저장 디렉토리 경로 |

---

## 3. 핵심 기능

### 3.1 AWS 데이터 수집 (`collect`)

수집 단계에서 다음 4가지 데이터를 자동으로 가져온다:

#### (1) ALB 목록
- AWS ELBv2 API로 전체 ALB 조회
- `k8s`로 시작하는 ALB는 자동 제외 (k8s ALB controller가 관리하므로)
- 수집 항목: ALB 이름, ARN, DNS 주소

#### (2) Route53 레코드
- 모든 Hosted Zone의 레코드를 한번에 조회
- ALB DNS를 가리키는 ALIAS/CNAME 레코드를 찾아 매칭
- 예: `dev-scms.nextpay.co.kr` → `alb-dev-signage-xxx.elb.amazonaws.com`

#### (3) Target Group 정보
- 각 ALB에 연결된 Target Group 조회 (10개 동시 처리)
- TG별 타겟 상태: 정상(healthy), 비정상(unhealthy), 미사용(unused), 전체 수

#### (4) HTTPS 헬스체크
- Route53에 연결된 도메인에 실제 HTTPS 요청 전송 (5개 동시 처리)
- 순서대로 `/`, `/health`, `/healthz` 시도
- 응답 코드 판정:

| 응답 | 의미 |
|------|------|
| 200~499 | 정상 (서버 응답 중) |
| 500~599 | 서버 에러 |
| -1 | 인증서 오류 |
| -2 | DNS 해석 실패 |
| 0 | 타임아웃 (10초) |

#### 수집 결과: ALB 상태 판정

수집된 데이터를 종합하여 각 ALB의 상태를 결정한다 (우선순위 순서):

| 상태 | 조건 | 의미 |
|------|------|------|
| `no_target` | Target Group에 등록된 타겟이 0개 | 트래픽을 받을 서버가 없음 |
| `no_record` | Route53 레코드가 없음 | 어떤 도메인도 이 ALB를 가리키지 않음 |
| `error` | 헬스체크 결과 500+, 인증서오류, DNS오류, 타임아웃 | 접속 불가 |
| `unhealthy` | TG에 비정상 타겟 존재 | 일부 서버 문제 |
| `healthy` | 위 조건에 해당 없음 | 정상 운영 중 |
| `unknown` | 판정 불가 | TG 조회 실패 등 |

---

### 3.2 자동 분류 (`classify`)

수집된 ALB를 3가지 기준으로 자동 분류한다.

#### (1) 솔루션 (프로젝트) 분류

ALB 이름에 포함된 키워드로 어떤 프로젝트의 ALB인지 판별한다:

| 솔루션 | 매칭 키워드 | 상태 |
|--------|------------|------|
| signage | signage, scms, swaiting, tizenweb, oss-cms | 운영중 |
| cms | cms-elb, dev-cms | 운영중 |
| aiagent | aiagent, knowledge-graph | 운영중 |
| nserise | nseries, nextpay-kiosk, nextpay-npos, webservice 등 | 중단 |
| bss | bss- | 중단 |
| ncount | ncount | 중단 |
| srt | srt | 중단 |
| wine | wine | 중단 |
| ws2025 | ws2025 | 중단 |
| dooh | dooh | 중단 |
| unknown | 키워드 매칭 실패 | 미분류 |

#### (2) 환경 분류

| 환경 | 판정 기준 |
|------|----------|
| `dev` | ALB 이름에 `-dev` 또는 `dev-` 포함 |
| `stg` | ALB 이름에 `staging` 또는 `-stg` 포함 |
| `prd` | 위에 해당하지 않으면 운영으로 간주 |

#### (3) 조치상태 추론

| 조건 | 조치 | 이유 |
|------|------|------|
| 상태가 no_target / no_record / error | **삭제** | 사용되지 않거나 접속 불가 |
| 중단된 프로젝트 (nserise, bss, srt, ws2025) | **미정** | 전략적 판단 필요 |
| 정상 운영 중 | **유지** | 현재 사용 중 |

#### (4) 티키타카 (미분류 ALB 대화형 질의)

자동 분류에 실패한 ALB(`unknown`)는 터미널에서 사용자에게 직접 질문한다:
1. "이 ALB는 어떤 솔루션에 속하나요?" → 11개 선택지 제시
2. "조치를 어떻게 할까요?" → 유지/합치기/삭제/미정 선택

---

### 3.3 웹 UI (`serve`)

브라우저에서 `http://127.0.0.1:8080` 접속하여 사용한다.

#### 화면 구성

```
┌─────────────────────────────────────────────────────────┐
│  ALB Resource Cleanup Tool        [Collect] [Classify]  │
├─────────────────────────────────────────────────────────┤
│  Solution [▼]  Status [▼]  Action [▼]  Env [▼]  🔍검색 │
├────┬──────────┬──────────┬─────┬────────┬──────┬────────┤
│ #  │ ALB Name │ Solution │ Env │ Status │Action│Records │
├────┼──────────┼──────────┼─────┼────────┼──────┼────────┤
│ 1  │ alb-dev  │ signage  │ dev │●정상   │ 유지 │ 2개    │
│ 2  │ alb-old  │ cms      │ prd │●에러   │ 삭제 │ 0개    │
│ ...│          │          │     │        │      │        │
├────┴──────────┴──────────┴─────┴────────┴──────┴────────┤
│  Total: 48  Showing: 48  유지:30 합치기:8 삭제:8 미정:2 │
└─────────────────────────────────────────────────────────┘
```

#### 주요 기능

| 기능 | 조작 방법 | 설명 |
|------|----------|------|
| **데이터 수집** | [Collect] 버튼 클릭 | AWS 데이터 재수집 (1~2분) |
| **자동 분류** | [Classify] 버튼 클릭 | 수집된 데이터 자동 분류 |
| **필터링** | 상단 드롭다운/검색창 | 솔루션, 상태, 조치, 환경, 이름으로 필터 |
| **정렬** | 컬럼 헤더 클릭 | 오름차순/내림차순 토글 |
| **솔루션 변경** | Solution 셀 클릭 | 드롭다운으로 솔루션 변경 (자동 저장) |
| **조치 변경** | Action 셀 클릭 | 드롭다운으로 조치 변경 (자동 저장) |
| **합치기 계획** | ALB 행을 다른 행으로 드래그 | 합병 계획 모달 표시 |

#### 합치기 계획 (Merge Plan)

ALB 행을 드래그하여 다른 ALB 위에 놓으면 합병 계획 모달이 표시된다:
- 소스 ALB의 Route53 레코드 목록
- 소스 ALB의 Target Group 목록
- 실제 AWS 변경은 수행하지 않음 (계획 확인용)

#### 상태 표시 색상

| 상태 | 색상 |
|------|------|
| healthy (정상) | 초록 |
| unhealthy (비정상) | 노랑 |
| no_target (타겟없음) | 빨강 |
| no_record (레코드없음) | 회색 |
| error (에러) | 빨강 |

---

## 4. API 엔드포인트

웹 UI 외에 curl 등으로 직접 API를 호출할 수 있다.

### GET /api/entries
ALB 목록 조회. 쿼리 파라미터로 필터링 가능.

```bash
# 전체 조회
curl http://127.0.0.1:8080/api/entries

# 필터링
curl "http://127.0.0.1:8080/api/entries?solution=signage&status=healthy&environment=dev"
```

| 파라미터 | 설명 |
|---------|------|
| `solution` | 솔루션 필터 |
| `status` | 상태 필터 (healthy, unhealthy, no_target, no_record, error, unknown) |
| `action` | 조치 필터 (유지, 합치기, 삭제, 미정) |
| `environment` | 환경 필터 (dev, stg, prd) |
| `search` | ALB 이름 검색 (부분 일치) |

### GET /api/entries/{name}
특정 ALB 상세 조회.

### PATCH /api/entries/{name}
ALB 정보 수동 수정.

```bash
curl -X PATCH http://127.0.0.1:8080/api/entries/alb-dev-signage \
  -H 'Content-Type: application/json' \
  -d '{"solution":"cms", "action":"합치기", "mergeTarget":"main-alb", "note":"메모"}'
```

| 필드 | 설명 |
|------|------|
| `solution` | 솔루션 변경 |
| `action` | 조치 변경 (유지/합치기/삭제/미정) |
| `environment` | 환경 변경 (dev/stg/prd) |
| `mergeTarget` | 합치기 대상 ALB 이름 |
| `note` | 메모 |

### POST /api/collect
AWS 데이터 수집 실행. 1~2분 소요.

### POST /api/classify
자동 분류 실행. 수집이 먼저 되어 있어야 한다.

### GET /api/report
현황 요약 리포트 (솔루션별/상태별/조치별/환경별 카운트).

---

## 5. 데이터 저장

- 파일: `data/entries.json` (JSON 배열)
- 수집/분류/수정 시 자동 저장
- 서버 재시작해도 데이터 유지
- 원자적 쓰기 (임시파일 → 이름변경)로 파일 손상 방지

---

## 6. 일반적인 사용 시나리오

### 시나리오 1: 처음 사용
```
1. go run . -mode=all          ← 수집 + 분류 + 웹 시작
2. 브라우저에서 127.0.0.1:8080 접속
3. 자동 분류 결과 확인
4. 'unknown'인 ALB → Solution 셀 클릭하여 수동 지정
5. 삭제/합치기 대상 확인
```

### 시나리오 2: 주기적 현황 파악
```
1. go run . -mode=serve        ← 웹 시작
2. 브라우저에서 [Collect] 클릭  ← 최신 데이터 수집
3. [Classify] 클릭              ← 재분류
4. 필터로 '삭제' 조치 ALB 확인
```

### 시나리오 3: 합치기 계획 수립
```
1. 웹 UI에서 Action 필터를 '합치기'로 설정
2. 합칠 ALB 행을 대상 ALB 위로 드래그
3. 모달에서 이전할 레코드/TG 확인
4. 실제 작업은 AWS 콘솔에서 수동 수행
```

---

## 7. 제약사항 및 주의사항

| 항목 | 내용 |
|------|------|
| AWS 리전 | ap-northeast-2 고정 |
| 접속 | localhost만 (127.0.0.1:8080) |
| 인증 | 없음 (로컬 전용) |
| 실제 변경 | 수행하지 않음 (계획 수립까지만) |
| k8s ALB | 자동 제외 (k8s controller 관리) |
| AWS 자격증명 | awscli 설정 필요 (이미 연결되어 있으면 별도 작업 불필요) |
| 동시 사용 | 단일 사용자 전용 |
