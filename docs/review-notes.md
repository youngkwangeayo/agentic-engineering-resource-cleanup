# 상호 검토 노트 (review-notes.md)

> 이 문서는 Planner, PM, Architect 간 상호 검토 피드백을 기록합니다.

---

## Phase 1 검토 (2026-04-09)

### plan.md에 대한 Planner 피드백

[Planner] plan.md는 JOB.md의 핵심 요구사항을 전반적으로 잘 반영하고 있음. 아래 몇 가지 보완/확인 사항을 제안함.

[Planner] **P1 — 500 에러 확인 방식 미정의**: JOB.md에 "500대 에러를 뱉는 프로젝트는 공지 후 삭제"라고 되어 있으나, plan.md Phase 3.4에서 "500 에러 감지"를 어떻게 수행하는지 구체적 방법이 명시되어 있지 않음. HTTP 직접 요청 vs CloudWatch 메트릭 등 방식을 architecture.md에서 명확히 정의해야 함.

[Planner] **P2 — 환경(environment) 분류 기준 부재**: JOB.md에 "환경-솔루션 엮어서" ALB를 합친다고 되어 있으나, 환경 분류 기준(dev/stg/prd)과 ALB 이름에서 환경을 식별하는 패턴이 정의되지 않음. Phase 4.1 자동 분류 시 환경 식별도 포함되어야 함.

[Planner] **P3 — 합치기 실행 범위 명확화**: plan.md Phase 5.4에서 "합치기 실행 시 필요 정보 출력"이라고 되어 있는데, 이것이 "계획 표시"인지 "실제 AWS 리소스 변경"인지 명확히 해야 함. JOB.md 맥락상 본 도구는 계획/시각화까지만 담당하는 것으로 보이며, 이를 plan.md에 명시적으로 기술하면 좋겠음. (리스크 테이블에는 언급되어 있으나 Phase 5.4 작업 설명에는 없음)

[Planner] **P4 — 티키타카가 CLI와 웹 UI 중 어디서 수행되는지**: Phase 4.3에서는 "CLI 인터페이스"로, Phase 8에서는 실제 분류를 수행하는데, 웹 UI에서도 분류 수정이 가능해야 하는지 확인 필요. 웹 UI 테이블에서 직접 솔루션/조치상태를 편집할 수 있으면 더 효율적일 수 있음.

[Planner] **긍정적 평가**: Phase 분리, 의존관계 다이어그램, 병렬 수행 가능 작업 정리, 리스크 테이블 모두 잘 구성되어 있음. 특히 "실제 실행은 별도 절차"라는 리스크 대응이 적절함.

### Planner 피드백에 대한 PM 의견

[PM] **P1 — 500 에러 확인 방식**: 동의. requirements.md 최종 확정으로 해결됨. plan.md rev.2에 Phase 3.5 "Route53 도메인 HTTPS 헬스체크" 작업을 신규 추가하여 반영 완료. 구체적 방법: Route53 레코드 도메인으로 HTTPS(443) 요청, 경로 `/` `/health` `/healthz` 순차 시도, 인증서 에러도 error로 통합. Architect가 architecture.md에 헬스체크 모듈 설계를 포함해야 함.

[PM] **P2 — 환경 분류 기준**: 동의. requirements.md에서 ALB 네이밍 패턴 `{aws리소스}-{환경}-{솔루션}-{서비스선택}`이 확정됨. plan.md Phase 4.1을 "솔루션/환경 매핑"으로 수정하고, 패턴 파싱 로직을 명시함. 단, 패턴 불일치 ALB도 존재하므로 Phase 4.3에서 티키타카 질의로 대응.

[PM] **P3 — 합치기 실행 범위**: 동의. requirements.md FR-13에 "계획/시각화까지만 담당, 실제 AWS 리소스 변경은 별도 절차"로 명확히 기술됨. plan.md Phase 5.4 설명도 이미 "필요 정보 출력"으로 되어 있어 혼란 소지가 적으나, 추가적으로 완료 기준에도 이를 반영해둠.

