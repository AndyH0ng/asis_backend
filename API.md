# Asis 백엔드 API 문서

**Base URL**: `https://asisbackend-production.up.railway.app`

---

## 엔드포인트 목록

### 1. 헬스 체크
서버 상태를 확인합니다.

**URL**: `/health`
**Method**: `GET`
**인증**: 불필요

#### 응답 예시
```json
{
  "status": "healthy"
}
```

#### 사용 예시
```bash
curl https://asisbackend-production.up.railway.app/health
```

```javascript
// JavaScript/Flextudio
fetch('https://asisbackend-production.up.railway.app/health')
  .then(response => response.json())
  .then(data => console.log(data));
```

---

### 2. 재료 목록 조회
Firestore에 저장된 모든 재료를 조회합니다.

**URL**: `/api/ingredients`
**Method**: `GET`
**인증**: 불필요

#### 응답 예시
```json
{
  "success": true,
  "count": 2,
  "ingredients": [
    {
      "_SUID": "nJ2Zh4xP",
      "createDateTime": "2025-11-19T11:47:44.884Z",
      "ingredient_group": "발효 식품",
      "ingredient_name": "김치",
      "ingredient_number": "1",
      "updateDateTime": "2025-11-19T11:47:44.884Z"
    },
    {
      "_SUID": "aB3Cd5eF",
      "createDateTime": "2025-11-19T12:30:00.000Z",
      "ingredient_group": "채소",
      "ingredient_name": "양파",
      "ingredient_number": "2",
      "updateDateTime": "2025-11-19T12:30:00.000Z"
    }
  ]
}
```

#### 사용 예시
```bash
curl https://asisbackend-production.up.railway.app/api/ingredients
```

```javascript
// JavaScript/Flextudio
fetch('https://asisbackend-production.up.railway.app/api/ingredients')
  .then(response => response.json())
  .then(data => {
    console.log(`총 ${data.count}개의 재료`);
    console.log(data.ingredients);
  });
```

```dart
// Flutter
import 'package:http/http.dart' as http;
import 'dart:convert';

Future<void> getIngredients() async {
  final response = await http.get(
    Uri.parse('https://asisbackend-production.up.railway.app/api/ingredients'),
  );

  if (response.statusCode == 200) {
    final data = json.decode(response.body);
    print('총 ${data['count']}개의 재료');
    print(data['ingredients']);
  }
}
```

---

### 3. 레시피 생성
보유한 재료를 기반으로 AI가 레시피를 생성하고 Firestore에 저장합니다.

**URL**: `/api/generate-recipe`
**Method**: `POST`
**인증**: 불필요
**Content-Type**: `application/json`

#### 요청
Body 없음 (Firestore에서 자동으로 재료 목록을 가져옴)

#### 응답 예시 (성공)
```json
{
  "success": true,
  "message": "Recipe successfully generated and saved with ID: abc123xyz",
  "recipe": {
    "recipe_name": "김치 볶음밥",
    "recipe_description": "간단하고 맛있는 김치 볶음밥 레시피입니다.",
    "recipe_difficulty": "●●○○○",
    "recipe_estimated_time": "15분",
    "recipe_ingredient_0": "김치 200g",
    "recipe_ingredient_1": "밥 1공기",
    "recipe_ingredient_2": "양파 1/2개",
    "recipe_ingredient_3": "참기름 1큰술",
    "recipe_ingredient_count": 4,
    "recipe_is_marked": false,
    "recipe_status": "0"
  }
}
```

#### 응답 예시 (실패)
```json
{
  "success": false,
  "message": "No ingredients found in Firestore"
}
```

#### 사용 예시
```bash
curl -X POST https://asisbackend-production.up.railway.app/api/generate-recipe
```

```javascript
// JavaScript/Flextudio
fetch('https://asisbackend-production.up.railway.app/api/generate-recipe', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json'
  }
})
  .then(response => response.json())
  .then(data => {
    if (data.success) {
      console.log('레시피 생성 완료:', data.recipe.recipe_name);
      console.log('소요 시간:', data.recipe.recipe_estimated_time);
    } else {
      console.error('에러:', data.message);
    }
  });
```

```dart
// Flutter
import 'package:http/http.dart' as http;
import 'dart:convert';

Future<void> generateRecipe() async {
  final response = await http.post(
    Uri.parse('https://asisbackend-production.up.railway.app/api/generate-recipe'),
    headers: {'Content-Type': 'application/json'},
  );

  if (response.statusCode == 200) {
    final data = json.decode(response.body);
    if (data['success']) {
      print('레시피: ${data['recipe']['recipe_name']}');
      print('난이도: ${data['recipe']['recipe_difficulty']}');
    }
  }
}
```

---

## 에러 코드

| 상태 코드 | 설명 |
|---------|------|
| 200 | 요청 성공 |
| 400 | 잘못된 요청 (예: 재료 없음) |
| 405 | 허용되지 않은 HTTP 메서드 |
| 500 | 서버 내부 오류 |

---

## CORS 정책

모든 오리진(`*`)에서 접근 가능하도록 CORS가 설정되어 있습니다.

- **Allow-Origin**: `*`
- **Allow-Methods**: `GET, POST, PUT, DELETE, OPTIONS`
- **Allow-Headers**: `Content-Type, Authorization`

---

## 주의사항

1. **레시피 생성 API**는 OpenAI GPT-4를 사용하므로 응답 시간이 5-15초 정도 소요될 수 있습니다.
2. Firestore에 재료가 없으면 레시피 생성이 실패합니다.
3. 생성된 레시피는 자동으로 Firestore의 `Recipes` 컬렉션에 저장됩니다.

---

## Flextudio 연결 예시

Flextudio에서 API를 호출할 때는 다음 정보를 사용하세요:

### 재료 목록 조회
- **URL**: `https://asisbackend-production.up.railway.app/api/ingredients`
- **Method**: GET
- **Response Path**: `ingredients` (배열)

### 레시피 생성
- **URL**: `https://asisbackend-production.up.railway.app/api/generate-recipe`
- **Method**: POST
- **Response Path**: `recipe` (객체)
