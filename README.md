# asis_backend

Firestore에 저장된 재료 정보를 기반으로 OpenAI GPT를 사용하여 레시피를 자동 생성하는 Go 서버입니다.

## 기능

- Firestore에서 재료 정보 조회
- OpenAI GPT-4를 사용한 레시피 생성
- 생성된 레시피를 Firestore에 자동 저장
- RESTful API 제공

## 설치 및 설정

### 1. 필요한 패키지 설치

```bash
go mod download
```

### 2. 환경 변수 설정

`.env.example` 파일을 복사하여 `.env` 파일을 생성하고 필요한 값을 입력합니다.

```bash
cp .env.example .env
```

필요한 환경 변수:
- `FIREBASE_CREDENTIALS_PATH`: Firebase 서비스 계정 키 JSON 파일 경로
- `OPENAI_API_KEY`: OpenAI API 키
- `PORT`: 서버 포트 (선택사항, 기본값: 8080)

### 3. Firebase 설정

1. Firebase Console에서 프로젝트 생성
2. 서비스 계정 키 JSON 파일 다운로드
3. JSON 파일을 프로젝트 루트에 저장하고 `.env` 파일의 경로 업데이트

## 실행

```bash
go run .
```

서버가 시작되면 기본적으로 `http://localhost:8080`에서 실행됩니다.

## API 엔드포인트

### 1. Health Check
```
GET /health
```

서버 상태를 확인합니다.

**응답 예시:**
```json
{
  "status": "healthy"
}
```

### 2. 재료 조회
```
GET /api/ingredients
```

Firestore에 저장된 모든 재료를 조회합니다.

**응답 예시:**
```json
{
  "success": true,
  "count": 3,
  "ingredients": [
    {
      "_SUID": "nJ2Zh4xP",
      "createDateTime": "2025-11-19T11:47:44.884Z",
      "ingredient_group": "발효 식품",
      "ingredient_name": "김치",
      "ingredient_number": "1",
      "updateDateTime": "2025-11-19T11:47:44.884Z"
    }
  ]
}
```

### 3. 레시피 생성
```
POST /api/generate-recipe
```

Firestore의 재료를 기반으로 새로운 레시피를 생성하고 저장합니다.

**요청:**
```bash
curl -X POST http://localhost:8080/api/generate-recipe
```

**응답 예시:**
```json
{
  "success": true,
  "message": "Recipe successfully generated and saved with ID: abc123",
  "recipe": {
    "recipe_name": "김치찌개",
    "recipe_description": "한국의 대표적인 찌개 요리로, 김치의 깊은 맛이 일품입니다...",
    "recipe_difficulty": "●●○○○",
    "recipe_estimated_time": "30분",
    "ingredients": {
      "recipe_ingredient_0": "김치 - 200g",
      "recipe_ingredient_1": "돼지고기 - 150g",
      "recipe_ingredient_2": "두부 - 1/2모"
    },
    "recipe_ingredient_count": 3,
    "recipe_is_marked": false,
    "recipe_status": "0",
    "createDateTime": "2025-11-19T12:00:00Z",
    "updateDateTime": "2025-11-19T12:00:00Z"
  }
}
```

## Firestore 데이터 구조

### Ingredients 컬렉션
```
{
  "_SUID": "string",
  "createDateTime": "ISO 8601 string",
  "ingredient_group": "string",
  "ingredient_name": "string",
  "ingredient_number": "string",
  "updateDateTime": "ISO 8601 string"
}
```

### Recipes 컬렉션
```
{
  "recipe_name": "string",
  "recipe_description": "string",
  "recipe_difficulty": "string (●●●○○ 형식)",
  "recipe_estimated_time": "string (분 단위)",
  "recipe_ingredient_0": "string",
  "recipe_ingredient_1": "string",
  ...
  "recipe_ingredient_count": "number",
  "recipe_is_marked": "boolean",
  "recipe_status": "string",
  "createDateTime": "ISO 8601 string",
  "updateDateTime": "ISO 8601 string"
}
```

## 프롬프트 설계

OpenAI에게 전달되는 프롬프트는 다음과 같이 최적화되어 있습니다:

1. **구조화된 입력**: 재료를 그룹별로 정리하여 제공
2. **명확한 지시사항**: 난이도, 조리 시간, 재료 양 등을 구체적으로 명시
3. **JSON 응답 형식**: 파싱이 쉬운 JSON 형식으로 응답 요청
4. **실용성 강조**: 실제로 만들 수 있는 레시피 생성

자세한 프롬프트는 `openai_service.go`의 `GenerateRecipePrompt` 함수를 참고하세요.

## 개발

### 프로젝트 구조
```
.
├── main.go                  # HTTP 서버 및 핸들러
├── models.go                # 데이터 모델 정의
├── firestore_service.go     # Firestore 연동 로직
├── openai_service.go        # OpenAI API 연동 로직
├── .env.example             # 환경 변수 템플릿
└── README.md                # 문서
```

## 라이센스

MIT