[PM] **P4 — 티키타카 수행 위치 (CLI vs 웹 UI)**: 좋은 지적. 현재 plan.md에서는 Phase 4.3이 CLI 기반 티키타카이고, 웹 UI(Phase 5)에서는 조회/필터/드래그앤드롭만 담당. 웹 UI에서 솔루션/조치상태 인라인 편집 기능을 추가하면 운영 Phase 8에서 편리할 것으로 동의함. 다만 이는 P2 우선순위로 분류하여, 핵심 기능 구현 후 여력이 있을 때 추가하는 것을 제안. Architect 설계 시 이를 고려한 API 설계(PATCH 엔드포인트 등)를 포함하도록 요청.

[PM] **긍정적 평가에 대한 감사**: Planner의 체계적 피드백 덕분에 plan.md를 보강할 수 있었음. requirements.md 최종 확정이 매우 충실하여 설계 단계로 넘어갈 준비가 됨.

---

## Phase 2 검토 (2026-04-09)

### architecture.md 작성 완료 — Architect 코멘트

[Architect] architecture.md 초안을 작성했습니다. Planner P1~P4 피드백을 모두 반영했습니다.

### requirements.md에 대한 Architect 피드백

[Architect] **A1 — Record 모델 확장**: requirements.md 데이터 모델의 `records`가 string 배열인데, 헬스체크 결과(healthCode)와 Zone 정보를 함께 저장해야 합니다. architecture.md에서는 Record를 구조체로 확장했습니다. requirements.md 데이터 모델 업데이트가 필요합니다.

[Architect] **A2 — HTTPS 헬스체크 400대 응답 판정**: requirements.md FR-08에 "에러 응답(500대, 인증서 에러 등)"이라 되어 있는데, 400대 응답(401, 403, 404 등)은 서버가 살아있는 것이므로 "정상"으로 분류하는 것이 맞는지 확인 필요. architecture.md에서는 200~499=정상, 500+=에러로 설계했습니다.

[Architect] **A3 — TG 조회 부분 실패 허용**: 특정 ALB의 TG 조회가 실패해도 전체 수집이 중단되면 안 됩니다. 해당 ALB만 status="unknown"으로 마킹하고 나머지는 계속 수집하도록 설계했습니다.

### plan.md에 대한 Architect 피드백

[Architect] **A4 — Phase 4/5 병렬 구현 시 선행 조건**: 분류와 웹 UI를 Phase 3 이후 병렬 구현 가능하다는 의견에 동의하나, 양쪽 모두 model/entry.go의 데이터 모델에 의존하므로, 모델 정의를 Phase 3.1에서 확정해야 합니다.

[Architect] **A5 — 외부 프레임워크 불필요**: 임시 도구이므로 net/http 표준 라이브러리 + golang.org/x/sync(errgroup)만 사용합니다. 웹 프레임워크(gin, echo 등)는 과잉입니다.

### Planner, PM 검토 요청 사항

[Architect] 다음 항목에 대해 3인 회의에서 합의가 필요합니다:
1. **A2**: 400대 응답을 "정상"으로 볼 것인지 (Architect 제안: 정상)
2. **A1**: requirements.md 데이터 모델을 architecture.md 기준으로 업데이트할지
3. **P4 후속**: 웹 UI 인라인 편집을 P2로 분류하되, API는 미리 설계에 포함 — 동의 여부

---

## Review Cycle - Architect 검토 (2026-04-09)

### requirements.md 검토

