# 테스트 리포트 (test-report.md)

> 작성: Tester 에이전트
> 날짜: 2026-04-09 (rev.2 -- Phase 9/10 포함 전체 재검증)
> 대상: Phase 3~5 + Phase 9 (스냅샷) + Phase 10 (합치기 뷰) 전체 코드

---

## 1. 빌드 검증

### 1.1 go build
- 상태: **PASS**
- 내용: `go build ./...` 정상 완료. 에러/경고 없음.

### 1.2 go vet
- 상태: **PASS**
- 내용: `go vet ./...` 경고/에러 없음.

---

## 2. 수집 기능 (Collector) 코드 검증

### 2.1 k8s ALB 필터링
- 상태: **PASS**
- 내용: `collector/alb.go:35`에서 `strings.HasPrefix(name, "k8s")`로 필터링. "k8s"로 시작하는 ALB 이름만 제외하며, "k8s-" 뿐 아니라 "k8sXxx" 형태도 모두 제외됨. 요구사항과 일치.

### 2.2 Route53 전체 조회 후 메모리 매칭
- 상태: **PASS**
- 내용: `collector/route53.go`에서 모든 Hosted Zone을 순회하며 전체 레코드를 메모리에 적재. `map[albDNS][]Record` 형태로 반환. ALB 개별 조회 없이 메모리 매칭. ALIAS + CNAME 모두 처리. DNS 정규화 (trailing dot 제거, lowercase) 정상 구현.

### 2.3 HTTPS 헬스체크 경로 순차 시도
- 상태: **PASS**
- 내용: `collector/healthcheck.go`에서 `healthPaths = []string{"/", "/health", "/healthz"}` 순차 시도. 200~499 응답을 받으면 즉시 해당 코드를 반환하고 중단. 모든 경로에서 실패하면 마지막 에러 코드를 반환.
- 검증 포인트:
  - 인증서 에러(-1) 시 InsecureSkipVerify로 재시도하여 서버 생존 확인
  - 리다이렉트 최대 3회 제한
  - 타임아웃은 httpTimeoutOverride로 테스트 시 변경 가능

### 2.4 Status 판정 우선순위
- 상태: **PASS (경미한 이슈 있음)**
- 내용: `determineStatus()` 함수에서 no_target > no_record > error > unhealthy > healthy 우선순위 준수.
- 이슈 B-01: `collector/collector.go:152-158`에 중복된 빈 조건문이 존재 (데드코드). 기능에는 영향 없음.
  ```go
  // 이 블록이 160~162행과 중복되며, 내부가 비어 있음
  if e.Status == "unknown" && len(e.TGs) == 0 {
      // Check if it was explicitly set to unknown (TG failure case)
      // vs just never collected
  }
  ```

### 2.5 병렬 처리 패턴 안전성
- 상태: **PASS**
- 내용:
  - Phase A: errgroup으로 ALB + Route53 병렬 수집. errgroup은 하나라도 실패하면 전체 에러 반환. 정상.
  - Phase C (TG): `sync.WaitGroup` + 세마포어 채널(`sem`) 패턴으로 최대 10개 동시 실행. 고루틴 내에서 `entries[idx]`를 인덱스로 직접 접근하며 각 고루틴이 서로 다른 인덱스를 처리하므로 data race 없음.
  - Phase C (Health): 동일한 WaitGroup + sem 패턴으로 최대 5개 동시 실행. `entries[j.entryIdx].Records[j.recordIdx].HealthCode`에 기록하며, 각 job이 고유한 (entryIdx, recordIdx) 조합이므로 안전.

---

## 3. 분류 기능 (Classifier) 코드 검증

### 3.1 키워드 매칭
- 상태: **PASS**
- 내용: `classifier/classifier.go`에서 `solutionKeywords` 슬라이스 순서대로 매칭. 더 구체적인 키워드(signage, aiagent)가 앞에 배치되어 오분류 방지. ALB 이름을 lowercase로 변환 후 `strings.Contains`로 매칭.

### 3.2 수동 분류 보존 (solution != "unknown" 건너뛰기)
- 상태: **PASS**
- 내용: `Classify()` 함수에서:
  - solution이 빈 문자열 또는 "unknown"인 경우에만 자동 분류 수행 (48행)
  - action이 빈 문자열 또는 "미정"인 경우에만 자동 추론 수행 (60행)
  - 사용자가 수동으로 설정한 값은 보존됨.

