# AWS ALB 현황 리포트

> 조사일: 2026-04-09
> 총 ALB 수: 57개 (k8s 제외 후: 48개)

---

## 1. ALB 전체 목록

| # | ALB 이름 | 솔루션 | 환경 | 리스너 | TG수 | 타겟수 | 헬스상태 | Route53 레코드 수 |
|---|---------|--------|------|--------|------|--------|---------|-------------------|
| 1 | BSSTizenWebSignage | signage | prd | 1 | 1 | 1 | unhealthy | 1 |
| 2 | BSSTizenWebSignageSub1 | signage | prd | 2 | 1 | 1 | healthy | 1 |
| 3 | BSSTizenWebSignageSub2 | signage | prd | 2 | 1 | 1 | healthy | 1 |
| 4 | LBTizenWebSignage | signage | prd | 1 | 1 | 1 | healthy | 0 |
| 5 | NCount-Live-ELB | ncount | prd | 1 | 1 | 1 | unused | 1 |
| 6 | aiagent-staging | aiagent | stg | 2 | 3 | 3 | mixed | 3 |
| 7 | alb-dev-signage | signage | dev | 2 | 3 | 3 | mixed | 3 |
| 8 | apiService | unknown | prd | 1 | 1 | 1 | unused | 1 |
| 9 | bss-cms | bss | prd | 2 | 1 | 1 | healthy | 1 |
| 10 | bss-kiosk | bss | prd | 3 | 2 | 2 | healthy | 1 |
| 11 | bss-order | bss | prd | 3 | 2 | 2 | healthy | 1 |
| 12 | cms-elb | cms | prd | 2 | 1 | 1 | unused | 2 |
| 13 | dev-cms-elb | cms | dev | 2 | 7 | 6 | mixed | 7 |
| 14 | dev-product-8080 | signage | dev | 2 | 1 | 1 | unhealthy | 1 |
| 15 | dev-waiting-8081 | signage | dev | 2 | 1 | 1 | healthy | 1 |
| 16 | dev-waiting-8082 | signage | dev | 2 | 1 | 1 | healthy | 1 |
| 17 | dooh-dev-elb | dooh | dev | 1 | 4 | 4 | healthy | 4 |
| 18 | elb-knowledge-graph | aiagent | prd | 2 | 1 | 1 | unhealthy | 2 |
| 19 | elb-srt-device-api | srt | dev | 2 | 1 | 1 | unhealthy | 3 |
| 20 | elb-srt-staging-device-api | srt | dev | 2 | 1 | 1 | unhealthy | 1 |
| 21 | elb-srt-staging-store-api | srt | stg | 2 | 1 | 1 | unhealthy | 3 |
| 22 | elb-srt-store-api | srt | prd | 2 | 1 | 1 | unhealthy | 2 |
| 23 | lb-dev-scms | signage | dev | 2 | 1 | 1 | unhealthy | 1 |
| 24 | lb-dev-swaiting | signage | dev | 2 | 1 | 1 | healthy | 1 |
| 25 | lb-scms | signage | prd | 2 | 1 | 1 | unhealthy | 2 |
| 26 | lb-signage-oss-cms | signage | prd | 2 | 1 | 0 | no_target | 1 |
| 27 | lb-signage-oss-waiting | signage | prd | 2 | 1 | 0 | no_target | 1 |
| 28 | lb-swaiting | signage | prd | 2 | 1 | 1 | healthy | 2 |
| 29 | ncount-dev | ncount | dev | 1 | 1 | 1 | unused | 2 |
| 30 | next-office-dev | unknown | dev | 3 | 2 | 2 | unused | 1 |
| 31 | nextpay-kiosk-dev | nserise | dev | 1 | 1 | 1 | unhealthy | 1 |
| 32 | nextpay-nkds-dev | nserise | dev | 1 | 1 | 1 | healthy | 1 |
| 33 | nextpay-norder-dev | nserise | dev | 2 | 1 | 1 | unhealthy | 1 |
| 34 | nextpay-npos-dev | nserise | dev | 2 | 1 | 1 | unhealthy | 1 |
| 35 | nextpay-order-ai | nserise | prd | 1 | 0 | 0 | no_tg | 1 |
| 36 | nextpay-shop8080 | nserise | prd | 1 | 1 | 1 | healthy | 1 |
| 37 | nseries-dev-elb | nserise | dev | 2 | 9 | 7 | mixed | 0 |
| 38 | nseries-elb | nserise | prd | 1 | 6 | 3 | mixed | 0 |
| 39 | product-8080 | signage | prd | 2 | 1 | 1 | unhealthy | 1 |
| 40 | socket-dev | nserise | dev | 2 | 1 | 0 | no_target | 1 |
| 41 | srt-device-api | srt | dev | 3 | 2 | 0 | no_target | 0 |
| 42 | srt-store-api | srt | prd | 1 | 1 | 0 | no_target | 0 |
| 43 | srt-system-api | srt | prd | 1 | 1 | 0 | no_target | 0 |
| 44 | waiting-8081 | signage | prd | 2 | 1 | 1 | healthy | 1 |
| 45 | waiting-8082 | signage | prd | 2 | 1 | 1 | healthy | 1 |
| 46 | webService | nserise | prd | 2 | 12 | 7 | mixed | 12 |
| 47 | wine-curation-elb | wine | prd | 1 | 1 | 1 | healthy | 1 |
| 48 | ws2025-lb | ws2025 | prd | 1 | 2 | 2 | unused | 2 |

---

## 2. ALB별 Route53 레코드 매핑

### BSSTizenWebSignage
- DNS: `BSSTizenWebSignage-1576074367.ap-northeast-2.elb.amazonaws.com`
- 연결된 레코드:
  - `bss-signage.nextpay.co.kr`

### BSSTizenWebSignageSub1
- DNS: `BSSTizenWebSignageSub1-972785156.ap-northeast-2.elb.amazonaws.com`
- 연결된 레코드:
  - `bss-signage-sub1.nextpay.co.kr`

### BSSTizenWebSignageSub2
- DNS: `BSSTizenWebSignageSub2-333608192.ap-northeast-2.elb.amazonaws.com`
- 연결된 레코드:
  - `bss-signage-sub2.nextpay.co.kr`

### NCount-Live-ELB
- DNS: `NCount-Live-ELB-967232354.ap-northeast-2.elb.amazonaws.com`
- 연결된 레코드:
  - `ncount-whisper.nextpay.co.kr`

