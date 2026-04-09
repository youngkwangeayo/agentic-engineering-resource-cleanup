# ALB 정리 도구 — 사용 가이드

## 시작하기

이 프로젝트는 **에이전틱 엔지니어링** 방식으로 진행됩니다.
7개의 AI 에이전트가 각자 역할을 맡아 협업하며, 사용자는 에이전트를 호출하고 판단을 내려주는 역할입니다.

---

## 에이전트 사용법

에이전트는 Claude Code에서 `@에이전트명`으로 호출합니다.

### 기획 단계 (Planner → PM → Architect)

```
@planner 요구사항을 정리해줘
```
- Planner가 JOB.md를 읽고 질문을 시작합니다
- 선택지 형태로 물어보니 번호로 답하면 됩니다
- 확정된 요구사항은 `docs/requirements.md`에 저장됩니다

```
@pm 작업 계획을 세워줘
```
- requirements.md 기반으로 작업을 분해합니다
- 각 작업에 담당 에이전트, 우선순위를 배정합니다
- 결과는 `docs/plan.md`에 저장됩니다

```
@architect 기술 설계를 해줘
```
- requirements.md + plan.md 기반으로 설계합니다
- 디렉토리 구조, 데이터 모델, API 전략 등
- 결과는 `docs/architecture.md`에 저장됩니다

### 상호 검토

```
/review-cycle
```
- Planner, PM, Architect가 서로의 산출물을 검토합니다
- 충돌이나 의견 차이가 있으면 사용자에게 판단을 요청합니다
- 피드백은 `docs/review-notes.md`에 기록됩니다

### 개발 단계

```
@developer architecture.md 기반으로 구현해줘
```
- 설계 문서를 따라 코드를 작성합니다
- 설계와 맞지 않는 부분이 있으면 보고합니다

### 빌드 & 테스트

```
/build
/build run
```
- Go 빌드를 실행하고 에러를 처리합니다
- `run` 옵션을 주면 빌드 후 바로 실행합니다

```
@tester 테스트해줘
```
- 빌드 검증, 기능 테스트, 엣지케이스를 확인합니다
- 결과는 `docs/test-report.md`에 기록됩니다

### 코드 리뷰

```
@reviewer 코드 리뷰해줘
```
- 설계 준수, 코드 품질, 보안을 점검합니다
- CRITICAL / WARNING / SUGGESTION 심각도로 분류합니다
- 결과는 `docs/review-report.md`에 기록됩니다

### AWS 조사

```
/collect
```
- awscli로 ALB, Route53, Target Group 데이터를 수집합니다
- 결과는 `docs/aws-report.md`에 기록됩니다

```
/classify
```
- 수집된 ALB를 솔루션/조치상태로 분류합니다
- 자동 추론이 어려운 항목은 선택지를 제시합니다:
  ```
  **alb-dev-test** DNS: xxx / 레코드: 없음 / 타겟: 0개
  서비스를 선택해주세요:
  1. signage  2. cms  3. nserise ...
  조치상태: 1. 유지  2. 합치기  3. 삭제
  ```

```
/report
/report signage
/report 삭제
```
- 현재 분류 현황을 요약합니다

---

## 권장 작업 순서

```
1단계: 기획
  @planner → @pm → @architect → /review-cycle

2단계: 개발
  @developer → /build → @tester → @reviewer
  (문제 있으면 @developer로 수정 → 반복)

3단계: 운영
  /collect → /classify → /report
  (분류 결과 확인 후 웹 UI로 공유)
```

---

## 산출물 위치

| 파일 | 작성자 | 내용 |
|------|--------|------|
| `docs/requirements.md` | Planner | 확정된 요구사항 |
| `docs/plan.md` | PM | 작업 분해, 우선순위 |
| `docs/architecture.md` | Architect | 기술 설계 |
| `docs/review-notes.md` | 공유 | 3인 상호 검토 피드백 |
| `docs/aws-report.md` | Inspector | AWS 리소스 현황 |
| `docs/test-report.md` | Tester | 테스트 결과 |
| `docs/review-report.md` | Reviewer | 코드 리뷰 결과 |
| `data/entries.json` | classify | ALB 분류 데이터 |

---

## 팁

- **티키타카**: 에이전트가 질문하면 번호로 간단히 답하세요
- **검토 사이클**: 기획 변경이 있으면 `/review-cycle`을 다시 돌려 정합성을 확인하세요
- **부분 작업**: 에이전트에게 특정 범위만 지시할 수 있습니다 (예: `@developer collector만 구현해줘`)
- **현황 확인**: 언제든 `/report`로 현재 상태를 볼 수 있습니다