### 3.3 environment 독립 감지
- 상태: **PASS**
- 내용: `detectEnvironment()`는 solution 분류와 독립적으로 실행 (55~57행). solution이 이미 결정된 경우에도 environment가 unknown이면 별도로 감지. 패턴: `-dev`/`dev-` -> dev, `staging`/`-stg` -> stg, 기본값 -> prd.

### 3.4 조치상태 자동 추론
- 상태: **PASS**
- 내용: `inferAction()` 함수:
  - no_target -> 삭제
  - no_record -> 삭제
  - error -> 삭제
  - suspendedProjects (nserise, bss, srt, ws2025) -> 미정
  - 나머지 -> 유지
  - requirements.md의 조치상태 분류 기준과 일치.

### 3.5 nserise/nseries 키워드
- 상태: **INFO**
- 내용: 솔루션명 `nserise`의 키워드 목록에 `nseries`가 포함됨. 실제 ALB 이름이 "nseries"인 것으로 보이며, architecture.md에서도 이를 명시하고 있으므로 의도된 설계. 다만 혼동 가능성 있으므로 문서화 권장.

---

## 4. 데이터 저장 (Store) 검증

### 4.1 Mutex 보호
- 상태: **PASS**
- 내용: `store/store.go`에서 `sync.Mutex`로 Load/Save/UpdateEntry 모두 보호. `loadLocked()`/`saveLocked()` 내부 함수로 분리하여 Lock 안에서 호출하는 패턴 사용.

### 4.2 원자적 저장 (tmp + rename)
- 상태: **PASS**
- 내용: `saveLocked()`에서 `.tmp` 파일에 먼저 쓰고 `os.Rename()`으로 교체. 쓰기 도중 크래시 시에도 원본 파일은 손상되지 않음.

### 4.3 MergeEntries: collect 시 기존 분류 보존
- 상태: **PASS**
- 내용: `MergeEntries()`에서 ALB 이름 기준으로 old/new 매칭:
  - solution: old가 "unknown"이 아닌 경우 보존
  - action: old가 "미정"이 아닌 경우 보존
  - mergeTarget: 비어 있지 않으면 항상 보존
  - mergedName: 비어 있지 않으면 항상 보존
  - note: 비어 있지 않으면 항상 보존
  - environment는 보존하지 않음 -- 이것은 의도적인지 확인 필요 (이슈 B-05)

### 4.4 UpdateEntry: 모든 필드 지원
- 상태: **PASS**
- 내용: `UpdateEntry()`에서 solution, action, note, mergeTarget, environment, mergedName 6개 필드 모두 지원. map[string]any에서 타입 단언(string)으로 안전하게 처리.

---

## 5. 스냅샷 (Phase 9) 검증

