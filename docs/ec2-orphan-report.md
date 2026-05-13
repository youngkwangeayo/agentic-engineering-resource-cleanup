# EC2 미아/중지 인스턴스 리포트

> 조사 시각: 2026-05-13 / 리전: ap-northeast-2 / 계정: 365485194891
> 최종 갱신: 2026-05-13 (Section 2/3 교차 검증 반영)
> 환율: 1 USD = 1,400 KRW

---

# 📌 핵심 요약 (Executive Summary)

## 1. 예상 절감 (ALB + EC2 + EBS 통합)

> 단가: ALB $16.20/mo (약 22,680원/월, $194.40/yr) · EC2 on-demand (ap-northeast-2) · EBS gp3 $0.0912/GB, gp2 $0.114/GB · 환율 1 USD = 1,400 KRW

### 1.1 카테고리별 절감

| 카테고리 | ALB 수 | EC2 수 | EBS(GB) | ALB 월 | EC2 월 | EBS 월 | **월 합계** | **연 합계 (USD)** | **연 합계 (KRW)** |
|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|
| **bss 폐기** (전체) | 6 | 4 (running 4) | 192 | $97.20 (약 136,080원) | $223.38 (약 312,730원) | $21.89 (약 30,650원) | **$342.47 (약 479,460원)** | **$4,109.64** | **약 575만원** |
| **ws2025 폐기** (전체) | 1 | 1 (stopped 1) | 8 | $16.20 (약 22,680원) | $0.00 (0원) | $0.73 (약 1,020원) | **$16.93 (약 23,700원)** | **$203.16** | **약 28만원** |
| **wine 폐기** (전체) | 1 | 2 (running 1, stopped 1) | 16 | $16.20 (약 22,680원) | $49.64 (약 69,500원) | $1.46 (약 2,040원) | **$67.30 (약 94,220원)** | **$807.60** | **약 113만원** |
| **nserise 폐기** (전체) | 10 | 11 (running 11) | 205 | $162.00 (약 226,800원) | $173.74 (약 243,240원) | $23.01 (약 32,210원) | **$358.75 (약 502,250원)** | **$4,305.00** | **약 603만원** |
| **미아 EC2 + 관련 ALB** (4 솔루션 외) | 5 | 8 (stopped 8) | 346 | $81.00 (약 113,400원) | $0.00 (0원) | $35.20 (약 49,280원) | **$116.20 (약 162,680원)** | **$1,394.40** | **약 195만원** |
| **기타 action=삭제 ALB** (signage 8 + srt 7) | 15 | - | - | $243.00 (약 340,200원) | $0.00 (0원) | $0.00 (0원) | **$243.00 (약 340,200원)** | **$2,916.00** | **약 408만원** |
| **총합** | **38** | **26** | **767** | **$615.60 (약 861,840원)** | **$446.76 (약 625,460원)** | **$82.29 (약 115,210원)** | **$1,144.65 (약 1,602,510원)** | **$13,735.80** | **약 1,923만원/년** |

> 정리 후 잔존 ALB: 48 - 38 = **10개** (signage 합치기/유지 6, cms 미정 1, aiagent 미정 1, 기타 미정 2)
> 추가로 signage 합치기 4→1 통합 시 ALB 3개 추가 절감 가능 (약 $583/년, 약 82만원/년)

### 1.2 핵심 숫자 (월 / 연)

| 항목 | 월 (USD) | 월 (KRW) | 연 (USD) | 연 (KRW) |
|---|---:|---:|---:|---:|
| 중지서비스 4종 폐기 소계 | $785.45 | 약 1,099,630원 | $9,425.40 | 약 1,320만원 |
| 미아 EC2 정리 | $116.20 | 약 162,680원 | $1,394.40 | 약 195만원 |
| 기타 ALB 삭제 | $243.00 | 약 340,200원 | $2,916.00 | 약 408만원 |
| **총 절감** | **$1,144.65** | **약 1,602,510원** | **$13,735.80** | **약 1,923만원** |

### 1.3 가정 및 주의

- ALB 비용은 고정 시간당 요금만 반영 (LCU 비용은 대부분 미미하여 제외)
- EC2는 running 인스턴스만 비용 계상, stopped는 EBS만 발생
- 중지서비스 폐기 시 해당 솔루션의 **모든 ALB(삭제/합치기/미정 포함)** 가 제거되는 시나리오
- Savings Plan / RI 적용 여부에 따라 실 EC2 절감은 변동될 수 있음
- 환율 1 USD = 1,400 KRW 기준 (시장 환율 변동 시 재계산 필요)

