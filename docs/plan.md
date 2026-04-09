# 작업 계획 (plan.md)

> 작성: PM 에이전트  
> 최종 갱신: 2026-04-09 (rev.3 — review-notes 액션 아이템 AI-04/AI-05/AI-08/AI-09 반영)

---

## 프로젝트 요약

AWS ALB 약 60개와 다수의 Route53 레코드를 수집/분류/정리하는 Go 기반 CLI + 웹 UI 도구를 개발한다.
- k8s ALB는 제외
- 솔루션(프로젝트)별로 ALB를 통합(합치기)하거나 삭제
- 분류가 어려운 항목은 사용자와 티키타카로 결정
- 웹 UI로 현황을 공유하고 드래그앤드롭으로 합칠 대상을 지정

---

## Phase 1: 기획 및 요구사항 확정

| # | 작업 | 담당 | 의존 | 우선순위 | 완료 기준 |
|---|------|------|------|----------|-----------|
| 1.1 | JOB.md 기반 요구사항 정리 및 티키타카 질의 | **Planner** | 없음 | P0 | `docs/requirements.md` 작성 완료, 사용자 확인 |
| 1.2 | 작업 분해 및 계획 수립 | **PM** | 1.1 | P0 | `docs/plan.md` 작성 완료 (본 문서) |

### 완료 기준
- requirements.md에 기능 목록, 솔루션 목록, 조치상태 분류 기준, 웹 UI 요구사항이 명확히 기술됨
- 사용자가 requirements.md 내용을 승인함

---

## Phase 2: 기술 설계

| # | 작업 | 담당 | 의존 | 우선순위 | 완료 기준 |
|---|------|------|------|----------|-----------|
| 2.1 | AWS 실데이터 조사 (ALB, Route53, TG 현황 파악) | **AWS Inspector** | 없음 | P0 | `docs/aws-report.md` 작성, ALB 목록/Route53 레코드/TG 현황 정리 |
| 2.2 | 기술 아키텍처 설계 | **Architect** | 1.1, 2.1 | P0 | `docs/architecture.md` 작성 완료 |
| 2.3 | 3인 상호 검토 (Planner-PM-Architect) | **Planner, PM, Architect** | 2.2 | P0 | `docs/review-notes.md` 작성, 충돌 해소 완료 |

### 설계에 포함될 핵심 항목
- **디렉토리 구조**: Go 프로젝트 레이아웃
- **데이터 모델**: ALB 엔트리 (albName, 솔루션, 이용상태, 조치상태 등)
- **수집 전략**: Route53 전체 조회 → 메모리 적재 → ALB 매칭 (개별 조회 금지)
- **병렬 처리**: Go goroutine 활용한 ALB/Route53 병렬 수집
- **파일 저장**: JSON 파일 기반 (`data/entries.json`)
- **웹 UI 구조**: HTML/JS 단일 파일, 드래그앤드롭 합치기 기능
- **API 설계**: 수집/분류/합치기/삭제 엔드포인트

### 완료 기준
- architecture.md에 디렉토리 구조, 데이터 모델, API, 시퀀스 다이어그램 포함
- review-notes.md에 3인 검토 피드백 기록, 미해결 충돌 없음

---

## Phase 3: 핵심 기능 구현 — 수집(Collector)

| # | 작업 | 담당 | 의존 | 우선순위 | 완료 기준 |
|---|------|------|------|----------|-----------|
| 3.1 | Go 프로젝트 초기 구조 생성 (go mod, 디렉토리, 데이터 모델) | **Developer** | 2.3 | P0 | `go build` 성공 + `model/entry.go` 데이터 모델 확정 (Entry, Record, TGInfo, MergePlan 구조체 정의 완료) |
| 3.2 | ALB 수집기 구현 (ELBv2 API, k8s 필터링) | **Developer** | 3.1 | P0 | k8s 제외한 ALB 목록 수집 가능 |
| 3.3 | Route53 레코드 수집 및 메모리 매칭 | **Developer** | 3.1 | P0 | 전체 레코드 조회 → ALB DNS와 매칭 |
| 3.4 | Target Group 상태 수집 (healthy/unhealthy/없음) | **Developer** | 3.1 | P0 | TG 상태 확인 |
| 3.5 | Route53 도메인 HTTPS 헬스체크 (에러 판정) | **Developer** | 3.3 | P0 | Route53 레코드 도메인으로 HTTPS(443) 요청. 경로: `/`, `/health`, `/healthz` 순차 시도. 모든 경로에서 500대 응답 또는 인증서 에러 시 status=error로 판정 |
| 3.6 | 병렬 수집 통합 및 JSON 저장 | **Developer** | 3.2, 3.3, 3.4, 3.5 | P0 | `data/entries.json`에 전체 결과 저장 |

