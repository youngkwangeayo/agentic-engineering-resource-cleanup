---
name: classify
description: ALB 분류 티키타카 - 수집된 ALB를 솔루션/조치상태로 분류, 미분류 항목은 사용자에게 선택지 질의
---

# Classify (ALB 분류 티키타카)

## 설명
수집된 ALB 데이터를 사용자와 대화하며 솔루션/조치상태를 분류한다.

## 트리거
사용자가 `/classify` 입력 시 실행

## 동작
1. `docs/aws-report.md`를 읽어 미분류 ALB를 파악한다
2. ALB 이름, DNS, 연결된 레코드 등 정보를 기반으로 솔루션을 자동 추론한다
3. 자동 추론이 가능한 항목은 추론 결과를 보여주고 확인받는다
4. 추론이 어려운 항목은 사용자에게 선택지를 제시한다:
   ```
   **[ALB이름]** DNS: [dns] / 레코드: [records] / 타겟: [n개]
   서비스를 선택해주세요:
   1. signage  2. cms  3. nserise  4. bss  5. ncount
   6. srt  7. aiagent  8. wine  9. ws2025  10. dooh  11. 기타
   
   조치상태를 선택해주세요:
   1. 유지  2. 합치기  3. 삭제
   ```
5. 분류 결과를 `data/entries.json`에 저장한다

## 사용법
```
/classify           # 미분류 항목부터 시작
/classify all       # 전체 재분류
```
