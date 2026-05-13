# ALB 리소스 정리 리포트

> 작성일: 2026-05-13
> 대상 리전: ap-northeast-2 (서울)
> 대상: k8s ALB Controller 관리 ALB를 제외한 전체 ALB

---

# 📌 핵심 요약 (Executive Summary)

## 1. 요약

| 항목 | 값 |
|---|---|
| 전체 ALB 수 | **48개** |
| 삭제 대상 | **26개** |
| 합치기 대상 | **15개** |
| 유지 (유지 2 + 미정 5) | **7개** |
| **정리 후 ALB 수** | **총 7개** |

> 합치기 그룹 통합 시 추가 감축 가능 (signage 합치기 4→1 등).

## 2. 예상 절감 효과

| 항목 | 값 |
|---|---|
| 정리되는 ALB 수 | **41개** (삭제 26 + 합치기 15) |
| 연간 절감 (USD) | 41 × $194.40 = **$7,970** |
| **연간 절감 (KRW, 1USD=1,400원)** | **약 1,116만원** |
| 월간 절감 | **약 93만원** |

> 위는 ALB 고정 시간당 요금만 반영했습니다.
> LCU(트래픽) 비용은 별도이며, 삭제 대상 대부분이 unused/error 상태이므로 추가 절감 여지가 있습니다.
> 실 절감액은 **연 1,200만원 이상**으로 추정됩니다.
> EC2/EBS까지 포함한 전체 절감액은 별도 리포트 `docs/ec2-orphan-report.md` 참조 (연 약 1,923만원).

---

> 📖 이하는 상세 분석입니다. 위 두 표만으로 충분하다면 여기까지 읽으시면 됩니다.

---

# 📊 상세 분석 (Detailed Breakdown)

## 3. 분포 분석

### 3.1 조치 분포

| 조치 | 개수 | 비중 |
|---|---:|---:|
| 삭제 | 26 | 54% |
| 합치기 | 15 | 31% |
| 미정 | 5 | 10% |
| 유지 | 2 | 4% |

### 3.2 환경 분포

| 환경 | 개수 |
|---|---:|
| prd | 26 |
| dev | 19 |
| stg | 3 |

### 3.3 상태 분포

| 상태 | 개수 | 설명 |
|---|---:|---|
| error | 14 | Route53 레코드는 있으나 5xx 응답 |
| healthy | 13 | 정상 응답 |
| no_record | 8 | Route53 레코드 없음 |
| unhealthy | 6 | Target Group unhealthy |
| no_target | 6 | Target Group에 타겟 없음 |
| unknown | 1 | 응답 확인 불가 |

---

## 4. 솔루션별 현황

| 솔루션 | 합계 | 삭제 | 합치기 | 유지 | 미정 |
|---|---:|---:|---:|---:|---:|
| signage | 14 | 8 | 4 | 2 | 0 |
| nserise | 10 | 5 | 3 | 0 | 2 |
| srt | 7 | 7 | 0 | 0 | 0 |
| bss | 6 | 1 | 4 | 0 | 1 |
| cms | 2 | 1 | 0 | 0 | 1 |
| ncount | 2 | 2 | 0 | 0 | 0 |
| aiagent | 2 | 0 | 1 | 0 | 1 |
| unknown | 2 | 2 | 0 | 0 | 0 |
| wine | 1 | 0 | 1 | 0 | 0 |
| ws2025 | 1 | 0 | 1 | 0 | 0 |
| dooh | 1 | 0 | 1 | 0 | 0 |

---

## 5. 조치 상세

### 5.1 삭제 대상 (26개)

#### signage (8개)
| ALB | 환경 | 상태 | 레코드 |
|---|---|---|---|
| LBTizenWebSignage | prd | no_record | - |
| product-8080 | prd | no_record | - |
| lb-signage-oss-cms | prd | no_target | x-signage.nextaistore.com |
| lb-signage-oss-waiting | prd | no_target | x-waiting.nextaistore.com |
| dev-product-8080 | dev | error | dev-ur.nextpay.co.kr |
| dev-waiting-8081 | dev | healthy | backup-dev-waiting.nextpay.co.kr |
| dev-waiting-8082 | dev | healthy | backup-dev-nsm.nextpay.co.kr |
| lb-dev-scms | dev | no_record | - |

#### srt (7개)
| ALB | 환경 | 상태 | 레코드 |
|---|---|---|---|
| srt-system-api | prd | no_target | - |
| srt-store-api | prd | no_target | - |
| elb-srt-store-api | prd | error | store-srtdev.nextpay.co.kr, sys-srtdev.nextpay.co.kr |
| srt-device-api | dev | no_target | - |
| elb-srt-device-api | dev | error | dvc-srtdev.nextpay.co.kr, str-srtdev.nextpay.co.kr, dvc.srtdev.nextpay.co.kr |
| elb-srt-staging-device-api | dev | unhealthy | dvc-srtstg.nextpay.co.kr |
| elb-srt-staging-store-api | stg | error | store-srtstg.nextpay.co.kr, str-srtstg.nextpay.co.kr, sys-srtstg.nextpay.co.kr |

#### nserise (5개)
| ALB | 환경 | 상태 | 레코드 |
|---|---|---|---|
| nextpay-order-ai | prd | unknown | norder-ai.nextpay.co.kr |
| nextpay-shop8080 | prd | error | shopvs.nextpay.co.kr |
| nseries-elb | prd | no_record | - |
| socket-dev | dev | no_target | socket-dev.nextpay.co.kr |
| nseries-dev-elb | dev | no_record | - |