---

## 2. 중지 솔루션별 인스턴스 (bss / ws2025 / wine / nserise)

### 2.1 bss
- 인스턴스 수: 4
- 월 EC2 비용 (running 합): **$223.38 (약 312,730원)**
- 월 EBS 비용: **$21.89 (약 30,650원)**

| ID | Name | 타입 | 상태 | EBS(GB) | EC2 월 | EBS 월 |
|---|---|---|---|---|---|---|
| i-082089d36ca5ad84b | bss-kiosk | t3.micro | running | 16 | $12.41 (약 17,370원) | $1.82 (약 2,550원) |
| i-01826b3f3e2fc1ae7 | bss-order | t3.micro | running | 16 | $12.41 (약 17,370원) | $1.82 (약 2,550원) |
| i-0ba059f4732a9c43f | bss_Signage | t3.large | running | 32 | $99.28 (약 138,990원) | $3.65 (약 5,110원) |
| i-094898923435487eb | bss_cms | t3.large | running | 128 | $99.28 (약 138,990원) | $14.59 (약 20,430원) |

### 2.2 ws2025
- 인스턴스 수: 1
- 월 EC2 비용 (running 합): **$0.00 (0원)**
- 월 EBS 비용: **$0.73 (약 1,020원)**

| ID | Name | 타입 | 상태 | EBS(GB) | EC2 월 | EBS 월 |
|---|---|---|---|---|---|---|
| i-0676dc1f436330caa | ws2025 | t2.micro | stopped | 8 | $0.00 (0원) | $0.73 (약 1,020원) |

### 2.3 wine
- 인스턴스 수: 2
- 월 EC2 비용 (running 합): **$49.64 (약 69,500원)**
- 월 EBS 비용: **$1.46 (약 2,040원)**

| ID | Name | 타입 | 상태 | EBS(GB) | EC2 월 | EBS 월 |
|---|---|---|---|---|---|---|
| i-0834a1522d1b35443 | wine-curation-dev | t2.small | stopped | 8 | $0.00 (0원) | $0.73 (약 1,020원) |
| i-04cf09d620cf1ec28 | wine-curation | t3.medium | running | 8 | $49.64 (약 69,500원) | $0.73 (약 1,020원) |

### 2.4 nserise
- 인스턴스 수: 11
- 월 EC2 비용 (running 합): **$173.74 (약 243,240원)**
- 월 EBS 비용: **$23.01 (약 32,210원)**

| ID | Name | 타입 | 상태 | EBS(GB) | EC2 월 | EBS 월 |
|---|---|---|---|---|---|---|
| i-01d12fac0e9097dfa | nextpay-npos-dev | t3.micro | running | 30 | $12.41 (약 17,370원) | $3.42 (약 4,790원) |
| i-0be06c6ccccc238f3 | n-kiosk-kinkos | t3.micro | running | 15 | $12.41 (약 17,370원) | $1.71 (약 2,390원) |
| i-09335e4f3bbf01908 | nextpay-kiosk-dev | t3.micro | running | 30 | $12.41 (약 17,370원) | $3.42 (약 4,790원) |
| i-05e80cfd2d77f9f84 | nextpay-shop | t3.micro | running | 16 | $12.41 (약 17,370원) | $1.46 (약 2,040원) |
| i-0999f0a959daa5507 | nextpay-kiosk-ai | t3.micro | running | 16 | $12.41 (약 17,370원) | $1.82 (약 2,550원) |
| i-0f33238e06e3d189b | nextpay-nkds-dev | t3.micro | running | 15 | $12.41 (약 17,370원) | $1.71 (약 2,390원) |
| i-0aa991fce295e1958 | nextpay-norder | t3.micro | running | 15 | $12.41 (약 17,370원) | $1.71 (약 2,390원) |
| i-04ef73dcb67e85b4f | nextpay-norder-dev | t3.micro | running | 30 | $12.41 (약 17,370원) | $3.42 (약 4,790원) |
| i-0afe06698963c9070 | nextpay-nkds-env | t3.micro | running | 15 | $12.41 (약 17,370원) | $1.71 (약 2,390원) |
| i-0ffc74bdb94236f55 | nextpay-kiosk | t3.micro | running | 8 | $12.41 (약 17,370원) | $0.91 (약 1,270원) |
| i-0a4457d68978b1e45 | nextpay-npos-01 | t3.medium | running | 15 | $49.64 (약 69,500원) | $1.71 (약 2,390원) |