### aiagent-staging
- DNS: `aiagent-staging-1061935423.ap-northeast-2.elb.amazonaws.com`
- 연결된 레코드:
  - `staging-api.nextaistore.co.kr`
  - `staging-system.nextaistore.co.kr`
  - `staging.nextaistore.co.kr`

### alb-dev-signage
- DNS: `alb-dev-signage-299073066.ap-northeast-2.elb.amazonaws.com`
- 연결된 레코드:
  - `dev-scms.nextpay.co.kr`
  - `dev-sw.nextpay.co.kr`
  - `signage-glitchtip.nextpay.co.kr`

### apiService
- DNS: `apiService-2020709608.ap-northeast-2.elb.amazonaws.com`
- 연결된 레코드:
  - `api.nextpay.co.kr`

### bss-cms
- DNS: `bss-cms-174951817.ap-northeast-2.elb.amazonaws.com`
- 연결된 레코드:
  - `bss.nextpay.co.kr`

### bss-kiosk
- DNS: `bss-kiosk-1230504848.ap-northeast-2.elb.amazonaws.com`
- 연결된 레코드:
  - `bss-kiosk.nextpay.co.kr`

### bss-order
- DNS: `bss-order-1305085170.ap-northeast-2.elb.amazonaws.com`
- 연결된 레코드:
  - `bss-order.nextpay.co.kr`

### cms-elb
- DNS: `cms-elb-393392806.ap-northeast-2.elb.amazonaws.com`
- 연결된 레코드:
  - `cms.nextpay.co.kr`
  - `lidar.ncount.nextpay.co.kr`

### dev-cms-elb
- DNS: `dev-cms-elb-621823357.ap-northeast-2.elb.amazonaws.com`
- 연결된 레코드:
  - `aiagent.nextpay.co.kr`
  - `aiot.nextpay.co.kr`
  - `dev-ncount.nextpay.co.kr`
  - `dev.nextpay.co.kr`
  - `dooh-api.nextpay.co.kr`
  - `staging.nextpay.co.kr`
  - `wine-curation-dev.nextpay.co.kr`

### dev-product-8080
- DNS: `dev-product-8080-1125016562.ap-northeast-2.elb.amazonaws.com`
- 연결된 레코드:
  - `dev-ur.nextpay.co.kr`

### dev-waiting-8081
- DNS: `dev-waiting-8081-1289716708.ap-northeast-2.elb.amazonaws.com`
- 연결된 레코드:
  - `backup-dev-waiting.nextpay.co.kr`

### dev-waiting-8082
- DNS: `dev-waiting-8082-2013353451.ap-northeast-2.elb.amazonaws.com`
- 연결된 레코드:
  - `backup-dev-nsm.nextpay.co.kr`

### dooh-dev-elb
- DNS: `dooh-dev-elb-1993049909.ap-northeast-2.elb.amazonaws.com`
- 연결된 레코드:
  - `backstage.dooh.nextpay.co.kr`
  - `broker.dooh.nextpay.co.kr`
  - `chat.dooh.nextpay.co.kr`
  - `dashboard.dooh.nextpay.co.kr`

### elb-knowledge-graph
- DNS: `elb-knowledge-graph-131700488.ap-northeast-2.elb.amazonaws.com`
- 연결된 레코드:
  - `graph.nextaistore.co.kr`
  - `graph.nextaistore.com`

### elb-srt-device-api
- DNS: `elb-srt-device-api-387181319.ap-northeast-2.elb.amazonaws.com`
- 연결된 레코드:
  - `dvc-srtdev.nextpay.co.kr`
  - `dvc.srtdev.nextpay.co.kr`
  - `str-srtdev.nextpay.co.kr`

### elb-srt-staging-device-api
- DNS: `elb-srt-staging-device-api-528593873.ap-northeast-2.elb.amazonaws.com`
- 연결된 레코드:
  - `dvc-srtstg.nextpay.co.kr`

### elb-srt-staging-store-api
- DNS: `elb-srt-staging-store-api-934084596.ap-northeast-2.elb.amazonaws.com`
- 연결된 레코드:
  - `store-srtstg.nextpay.co.kr`
  - `str-srtstg.nextpay.co.kr`
  - `sys-srtstg.nextpay.co.kr`

### elb-srt-store-api
- DNS: `elb-srt-store-api-1786817859.ap-northeast-2.elb.amazonaws.com`
- 연결된 레코드:
  - `store-srtdev.nextpay.co.kr`
  - `sys-srtdev.nextpay.co.kr`

### lb-dev-scms
- DNS: `lb-dev-scms-698973874.ap-northeast-2.elb.amazonaws.com`
- 연결된 레코드:
  - `dev-nsm.nextpay.co.kr`

### lb-dev-swaiting
- DNS: `lb-dev-swaiting-1347573332.ap-northeast-2.elb.amazonaws.com`
- 연결된 레코드:
  - `dev-waiting.nextpay.co.kr`

### lb-scms
- DNS: `lb-scms-390862235.ap-northeast-2.elb.amazonaws.com`
- 연결된 레코드:
  - `nsm.nextpay.co.kr`
  - `scms.nextpay.co.kr`

### lb-signage-oss-cms
- DNS: `lb-signage-oss-cms-701641332.ap-northeast-2.elb.amazonaws.com`
- 연결된 레코드:
  - `x-signage.nextaistore.com`

### lb-signage-oss-waiting
- DNS: `lb-signage-oss-waiting-1181203740.ap-northeast-2.elb.amazonaws.com`
- 연결된 레코드:
  - `x-waiting.nextaistore.com`

### lb-swaiting
- DNS: `lb-swaiting-581779347.ap-northeast-2.elb.amazonaws.com`
- 연결된 레코드:
  - `sw.nextpay.co.kr`
  - `waiting.nextpay.co.kr`

### ncount-dev
- DNS: `ncount-dev-1417780276.ap-northeast-2.elb.amazonaws.com`
- 연결된 레코드:
  - `hybo.lidar.ncount.nextpay.co.kr`
  - `poc.hybo.lidar.ncount.nextpay.co.kr`

### next-office-dev
- DNS: `next-office-dev-1163943400.ap-northeast-2.elb.amazonaws.com`
- 연결된 레코드:
  - `dev-nextoffice.nextpay.co.kr`

### nextpay-kiosk-dev
- DNS: `nextpay-kiosk-dev-767343648.ap-northeast-2.elb.amazonaws.com`
- 연결된 레코드:
  - `nkiosk-dev.nextpay.co.kr`

### nextpay-nkds-dev
- DNS: `nextpay-nkds-dev-2073863581.ap-northeast-2.elb.amazonaws.com`
- 연결된 레코드:
  - `nkds-dev.nextpay.co.kr`