- [충족]: **FR-01~FR-06 (수집)** — architecture.md 3장(Collector)에서 ALB 수집, k8s 필터링, Route53 전체 조회 후 메모리 매칭, TG 상태 수집, goroutine 병렬 수집, JSON 저장을 모두 설계에 반영함. 특히 FR-03(개별 조회 금지)은 3.2 핵심 설계 결정 테이블에서 명시적으로 근거를 제시함.
- [충족]: **FR-07~FR-10 (분류)** — architecture.md 4장(Classifier)에서 ALB 이름 파싱, 조치상태 자동 추론 규칙, 티키타카 질의 흐름, JSON 업데이트를 모두 반영함. 네이밍 패턴 불일치 시 fallback으로 티키타카 질의를 수행하는 설계가 FR-09와 정확히 일치.
- [충족]: **FR-11~FR-13 (합치기)** — architecture.md 5장에서 MergePlan 모델, 합치기 계획 생성 로직, "계획/시각화까지만 담당" 원칙을 명시. FR-12의 레코드/리스너/TG 정보 표시도 MergePlan 구조체에 포함.
- [충족]: **FR-14~FR-19 (웹 UI)** — architecture.md 7장에서 테이블 뷰, 필터링, 드래그앤드롭, 단일 HTML/JS, 파일 기반 저장을 모두 반영. localhost 바인딩은 server.go에서 처리 예정.
- [충족]: **TC-01~TC-05 (기술 제약)** — Go 언어, AWS SDK v2, rate limit 대응(세마포어+재시도), 파일 기반 JSON 저장 모두 반영.
- [보완필요]: **데이터 모델 불일치** — requirements.md의 `records`와 `targetGroups`가 `string` 배열인데, architecture.md에서는 `Record` 구조체와 `TGInfo` 구조체로 확장함. 이는 헬스체크 결과(healthCode)와 TG 상세 정보를 저장하기 위해 필수적인 확장이므로, requirements.md 데이터 모델을 architecture.md 기준으로 업데이트해야 함 (기존 피드백 A1 재확인).
- [보완필요]: **status "unknown" 누락** — architecture.md 에러 처리(9장)에서 TG 조회 실패 시 `status="unknown"`으로 마킹하는데, requirements.md의 status 필드에는 `healthy | unhealthy | no_target | no_record | error`만 있고 `unknown`이 없음. `unknown`을 추가하거나, 해당 케이스를 `error`로 통합할지 합의 필요.

### plan.md 검토

- [구현가능]: **Phase 분리 및 의존관계** — Phase 3(수집) -> Phase 4(분류) + Phase 5(웹 UI) 병렬 -> Phase 6(테스트)의 흐름이 architecture.md의 설계와 정확히 일치함. 특히 Phase 3.5(HTTPS 헬스체크)가 3.3(Route53) 의존으로 설정된 것이 적절.
- [구현가능]: **Phase 3 세부 작업 분해** — 3.1(프로젝트 초기화) -> 3.2/3.3/3.4 병렬 -> 3.5 -> 3.6 통합의 흐름이 architecture.md 8장의 병렬 처리 설계(Phase 1: ALB+Route53 병렬, Phase 2: 매칭, Phase 3: TG+헬스체크 워커풀)와 부합.
- [구현가능]: **Phase 5 웹 UI 구현 순서** — 5.1(서버) -> 5.2(테이블) -> 5.3(드래그앤드롭) -> 5.4(합치기 정보 표시)가 점진적이며 합리적. architecture.md API 6장의 엔드포인트가 이 순서대로 구현 가능.
- [리스크]: **Phase 4/5 병렬 구현 시 모델 확정 시점** — plan.md에서 Phase 4와 5를 Phase 3 이후 병렬로 진행 가능하다고 했는데, 양쪽 모두 `model/entry.go`에 강하게 의존함. Phase 3.1에서 모델을 확정하지 않으면 병렬 구현 시 충돌 발생 가능. 기존 피드백 A4를 반영하여, Phase 3.1의 완료 기준에 "model/entry.go 데이터 모델 확정"을 명시적으로 추가해야 함.
- [리스크]: **Phase 8 운영 단계의 엔드포인트 불일치** — plan.md Phase 8에서 `/collect`, `/classify`, `/report` 엔드포인트를 언급하는데, architecture.md API 설계(6장)에는 `/api/collect`와 `/api/summary`만 있고 `/classify`와 `/report`는 없음. 분류는 CLI 기반(Phase 4.3) + 웹 UI PUT 엔드포인트로 처리하는 설계이므로, plan.md Phase 8.2~8.3의 기술적 수행 방식을 architecture.md API에 맞게 보정해야 함.
- [보완필요]: **Phase 6 테스트 범위** — 테스트 계획에 HTTPS 헬스체크의 인증서 에러 케이스, 타임아웃 케이스, 400대 응답 정상 판정 케이스가 명시되어 있지 않음. architecture.md 3.4의 복잡한 판정 로직을 고려하면 이 부분의 테스트 케이스를 Phase 6.5에 구체적으로 추가하는 것을 권장.

