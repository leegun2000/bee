# 이슈 관리 API

해당 프로젝트는 git copilot Claude Sonnet4 를 이용하여 만든 vibe coding 입니다.

Go로 구현된 이슈 관리 API 서버입니다. 인터페이스 기반 아키텍처를 사용하여 외부 프레임워크에 직접 의존하지 않는 클린 아키텍처를 구현했습니다.

DB를 사용하지 않고 Memory로만 관리하도록 구성했습니다.

## 주요 특징

- **인터페이스 기반 설계**: 외부 모듈(Gin)에 직접 의존하지 않고, 인터페이스를 통한 추상화 구현
- **TDD 개발**: 테스트 주도 개발로 모든 기능 검증
- **Go 표준 레이아웃**: `internal`과 `pkg` 구조로 프로젝트 조직화
- **웹 프레임워크 독립적**: 서버 추상화를 통해 다양한 웹 프레임워크 지원 가능
- **그레이스풀 셧다운**: 안전한 서버 종료 지원
- 상태별 이슈 필터링
- **Graceful Shutdown**: SIGINT, SIGTERM 신호를 받으면 안전하게 서버 종료

## 실행 방법

### 1. 사전 요구사항

- Go 1.24.0 이상

### 2. 실행

```bash
# 저장소 클론
git clone <repository-url>
cd aoroa

# 의존성 설치
go mod tidy

# 서버 실행
go run .
```

서버는 포트 8080에서 실행됩니다.

### 3. 헬스 체크

```bash
curl http://localhost:8080/health
```

### 4. Graceful Shutdown 테스트

서버 종료 테스트:
```bash
# 서버 실행
go run .

# 다른 터미널에서 SIGTERM 신호 보내기
kill -TERM <PID>

# 또는 Ctrl+C로 종료
# 서버가 안전하게 종료되는 것을 확인할 수 있습니다
```

## API 테스트 방법

### 1. 이슈 생성 (POST /issue)

담당자 없는 이슈 생성:
```bash
curl -X POST http://localhost:8080/issue \
  -H "Content-Type: application/json" \
  -d '{
    "title": "버그 수정 필요",
    "description": "로그인 페이지에서 오류 발생"
  }'
```

담당자 있는 이슈 생성:
```bash
curl -X POST http://localhost:8080/issue \
  -H "Content-Type: application/json" \
  -d '{
    "title": "UI 개선",
    "description": "메인 페이지 디자인 개선",
    "userId": 1
  }'
```

### 2. 이슈 목록 조회 (GET /issues)

전체 이슈 조회:
```bash
curl http://localhost:8080/issues
```

상태별 필터링:
```bash
curl "http://localhost:8080/issues?status=PENDING"
curl "http://localhost:8080/issues?status=IN_PROGRESS"
curl "http://localhost:8080/issues?status=COMPLETED"
curl "http://localhost:8080/issues?status=CANCELLED"
```

### 3. 이슈 상세 조회 (GET /issue/:id)

```bash
curl http://localhost:8080/issue/1
```

### 4. 이슈 수정 (PATCH /issue/:id)

제목과 설명 수정:
```bash
curl -X PATCH http://localhost:8080/issue/1 \
  -H "Content-Type: application/json" \
  -d '{
    "title": "로그인 버그 수정",
    "description": "로그인 페이지 오류 해결"
  }'
```

담당자 할당:
```bash
curl -X PATCH http://localhost:8080/issue/1 \
  -H "Content-Type: application/json" \
  -d '{
    "userId": 2
  }'
```

상태 변경:
```bash
curl -X PATCH http://localhost:8080/issue/1 \
  -H "Content-Type: application/json" \
  -d '{
    "status": "COMPLETED"
  }'
```

담당자 제거:
```bash
curl -X PATCH http://localhost:8080/issue/1 \
  -H "Content-Type: application/json" \
  -d '{
    "userId": null
  }'
```

## 데이터 모델

### User
```json
{
  "id": 1,
  "name": "김개발"
}
```

### Issue
```json
{
  "id": 1,
  "title": "버그 수정 필요",
  "description": "로그인 페이지에서 오류 발생",
  "status": "PENDING",
  "user": {
    "id": 1,
    "name": "김개발"
  },
  "createdAt": "2025-07-11T10:00:00Z",
  "updatedAt": "2025-07-11T10:00:00Z"
}
```

