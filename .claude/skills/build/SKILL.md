---
name: build
description: Go 프로젝트 빌드 및 실행 - go build 후 에러 처리, /build run으로 빌드+실행
---

# Build (빌드 및 실행)

## 설명
Go 프로젝트를 빌드하고 실행한다.

## 트리거
사용자가 `/build` 입력 시 실행

## 동작
1. `go build -o new-lb .`로 빌드한다
2. 빌드 에러가 있으면 에러를 분석하고 Developer 에이전트에게 수정을 지시한다
3. 빌드 성공 시 실행 여부를 사용자에게 확인한다
4. 실행 시 `./new-lb`를 실행하고 웹 서버 주소를 안내한다

## 사용법
```
/build              # 빌드만
/build run          # 빌드 + 실행
```