#### ncount (2개)
| ALB | 환경 | 상태 | 레코드 |
|---|---|---|---|
| NCount-Live-ELB | prd | error | ncount-whisper.nextpay.co.kr |
| ncount-dev | dev | no_record | - |

#### unknown (2개)
| ALB | 환경 | 상태 | 레코드 |
|---|---|---|---|
| apiService | prd | no_record | - |
| next-office-dev | dev | error | dev-nextoffice.nextpay.co.kr |

#### cms (1개) / bss (1개)
| ALB | 솔루션 | 환경 | 상태 | 레코드 |
|---|---|---|---|---|
| cms-elb | cms | prd | error | cms.nextpay.co.kr, lidar.ncount.nextpay.co.kr |
| BSSTizenWebSignage | bss | prd | error | bss-signage.nextpay.co.kr |

---

### 5.2 합치기 대상 (15개)

#### signage (4개) — 통합 후보
| ALB | 환경 | 상태 | 레코드 |
|---|---|---|---|
| lb-swaiting | prd | healthy | sw.nextpay.co.kr, waiting.nextpay.co.kr |
| waiting-8081 | prd | healthy | backup-waiting.nextpay.co.kr |
| waiting-8082 | prd | healthy | backup-nsm.nextpay.co.kr |
| lb-dev-swaiting | dev | healthy | dev-waiting.nextpay.co.kr |

#### bss (4개) — 통합 후보
| ALB | 환경 | 상태 | 레코드 |
|---|---|---|---|
| bss-order | prd | healthy | bss-order.nextpay.co.kr |
| bss-kiosk | prd | healthy | bss-kiosk.nextpay.co.kr |
| BSSTizenWebSignageSub1 | prd | healthy | bss-signage-sub1.nextpay.co.kr |
| BSSTizenWebSignageSub2 | prd | healthy | bss-signage-sub2.nextpay.co.kr |

#### nserise (3개) — 통합 후보
| ALB | 환경 | 상태 | 레코드 |
|---|---|---|---|
| nextpay-norder-dev | dev | unhealthy | norder-dev.nextpay.co.kr |
| nextpay-npos-dev | dev | unhealthy | npos-dev.nextpay.co.kr |
| nextpay-nkds-dev | dev | healthy | nkds-dev.nextpay.co.kr |

#### 기타 (4개)
| ALB | 솔루션 | 환경 | 상태 | 레코드 |
|---|---|---|---|---|
| wine-curation-elb | wine | prd | healthy | wine-curation.nextpay.co.kr |
| ws2025-lb | ws2025 | prd | error | ws2025-api.nextpay.co.kr, ws2025-web.nextpay.co.kr |
| aiagent-staging | aiagent | stg | error | staging-api.nextaistore.co.kr, staging-system.nextaistore.co.kr, staging.nextaistore.co.kr |
| dooh-dev-elb | dooh | dev | no_record | - |

---

### 5.3 유지 (7개) — 담당자 협의 필요

> 유지 2건(unhealthy 상태로 정상화 검토 필요) + 미정 5건(조치 미결정) = **정리 후 잔존 ALB 7개**

#### 유지 — 정상화 검토 (2개)
| ALB | 솔루션 | 환경 | 상태 | 레코드 | 확인 포인트 |
|---|---|---|---|---|---|
| lb-scms | signage | prd | unhealthy | nsm.nextpay.co.kr, scms.nextpay.co.kr | 트래픽 살아있음 — Target 정상화 필요 |
| alb-dev-signage | signage | dev | unhealthy | dev-scms.nextpay.co.kr, dev-sw.nextpay.co.kr, signage-glitchtip.nextpay.co.kr | Target 정상화 필요 |

#### 미정 — 조치 결정 필요 (5개)
| ALB | 솔루션 | 환경 | 상태 | 레코드 | 확인 포인트 |
|---|---|---|---|---|---|
| dev-cms-elb | cms | dev | error | dev/aiagent/aiot/ncount/dooh-api/staging/wine-curation-dev.nextpay.co.kr 등 7건 | 다수 도메인 묶음 — 분리/이관 필요 |
| webService | nserise | prd | error | back-office/ca/ncms/nkds(-new)/nkiosk/norder/npos/salesking/shop/socket.nextpay.co.kr 등 12건 | 운영 핵심 도메인 묶음 — 영향도 큼 |
| nextpay-kiosk-dev | nserise | dev | unhealthy | nkiosk-dev.nextpay.co.kr | 합치기 그룹 편입 여부 |
| bss-cms | bss | stg | healthy | bss.nextpay.co.kr | bss 통합 그룹 편입 여부 |
| elb-knowledge-graph | aiagent | prd | error | graph.nextaistore.com, graph.nextaistore.co.kr | aiagent-staging과 통합 여부 |

---

## 6. 진행 절차 (제안)

1. **미정 5건 협의** — 운영/개발팀과 도메인 영향도 확인
2. **합치기 그룹 설계** — 솔루션별 통합 ALB 구성 (signage, bss, nserise, 기타)
3. **레코드 이관** — Route53 레코드를 통합 ALB로 변경, TTL 단축 후 절체
4. **검증 기간** — 1~2주 모니터링
5. **삭제 실행** — 무트래픽 확인된 ALB 순차 삭제
6. **Target Group / 미사용 SG 정리**

---

## 참고

- 데이터 출처: `data/entries.json` (수집 + 분류 결과)
- 수집 도구: 사내 `new-lb` Go 도구 (ALB / Route53 / Target Group 통합 수집)
- k8s ALB Controller가 관리하는 ALB(이름 `k8s-` 시작)는 본 리포트 대상에서 제외
- 관련 리포트: `docs/ec2-orphan-report.md` (EC2 미아/중지 인스턴스, EC2+EBS+ALB 통합 절감)