**솔루션 폐기 시 영향도 추정**

- **bss**: ALB 다수(`BSSTizen*`, `bss-cms`, `bss-order` 등)가 합치기/유지로 표시되어 솔루션은 살아있음. 인스턴스 즉시 폐기는 권장하지 않음.
- **ws2025**: `ws2025` 인스턴스가 stopped 이며 TG도 모두 unused → 사실상 폐기 상태. EC2 종료 + EBS 삭제로 즉시 정리 가능.
- **wine**: `wine-curation-dev` stopped. ALB(`wine-curation-elb`)는 합치기 대상. dev 인스턴스는 폐기 검토.
- **nserise**: nseries 계열 ALB 다수 삭제 대상. 그러나 인스턴스(`nextpay-shop`, `bss-order` 등)는 다른 nserise/bss ALB에도 묶여 있어 인스턴스 단위 폐기는 신중 필요.

---

> 📖 이하는 상세 분석입니다. 위 두 섹션만으로 충분하다면 여기까지 읽으시면 됩니다.

---

# 📊 상세 분석 (Detailed Breakdown)

## 3. 검증 노트

이전 리포트의 "502/503 룰 삭제로 미아 후보"는 "TG 헬스 상태"만으로 판정했기 때문에, 해당 EC2가 다른 ALB의 TG에 healthy로 등록되어 있는지를 확인하지 않았다. 이를 보완하기 위해 18개 후보 인스턴스 각각에 대해 instance-type TG 103개를 전수 스캔하여 "EC2 → 묶인 모든 TG → TG가 연결된 ALB → 해당 ALB의 action(`data/entries.json`)"을 매핑했다. 그 결과 (1) 모든 후보 인스턴스에서 healthy 상태 TG 등록은 0건이었으므로 "부분 유지"로 빠지는 인스턴스는 없으며, (2) TG가 연결된 ALB의 action이 전부 `삭제` 또는 ALB 미연결이면 **완전 미아**, action이 `합치기`/`미정` ALB가 하나라도 포함되면 **그레이존**(합치기/미정 결과에 따라 처리 방향 결정 필요)으로 재분류했다.

## 4. 인스턴스 분류 요약

| 구분 | 인스턴스 수 | running | stopped | EC2 월비용 | EBS 월비용 |
|---|---|---|---|---|---|
| 완전 미아 (Section 5) | 8 | 0 | 8 | $0.00 (0원) | $35.20 (약 49,280원) |
| 그레이존 — 합치기/미정 ALB에 unhealthy/unused (Section 6) | 10 | 7 | 3 | $111.69 (약 156,370원) | $48.88 (약 68,430원) |
| 부분 유지 (다른 ALB에 healthy로 살아있음) | 0 | 0 | 0 | $0.00 (0원) | $0.00 (0원) |
| 현재 stopped 전체 (Section 7) | 20 | 0 | 20 | $0.00 (0원) | $78.71 (약 110,190원) |

- **합산 절감 가능 (running EC2 종료 + 미아 EBS 정리)**: 약 **$225.60/월 (약 315,840원/월)**
- 부분 유지 항목 없음 → 미아 인스턴스 수는 검증 전후 동일. 다만 Section 6 10건은 합치기/미정 ALB에 발이 걸쳐 있어 **즉시 폐기 전 ALB 합치기/결정 후행**이 필요한 그레이존으로 격하.

## 5. 완전 미아 EC2

> 묶인 모든 TG가 (a) 어떤 ALB에도 연결되지 않음 OR (b) `action=삭제` ALB에만 연결됨. 모든 등록 상태도 unhealthy/unused. 합치기/유지 ALB에 발이 걸쳐 있지 않으므로 ALB 정리 시 안전하게 함께 폐기 가능.

