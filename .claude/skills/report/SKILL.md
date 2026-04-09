---
name: report
description: ALB 현황 요약 리포트 - 솔루션별/조치상태별/상태별 ALB 현황을 테이블로 출력
---

# Report (현황 리포트)

## 설명
현재 ALB 정리 현황을 요약 리포트로 출력한다.

## 트리거
사용자가 `/report` 입력 시 실행

## 동작
1. `data/entries.json`을 읽는다
2. 다음 기준으로 요약한다:
   - 솔루션별 ALB 수
   - 조치상태별 ALB 수 (유지/합치기/삭제)
   - 상태별 ALB 수 (active/500/no-target/no-record)
   - 미분류 ALB 수
   - 합치기 대상 그룹별 목록
3. 터미널에 테이블 형태로 출력한다

## 사용법
```
/report             # 전체 요약
/report signage     # 특정 솔루션만
/report 삭제        # 특정 조치상태만
```