---

## Review Cycle - Planner 검토 (2026-04-09)

### plan.md 검토

- [부합]: **프로젝트 요약 및 Phase 구조** -- requirements.md의 프로젝트 목표(수집/분류/정리, 웹 UI 현황 공유, 드래그앤드롭 합치기)를 plan.md가 충실히 반영하고 있음. 8개 Phase로 분해한 것이 적절하며 의존관계도 논리적임.
- [부합]: **기능 요구사항 커버리지** -- FR-01~FR-19가 Phase 3~5에 빠짐없이 매핑됨. 특히 FR-03(Route53 전체 조회 후 메모리 매칭, 개별 조회 금지)이 Phase 3.3에, FR-05(goroutine 병렬 수집)이 Phase 3.6에 명확히 반영됨.
- [부합]: **솔루션 목록 및 프로젝트 상태** -- Phase 8 완료 기준에 중지 프로젝트(nserise, bss, srt, ws2025)의 전략팀 컨펌, 중지예상 프로젝트(wine, ncount, dooh)의 개발팀 컨펌 반영이 명시되어 있어 requirements.md 2장과 일치.
- [부합]: **HTTPS 헬스체크** -- Phase 3.5에 "Route53 레코드 도메인으로 HTTPS(443) 요청, 경로 `/`, `/health`, `/healthz` 순차 시도, 인증서 에러 포함"이 명확히 기술되어 FR-08 에러 확인 방법과 정확히 일치.
- [부합]: **기술 제약사항** -- TC-01(Go), TC-03(AWS SDK v2), TC-04(rate limit), TC-05(파일 기반 JSON)가 plan.md 전반에 반영됨.
- [부합]: **비기능 요구사항** -- NFR-02(과도한 설계 지양)가 "임시용 도구" 톤으로 잘 반영됨. NFR-03(실제 리소스 변경은 범위 밖)이 리스크 테이블과 Phase 5.4에 반영됨.
- [보완필요]: **FR-09 티키타카 질의 세부사항** -- requirements.md에 정의된 구체적 선택지 형식("이 ALB의 솔루션을 선택해주세요: 1. signage, 2. cms ...")이 plan.md Phase 4.3에서는 "선택지 제시, 사용자 입력 수용"으로만 기술되어 있음. 구현 시 혼란 방지를 위해 선택지 목록(솔루션 10개, 환경 3개, 조치상태 3개)을 Phase 4.3 완료 기준에 명시하면 좋겠음. 단, 이는 architecture.md에서 상세히 다루고 있으므로 심각한 불일치는 아님.
- [보완필요]: **FR-17 단일 HTML/JS 파일 제약** -- requirements.md에 "단순 HTML/JS 단일 파일"로 명시되어 있으나, plan.md Phase 5에서는 이 제약을 명시적으로 언급하지 않음. Phase 5.2~5.4의 완료 기준에 "단일 HTML 파일 내 구현"을 명기하면 설계 의도가 더 명확해짐.

### architecture.md 검토