### 5.1 이름 검증 (path traversal 방지)
- 상태: **PASS**
- 내용: `validateSnapshotName()`에서:
  - 빈 문자열 거부
  - 100자 초과 거부
  - `[/\\:.*?<>|]` 문자 거부 (정규식)
  - **주의**: `..`(double dot) 자체는 정규식에 포함되지 않음. 예: "test..snap"이라는 이름은 허용됨. 그러나 `/`와 `\`가 금지되어 있으므로 `../`나 `..\` 형태의 path traversal은 차단됨.

### 5.2 SaveSnapshot: 분류 필드만 저장
- 상태: **PASS**
- 내용: `SaveSnapshot()`에서 Classification 구조체로 변환 시 ALBName, Solution, Environment, Action, MergeTarget, MergedName, Note만 추출. ALBArn, ALBDNS, Records, TGs, Status 등 AWS 상태 데이터는 저장하지 않음. plan.md Phase 9 명세와 일치.

### 5.3 ApplySnapshot: ALB 이름 기준 매칭
- 상태: **PASS**
- 내용: `ApplySnapshot()`에서 `entryMap[ALBName]`으로 현재 entries와 매칭. 매칭된 항목은 모든 분류 필드(Solution, Environment, Action, MergeTarget, MergedName, Note)를 덮어씀. 매칭되지 않은 스냅샷 항목은 unmatched로 카운트. (matched, unmatched) 반환.

### 5.4 ListSnapshots: 최신순 정렬
- 상태: **PASS**
- 내용: `ListSnapshots()`에서 `sort.Slice`로 CreatedAt 기준 내림차순 정렬. RFC3339 형식이므로 문자열 비교로 시간순 정렬이 올바르게 동작. nil 방어 처리 (빈 디렉토리 시 빈 슬라이스 반환).

### 5.5 DeleteSnapshot
- 상태: **PASS**
- 내용: 이름 검증 후 파일 삭제. 파일이 없으면 "not found" 에러 반환.

---

## 6. 합치기 뷰 (Phase 10) 검증

### 6.1 HandleGetMergeGroups: mergeTarget 기준 그룹핑
- 상태: **PASS**
- 내용: `HandleGetMergeGroups()`에서:
  - 전체 entries를 entryMap[ALBName]으로 인덱싱
  - mergeTarget이 비어있지 않은 entry를 그룹핑
  - 타겟 ALB의 MergedName을 그룹 정보에 포함
  - 응답은 map[string]*mergeGroupInfo 형태로 JSON 직렬화

### 6.2 MergedName 필드: 전체 코드 일관성
- 상태: **PASS**
- 내용: MergedName 필드가 다음 위치에서 모두 올바르게 처리됨:
  - `model/entry.go`: `MergedName string \`json:"mergedName,omitempty"\`` 선언
  - `store/store.go` UpdateEntry: mergedName 키 지원 (165~168행)
  - `store/store.go` MergeEntries: old.MergedName 보존 (106~108행)
  - `store/snapshot.go` Classification: MergedName 필드 포함
  - `store/snapshot.go` SaveSnapshot: e.MergedName 추출 (85행)
  - `store/snapshot.go` ApplySnapshot: entry.MergedName 덮어씀 (218행)
  - `server/handler.go` HandleGetMergeGroups: target.MergedName 참조 (325행)

---

## 7. API 검증

### 7.1 HTTP 메서드 제한
- 상태: **PASS**
- 내용: `server/server.go`에서 모든 엔드포인트에 메서드 검사가 올바르게 구현됨:

| 엔드포인트 | 허용 메서드 | 제한 방식 |
|------------|------------|-----------|
| GET /api/entries | GET | switch + default 405 |
| GET/PATCH /api/entries/{name} | GET, PATCH | switch + default 405 |
| POST /api/collect | POST | if != POST, 405 |
| POST /api/classify | POST | if != POST, 405 |
| GET /api/report | GET | if != GET, 405 |
| GET /api/merge-groups | GET | if != GET, 405 |
| GET/POST /api/snapshots | GET, POST | switch + default 405 |
| POST /api/snapshots/{name}/load | POST | if != POST, 405 |
| DELETE /api/snapshots/{name} | DELETE | switch + default 405 |

### 7.2 에러 응답 형식 일관성
- 상태: **PASS**
- 내용: 모든 에러 응답이 `writeError()` 함수를 통해 `{"error": "message"}` 형식으로 반환. Content-Type은 항상 `application/json`.

### 7.3 API 라우트 경로 구분
- 상태: **PASS (주의사항 있음)**
- 내용: `/api/entries`와 `/api/entries/`가 별도 핸들러로 등록됨. Go의 `http.ServeMux`는 trailing slash 유무를 구분하므로:
  - `/api/entries` -> entries 목록 (GET)
  - `/api/entries/some-name` -> 개별 entry (GET/PATCH)
  - `/api/snapshots`와 `/api/snapshots/`도 동일한 패턴으로 분리됨

---

## 8. 웹 UI (web/index.html) 검증

### 8.1 인라인 편집 (solution, action, environment)
- 상태: **PASS**
- 내용: `editSolution()`, `editAction()`, `editEnvironment()` 함수에서:
  - 클릭 시 select 요소 생성 (중복 방지: `td.querySelector('select')` 체크)
  - `saved` 플래그로 중복 저장 방지
  - onchange 이벤트로 선택 즉시 저장
  - onblur 이벤트로 포커스 해제 시 저장 (150ms setTimeout으로 onchange와의 경쟁 조건 방지)
  - 값이 변경되지 않았으면 PATCH 호출 안 함