### 완료 기준
- `go run .` 실행 시 ALB/Route53/TG 데이터를 병렬 수집하여 JSON으로 저장
- k8s ALB 제외 확인
- Route53 개별 조회 없이 메모리 매칭 확인
- Route53 레코드 도메인으로 HTTPS 헬스체크 수행, 인증서 에러 포함한 에러 판정 확인

---

## Phase 4: 핵심 기능 구현 — 분류(Classifier) + 티키타카

| # | 작업 | 담당 | 의존 | 우선순위 | 완료 기준 |
|---|------|------|------|----------|-----------|
| 4.1 | 자동 분류 로직 (ALB 이름 → 솔루션/환경 매핑) | **Developer** | 3.6 | P1 | ALB 네이밍 패턴 `{aws리소스}-{환경}-{솔루션}-{서비스선택}` 파싱으로 솔루션+환경 자동 배정 |
| 4.2 | 조치상태 자동 추론 (타겟 없음→삭제, 레코드 없음→삭제, 에러→삭제후보 등) | **Developer** | 3.6 | P1 | 규칙 기반 자동 분류. status=error(500+인증서)인 경우도 삭제 후보로 분류 |
| 4.3 | 패턴 불일치 ALB 티키타카 질의 | **Developer** | 4.1, 4.2 | P1 | 네이밍 패턴 불일치 ALB에 대해 솔루션/환경/조치상태 선택지 제시, 사용자 입력 수용 |
| 4.4 | 분류 결과 JSON 업데이트 | **Developer** | 4.3 | P1 | entries.json에 솔루션/환경/조치상태 반영 |

### 완료 기준
- 자동 분류 가능한 항목은 자동 배정
- 미분류 항목은 선택지 형태로 질문
- 결과가 entries.json에 반영

---

## Phase 5: 웹 UI 구현

| # | 작업 | 담당 | 의존 | 우선순위 | 완료 기준 |
|---|------|------|------|----------|-----------|
| 5.1 | 웹서버 및 API 엔드포인트 구현 (localhost only) | **Developer** | 3.6 | P1 | HTTP 서버 localhost 바인딩, JSON API 응답 |
| 5.2 | ALB 현황 테이블 뷰 (조회/필터링) | **Developer** | 5.1 | P1 | 솔루션별/상태별 필터링 가능한 테이블. 단일 HTML 파일(`web/index.html`) 내 구현 |
| 5.3 | 드래그앤드롭 합치기 UI | **Developer** | 5.2 | P2 | 합칠 대상 ALB에 드래그하면 관련 정보(레코드, TG) 표시. 단일 HTML 파일 내 구현 |
| 5.4 | 합치기 실행 시 필요 정보 출력 (레코드 변경 사항, ALB 설정 등) | **Developer** | 5.3 | P2 | 합치기 대상의 레코드/리스너 룰 이동 정보 표시. 단일 HTML 파일 내 구현. 실제 AWS 리소스 변경은 범위 밖 |

### 완료 기준
- 브라우저에서 ALB 현황 테이블 확인 가능
- 솔루션별, 조치상태별 필터링 동작
- 드래그앤드롭으로 합칠 대상 지정 시 관련 정보 표시

---

## Phase 6: 빌드 검증 및 테스트