- [부합]: **데이터 모델** -- requirements.md 4장의 ALB 엔트리 JSON 구조와 architecture.md 2.1의 Entry 구조체가 필드명과 타입 모두 일치함(albName, albArn, albDns, solution, environment, status, action, mergeTarget, note).
- [부합]: **수집 전략** -- FR-03(Route53 전체 조회 -> 메모리 매칭)이 3.2~3.3에 정확히 설계됨. ALIAS/CNAME 매칭 로직이 구체적이고 실용적임.
- [부합]: **HTTPS 헬스체크 설계** -- 3.4의 헬스체크 로직이 FR-08 및 Q1/Q4/Q6 확인사항과 완전히 일치. Route53 도메인으로 요청, 경로 순차 시도, 인증서 에러(-1) 통합 처리 모두 반영됨.
- [부합]: **분류 전략** -- 4.1 ALB 이름 파싱이 FR-07의 네이밍 패턴과 일치. 4.2 조치상태 자동 추론이 requirements.md 5장의 분류 기준 테이블과 일치.
- [부합]: **합치기 계획** -- FR-11~FR-13 요구사항이 5장에 충실히 반영됨. "계획/시각화까지만 담당, 실행 버튼 없음"이 NFR-03과 일치.
- [부합]: **웹 UI** -- FR-14(테이블 뷰), FR-15(필터링), FR-16(드래그앤드롭), FR-17(단일 HTML/JS), FR-19(localhost)가 7장에 모두 반영됨.
- [부합]: **병렬 처리** -- FR-05 및 TC-04가 8장에 errgroup + 세마포어 패턴으로 구체적으로 설계됨.
- [부합]: **기술 스택** -- TC-01(Go), TC-03(AWS SDK v2), TC-05(파일 기반)이 10장에 정확히 반영됨.
- [보완필요]: **A1 -- Record 모델 확장 건 동의**: requirements.md의 `records: ["string"]`을 architecture.md의 Record 구조체(name, zoneId, zoneName, type, healthCode)로 업데이트해야 함. architecture.md의 설계가 더 실용적이므로 requirements.md를 수정하는 것에 동의.
- [보완필요]: **A2 -- 400대 응답 판정 건 동의**: 400대(401, 403, 404)는 서버가 응답하고 있으므로 "정상"으로 분류하는 Architect 제안에 동의. requirements.md FR-08에 이 판정 기준을 추가 명시 필요.
- [보완필요]: **A3 -- TG 조회 부분 실패 허용 건 동의**: 특정 ALB의 TG 조회 실패 시 해당 ALB만 unknown 마킹하는 설계가 합리적. 다만 requirements.md 데이터 모델의 status 필드에 "unknown" 값이 정의되어 있지 않으므로 추가 필요.
- [보완필요]: **FR-18 파일 기반 데이터 저장의 동시성** -- 웹 UI에서 PUT으로 수정하면서 동시에 재수집(POST /api/collect)이 진행될 경우 data/entries.json의 동시 접근 문제 발생 가능. architecture.md 9장의 "임시 파일 -> rename" 원자적 저장은 있으나, 읽기-수정-쓰기 경합에 대한 대응이 없음. sync.Mutex 등 간단한 대응 추가 설계 권장.
- [보완필요]: **API 엔드포인트의 classify 경로 부재** -- plan.md Phase 8.2에 `/classify`가 언급되어 있으나, architecture.md 6.1 API 목록에 분류 트리거 API가 없음. `POST /api/classify` 엔드포인트를 추가하거나, CLI 전용임을 명시해야 함.

---

## Review Cycle - PM 검토 (2026-04-09)

### requirements.md 검토

- [커버됨] **FR-01~FR-06 (수집 기능)**: plan.md Phase 3(작업 3.1~3.6)에서 ALB 수집, k8s 필터링, Route53 전체 조회/메모리 매칭, TG 상태 수집, 병렬 수집, JSON 저장이 모두 대응됨. 빠짐 없음.
- [커버됨] **FR-07~FR-10 (분류 기능)**: Phase 4(작업 4.1~4.4)에서 ALB 이름 파싱(솔루션+환경), 조치상태 자동 추론, 티키타카 질의, JSON 업데이트가 모두 대응됨.
- [커버됨] **FR-11~FR-13 (합치기 계획)**: Phase 5(작업 5.3~5.4)에서 드래그앤드롭 합치기 UI와 필요 정보 출력을 커버. "계획/시각화까지만" 범위 제한도 plan.md 리스크 테이블과 Phase 5.4 설명에 반영됨.
- [커버됨] **FR-14~FR-19 (웹 UI)**: Phase 5(작업 5.1~5.4)에서 테이블 뷰, 필터링, 드래그앤드롭, 단일 HTML/JS, 파일 기반 저장, localhost 바인딩 모두 대응됨.
- [커버됨] **TC-01~TC-05 (기술 제약)**: Go 언어, AWS SDK v2, rate limit 대응, 파일 기반 저장이 plan.md와 architecture.md 양쪽에 명시됨.
- [커버됨] **NFR-01~NFR-03 (비기능 요구사항)**: 병렬 수집으로 합리적 시간 처리(NFR-01), 과도한 설계 지양(NFR-02), 실제 리소스 변경 범위 밖(NFR-03) 모두 반영됨.
- [보완필요] **A1 -- requirements.md 데이터 모델 업데이트**: requirements.md의 `records` 필드가 `string[]`이나, architecture.md에서는 `Record` 구조체(name, zoneId, zoneName, type, healthCode 포함)로 확장됨. 헬스체크 결과 저장을 위해 필수적 확장이므로 requirements.md 갱신 필요. **PM 의견: Architect 제안에 동의. Planner가 requirements.md 데이터 모델 섹션을 갱신해 주기 바람.**
- [보완필요] **A2 -- 400대 응답 판정 기준**: requirements.md FR-08에 400대 응답에 대한 명시적 기준이 없음. architecture.md에서 200~499=정상, 500+=에러로 설계한 것이 합리적 (401/403/404는 서버가 응답하고 있으므로 정상). **PM 의견: Architect 제안(400대=정상)에 동의. requirements.md FR-08에 이 기준을 명시하도록 Planner에 요청.**