| 인스턴스 ID | Name | 상태 | 타입 | AZ | EBS합(GB) | 묶인 TG → ALB → action | 솔루션 |
|---|---|---|---|---|---|---|---|
| i-058e2739a704c0e2a | next-office-dev | stopped | t3.large | ap-northeast-2b | 128 | next-office-dev(unused) → next-office-dev → 삭제 / office-dev-monitoring(unused) → next-office-dev → 삭제 | unknown |
| i-0e9c0536c4b85f996 | cms-web-11 | stopped | c5.large | ap-northeast-2a | 32 | cms-web-tg(unused) → (ALB 미연결) | cms |
| i-0ea629ff65a69843f | cms-web-22 | stopped | c5.large | ap-northeast-2a | 32 | cms-web-tg(unused) → (ALB 미연결) | cms |
| i-04fa271011856d980 | cms-web-33 | stopped | c5.large | ap-northeast-2a | 32 | cms-web-tg(unused) → (ALB 미연결) | cms |
| i-0e179dcb732712136 | api-dev | stopped | t2.micro | ap-northeast-2a | 32 | apiDev(unused) → apiService → 삭제 | unknown |
| i-04e5c1e0c0770cdea | NCount Lidar Live | stopped | t4g.small | ap-northeast-2c | 30 | NCount(unused) → cms-elb → 삭제 | ncount |
| i-029126fd8b34d4257 | OpenAI-Whisper | stopped | t4g.small | ap-northeast-2c | 30 | NCountWhisper-TG(unused) → NCount-Live-ELB → 삭제 | ncount |
| i-0be72303280d93ac1 | ncount-lidar-hybo-poc | stopped | t4g.nano | ap-northeast-2c | 30 | ncount-lidar-hybo-poc(unused) → ncount-dev → 삭제 | ncount |

**비고**

- 8개 모두 현재 stopped — running EC2 비용 절감 효과는 없으나 EBS 비용은 계속 발생 중.
- `cms-web-11/22/33` (c5.large 3대)는 TG=`cms-web-tg`에 묶여있지만 이 TG는 이미 어떤 ALB에도 연결되어 있지 않음 → ALB 삭제와 무관하게 이미 미아 상태.
- `next-office-dev`는 t3.large + EBS 128GB로 EBS 절감 효과가 가장 큼.
- ncount 계열 3대(`NCount Lidar Live`, `OpenAI-Whisper`, `ncount-lidar-hybo-poc`)는 동일 비즈니스(라이다 POC) 묶음으로 보임 → 일괄 폐기 검토.

## 6. 그레이존 EC2 (합치기/미정 ALB에 unhealthy 등록)

> 묶인 TG 중 하나 이상이 `action=합치기` 또는 `action=미정` ALB에 연결되어 있으나, 모든 상태가 unhealthy/unused. 합치기 후 신 ALB가 이 EC2를 계속 사용할 의도인지(=서비스 복구 대상), 아니면 함께 폐기할 대상인지 결정 필요. (모든 후보에서 healthy 등록은 0건이었으므로 "부분 유지"로 빠진 EC2는 없음.)