### nextpay-norder-dev
- DNS: `nextpay-norder-dev-36629141.ap-northeast-2.elb.amazonaws.com`
- 연결된 레코드:
  - `norder-dev.nextpay.co.kr`

### nextpay-npos-dev
- DNS: `nextpay-npos-dev-1407072045.ap-northeast-2.elb.amazonaws.com`
- 연결된 레코드:
  - `npos-dev.nextpay.co.kr`

### nextpay-order-ai
- DNS: `nextpay-order-ai-894286297.ap-northeast-2.elb.amazonaws.com`
- 연결된 레코드:
  - `norder-ai.nextpay.co.kr`

### nextpay-shop8080
- DNS: `nextpay-shop8080-1474328396.ap-northeast-2.elb.amazonaws.com`
- 연결된 레코드:
  - `shopvs.nextpay.co.kr`

### product-8080
- DNS: `product-8080-1041987207.ap-northeast-2.elb.amazonaws.com`
- 연결된 레코드:
  - `ur.nextpay.co.kr`

### socket-dev
- DNS: `socket-dev-587569175.ap-northeast-2.elb.amazonaws.com`
- 연결된 레코드:
  - `socket-dev.nextpay.co.kr`

### waiting-8081
- DNS: `waiting-8081-442550494.ap-northeast-2.elb.amazonaws.com`
- 연결된 레코드:
  - `backup-waiting.nextpay.co.kr`

### waiting-8082
- DNS: `waiting-8082-1157880650.ap-northeast-2.elb.amazonaws.com`
- 연결된 레코드:
  - `backup-nsm.nextpay.co.kr`

### webService
- DNS: `webService-2021174328.ap-northeast-2.elb.amazonaws.com`
- 연결된 레코드:
  - `back-office.nextpay.co.kr`
  - `ca.nextpay.co.kr`
  - `ncms.nextpay.co.kr`
  - `nkds-new.nextpay.co.kr`
  - `nkds.nextpay.co.kr`
  - `nkiosk.nextpay.co.kr`
  - `nkioskkinkos.nextpay.co.kr`
  - `norder.nextpay.co.kr`
  - `npos.nextpay.co.kr`
  - `salesking.nextpay.co.kr`
  - `shop.nextpay.co.kr`
  - `socket.nextpay.co.kr`

### wine-curation-elb
- DNS: `wine-curation-elb-1718355595.ap-northeast-2.elb.amazonaws.com`
- 연결된 레코드:
  - `wine-curation.nextpay.co.kr`

### ws2025-lb
- DNS: `ws2025-lb-251643516.ap-northeast-2.elb.amazonaws.com`
- 연결된 레코드:
  - `ws2025-api.nextpay.co.kr`
  - `ws2025-web.nextpay.co.kr`

---

## 3. 타겟이 없는 ALB (삭제 후보)

타겟그룹이 없거나, 타겟그룹은 있지만 등록된 타겟이 없는 ALB입니다.

| ALB 이름 | 솔루션 | 상태 | 타겟그룹 |
|---------|--------|------|---------|
| lb-signage-oss-cms | signage | no_target | tg-signage-oss-cms |
| lb-signage-oss-waiting | signage | no_target | tg-signage-oss-waiting |
| nextpay-order-ai | nserise | no_target_group | (없음) |
| socket-dev | nserise | no_target | socket-dev-tg |
| srt-device-api | srt | no_target | srt-dev-monitoring, srt-device-api |
| srt-store-api | srt | no_target | srt-store-api |
| srt-system-api | srt | no_target | srt-system-api |

---

## 4. Route53 레코드가 연결되지 않은 ALB (삭제 후보)

| ALB 이름 | 솔루션 | 헬스상태 |
|---------|--------|---------|
| LBTizenWebSignage | signage | healthy |
| nseries-dev-elb | nserise | mixed |
| nseries-elb | nserise | mixed |
| srt-device-api | srt | no_target |
| srt-store-api | srt | no_target |
| srt-system-api | srt | no_target |

---

## 5. Unhealthy 타겟이 있는 ALB

| ALB 이름 | 솔루션 | 환경 | Healthy | Unhealthy | Unused |
|---------|--------|------|---------|-----------|--------|
| BSSTizenWebSignage | signage | prd | 0 | 1 | 0 |
| aiagent-staging | aiagent | stg | 1 | 2 | 0 |
| alb-dev-signage | signage | dev | 2 | 1 | 0 |
| dev-product-8080 | signage | dev | 0 | 1 | 0 |
| elb-knowledge-graph | aiagent | prd | 0 | 1 | 0 |
| elb-srt-device-api | srt | dev | 0 | 1 | 0 |
| elb-srt-staging-device-api | srt | dev | 0 | 1 | 0 |
| elb-srt-staging-store-api | srt | stg | 0 | 1 | 0 |
| elb-srt-store-api | srt | prd | 0 | 1 | 0 |
| lb-dev-scms | signage | dev | 0 | 1 | 0 |
| lb-scms | signage | prd | 0 | 1 | 0 |
| nextpay-kiosk-dev | nserise | dev | 0 | 1 | 0 |
| nextpay-norder-dev | nserise | dev | 0 | 1 | 0 |
| nextpay-npos-dev | nserise | dev | 0 | 1 | 0 |
| nseries-dev-elb | nserise | dev | 4 | 2 | 1 |
| nseries-elb | nserise | prd | 1 | 2 | 0 |
| product-8080 | signage | prd | 0 | 1 | 0 |
| webService | nserise | prd | 4 | 3 | 0 |

---

## 6. 솔루션별 ALB 현황

### signage (활성) - 17개

| ALB 이름 | 환경 | 헬스상태 | 레코드수 |
|---------|------|---------|---------|
| BSSTizenWebSignage | prd | unhealthy | 1 |
| BSSTizenWebSignageSub1 | prd | healthy | 1 |
| BSSTizenWebSignageSub2 | prd | healthy | 1 |
| LBTizenWebSignage | prd | healthy | 0 |
| alb-dev-signage | dev | mixed | 3 |
| dev-product-8080 | dev | unhealthy | 1 |
| dev-waiting-8081 | dev | healthy | 1 |
| dev-waiting-8082 | dev | healthy | 1 |
| lb-dev-scms | dev | unhealthy | 1 |
| lb-dev-swaiting | dev | healthy | 1 |
| lb-scms | prd | unhealthy | 2 |
| lb-signage-oss-cms | prd | no_target | 1 |
| lb-signage-oss-waiting | prd | no_target | 1 |
| lb-swaiting | prd | healthy | 2 |
| product-8080 | prd | unhealthy | 1 |
| waiting-8081 | prd | healthy | 1 |
| waiting-8082 | prd | healthy | 1 |

