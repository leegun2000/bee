# 백엔드 개발자 채용 과제 - 이슈 관리 API

## 과제 개요

- 제출 형식: 제출물을 깃허브에 업로드 후 링크 공유 (public repository)
- 실행 환경: Go 언어, 포트 8080번 사용

## 평가 기준

- REST API 설계 및 구현 능력
- 비즈니스 로직 구현 능력
- 코드 구조화 및 가독성
- 에러 핸들링 및 데이터 검증

## 구현 요구사항

### 1. 프로젝트 구조

- 실행 방법이 포함된 README.md 파일을 반드시 포함해주세요.
- API 테스트 방법을 README.md에 명시해주세요.

### 2. 데이터 모델

API를 구현 시 아래 구조체를 참고해주세요(수정 가능):

```go
type User struct {
    ID   uint   `json:"id"`
    Name string `json:"name"`
}

type Issue struct {
    ID          uint      `json:"id"`
    Title       string    `json:"title"`
    Description string    `json:"description"`
    Status      string    `json:"status"`
    User        *User     `json:"user,omitempty"` 
    CreatedAt   time.Time `json:"createdAt"`
    UpdatedAt   time.Time `json:"updatedAt"`
}
```

### 3. 비즈니스 규칙

#### 3.1 이슈 상태 (Status)

- 유효한 상태값: `PENDING`, `IN_PROGRESS`, `COMPLETED`, `CANCELLED`

#### 3.2 사용자 관리

- 적어도 아래 3명의 사용자는 시스템에 존재해야 합니다:

```json
[
    { "id": 1, "name": "김개발" },
    { "id": 2, "name": "이디자인" },
    { "id": 3, "name": "박기획" }
]
```

### 4. API 명세

#### API 공통

- 필수 파라미터가 누락되었거나, 미리 정의된 값이 아닌 파라미터가 전달되었거나, 유효하지 않은 데이터가 요청된 경우 에러 응답을 반환합니다.

#### 4.1 이슈 생성 [POST] /issue

- 담당자(userId)가 있는 경우: 상태를 `IN_PROGRESS`로 설정
  - 존재하지 않는 사용자를 담당자로 지정할 수 없습니다.
- 담당자(userId)가 없는 경우: 상태를 `PENDING`으로 설정

요청 예시:

```json
{
    "title": "버그 수정 필요", // 필수
    "description": "로그인 페이지에서 오류 발생", // 선택
    "userId": 1 // 선택
}
```

응답 예시 (201 Created):

```json
{
    "id": 1,
    "title": "버그 수정 필요",
    "description": "로그인 페이지에서 오류 발생",
    "status": "IN_PROGRESS",
    "user": { "id": 1, "name": "김개발" },
    "createdAt": "2025-06-02T10:00:00Z",
    "updatedAt": "2025-06-02T10:00:00Z"
}
```

#### 4.2 이슈 목록 조회 [GET] /issues

- 상태별 필터링 지원 (쿼리 파라미터: status)
- status 파라미터가 없는 경우 전체 이슈 조회

응답 예시 (200 OK):

```json
{
    "issues": [
        {
            "id": 1,
            "title": "버그 수정 필요",
            "description": "로그인 페이지에서 오류 발생",
            "status": "PENDING",
            "createdAt": "2025-06-02T10:00:00Z",
            "updatedAt": "2025-06-02T10:05:00Z"
        }
    ]
}
```

#### 4.3 이슈 상세 조회 [GET] /issue/:id

응답 예시 (200 OK):

```json
{
    "id": 1,
    "title": "버그 수정 필요",
    "description": "로그인 페이지에서 오류 발생",
    "status": "PENDING",
    "createdAt": "2025-06-02T10:00:00Z",
    "updatedAt": "2025-06-02T10:05:00Z"
}
```

#### 4.4 이슈 수정 [PATCH] /issue/:id

- 제목, 설명, 상태, 담당자 변경 가능
- 요청 data에 명시되지 않은 필드는 업데이트하지 않음
- 업데이트 시, 담당자를 지정하지 않고 `PENDING`, `CANCELLED` 이외 다른 상태로 변경 불가
- 상태가 `PENDING`일 때 담당자 할당 시 상태를 변경
  - 따로 상태를 지정하지 않는 경우 `IN_PROGRESS`로 변경되어야 함
- 담당자 제거(userId -> null) 시 상태는 `PENDING` 으로 변경
- `COMPLETED` 또는 `CANCELLED` 상태에서는 해당 issue 업데이트 불가

요청 예시:

```json
{
    "title": "로그인 버그 수정",
    "status": "IN_PROGRESS",
    "userId": 2
}
```

응답 예시 (200 OK):

```json
{
    "id": 1,
    "title": "로그인 버그 수정",
    "description": "로그인 페이지에서 오류 발생",
    "status": "IN_PROGRESS",
    "user": { "id": 2, "name": "이디자인" },
    "createdAt": "2025-06-02T10:00:00Z",
    "updatedAt": "2025-06-02T10:10:00Z"
}
```

### 5. 에러 응답

API는 적절한 HTTP 상태 코드, 메시지와 함께 에러 응답을 반환해야 합니다.
에러 응답은 아래 형식을 따릅니다:

```json
{
    "error": "에러 메시지",
    "code": 400
}
```

응답에 포함되는 HTTP 상태 코드는 REST API 설계 원칙에 따라 적절하게 선택해주세요.