| 인스턴스 ID | Name | 상태 | 타입 | AZ | EBS합(GB) | 묶인 TG → ALB → action | 솔루션 |
|---|---|---|---|---|---|---|---|
| i-0676dc1f436330caa | ws2025 | stopped | t2.micro | ap-northeast-2a | 8 | ws2025-api(unused) → ws2025-lb → 합치기 / ws2025-tg(unused) → ws2025-lb → 합치기 | ws2025 |
| i-09335e4f3bbf01908 | nextpay-kiosk-dev | running | t3.micro | ap-northeast-2c | 30 | nextpay-kiosk-dev(unhealthy) → nextpay-kiosk-dev → 미정 / nkiosk-dev-tg(unhealthy) → nseries-dev-elb → 삭제 | nserise |
| i-0a4457d68978b1e45 | nextpay-npos-01 | running | t3.medium | ap-northeast-2a | 15 | nextpay-npos(unhealthy) → webService → 미정 | nserise |
| i-071e75b8b961eb20a | srt-staging-ec2 | running | m4.xlarge | ap-northeast-2a | 200 | aiagent-api(unhealthy) → aiagent-staging → 합치기 / aiagent-front(unhealthy) → aiagent-staging → 합치기 / aiagent-system(unhealthy) → aiagent-staging → 합치기 / aiagent-weraser(unused) → (ALB 미연결) / elb-knowledge-graph(unhealthy) → elb-knowledge-graph → 미정 / elb-srt-staging-device-api(unhealthy) → elb-srt-staging-device-api → 삭제 / elb-srt-staging-store-api(unhealthy) → elb-srt-staging-store-api → 삭제 | srt |
| i-01d12fac0e9097dfa | nextpay-npos-dev | running | t3.micro | ap-northeast-2b | 30 | nextpay-npos-dev-tg(unhealthy) → nextpay-npos-dev → 합치기 / npos-dev-tg(unused) → nseries-dev-elb → 삭제 | nserise |
| i-0ffc74bdb94236f55 | nextpay-kiosk | running | t3.micro | ap-northeast-2a | 8 | nextpay-kiosk(unhealthy) → webService → 미정 / nkiosk-tg(unhealthy) → nseries-elb → 삭제 | nserise |
| i-04ef73dcb67e85b4f | nextpay-norder-dev | running | t3.micro | ap-northeast-2c | 30 | nextpay-norder-dev(unhealthy) → nextpay-norder-dev → 합치기 / norder-dev-tg(unhealthy) → nseries-dev-elb → 삭제 | nserise |
| i-0ac99b083501baf48 | cms-web-dev | stopped | t3.large | ap-northeast-2c | 128 | dev-cms-web-tg(unused) → dev-cms-elb → 미정 / cms-dev-monitoring(unused) → (ALB 미연결) | cms |
| i-0aa991fce295e1958 | nextpay-norder | running | t3.micro | ap-northeast-2c | 15 | nextpay-norder-tg(unhealthy) → webService → 미정 / norder-tg(unhealthy) → nseries-elb → 삭제 | nserise |
| i-0834a1522d1b35443 | wine-curation-dev | stopped | t2.small | ap-northeast-2a | 8 | wine-curation-dev-tg(unused) → dev-cms-elb → 미정 | wine |

**비고**

- 본 섹션 10건은 모두 unhealthy/unused 상태이므로 트래픽을 받고 있지 않다. 그러나 합치기/미정 ALB에 발이 걸쳐 있어 "ALB 합치기 완료 후 신 ALB가 이 EC2를 살릴 계획인지"를 솔루션 담당자에게 확인해야 한다.
- `srt-staging-ec2` (i-071e75b8b961eb20a)는 7개 TG에 걸쳐 있으며 그 중 3개가 `aiagent-staging`(합치기) 소속 → aiagent 합치기 후속 작업 결과 결정.
- nserise 계열 6대(`nextpay-kiosk-dev`, `nextpay-npos-01`, `nextpay-npos-dev`, `nextpay-kiosk`, `nextpay-norder-dev`, `nextpay-norder`)는 동일 패턴 — 개별 ALB(합치기/미정) + `nseries-*-elb`(삭제) 이중 등록. `webService`/`nseries-*-elb` 합치기/미정 결정 시 일괄 처리 가능.
- 신 ALB가 이 EC2를 사용할 계획이 없다면 → 완전 미아로 격상하여 함께 폐기.

## 7. 현재 stopped 상태 EC2 전체