| # | 작업 | 담당 | 의존 | 우선순위 | 완료 기준 |
|---|------|------|------|----------|-----------|
| 6.1 | 빌드 검증 (go build, go vet, 린트) | **Tester** | Phase 3~5 | P1 | 빌드 성공, vet 경고 없음 |
| 6.2 | 수집 기능 테스트 (실제 AWS 데이터) | **Tester** | 6.1 | P1 | 실제 ALB/Route53 수집 정상 동작 |
| 6.3 | 분류 로직 테스트 (자동 분류 정확도) | **Tester** | 6.1 | P1 | 알려진 ALB 이름에 대해 올바른 솔루션 배정 |
| 6.4 | 웹 UI 기능 테스트 | **Tester** | 6.1 | P1 | 테이블 표시, 필터링, 드래그앤드롭 동작 확인 |
| 6.5 | 엣지케이스 테스트 (타겟 없는 ALB, 레코드 없는 ALB, HTTPS 헬스체크 등) | **Tester** | 6.2, 6.3 | P1 | 엣지케이스 정상 처리 확인. HTTPS 헬스체크 테스트 케이스 포함: (1) 인증서 에러(healthCode=-1) 판정, (2) 타임아웃(healthCode=0) 판정, (3) 400대 응답(401/403/404)이 정상으로 판정되는지 확인, (4) DNS 해석 실패(healthCode=-2) 판정 |

### 완료 기준
- `docs/test-report.md` 작성 완료
- PASS/FAIL 결과 명시, FAIL 항목에 대한 원인과 조치 방안 기록

---

## Phase 7: 코드 리뷰 및 최종 점검

| # | 작업 | 담당 | 의존 | 우선순위 | 완료 기준 |
|---|------|------|------|----------|-----------|
| 7.1 | 설계 준수 여부 확인 (architecture.md 대비) | **Reviewer** | Phase 6 | P1 | 설계-구현 불일치 없음 |
| 7.2 | 코드 품질 점검 (에러 핸들링, 병렬 안전성) | **Reviewer** | Phase 6 | P1 | CRITICAL 이슈 없음 |
| 7.3 | 보안 점검 (AWS 자격증명 노출, 입력값 검증) | **Reviewer** | Phase 6 | P1 | 보안 취약점 없음 |
| 7.4 | 리팩토링/개선 제안 | **Reviewer** | 7.1~7.3 | P2 | 개선 사항 목록화 |

### 완료 기준
- `docs/review-report.md` 작성 완료
- CRITICAL 이슈 0건, WARNING 이슈 조치 완료

---

## Phase 8: 운영 — 실제 분류 실행

| # | 작업 | 담당 | 의존 | 우선순위 | 완료 기준 |
|---|------|------|------|----------|-----------|
| 8.1 | `/collect`로 최신 AWS 데이터 수집 | **AWS Inspector** | Phase 7 | P1 | 최신 데이터 수집 완료 |
| 8.2 | `POST /api/classify`로 자동 분류 실행 + 웹 UI `PATCH /api/entries/:name`으로 미분류(unknown) 항목 개별 수정 (티키타카 포함) | **Planner** | 8.1 | P1 | 전체 ALB 분류 완료 |
| 8.3 | `GET /api/report`로 현황 리포트 조회 | **PM** | 8.2 | P1 | 솔루션별/상태별 요약 리포트 |
| 8.4 | 웹 UI로 결과 공유 및 합치기 대상 확정 | **PM** | 8.3 | P1 | 관계자 검토 완료 |

### 완료 기준
- 전체 ALB에 대해 솔루션과 조치상태가 확정됨
- 중지 프로젝트(nserise, bss, srt, ws2025)는 전략팀 컨펌 반영
- 중지예상 프로젝트(wine, ncount, dooh)는 개발팀 컨펌 반영

---

## 의존관계 다이어그램

```
Phase 1 (기획)
  └─→ Phase 2 (설계) ← Phase 2.1 (AWS 조사, 병렬 가능)
        └─→ Phase 3 (수집 구현)
              ├─→ Phase 4 (분류 구현)
              └─→ Phase 5 (웹 UI 구현)
                    └─→ Phase 6 (테스트)
                          └─→ Phase 7 (코드 리뷰)
                                └─→ Phase 8 (운영)
```

---

## 우선순위 정의

| 등급 | 의미 |
|------|------|
| P0 | 필수. 이것 없이 다음 단계 진행 불가 |
| P1 | 핵심 기능. 프로젝트 목표 달성에 필수 |
| P2 | 부가 기능. 없어도 최소 동작 가능 |

---

## 병렬 수행 가능 작업