### cms (활성) - 2개

| ALB 이름 | 환경 | 헬스상태 | 레코드수 |
|---------|------|---------|---------|
| cms-elb | prd | unused | 2 |
| dev-cms-elb | dev | mixed | 7 |

### aiagent (활성) - 2개

| ALB 이름 | 환경 | 헬스상태 | 레코드수 |
|---------|------|---------|---------|
| aiagent-staging | stg | mixed | 3 |
| elb-knowledge-graph | prd | unhealthy | 2 |

### nserise (중지 예정) - 10개

| ALB 이름 | 환경 | 헬스상태 | 레코드수 |
|---------|------|---------|---------|
| nextpay-kiosk-dev | dev | unhealthy | 1 |
| nextpay-nkds-dev | dev | healthy | 1 |
| nextpay-norder-dev | dev | unhealthy | 1 |
| nextpay-npos-dev | dev | unhealthy | 1 |
| nextpay-order-ai | prd | no_target_group | 1 |
| nextpay-shop8080 | prd | healthy | 1 |
| nseries-dev-elb | dev | mixed | 0 |
| nseries-elb | prd | mixed | 0 |
| socket-dev | dev | no_target | 1 |
| webService | prd | mixed | 12 |

### bss (중지 예정) - 3개

| ALB 이름 | 환경 | 헬스상태 | 레코드수 |
|---------|------|---------|---------|
| bss-cms | prd | healthy | 1 |
| bss-kiosk | prd | healthy | 1 |
| bss-order | prd | healthy | 1 |

### ncount (중지 예상) - 2개

| ALB 이름 | 환경 | 헬스상태 | 레코드수 |
|---------|------|---------|---------|
| NCount-Live-ELB | prd | unused | 1 |
| ncount-dev | dev | unused | 2 |

### srt (중지 예정) - 7개

| ALB 이름 | 환경 | 헬스상태 | 레코드수 |
|---------|------|---------|---------|
| elb-srt-device-api | dev | unhealthy | 3 |
| elb-srt-staging-device-api | dev | unhealthy | 1 |
| elb-srt-staging-store-api | stg | unhealthy | 3 |
| elb-srt-store-api | prd | unhealthy | 2 |
| srt-device-api | dev | no_target | 0 |
| srt-store-api | prd | no_target | 0 |
| srt-system-api | prd | no_target | 0 |

### wine (중지 예상) - 1개

| ALB 이름 | 환경 | 헬스상태 | 레코드수 |
|---------|------|---------|---------|
| wine-curation-elb | prd | healthy | 1 |

### ws2025 (중지 예정) - 1개

| ALB 이름 | 환경 | 헬스상태 | 레코드수 |
|---------|------|---------|---------|
| ws2025-lb | prd | unused | 2 |

### dooh (중지 예상) - 1개

| ALB 이름 | 환경 | 헬스상태 | 레코드수 |
|---------|------|---------|---------|
| dooh-dev-elb | dev | healthy | 4 |

### unknown (-) - 2개

| ALB 이름 | 환경 | 헬스상태 | 레코드수 |
|---------|------|---------|---------|
| apiService | prd | unused | 1 |
| next-office-dev | dev | unused | 1 |

---

## 7. ALB 네이밍 패턴 분석

기대 패턴: `{aws리소스}-{환경}-{솔루션}-{서비스}` (예: `alb-dev-signage-xxx`)

### 패턴 준수 (3개)

- `alb-dev-signage`
- `lb-dev-scms`
- `lb-dev-swaiting`

### 패턴 미준수 (45개)

| ALB 이름 | 현재 패턴 | 솔루션 추정 |
|---------|---------|------------|
| BSSTizenWebSignage | BSSTizenWebSignage | signage |
| BSSTizenWebSignageSub1 | BSSTizenWebSignageSub1 | signage |
| BSSTizenWebSignageSub2 | BSSTizenWebSignageSub2 | signage |
| LBTizenWebSignage | LBTizenWebSignage | signage |
| NCount-Live-ELB | NCount-Live-ELB | ncount |
| aiagent-staging | aiagent-staging | aiagent |
| apiService | apiService | unknown |
| bss-cms | bss-cms | bss |
| bss-kiosk | bss-kiosk | bss |
| bss-order | bss-order | bss |
| cms-elb | cms-elb | cms |
| dev-cms-elb | dev-cms-elb | cms |
| dev-product-8080 | dev-product-8080 | signage |
| dev-waiting-8081 | dev-waiting-8081 | signage |
| dev-waiting-8082 | dev-waiting-8082 | signage |
| dooh-dev-elb | dooh-dev-elb | dooh |
| elb-knowledge-graph | elb-knowledge-graph | aiagent |
| elb-srt-device-api | elb-srt-device-api | srt |
| elb-srt-staging-device-api | elb-srt-staging-device-api | srt |
| elb-srt-staging-store-api | elb-srt-staging-store-api | srt |
| elb-srt-store-api | elb-srt-store-api | srt |
| lb-scms | lb-scms | signage |
| lb-signage-oss-cms | lb-signage-oss-cms | signage |
| lb-signage-oss-waiting | lb-signage-oss-waiting | signage |
| lb-swaiting | lb-swaiting | signage |
| ncount-dev | ncount-dev | ncount |
| next-office-dev | next-office-dev | unknown |
| nextpay-kiosk-dev | nextpay-kiosk-dev | nserise |
| nextpay-nkds-dev | nextpay-nkds-dev | nserise |
| nextpay-norder-dev | nextpay-norder-dev | nserise |
| nextpay-npos-dev | nextpay-npos-dev | nserise |
| nextpay-order-ai | nextpay-order-ai | nserise |
| nextpay-shop8080 | nextpay-shop8080 | nserise |
| nseries-dev-elb | nseries-dev-elb | nserise |
| nseries-elb | nseries-elb | nserise |
| product-8080 | product-8080 | signage |
| socket-dev | socket-dev | nserise |
| srt-device-api | srt-device-api | srt |
| srt-store-api | srt-store-api | srt |
| srt-system-api | srt-system-api | srt |
| waiting-8081 | waiting-8081 | signage |
| waiting-8082 | waiting-8082 | signage |
| webService | webService | nserise |
| wine-curation-elb | wine-curation-elb | wine |
| ws2025-lb | ws2025-lb | ws2025 |