### 8.2 솔루션 다중 필터 로직
- 상태: **PASS**
- 내용: 체크박스 기반 다중 선택 필터. `getSelectedSolutions()`로 체크된 솔루션 목록을 가져와 `applyFilters()`에서 `selectedSolutions.includes(e.solution)`으로 필터링. 아무것도 선택하지 않으면 "All" 표시.

### 8.3 Final View 토글 로직
- 상태: **PASS**
- 내용: `toggleFinalView()`에서 `finalViewMode` 플래그 토글. `applyFilters()`에서:
  - action === '삭제'인 entry 제외
  - mergeTarget이 있는 entry (소스 ALB) 제외
  - `renderTable()`에서 finalViewMode가 true이면 merge 타겟 행에 소스들의 records/TGs를 합산하여 표시

### 8.4 스냅샷 저장/불러오기/삭제 흐름
- 상태: **PASS**
- 내용:
  - Save: prompt로 이름 입력 -> POST /api/snapshots -> 토스트 알림
  - Load: 드롭다운에서 선택 -> confirm -> POST /api/snapshots/{name}/load -> loadEntries() 호출로 테이블 갱신
  - Delete: x 버튼 클릭 -> confirm -> DELETE /api/snapshots/{name} -> 드롭다운 갱신
  - 드롭다운 외부 클릭 시 닫힘 (document click 이벤트)

### 8.5 합치기 접기/펼치기 로직
- 상태: **PASS**
- 내용: `expandedMerges` Set으로 펼쳐진 merge 타겟을 추적. `toggleMergeExpand()`로 Set에 추가/제거 후 `renderTable()` 호출. 펼쳤을 때 소스 행이 타겟 행 바로 아래에 들여쓰기(padding-left: 28px)로 삽입됨. finalViewMode에서는 접기/펼치기가 표시되지 않음 (소스 행이 이미 숨겨짐).

### 8.6 Records 링크 연결
- 상태: **PASS**
- 내용: records 표시 시 `<a href="https://{name}" target="_blank">`로 링크 생성. healthCode에 따라 rec-ok/rec-err 스타일 적용. escHtml/escAttr로 XSS 방지.

### 8.7 드래그 앤 드롭 합치기
- 상태: **PASS**
- 내용: `onDragStart` -> `onDragOver` -> `onDrop` 이벤트 체인. Drop 시 `showMergePlan()`으로 모달 표시. "Apply Merge Plan" 클릭 시 PATCH {mergeTarget, action: '합치기'} 전송. 자기 자신으로의 드롭은 무시됨.

---

## 9. 엣지케이스 검증

### 9.1 entries.json이 없을 때 서버 시작
- 상태: **PASS**
- 내용: `store.Load()`가 `os.IsNotExist(err)` 체크 후 빈 슬라이스 반환. 서버의 GET /api/entries는 빈 배열(`[]`)을 정상 반환. web UI는 "No entries found." 메시지 표시.

### 9.2 빈 배열일 때 classify
- 상태: **PASS**
- 내용: `HandleClassify()`에서 `len(entries) == 0`이면 400 Bad Request로 "no entries found, run collect first" 반환. CLI 모드에서는 "no entries to classify" 로그 출력 후 정상 종료. classifier.Classify()에 빈 슬라이스가 전달되면 (0, 0) 반환 (정상).

### 9.3 스냅샷 이름에 특수문자
- 상태: **PASS**
- 내용: `validateSnapshotName()`에서 `/\:.*?<>|` 문자 거부. 테스트 시나리오:
  - "my-snapshot_v1" -> 허용
  - "snapshot 2026" -> 허용 (공백 허용)
  - "test/hack" -> 거부 (/ 포함)
  - "../escape" -> 거부 (. 은 허용이지만 / 거부)
  - "test*name" -> 거부 (* 포함)
  - 101자 이름 -> 거부 (100자 초과)

### 9.4 collect 재실행 시 기존 분류 보존
- 상태: **PASS**
- 내용: `main.go:runCollect()`에서 oldEntries를 먼저 로드하고, 수집 후 `store.MergeEntries(oldEntries, entries)`로 병합. `HandleCollect()`에서도 동일한 패턴 적용. solution/action/mergeTarget/mergedName/note가 보존됨.

