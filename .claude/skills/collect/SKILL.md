---
name: collect
description: AWS 데이터 수집 - ALB, Route53, Target Group 현황을 awscli로 수집하여 docs/aws-report.md에 정리
---

# Collect (AWS 데이터 수집)

## 설명
AWS Inspector 에이전트를 호출하여 ALB, Route53, Target Group 데이터를 수집한다.

## 트리거
사용자가 `/collect` 입력 시 실행

## 동작
1. AWS Inspector 에이전트를 실행한다
2. awscli로 다음을 수집한다:
   - `aws elbv2 describe-load-balancers` — ALB 목록 (k8s 제외)
   - `aws elbv2 describe-listeners` — ALB별 리스너
   - `aws elbv2 describe-target-groups` — 타겟그룹
   - `aws elbv2 describe-target-health` — 타겟 헬스
   - `aws route53 list-hosted-zones` + `list-resource-record-sets` — 레코드
3. 수집 결과를 `docs/aws-report.md`에 정리한다
4. 분류가 어려운 항목은 사용자에게 질문한다

## 사용법
```
/collect
```