---

## 8. 종합 통계

### 헬스 상태별

| 상태 | 개수 |
|------|------|
| healthy | 16 |
| mixed | 6 |
| unhealthy | 13 |
| unused | 6 |
| no_target | 6 |
| no_target_group | 1 |

### 환경별

| 환경 | 개수 |
|------|------|
| dev | 19 |
| stg | 2 |
| prd | 27 |

### 삭제 후보 ALB (타겟 없음 또는 레코드 없음)

총 10개:

- **LBTizenWebSignage** (signage) - 레코드 없음
- **lb-signage-oss-cms** (signage) - 타겟 없음
- **lb-signage-oss-waiting** (signage) - 타겟 없음
- **nextpay-order-ai** (nserise) - 타겟 없음
- **nseries-dev-elb** (nserise) - 레코드 없음
- **nseries-elb** (nserise) - 레코드 없음
- **socket-dev** (nserise) - 타겟 없음
- **srt-device-api** (srt) - 타겟 없음, 레코드 없음
- **srt-store-api** (srt) - 타겟 없음, 레코드 없음
- **srt-system-api** (srt) - 타겟 없음, 레코드 없음

---

## 9. ALB별 상세 타겟그룹 및 타겟 헬스

### BSSTizenWebSignage
- 솔루션: signage | 환경: prd
- DNS: `BSSTizenWebSignage-1576074367.ap-northeast-2.elb.amazonaws.com`
- 리스너: HTTPS:443
- Route53: `bss-signage.nextpay.co.kr`
- 타겟그룹:
  - `BSSTizenWebSignage`: 총 1개 (healthy:0, unhealthy:1, unused:0)

### BSSTizenWebSignageSub1
- 솔루션: signage | 환경: prd
- DNS: `BSSTizenWebSignageSub1-972785156.ap-northeast-2.elb.amazonaws.com`
- 리스너: HTTP:80, HTTPS:443
- Route53: `bss-signage-sub1.nextpay.co.kr`
- 타겟그룹:
  - `BSSTizenWebSignageSub1`: 총 1개 (healthy:1, unhealthy:0, unused:0)

### BSSTizenWebSignageSub2
- 솔루션: signage | 환경: prd
- DNS: `BSSTizenWebSignageSub2-333608192.ap-northeast-2.elb.amazonaws.com`
- 리스너: HTTP:80, HTTPS:443
- Route53: `bss-signage-sub2.nextpay.co.kr`
- 타겟그룹:
  - `BSSTizenWebSignageSub2`: 총 1개 (healthy:1, unhealthy:0, unused:0)

### LBTizenWebSignage
- 솔루션: signage | 환경: prd
- DNS: `LBTizenWebSignage-1170607247.ap-northeast-2.elb.amazonaws.com`
- 리스너: HTTPS:443
- Route53: (없음)
- 타겟그룹:
  - `RTTizenWebSignage`: 총 1개 (healthy:1, unhealthy:0, unused:0)

### NCount-Live-ELB
- 솔루션: ncount | 환경: prd
- DNS: `NCount-Live-ELB-967232354.ap-northeast-2.elb.amazonaws.com`
- 리스너: HTTPS:443
- Route53: `ncount-whisper.nextpay.co.kr`
- 타겟그룹:
  - `NCountWhisper-TG`: 총 1개 (healthy:0, unhealthy:0, unused:1)

### aiagent-staging
- 솔루션: aiagent | 환경: stg
- DNS: `aiagent-staging-1061935423.ap-northeast-2.elb.amazonaws.com`
- 리스너: HTTPS:443, HTTP:80
- Route53: `staging-api.nextaistore.co.kr`, `staging-system.nextaistore.co.kr`, `staging.nextaistore.co.kr`
- 타겟그룹:
  - `aiagent-api`: 총 1개 (healthy:0, unhealthy:1, unused:0)
  - `aiagent-front`: 총 1개 (healthy:0, unhealthy:1, unused:0)
  - `aiagent-system`: 총 1개 (healthy:1, unhealthy:0, unused:0)

### alb-dev-signage
- 솔루션: signage | 환경: dev
- DNS: `alb-dev-signage-299073066.ap-northeast-2.elb.amazonaws.com`
- 리스너: HTTPS:443, HTTP:80
- Route53: `dev-scms.nextpay.co.kr`, `dev-sw.nextpay.co.kr`, `signage-glitchtip.nextpay.co.kr`
- 타겟그룹:
  - `tg-dev-signage-glitchtip`: 총 1개 (healthy:1, unhealthy:0, unused:0)
  - `tg-dev-signage-scms`: 총 1개 (healthy:0, unhealthy:1, unused:0)
  - `tg-dev-signage-sw`: 총 1개 (healthy:1, unhealthy:0, unused:0)

### apiService
- 솔루션: unknown | 환경: prd
- DNS: `apiService-2020709608.ap-northeast-2.elb.amazonaws.com`
- 리스너: HTTPS:443
- Route53: `api.nextpay.co.kr`
- 타겟그룹:
  - `apiDev`: 총 1개 (healthy:0, unhealthy:0, unused:1)

### bss-cms
- 솔루션: bss | 환경: prd
- DNS: `bss-cms-174951817.ap-northeast-2.elb.amazonaws.com`
- 리스너: HTTPS:443, HTTP:80
- Route53: `bss.nextpay.co.kr`
- 타겟그룹:
  - `bss-cms`: 총 1개 (healthy:1, unhealthy:0, unused:0)

### bss-kiosk
- 솔루션: bss | 환경: prd
- DNS: `bss-kiosk-1230504848.ap-northeast-2.elb.amazonaws.com`
- 리스너: HTTPS:9100, HTTPS:443, HTTP:80
- Route53: `bss-kiosk.nextpay.co.kr`
- 타겟그룹:
  - `bss-kiosk`: 총 1개 (healthy:1, unhealthy:0, unused:0)
  - `bss-kiosk-monitoring`: 총 1개 (healthy:1, unhealthy:0, unused:0)

### bss-order
- 솔루션: bss | 환경: prd
- DNS: `bss-order-1305085170.ap-northeast-2.elb.amazonaws.com`
- 리스너: HTTPS:9100, HTTPS:443, HTTP:80
- Route53: `bss-order.nextpay.co.kr`
- 타겟그룹:
  - `bss-order`: 총 1개 (healthy:1, unhealthy:0, unused:0)
  - `bss-order-monitoring`: 총 1개 (healthy:1, unhealthy:0, unused:0)