## 비즈니스 규칙

### 이슈 상태
- `PENDING`: 대기 중
- `IN_PROGRESS`: 진행 중
- `COMPLETED`: 완료
- `CANCELLED`: 취소

### 상태 변경 규칙
- 담당자가 할당되면 `PENDING` → `IN_PROGRESS`
- 담당자가 제거되면 상태는 `PENDING`으로 변경
- `COMPLETED` 또는 `CANCELLED` 상태의 이슈는 수정 불가
- 담당자 없이는 `PENDING`, `CANCELLED` 외의 상태로 변경 불가

### 기본 사용자
시스템에는 기본적으로 다음 사용자들이 존재합니다:
- ID 1: 김개발
- ID 2: 이디자인  
- ID 3: 박기획

## 에러 응답

API는 다음 형식으로 에러를 반환합니다:

```json
{
  "error": "에러 메시지",
  "code": 400
}
```

주요 HTTP 상태 코드:
- `400 Bad Request`: 잘못된 요청 데이터
- `404 Not Found`: 리소스를 찾을 수 없음
- `409 Conflict`: 비즈니스 규칙 위반
- `201 Created`: 리소스 생성 성공
- `200 OK`: 요청 처리 성공
`````markdown
## 프로젝트 구조

```
aoroa/
├── main.go                     # 애플리케이션 진입점
├── internal/                   # 비즈니스 로직 (외부 접근 불가)
│   ├── domain/                 # 도메인 타입 및 상수
│   │   └── types.go
│   ├── models/                 # 데이터 모델
│   │   └── models.go
│   ├── service/                # 비즈니스 로직
│   │   ├── user_service.go
│   │   ├── issue_service.go
│   │   └── issue_service_test.go
│   ├── handler/                # HTTP 핸들러 (인터페이스 기반)
│   │   ├── issue_handler.go    # 핵심 핸들러 인터페이스
│   │   ├── gin_wrapper.go      # Gin 프레임워크 어댑터
│   │   ├── http_handler.go     # 표준 HTTP 핸들러
│   │   └── handler_test.go     # 핸들러 테스트
│   └── server/                 # 서버 초기화
│       └── server.go
├── pkg/                        # 외부에서 사용 가능한 패키지
│   ├── server/                 # 서버 추상화
│   │   ├── interfaces.go       # 서버 인터페이스 정의
│   │   ├── abstract_server.go  # 프레임워크 독립적 서버
│   │   └── gin_adapter.go      # Gin 프레임워크 어댑터
│   └── utils/                  # 유틸리티
│       ├── http.go             # HTTP 컨텍스트 인터페이스
│       ├── gin_adapter.go      # Gin 컨텍스트 어댑터
│       └── http_adapter.go     # 표준 HTTP 어댑터
└── go.mod
```

## 아키텍처

### 인터페이스 기반 설계

1. **HTTPContext 인터페이스**: HTTP 요청/응답을 추상화
   - Gin, 표준 HTTP 등 다양한 프레임워크 지원
   - 테스트 가능한 Mock 구현

2. **WebFramework 인터페이스**: 웹 프레임워크 추상화
   - 라우팅 및 서버 시작 기능 추상화
   - 프레임워크 교체 시 최소한의 변경

3. **IssueHandlerInterface**: 비즈니스 로직 핸들러 인터페이스
   - 프레임워크에 독립적인 핸들러 구현
   - 테스트 및 Mock 생성 용이

### 계층 구조

```
┌─────────────────┐
│   main.go       │ # 애플리케이션 시작점
└─────────────────┘
         │
┌─────────────────┐
│ Server Layer    │ # 서버 추상화 및 프레임워크 어댑터
└─────────────────┘
         │
┌─────────────────┐
│ Handler Layer   │ # HTTP 핸들러 (인터페이스 기반)
└─────────────────┘
         │
┌─────────────────┐
│ Service Layer   │ # 비즈니스 로직
└─────────────────┘
         │
┌─────────────────┐
│ Model Layer     │ # 데이터 모델
└─────────────────┘
```

## 기능