- **Phase 2.1** (AWS 조사)은 Phase 1 완료 없이 바로 시작 가능 (실데이터 파악은 독립적)
- **Phase 3.2, 3.3, 3.4** (ALB/Route53/TG 수집)은 3.1 이후 병렬 구현 가능
- **Phase 3.5** (HTTPS 헬스체크)는 3.3 완료 후 수행 (Route53 레코드 도메인이 필요하므로)
- **Phase 4** (분류)와 **Phase 5** (웹 UI)는 Phase 3 이후 병렬 구현 가능

---

## 리스크 및 블로커

| 리스크 | 영향 | 대응 |
|--------|------|------|
| 중지 프로젝트 컨펌 지연 (전략팀/개발팀) | 분류 확정 불가 | 컨펌 전까지 "미확정" 상태로 유지, 도구 개발은 계속 진행 |
| ALB 60개 + 레코드 다수로 API rate limit | 수집 실패 | 병렬도 조절, 재시도 로직 추가 |
| 솔루션 자동 분류 정확도 낮음 | 티키타카 질문 과다 | ALB 이름 패턴 분석 후 매핑 규칙 강화 |
| 합치기 실행 시 서비스 영향 | 장애 위험 | 본 도구는 "계획/시각화"까지만 담당, 실제 실행은 별도 절차 |

---

## Phase 9: 스냅샷 저장/불러오기 기능

> 추가: 2026-04-09 -- 사용자 요청에 의한 신규 기능

### 배경
- 현재 MergeEntries로 collect 시 분류를 보존하지만, 명시적 저장/복원이 불가
- 사용자가 수동 분류 작업을 "이름 붙여 저장"하고 나중에 "불러오기"할 수 있어야 함
- collect/classify를 자유롭게 실행해도 이전 분류를 이름으로 복원 가능

### 작업 분해

| # | 작업 | 담당 | 의존 | 우선순위 | 완료 기준 |
|---|------|------|------|----------|-----------|
| 9.1 | `store/snapshot.go` 신규 작성 (Snapshot 타입, Save/List/Load/Delete/Apply 함수) | **Developer** | 없음 | P0 | 스냅샷 CRUD + Apply 동작 |
| 9.2 | `server/handler.go` 핸들러 4개 추가 + `server/server.go` 라우트 추가 | **Developer** | 9.1 | P0 | API 4개 정상 동작 |
| 9.3 | `web/index.html` 스냅샷 UI (Save/Load/Delete 버튼+드롭다운) | **Developer** | 9.2 | P0 | 웹에서 저장/불러오기/삭제 가능 |
| 9.4 | 빌드 검증 | **Developer** | 9.3 | P0 | go build 성공, 수동 테스트 |

### 구현 명세

#### 파일 구조
```
data/
  entries.json
  snapshots/
    {name}.json
```

#### 스냅샷 데이터 모델 (`store/snapshot.go`)
```go
type Classification struct {
    ALBName     string `json:"albName"`
    Solution    string `json:"solution"`
    Environment string `json:"environment"`
    Action      string `json:"action"`
    MergeTarget string `json:"mergeTarget"`
    Note        string `json:"note"`
}

type Snapshot struct {
    Name            string           `json:"name"`
    CreatedAt       string           `json:"createdAt"`
    EntryCount      int              `json:"entryCount"`
    Classifications []Classification `json:"classifications"`
}

type SnapshotMeta struct {
    Name       string `json:"name"`
    CreatedAt  string `json:"createdAt"`
    EntryCount int    `json:"entryCount"`
}
```

- 스냅샷에는 **분류 필드만** 저장 (solution, environment, action, mergeTarget, note)
- AWS 상태 데이터(arn, dns, records, targetGroups, status)는 저장하지 않음

#### 함수 목록
- `SaveSnapshot(dataDir, entries, name)` -- 분류 필드 추출 후 `data/snapshots/{name}.json` 저장
- `ListSnapshots(dataDir)` -- snapshots 디렉토리의 메타 목록 (최신순)
- `LoadSnapshot(dataDir, name)` -- 파일 읽기
- `DeleteSnapshot(dataDir, name)` -- 파일 삭제
- `ApplySnapshot(snapshot, entries) (matched, unmatched)` -- ALB 이름 기준 매칭하여 분류 덮어쓰기