### cms-elb
- 솔루션: cms | 환경: prd
- DNS: `cms-elb-393392806.ap-northeast-2.elb.amazonaws.com`
- 리스너: HTTP:80, HTTPS:443
- Route53: `cms.nextpay.co.kr`, `lidar.ncount.nextpay.co.kr`
- 타겟그룹:
  - `NCount`: 총 1개 (healthy:0, unhealthy:0, unused:1)

### dev-cms-elb
- 솔루션: cms | 환경: dev
- DNS: `dev-cms-elb-621823357.ap-northeast-2.elb.amazonaws.com`
- 리스너: HTTPS:443, HTTP:80
- Route53: `aiagent.nextpay.co.kr`, `aiot.nextpay.co.kr`, `dev-ncount.nextpay.co.kr`, `dev.nextpay.co.kr`, `dooh-api.nextpay.co.kr`, `staging.nextpay.co.kr`, `wine-curation-dev.nextpay.co.kr`
- 타겟그룹:
  - `Ncount-api-dev`: 총 1개 (healthy:1, unhealthy:0, unused:0)
  - `aiot-dashboard`: 총 1개 (healthy:1, unhealthy:0, unused:0)
  - `dev-cms-web-tg`: 총 1개 (healthy:0, unhealthy:0, unused:1)
  - `dooh-api`: 총 1개 (healthy:1, unhealthy:0, unused:0)
  - `staging-cms-web-tg`: 총 0개 (healthy:0, unhealthy:0, unused:0)
  - `tg-dev-cms-reverse-proxy`: 총 1개 (healthy:0, unhealthy:0, unused:0)
  - `wine-curation-dev-tg`: 총 1개 (healthy:0, unhealthy:0, unused:1)

### dev-product-8080
- 솔루션: signage | 환경: dev
- DNS: `dev-product-8080-1125016562.ap-northeast-2.elb.amazonaws.com`
- 리스너: HTTP:80, HTTPS:443
- Route53: `dev-ur.nextpay.co.kr`
- 타겟그룹:
  - `dev-product-8080`: 총 1개 (healthy:0, unhealthy:1, unused:0)

### dev-waiting-8081
- 솔루션: signage | 환경: dev
- DNS: `dev-waiting-8081-1289716708.ap-northeast-2.elb.amazonaws.com`
- 리스너: HTTP:80, HTTPS:443
- Route53: `backup-dev-waiting.nextpay.co.kr`
- 타겟그룹:
  - `dev-waiting-8081`: 총 1개 (healthy:1, unhealthy:0, unused:0)

### dev-waiting-8082
- 솔루션: signage | 환경: dev
- DNS: `dev-waiting-8082-2013353451.ap-northeast-2.elb.amazonaws.com`
- 리스너: HTTP:80, HTTPS:443
- Route53: `backup-dev-nsm.nextpay.co.kr`
- 타겟그룹:
  - `grp-dev-waiting-8082`: 총 1개 (healthy:1, unhealthy:0, unused:0)

### dooh-dev-elb
- 솔루션: dooh | 환경: dev
- DNS: `dooh-dev-elb-1993049909.ap-northeast-2.elb.amazonaws.com`
- 리스너: HTTPS:443
- Route53: `backstage.dooh.nextpay.co.kr`, `broker.dooh.nextpay.co.kr`, `chat.dooh.nextpay.co.kr`, `dashboard.dooh.nextpay.co.kr`
- 타겟그룹:
  - `adbot-dooh-dev`: 총 1개 (healthy:1, unhealthy:0, unused:0)
  - `backstage-dooh-dev`: 총 1개 (healthy:1, unhealthy:0, unused:0)
  - `broker-dooh-dev`: 총 1개 (healthy:1, unhealthy:0, unused:0)
  - `dooh-dev-dashboard-8083`: 총 1개 (healthy:1, unhealthy:0, unused:0)

### elb-knowledge-graph
- 솔루션: aiagent | 환경: prd
- DNS: `elb-knowledge-graph-131700488.ap-northeast-2.elb.amazonaws.com`
- 리스너: HTTPS:443, HTTP:80
- Route53: `graph.nextaistore.co.kr`, `graph.nextaistore.com`
- 타겟그룹:
  - `elb-knowledge-graph`: 총 1개 (healthy:0, unhealthy:1, unused:0)

### elb-srt-device-api
- 솔루션: srt | 환경: dev
- DNS: `elb-srt-device-api-387181319.ap-northeast-2.elb.amazonaws.com`
- 리스너: HTTPS:443, HTTP:80
- Route53: `dvc-srtdev.nextpay.co.kr`, `dvc.srtdev.nextpay.co.kr`, `str-srtdev.nextpay.co.kr`
- 타겟그룹:
  - `elb-srt-device-api`: 총 1개 (healthy:0, unhealthy:1, unused:0)

### elb-srt-staging-device-api
- 솔루션: srt | 환경: dev
- DNS: `elb-srt-staging-device-api-528593873.ap-northeast-2.elb.amazonaws.com`
- 리스너: HTTPS:443, HTTP:80
- Route53: `dvc-srtstg.nextpay.co.kr`
- 타겟그룹:
  - `elb-srt-staging-device-api`: 총 1개 (healthy:0, unhealthy:1, unused:0)

### elb-srt-staging-store-api
- 솔루션: srt | 환경: stg
- DNS: `elb-srt-staging-store-api-934084596.ap-northeast-2.elb.amazonaws.com`
- 리스너: HTTPS:443, HTTP:80
- Route53: `store-srtstg.nextpay.co.kr`, `str-srtstg.nextpay.co.kr`, `sys-srtstg.nextpay.co.kr`
- 타겟그룹:
  - `elb-srt-staging-store-api`: 총 1개 (healthy:0, unhealthy:1, unused:0)

### elb-srt-store-api
- 솔루션: srt | 환경: prd
- DNS: `elb-srt-store-api-1786817859.ap-northeast-2.elb.amazonaws.com`
- 리스너: HTTPS:443, HTTP:80
- Route53: `store-srtdev.nextpay.co.kr`, `sys-srtdev.nextpay.co.kr`
- 타겟그룹:
  - `elb-srt-store-api`: 총 1개 (healthy:0, unhealthy:1, unused:0)

### lb-dev-scms
- 솔루션: signage | 환경: dev
- DNS: `lb-dev-scms-698973874.ap-northeast-2.elb.amazonaws.com`
- 리스너: HTTP:80, HTTPS:443
- Route53: `dev-nsm.nextpay.co.kr`
- 타겟그룹:
  - `tg-dev-scms`: 총 1개 (healthy:0, unhealthy:1, unused:0)