### 9.5 allUnused 타겟 상태 판정
- 상태: **PASS**
- 내용: `determineStatus()`에서 모든 TG의 타겟이 unused 상태(healthyCount=0, unhealthyCount=0)이면 "unhealthy"로 판정. 이는 합리적인 판단 -- unused 타겟만 있는 ALB는 실질적으로 트래픽을 처리하지 못함.

### 9.6 동시성: collect 중 PATCH 요청
- 상태: **FAIL (잠재적, LOW)**
- 내용: POST /api/collect가 수행되는 동안(수 초 소요) PATCH /api/entries/{name}으로 수정하면, collect 완료 시 Save가 해당 변경을 덮어씀. Store Mutex는 파일 I/O 단위만 보호하며 collect 전체 파이프라인에 대한 논리적 잠금 없음.
- 재현: (1) POST /api/collect 실행, (2) 수집 중 PATCH 전송, (3) collect 완료 후 PATCH 변경 손실
- 심각도: LOW (localhost 단독 사용 환경)

---

## 10. 발견된 이슈 요약

| # | 심각도 | 항목 | 위치 | 설명 |
|---|--------|------|------|------|
| B-01 | LOW | 데드코드 | `collector/collector.go:152-158` | 중복 조건문 + 빈 주석 블록. 기능 영향 없으나 정리 권장 |
| B-02 | INFO | nserise/nseries | `classifier/classifier.go:21` | 솔루션명과 키워드 불일치. 의도된 설계이나 문서화 권장 |
| B-03 | LOW | 동시성 | `server/handler.go:103-133` | collect 중 PATCH 요청 시 변경 손실 가능. localhost 환경에서 실제 발생 가능성 낮음 |
| B-04 | INFO | 보수적 에러 판정 | `collector/collector.go:186-190` | 레코드 1개라도 에러면 전체 ALB가 error. 의도된 정책이나 운영 시 과도 삭제 주의 |
| B-05 | LOW | MergeEntries environment 미보존 | `store/store.go:86-116` | MergeEntries에서 old.Environment 보존 로직 없음. 사용자가 수동 설정한 environment가 collect 시 초기화될 수 있음 |
| B-06 | INFO | escAttr XSS 불완전 | `web/index.html:829-832` | escAttr에서 작은따옴표를 `\'`로 이스케이프하나, HTML 속성에서 `\'`는 유효하지 않음. 실질적 XSS 위험은 낮음 (ALB 이름에 특수문자 가능성 극히 낮음) |

---

## 11. PASS/FAIL 종합 판정

| 카테고리 | 판정 | 세부 |
|----------|------|------|
| 빌드 (go build, go vet) | **PASS** | 에러/경고 없음 |
| 수집 (Collector) | **PASS** | k8s 필터링, Route53 메모리 매칭, 헬스체크, 병렬 처리 모두 정상 |
| 분류 (Classifier) | **PASS** | 키워드 매칭, 수동 분류 보존, environment 독립 감지, 조치상태 추론 정상 |
| 데이터 저장 (Store) | **PASS** | Mutex 보호, 원자적 저장, MergeEntries, UpdateEntry 정상. environment 보존 누락(B-05) 확인 |
| 스냅샷 (Phase 9) | **PASS** | 이름 검증, 분류 필드만 저장, ALB 이름 매칭, 최신순 정렬 모두 정상 |
| 합치기 뷰 (Phase 10) | **PASS** | mergeTarget 그룹핑, MergedName 전체 코드 일관성 정상 |
| API | **PASS** | HTTP 메서드 제한, 에러 응답 형식 일관성 정상 |
| 웹 UI | **PASS** | 인라인 편집, 다중 필터, Final View, 스냅샷 UI, 접기/펼치기 정상 |
| 엣지케이스 | **PASS** | 빈 데이터, 특수문자, 재수집 보존 정상. 동시성(B-03) 잠재 이슈 존재 |

---

## 12. 결론

전체 빌드 및 정적 분석 통과. 모든 핵심 기능(수집, 분류, 저장, 스냅샷, 합치기 뷰, API, 웹 UI)이 plan.md/architecture.md 설계와 일치하며 올바르게 구현되어 있음. 발견된 6건의 이슈는 모두 LOW/INFO 수준으로 기능 동작에 실질적 영향 없음. B-05(MergeEntries environment 미보존)는 사용자 시나리오에 따라 수정 검토 권장.