| 인스턴스 ID | Name | 타입 | AZ | EBS합(GB) | 중지 시점 (UTC) | EBS 월비용 | 솔루션 |
|---|---|---|---|---|---|---|---|
| i-0ac99b083501baf48 | cms-web-dev | t3.large | ap-northeast-2c | 128 | 2025-12-02 08:17:15 GMT | $14.59 (약 20,430원) | cms |
| i-058e2739a704c0e2a | next-office-dev | t3.large | ap-northeast-2b | 128 | 2024-09-04 09:27:16 GMT | $14.59 (약 20,430원) | unknown |
| i-07406a1cc16cf2c5a | [사이니지 개발서버 II - Docker] - 일시정지 | t2.micro | ap-northeast-2a | 64 | 2025-03-04 02:41:24 GMT | $7.30 (약 10,220원) | unknown |
| i-053e5e15a87ad0ef7 | PG-AP01 | t2.micro | ap-northeast-2c | 32 | 2023-01-30 08:33:23 GMT | $3.65 (약 5,110원) | unknown |
| i-0a30d192f464f7539 | PG-DB01 | t2.micro | ap-northeast-2c | 32 | 2023-01-30 08:33:23 GMT | $3.65 (약 5,110원) | unknown |
| i-0e179dcb732712136 | api-dev | t2.micro | ap-northeast-2a | 32 | 2022-06-09 07:02:04 GMT | $3.65 (약 5,110원) | unknown |
| i-0e9c0536c4b85f996 | cms-web-11 | c5.large | ap-northeast-2a | 32 | 2025-12-02 08:11:33 GMT | $2.92 (약 4,090원) | cms |
| i-0ea629ff65a69843f | cms-web-22 | c5.large | ap-northeast-2a | 32 | 2025-12-02 08:12:39 GMT | $2.92 (약 4,090원) | cms |
| i-04fa271011856d980 | cms-web-33 | c5.large | ap-northeast-2a | 32 | 2025-12-02 08:13:07 GMT | $2.92 (약 4,090원) | cms |
| i-04e5c1e0c0770cdea | NCount Lidar Live | t4g.small | ap-northeast-2c | 30 | 2025-06-27 03:28:10 GMT | $2.74 (약 3,840원) | ncount |
| i-0e963ae9e5aa00cbc | NCount Lidar Live DB | t4g.small | ap-northeast-2c | 30 | 2025-06-27 03:28:10 GMT | $2.74 (약 3,840원) | ncount |
| i-029126fd8b34d4257 | OpenAI-Whisper | t4g.small | ap-northeast-2c | 30 | 2025-06-27 03:28:10 GMT | $2.74 (약 3,840원) | ncount |
| i-0df0d40daa44d8e5b | [매직인포 윈도우 서버]운영서버 | t2.micro | ap-northeast-2c | 30 | 2024-08-29 06:48:38 GMT | $3.42 (약 4,790원) | unknown |
| i-0be72303280d93ac1 | ncount-lidar-hybo-poc | t4g.nano | ap-northeast-2c | 30 | 2025-06-27 03:28:10 GMT | $2.74 (약 3,840원) | ncount |
| i-057f34714835e46b1 | nextpay-salesking-dev | t3.micro | ap-northeast-2c | 30 | 2024-09-05 01:50:35 GMT | $3.42 (약 4,790원) | unknown |
| i-03c0b900e85b76321 | server_monitoring | t2.micro | ap-northeast-2a | 16 | 2024-09-03 09:37:52 GMT | $1.46 (약 2,040원) | unknown |
| i-0aee6750bf3371ba9 | MTOUCH-Test | t2.micro | ap-northeast-2a | 8 | 2024-08-29 07:00:20 GMT | $0.91 (약 1,270원) | unknown |
| i-0834a1522d1b35443 | wine-curation-dev | t2.small | ap-northeast-2a | 8 | 2025-12-10 00:24:54 GMT | $0.73 (약 1,020원) | wine |
| i-0676dc1f436330caa | ws2025 | t2.micro | ap-northeast-2a | 8 | 2026-01-20 08:58:23 GMT | $0.73 (약 1,020원) | ws2025 |
| i-02971c39d7441d75e | 유니버셜로봇-운영 | t2.micro | ap-northeast-2c | 8 | 2024-08-29 06:50:26 GMT | $0.91 (약 1,270원) | unknown |

**합계: 20대, EBS 총 740GB, 월 EBS 비용 약 $78.71 (약 110,190원)**

- 가장 오래된 중지: `api-dev` (2022-06-09), `PG-AP01/PG-DB01` (2023-01-30) — 3년 이상 미사용.
- 가장 큰 EBS: `next-office-dev`, `cms-web-dev` (각 128GB).
- AMI/스냅샷으로 백업 후 인스턴스 종료(terminate) + 미사용 EBS 삭제 권고.

---

## 데이터 출처

- `aws elbv2 describe-target-groups --region ap-northeast-2` → 142개 TG (instance 타입 103개 health 조회)
- `aws elbv2 describe-target-health --target-group-arn <arn>` (103건 병렬 호출)
- `aws ec2 describe-instances --region ap-northeast-2` → 70개
- `aws ec2 describe-volumes --region ap-northeast-2` → 91개
- 분류 입력: `data/entries.json` (48개 ALB, action=삭제 26개 기준)
- 조사 시각: 2026-05-13 (Asia/Seoul)
- 관련 리포트: `docs/alb-cleanup-report.md` (ALB 정리 분류 결과)