#### API 설계
| Method | Path | 설명 |
|--------|------|------|
| POST | `/api/snapshots` | 저장. body: `{"name": "..."}` |
| GET | `/api/snapshots` | 목록 조회 |
| POST | `/api/snapshots/{name}/load` | 불러오기 (현재 entries에 적용) |
| DELETE | `/api/snapshots/{name}` | 삭제 |

#### 웹 UI
- 헤더에 [Save Snapshot] 버튼, [Load Snapshot] 드롭다운 추가
- Save: prompt로 이름 입력 -> POST /api/snapshots
- Load: 드롭다운에서 선택 -> confirm -> POST /api/snapshots/{name}/load -> 테이블 갱신
- Delete: 드롭다운 항목에 x 버튼 -> DELETE /api/snapshots/{name}

#### 이름 검증
- 허용: 한글, 영문, 숫자, `-`, `_`, `.`, 공백
- 금지: `/`, `\`, `..`, `:`, `*`, `?`, `<`, `>`, `|`
- 최대 길이: 100자

#### MergeEntries와의 관계
- MergeEntries는 **그대로 유지** (collect 시 자동 보존)
- Snapshot은 **명시적 저장/복원** (사용자가 원할 때)
- 역할이 다르므로 공존

---

## Phase 10: 합치기 뷰 및 시나리오 관리

> 추가: 2026-04-09 -- 사용자 요청에 의한 신규 기능
> 3인 논의(Planner/PM/Architect) 완료

### 배경
- 여러 사람이 합치기/삭제 시나리오를 시도해보고 결과를 확인하고 싶음
- 합쳐진 ALB를 "최종 뷰"로 보여줘야 함 (3개->1개로 축소된 상태)
- 합쳐진 것을 펼치거나 해제하는 기능 필요
- 초기화/되돌리기는 Phase 9 스냅샷으로 이미 해결됨 (추가 구현 불필요)

### 작업 분해

| # | 작업 | 담당 | 의존 | 우선순위 | 완료 기준 |
|---|------|------|------|----------|-----------|
| 10.1 | `model/entry.go`에 MergedName 필드 추가 | **Developer** | 없음 | P0 | 빌드 성공 |
| 10.2 | `store/store.go` UpdateEntry에 mergedName 지원 추가 | **Developer** | 10.1 | P0 | PATCH로 mergedName 변경 가능 |
| 10.3 | `store/snapshot.go` Classification에 MergedName 추가 | **Developer** | 10.1 | P0 | 스냅샷에 mergedName 포함 |
| 10.4 | `server/handler.go`에 GET /api/merge-groups 핸들러 추가 | **Developer** | 10.1 | P0 | 합치기 그룹 조회 API 동작 |
| 10.5 | `server/server.go`에 라우트 추가 | **Developer** | 10.4 | P0 | 라우트 연결 |
| 10.6 | `web/index.html` 최종 뷰 토글 + Merge Info 컬럼 + 접기/펼치기 + 이름 변경 + 합치기 해제 | **Developer** | 10.2, 10.4 | P0 | 웹에서 전체 기능 동작 |
| 10.7 | 빌드 검증 | **Developer** | 10.6 | P0 | go build 성공, 수동 테스트 |

### 구현 명세

#### 10.1 model/entry.go
```go
MergedName string `json:"mergedName,omitempty"` // 합쳐진 후 표시 이름 (타겟 ALB에만 설정)
```

#### 10.2 store/store.go - UpdateEntry 확장
```go
if v, ok := updates["mergedName"]; ok {
    if s, ok := v.(string); ok {
        found.MergedName = s
    }
}
```

#### 10.3 store/snapshot.go - Classification 확장
```go
type Classification struct {
    // ... 기존 필드 ...
    MergedName string `json:"mergedName"`
}
```
SaveSnapshot, ApplySnapshot에서도 MergedName 처리.

#### 10.4 server/handler.go - HandleGetMergeGroups
entries 순회하여 mergeTarget 기준 그룹핑. 응답:
```json
{
  "target-alb": {
    "mergedName": "new-name",
    "sources": ["src1", "src2"],
    "sourceRecords": [...],
    "sourceTGs": [...]
  }
}
```

#### 10.5 server/server.go
```go
mux.HandleFunc("GET /api/merge-groups", h.HandleGetMergeGroups)
```

#### 10.6 web/index.html 변경 상세

**A. 필터 영역 - "Final View" 토글 버튼 추가**
- `let finalViewMode = false;`
- 토글 시 applyFilters에서: 삭제/합치기(소스) 행 숨김
- 합치기 타겟 행에는 소스들의 records/TGs 합산 표시

**B. 테이블 - "Merge Info" 컬럼 추가 (TGs 다음)**
- 소스 ALB: "-> {target}" 표시
- 타겟 ALB: "{N}개 합침 [펼치기]" 표시 + 클릭 시 하위에 소스 행 삽입
- mergedName 있으면: "(새이름: {mergedName})" + 변경 아이콘

**C. 접기/펼치기**
- JS에서 mergeTarget 역방향 맵 계산 (allEntries 기반)
- 타겟 행 클릭 시 소스 행을 바로 아래에 들여쓰기로 삽입 (배경색: #f0f8ff)
- 소스 행에 "합치기 해제" 버튼: PATCH { action: "미정", mergeTarget: "" }

**D. 합쳐진 ALB 이름 변경**
- Merge Info 영역 더블클릭 -> prompt -> PATCH { mergedName: "..." }

**E. 최종 뷰 상세 로직**
```javascript
if (finalViewMode) {
  // 1. 삭제 entry 제외
  // 2. mergeTarget이 있는 entry (소스) 제외
  // 3. 남은 entry 중 합치기 타겟은 소스들의 records/TGs를 합산하여 표시
}
```

#### MergeEntries/스냅샷과의 관계
- MergeEntries: 그대로 유지
- 스냅샷: MergedName 필드 추가로 확장
- 초기화: 스냅샷 불러오기로 해결 (별도 구현 불필요)
- 여러 사람 시나리오: 각자 다른 이름으로 스냅샷 저장 (별도 구현 불필요)

---

## 현재 진행 상황 (2026-04-09)

| Phase | 상태 | 비고 |
|-------|------|------|
| Phase 1 (기획) | **완료** | requirements.md 최종 확정, plan.md rev.3 반영 |
| Phase 2 (설계) | **진행중** | 2.1 AWS 조사 완료, 2.2 아키텍처 설계 완료(rev.2), 2.3 상호 검토 진행중 (AI-01~03 Planner 미완료) |
| Phase 3~8 | 미착수 | Phase 2.3 완료 후 진행 |
| **Phase 9 (스냅샷)** | **완료** | store/snapshot.go, handler, UI 모두 구현됨 |
| **Phase 10 (합치기 뷰)** | **즉시 착수 가능** | 기존 코드 위에 독립적으로 추가 가능 |

## 다음 액션

> Phase 1 완료. 아래 두 작업을 **병렬**로 즉시 시작한다.

1. **병렬 A** — `@aws-inspector`: AWS 실데이터 조사 (Phase 2.1)
   - ALB 목록, Route53 Hosted Zone/레코드 현황, Target Group 상태 파악
   - 산출물: `docs/aws-report.md`
2. **병렬 B** — `@architect`: 기술 아키텍처 설계 (Phase 2.2)
   - requirements.md 기반으로 디렉토리 구조, 데이터 모델, API, 시퀀스 다이어그램 설계
   - 특히 반영 필수 사항:
     - HTTPS 헬스체크 모듈 설계 (Route53 도메인 -> HTTPS 요청, 경로 `/` `/health` `/healthz`, 인증서 에러 포함)
     - ALB 네이밍 패턴 파서 (`{aws리소스}-{환경}-{솔루션}-{서비스선택}`) + 패턴 불일치 시 fallback 로직
     - 웹 UI localhost 바인딩
     - status 필드: healthy | unhealthy | no_target | no_record | error
   - 산출물: `docs/architecture.md`
3. **병렬 A+B 완료 후** — `@planner` + `@pm` + `@architect`: 3인 상호 검토 (Phase 2.3)
   - aws-report.md의 실데이터를 architecture.md에 반영
   - review-notes.md에 피드백 기록
4. **Phase 2.3 완료 후** — `@developer`: Phase 3 구현 착수