### architecture.md 검토

- [실현가능] **디렉토리 구조**: collector/classifier/model/store/server 5개 패키지 분리는 Go 프로젝트 표준에 부합하고, 임시 도구 규모에 적절. model 패키지를 독립시켜 순환 의존을 방지한 설계가 좋음.
- [실현가능] **데이터 모델**: Entry/Record/TGInfo/MergePlan 4개 구조체가 requirements.md의 기능 요구사항을 충실히 반영. targetGroups를 TGInfo 구조체로 확장하여 healthy/unhealthy 카운트를 포함한 것은 조치상태 자동 추론에 유용.
- [실현가능] **수집 전략**: ALB+Route53 병렬 수집 후 매칭, 이후 TG/헬스체크를 워커풀로 수행하는 3단계 파이프라인은 합리적. errgroup + 세마포어 패턴은 Go에서 검증된 방식.
- [실현가능] **HTTPS 헬스체크**: Route53 도메인 기준 HTTPS 요청, 경로 순차 시도, 인증서 에러(healthCode=-1), 타임아웃(healthCode=0) 처리가 requirements.md FR-08과 정확히 일치. 동시성 5개(maxHealthWorkers)는 외부 서버 부하를 고려한 보수적 설정으로 적절.
- [실현가능] **API 설계**: 7개 엔드포인트가 웹 UI 기능을 충분히 지원. PUT /api/entries/{albName}으로 웹 UI 인라인 편집(P4 후속)도 추가 설계 변경 없이 가능.
- [실현가능] **웹 UI**: 단일 HTML/JS 파일, 외부 CDN 최소화 방침은 NFR-02(과도한 설계 지양)에 부합. 레이아웃 와이어프레임이 구체적이어 개발자 가이드로 충분.
- [실현가능] **병렬 처리 및 Rate Limit**: maxTGWorkers=10, maxHealthWorkers=5 동시성 제한 + AWS SDK v2 내장 retry + 지수 백오프는 ALB 60개 규모에 충분.
- [리스크-낮음] **A3 -- status "unknown" 누락**: TG 조회 실패 시 status="unknown" 마킹이 architecture.md에는 있으나 requirements.md status 필드에 정의되어 있지 않음. **PM 의견: Planner가 requirements.md status 필드에 "unknown" 값을 추가해 주기 바람.**
- [리스크-낮음] **A4 -- 모델 정의 선행**: Phase 4/5 병렬 구현 시 model/entry.go가 선행 확정되어야 한다는 Architect 의견에 동의. **PM 의견: plan.md Phase 3.1의 완료 기준에 "model/entry.go 데이터 모델 확정"을 추가할 것.**
- [리스크-중간] **Phase 8 엔드포인트 불일치**: plan.md Phase 8에서 `/collect`, `/classify`, `/report`를 언급하나, architecture.md API에는 `/api/collect`와 `/api/summary`만 있고 `/classify`, `/report`는 없음. 분류는 CLI + 웹 UI PUT으로 처리하는 설계이므로, plan.md Phase 8.2~8.3의 수행 방식을 architecture.md API에 맞게 보정 필요. Architect/Planner 검토에서도 동일 발견.
- [실현가능] **일정 판단**: 전체 8 Phase의 의존 관계가 일관되고, Phase 3~5 핵심 구현의 병렬 가능성으로 일정 단축 여지 있음. 임시 도구 규모에 비해 Phase가 다소 상세하나 체계적 진행에 적절.
- [보완필요] **Planner 신규 발견 -- 파일 동시성 문제 동의**: Planner가 지적한 PUT 수정과 POST /api/collect 동시 실행 시 entries.json 경합 문제에 동의. sync.Mutex로 간단히 대응 가능하며, architecture.md store 패키지 설계에 "파일 접근 mutex" 항목 추가를 Architect에 요청.

