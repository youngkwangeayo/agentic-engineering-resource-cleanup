---
name: review-cycle
description: Planner/PM/Architect 3인 상호 검토 사이클 - 산출물 교차 검토 후 피드백을 review-notes.md에 기록
---

# Review Cycle (상호 검토 사이클)

## 설명
Planner, PM, Architect 3인의 상호 검토 사이클을 실행한다.
각 에이전트가 서로의 산출물을 읽고 피드백을 남긴다.

## 트리거
사용자가 `/review-cycle` 입력 시 실행

## 동작
1. 현재 `docs/` 디렉토리의 산출물 상태를 확인한다:
   - `docs/requirements.md` (Planner)
   - `docs/plan.md` (PM)
   - `docs/architecture.md` (Architect)
2. 각 산출물이 있는 에이전트부터 검토를 시작한다
3. 검토 순서:
   a. Planner → PM/Architect 산출물 검토 → 피드백
   b. PM → Planner/Architect 산출물 검토 → 피드백
   c. Architect → Planner/PM 산출물 검토 → 피드백
4. 모든 피드백을 `docs/review-notes.md`에 기록한다
5. 충돌이나 미해결 사항이 있으면 사용자에게 판단을 요청한다

## 사용법
```
/review-cycle                # 전체 검토 사이클
/review-cycle planner        # Planner 산출물만 검토
/review-cycle architecture   # 설계만 검토
```
