# ALB Resource Cleanup Tool

## 프로젝트 개요
AWS ALB 리소스를 정리하는 도구. 약 60개의 ALB와 다수의 Route53 레코드를 수집/분류/관리한다.
에이전틱 엔지니어링 방식으로 개발한다 — 7개 에이전트가 역할을 분담하여 협업.

## 기술 스택
- Go (병렬 처리)
- AWS SDK v2 (ALB, Route53, ELBv2)
- 웹 UI (임시용, 단순 HTML/JS)
- 파일 기반 데이터 저장 (JSON)

## 핵심 규칙
- k8s로 시작하는 ALB 이름은 제외 (k8s ALB controller가 관리)
- Route53 레코드는 한번에 전체 조회 후 메모리에서 ALB 매칭
- 솔루션(프로젝트): signage, cms, nserise, bss, ncount, srt, aiagent, wine, ws2025, dooh
- 조치상태: 유지, 합치기, 삭제
- 분류가 어려운 ALB는 사용자에게 선택지를 제시하며 질문 (티키타카)

## 에이전트 구성 (.claude/agents/)
| 에이전트 | 역할 | 산출물 |
|---------|------|--------|
| Planner | 요구사항 수렴, 티키타카 질의 | docs/requirements.md |
| PM | 작업 분해, 우선순위, 에이전트 조율 | docs/plan.md |
| Architect | 기술 설계, 구조, 데이터 모델 | docs/architecture.md |
| Developer | 설계 기반 코드 구현 | 소스 코드 |
| AWS Inspector | awscli로 실데이터 조사 | docs/aws-report.md |
| Tester | 빌드/기능 검증, 엣지케이스 | docs/test-report.md |
| Reviewer | 코드 리뷰, 품질/보안 체크 | docs/review-report.md |

### 상호 검토 구조
- Planner ↔ PM ↔ Architect 는 서로의 산출물을 검토한다
- 피드백은 `docs/review-notes.md`에 `[역할명] ...` 형식으로 남긴다
- 충돌 발견 시 대안과 함께 사용자에게 판단을 요청한다

## 스킬 (.claude/skills/)
| 스킬 | 명령어 | 설명 |
|------|--------|------|
| collect | `/collect` | AWS 데이터 수집 (ALB, Route53, Target Group) |
| classify | `/classify` | ALB 분류 티키타카 (솔루션/조치상태 결정) |
| report | `/report` | ALB 현황 요약 리포트 |
| build | `/build` | Go 빌드 및 실행 |
| review-cycle | `/review-cycle` | Planner↔PM↔Architect 상호 검토 사이클 |

## 작업 흐름 (권장 순서)
1. Planner로 요구사항 확정 → `docs/requirements.md`
2. PM으로 작업 분해 → `docs/plan.md`
3. Architect로 기술 설계 → `docs/architecture.md`
4. `/review-cycle`로 3인 상호 검토
5. Developer로 구현
6. `/build`로 빌드 검증
7. Tester로 테스트 → `docs/test-report.md`
8. Reviewer로 코드 리뷰 → `docs/review-report.md`
9. `/collect`로 AWS 데이터 수집
10. `/classify`로 ALB 분류
11. `/report`로 현황 확인

## 빌드 & 실행
```bash
go run .              # 수집 + 웹서버 시작
go build -o new-lb .  # 빌드
```

## AWS CLI
- awscli 이미 연결되어 있음. 별도 credential 설정 불필요.