### Architect 합의 요청에 대한 PM 최종 의견

1. **A2 (400대 응답 판정)**: **동의**. 3인 모두 동의 -- 합의 완료.
2. **A1 (데이터 모델 업데이트)**: **동의**. 3인 모두 동의 -- 합의 완료.
3. **P4 후속 (웹 UI 인라인 편집)**: **동의**. 3인 모두 동의 -- 합의 완료. API는 이미 설계에 포함됨.

### 미해결 액션 아이템 (통합)

| # | 내용 | 담당 | 상태 |
|---|------|------|------|
| AI-01 | requirements.md 데이터 모델에서 records를 Record 구조체로, targetGroups를 TGInfo 구조체로 업데이트 | Planner | **완료** |
| AI-02 | requirements.md FR-08에 "200~499=정상, 500+=에러" 기준 명시 | Planner | **완료** |
| AI-03 | requirements.md status 필드에 "unknown" 값 추가 | Planner | **완료** |
| AI-04 | plan.md Phase 3.1 완료 기준에 "model/entry.go 데이터 모델 확정" 추가 | PM | **완료** |
| AI-05 | plan.md Phase 8.2~8.3의 엔드포인트를 architecture.md API 설계에 맞게 보정 | PM | **완료** |
| AI-06 | architecture.md store 패키지에 파일 접근 mutex 설계 추가 (Planner 발견) | Architect | **완료** (rev.2) |
| AI-07 | architecture.md API에 POST /api/classify 추가 (Planner 발견) | Architect | **완료** (rev.2) |
| AI-08 | plan.md Phase 5.2~5.4 완료 기준에 "단일 HTML 파일 내 구현" 명기 (Planner 발견) | PM | **완료** |
| AI-09 | plan.md Phase 6.5에 HTTPS 헬스체크 테스트 케이스(인증서 에러, 타임아웃, 400대 정상 판정) 추가 (Architect 발견) | PM | **완료** |
| AI-10 | architecture.md 5.6 티키타카 대상을 10개→2개로 갱신 | Architect | **완료** (rev.3) |
| AI-11 | architecture.md 5.3 키워드 매핑에 product/waiting(signage), socket/webService(nserise) 추가 | Architect | **완료** (rev.3) |

---

## Phase 3 검토 — unknown ALB 분류 반영 (2026-04-09)

### 검토 배경

aws-report.md에서 unknown ALB 10개 중 8개가 사용자 확인으로 분류 완료됨:
- signage (6개): dev-product-8080, dev-waiting-8081, dev-waiting-8082, product-8080, waiting-8081, waiting-8082
- nserise (2개): socket-dev, webService
- unknown 유지 (2개): apiService, next-office-dev

### Planner 검토

[Planner] **AI-01~03 완료**: requirements.md 데이터 모델을 architecture.md와 일치시킴. records→Record 구조체, targetGroups→TGInfo 구조체로 변경. FR-08에 "200~499=정상, 500+=에러" 판정 기준 명시. status 필드에 unknown 추가.

[Planner] **PM 미완료 항목 확인(AI-04,05,08,09)**: PM이 plan.md를 보정해야 Phase 2.3이 완료됨.

[Planner] **AI-10 신규**: architecture.md 5.6 티키타카 대상 테이블이 여전히 10개로 되어 있음. 2개로 갱신 필요.