### lb-dev-swaiting
- 솔루션: signage | 환경: dev
- DNS: `lb-dev-swaiting-1347573332.ap-northeast-2.elb.amazonaws.com`
- 리스너: HTTPS:443, HTTP:80
- Route53: `dev-waiting.nextpay.co.kr`
- 타겟그룹:
  - `tg-dev-swaiting`: 총 1개 (healthy:1, unhealthy:0, unused:0)

### lb-scms
- 솔루션: signage | 환경: prd
- DNS: `lb-scms-390862235.ap-northeast-2.elb.amazonaws.com`
- 리스너: HTTP:80, HTTPS:443
- Route53: `nsm.nextpay.co.kr`, `scms.nextpay.co.kr`
- 타겟그룹:
  - `tg-scms`: 총 1개 (healthy:0, unhealthy:1, unused:0)

### lb-signage-oss-cms
- 솔루션: signage | 환경: prd
- DNS: `lb-signage-oss-cms-701641332.ap-northeast-2.elb.amazonaws.com`
- 리스너: HTTPS:443, HTTP:80
- Route53: `x-signage.nextaistore.com`
- 타겟그룹:
  - `tg-signage-oss-cms`: 총 0개 (healthy:0, unhealthy:0, unused:0)

### lb-signage-oss-waiting
- 솔루션: signage | 환경: prd
- DNS: `lb-signage-oss-waiting-1181203740.ap-northeast-2.elb.amazonaws.com`
- 리스너: HTTPS:443, HTTP:80
- Route53: `x-waiting.nextaistore.com`
- 타겟그룹:
  - `tg-signage-oss-waiting`: 총 0개 (healthy:0, unhealthy:0, unused:0)

### lb-swaiting
- 솔루션: signage | 환경: prd
- DNS: `lb-swaiting-581779347.ap-northeast-2.elb.amazonaws.com`
- 리스너: HTTP:80, HTTPS:443
- Route53: `sw.nextpay.co.kr`, `waiting.nextpay.co.kr`
- 타겟그룹:
  - `tg-waiting`: 총 1개 (healthy:1, unhealthy:0, unused:0)

### ncount-dev
- 솔루션: ncount | 환경: dev
- DNS: `ncount-dev-1417780276.ap-northeast-2.elb.amazonaws.com`
- 리스너: HTTPS:443
- Route53: `hybo.lidar.ncount.nextpay.co.kr`, `poc.hybo.lidar.ncount.nextpay.co.kr`
- 타겟그룹:
  - `ncount-lidar-hybo-poc`: 총 1개 (healthy:0, unhealthy:0, unused:1)

### next-office-dev
- 솔루션: unknown | 환경: dev
- DNS: `next-office-dev-1163943400.ap-northeast-2.elb.amazonaws.com`
- 리스너: HTTPS:9100, HTTP:80, HTTPS:443
- Route53: `dev-nextoffice.nextpay.co.kr`
- 타겟그룹:
  - `next-office-dev`: 총 1개 (healthy:0, unhealthy:0, unused:1)
  - `office-dev-monitoring`: 총 1개 (healthy:0, unhealthy:0, unused:1)

### nextpay-kiosk-dev
- 솔루션: nserise | 환경: dev
- DNS: `nextpay-kiosk-dev-767343648.ap-northeast-2.elb.amazonaws.com`
- 리스너: HTTPS:443
- Route53: `nkiosk-dev.nextpay.co.kr`
- 타겟그룹:
  - `nextpay-kiosk-dev`: 총 1개 (healthy:0, unhealthy:1, unused:0)

### nextpay-nkds-dev
- 솔루션: nserise | 환경: dev
- DNS: `nextpay-nkds-dev-2073863581.ap-northeast-2.elb.amazonaws.com`
- 리스너: HTTPS:443
- Route53: `nkds-dev.nextpay.co.kr`
- 타겟그룹:
  - `nextpay-nkds-dev`: 총 1개 (healthy:1, unhealthy:0, unused:0)

### nextpay-norder-dev
- 솔루션: nserise | 환경: dev
- DNS: `nextpay-norder-dev-36629141.ap-northeast-2.elb.amazonaws.com`
- 리스너: HTTP:80, HTTPS:443
- Route53: `norder-dev.nextpay.co.kr`
- 타겟그룹:
  - `nextpay-norder-dev`: 총 1개 (healthy:0, unhealthy:1, unused:0)

### nextpay-npos-dev
- 솔루션: nserise | 환경: dev
- DNS: `nextpay-npos-dev-1407072045.ap-northeast-2.elb.amazonaws.com`
- 리스너: HTTPS:443, HTTP:80
- Route53: `npos-dev.nextpay.co.kr`
- 타겟그룹:
  - `nextpay-npos-dev-tg`: 총 1개 (healthy:0, unhealthy:1, unused:0)

### nextpay-order-ai
- 솔루션: nserise | 환경: prd
- DNS: `nextpay-order-ai-894286297.ap-northeast-2.elb.amazonaws.com`
- 리스너: HTTP:80
- Route53: `norder-ai.nextpay.co.kr`
- 타겟그룹:
  - (없음)

### nextpay-shop8080
- 솔루션: nserise | 환경: prd
- DNS: `nextpay-shop8080-1474328396.ap-northeast-2.elb.amazonaws.com`
- 리스너: HTTP:8080
- Route53: `shopvs.nextpay.co.kr`
- 타겟그룹:
  - `nextpay-shop-8080`: 총 1개 (healthy:1, unhealthy:0, unused:0)

### nseries-dev-elb
- 솔루션: nserise | 환경: dev
- DNS: `nseries-dev-elb-1604744663.ap-northeast-2.elb.amazonaws.com`
- 리스너: HTTPS:9100, HTTPS:443
- Route53: (없음)
- 타겟그룹:
  - `ncms-dev-tg`: 총 0개 (healthy:0, unhealthy:0, unused:0)
  - `nkds-dev-tg`: 총 1개 (healthy:1, unhealthy:0, unused:0)
  - `nkiosk-bss-tg`: 총 1개 (healthy:1, unhealthy:0, unused:0)
  - `nkiosk-dev-tg`: 총 1개 (healthy:0, unhealthy:1, unused:0)
  - `norder-bss-tg`: 총 1개 (healthy:1, unhealthy:0, unused:0)
  - `norder-dev-tg`: 총 1개 (healthy:0, unhealthy:1, unused:0)
  - `norder-monitoring-bss-tg`: 총 1개 (healthy:1, unhealthy:0, unused:0)
  - `npos-dev-tg`: 총 1개 (healthy:0, unhealthy:0, unused:1)
  - `socket-io-dev-tg`: 총 0개 (healthy:0, unhealthy:0, unused:0)