[Planner] **AI-11 신규**: architecture.md 5.3 키워드 매핑에 product/waiting(signage), socket/webService(nserise) 키워드 추가 필요. 자동 분류기가 이 ALB들을 올바르게 분류하려면 필수.

### PM 검토

[PM] **AI-04,05,08,09 완료**: plan.md rev.3 반영. Phase 3.1에 모델 확정 기준 추가, Phase 8.2~8.3 엔드포인트를 API 설계에 맞게 보정, Phase 5에 단일 HTML 제약 명기, Phase 6.5에 HTTPS 헬스체크 테스트 케이스 4종 추가.

[PM] **unknown 축소 영향**: plan.md에 직접적 수정은 불필요 — 작업 정의가 개수에 구애받지 않는 일반적 기술로 되어 있음. Phase 4.3 티키타카 공수는 대폭 감소.

[PM] **NF-01**: architecture.md 5.3 키워드 매핑 보강 필요 (Planner AI-11과 동일 발견).

[PM] **NF-02**: architecture.md 5.6 티키타카 대상 테이블 10→2 갱신 필요 (Planner AI-10과 동일 발견).

### Architect 검토

[Architect] **AI-10,11 완료**: architecture.md rev.3 반영. 5.1 수치 갱신(추정 가능 38→46, unknown 10→2), 5.3 키워드 매핑에 product-8080/waiting-808(signage), socket-dev/webService(nserise) 추가, 5.6 티키타카 대상 10→2개로 축소.

[Architect] **A10 — 키워드 오탐 리스크**: product-8080, waiting-808은 숫자 포함 특이 패턴. 현재 48개 범위에서 충돌 없으나, 구현 시 contains 매칭으로 충분. 리스크 낮음.

[Architect] **A11 — socket-dev/webService 분류 방식**: 이들은 범용 단어이므로 키워드 매칭보다 "완전 일치 우선, 부분 일치 후순위" 2단계 매칭이 안전. 또는 사전 분류 데이터(pre-classified)로 처리 권장.

[Architect] **A12 — 미완료 항목 일괄 처리 제안**: 이번 라운드에서 AI-01~09 전량 + AI-10,11 신규 모두 해소됨. Phase 2.3 상호 검토 완료 조건 충족.

### 최종 액션 아이템 상태

| # | 내용 | 담당 | 상태 |
|---|------|------|------|
| AI-01 | requirements.md 데이터 모델 구조체 업데이트 | Planner | **완료** |
| AI-02 | requirements.md FR-08에 판정 기준 명시 | Planner | **완료** |
| AI-03 | requirements.md status 필드에 unknown 추가 | Planner | **완료** |
| AI-04 | plan.md Phase 3.1에 모델 확정 기준 추가 | PM | **완료** |
| AI-05 | plan.md Phase 8.2~8.3 엔드포인트 보정 | PM | **완료** |
| AI-06 | architecture.md store mutex 설계 | Architect | **완료** |
| AI-07 | architecture.md POST /api/classify 추가 | Architect | **완료** |
| AI-08 | plan.md Phase 5 단일 HTML 명기 | PM | **완료** |
| AI-09 | plan.md Phase 6.5 헬스체크 테스트 추가 | PM | **완료** |
| AI-10 | architecture.md 5.6 티키타카 대상 갱신 | Architect | **완료** |
| AI-11 | architecture.md 5.3 키워드 매핑 보강 | Architect | **완료** |

### Phase 2.3 완료 판정

**Phase 2.3 상호 검토: 완료** — AI-01~AI-11 전량 해소. requirements.md, plan.md, architecture.md 3개 산출물 간 일관성 확보. 충돌 없음.

### 구현 시 참고사항 (Architect)

1. 키워드 매칭 우선순위: 완전 일치 > 부분 일치 (socket-dev, webService 오탐 방지)
2. product-8080, waiting-808 패턴은 숫자 포함 — contains 매칭으로 충분
3. store 패키지 sync.Mutex로 파일 동시성 제어

### 다음 단계

> **Phase 3 구현 착수 가능**. Developer에게 architecture.md 기반 구현을 위임한다.