### nseries-elb
- 솔루션: nserise | 환경: prd
- DNS: `nseries-elb-894077439.ap-northeast-2.elb.amazonaws.com`
- 리스너: HTTPS:443
- Route53: (없음)
- 타겟그룹:
  - `ncms-tg`: 총 0개 (healthy:0, unhealthy:0, unused:0)
  - `nkds-tg`: 총 1개 (healthy:1, unhealthy:0, unused:0)
  - `nkiosk-tg`: 총 1개 (healthy:0, unhealthy:1, unused:0)
  - `norder-tg`: 총 1개 (healthy:0, unhealthy:1, unused:0)
  - `npos-tg`: 총 0개 (healthy:0, unhealthy:0, unused:0)
  - `socket-io-tg`: 총 0개 (healthy:0, unhealthy:0, unused:0)

### product-8080
- 솔루션: signage | 환경: prd
- DNS: `product-8080-1041987207.ap-northeast-2.elb.amazonaws.com`
- 리스너: HTTPS:443, HTTP:80
- Route53: `ur.nextpay.co.kr`
- 타겟그룹:
  - `product-8080`: 총 1개 (healthy:0, unhealthy:1, unused:0)

### socket-dev
- 솔루션: nserise | 환경: dev
- DNS: `socket-dev-587569175.ap-northeast-2.elb.amazonaws.com`
- 리스너: HTTPS:443, HTTP:80
- Route53: `socket-dev.nextpay.co.kr`
- 타겟그룹:
  - `socket-dev-tg`: 총 0개 (healthy:0, unhealthy:0, unused:0)

### srt-device-api
- 솔루션: srt | 환경: dev
- DNS: `srt-device-api-1442175753.ap-northeast-2.elb.amazonaws.com`
- 리스너: HTTPS:443, HTTPS:9323, HTTP:80
- Route53: (없음)
- 타겟그룹:
  - `srt-dev-monitoring`: 총 0개 (healthy:0, unhealthy:0, unused:0)
  - `srt-device-api`: 총 0개 (healthy:0, unhealthy:0, unused:0)

### srt-store-api
- 솔루션: srt | 환경: prd
- DNS: `srt-store-api-1085020352.ap-northeast-2.elb.amazonaws.com`
- 리스너: HTTPS:443
- Route53: (없음)
- 타겟그룹:
  - `srt-store-api`: 총 0개 (healthy:0, unhealthy:0, unused:0)

### srt-system-api
- 솔루션: srt | 환경: prd
- DNS: `srt-system-api-932950302.ap-northeast-2.elb.amazonaws.com`
- 리스너: HTTPS:443
- Route53: (없음)
- 타겟그룹:
  - `srt-system-api`: 총 0개 (healthy:0, unhealthy:0, unused:0)

### waiting-8081
- 솔루션: signage | 환경: prd
- DNS: `waiting-8081-442550494.ap-northeast-2.elb.amazonaws.com`
- 리스너: HTTPS:443, HTTP:80
- Route53: `backup-waiting.nextpay.co.kr`
- 타겟그룹:
  - `waiting-8081`: 총 1개 (healthy:1, unhealthy:0, unused:0)

### waiting-8082
- 솔루션: signage | 환경: prd
- DNS: `waiting-8082-1157880650.ap-northeast-2.elb.amazonaws.com`
- 리스너: HTTP:80, HTTPS:443
- Route53: `backup-nsm.nextpay.co.kr`
- 타겟그룹:
  - `grp-waiting-8082`: 총 1개 (healthy:1, unhealthy:0, unused:0)

### webService
- 솔루션: nserise | 환경: prd
- DNS: `webService-2021174328.ap-northeast-2.elb.amazonaws.com`
- 리스너: HTTPS:443, HTTP:80
- Route53: `back-office.nextpay.co.kr`, `ca.nextpay.co.kr`, `ncms.nextpay.co.kr`, `nkds-new.nextpay.co.kr`, `nkds.nextpay.co.kr`, `nkiosk.nextpay.co.kr`, `nkioskkinkos.nextpay.co.kr`, `norder.nextpay.co.kr`, `npos.nextpay.co.kr`, `salesking.nextpay.co.kr`, `shop.nextpay.co.kr`, `socket.nextpay.co.kr`
- 타겟그룹:
  - `ca-site`: 총 0개 (healthy:0, unhealthy:0, unused:0)
  - `n-kiosk-kinkos`: 총 1개 (healthy:1, unhealthy:0, unused:0)
  - `next-office`: 총 0개 (healthy:0, unhealthy:0, unused:0)
  - `nextpay-kiosk`: 총 1개 (healthy:0, unhealthy:1, unused:0)
  - `nextpay-ncms`: 총 0개 (healthy:0, unhealthy:0, unused:0)
  - `nextpay-nkds`: 총 1개 (healthy:1, unhealthy:0, unused:0)
  - `nextpay-nkds-new`: 총 1개 (healthy:1, unhealthy:0, unused:0)
  - `nextpay-norder-tg`: 총 1개 (healthy:0, unhealthy:1, unused:0)
  - `nextpay-npos`: 총 1개 (healthy:0, unhealthy:1, unused:0)
  - `nextpay-salesking`: 총 0개 (healthy:0, unhealthy:0, unused:0)
  - `nextpay-shop`: 총 1개 (healthy:1, unhealthy:0, unused:0)
  - `socket-tg`: 총 0개 (healthy:0, unhealthy:0, unused:0)

### wine-curation-elb
- 솔루션: wine | 환경: prd
- DNS: `wine-curation-elb-1718355595.ap-northeast-2.elb.amazonaws.com`
- 리스너: HTTPS:443
- Route53: `wine-curation.nextpay.co.kr`
- 타겟그룹:
  - `wine-curation-tg`: 총 1개 (healthy:1, unhealthy:0, unused:0)

### ws2025-lb
- 솔루션: ws2025 | 환경: prd
- DNS: `ws2025-lb-251643516.ap-northeast-2.elb.amazonaws.com`
- 리스너: HTTPS:443
- Route53: `ws2025-api.nextpay.co.kr`, `ws2025-web.nextpay.co.kr`
- 타겟그룹:
  - `ws2025-api`: 총 1개 (healthy:0, unhealthy:0, unused:1)
  - `ws2025-tg`: 총 1개 (healthy:0, unhealthy:0, unused:1)
