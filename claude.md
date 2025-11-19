## 명세서
GPT API에게 프롬프트와 firestore에 저장된 값을 보낸 뒤에 값을 받는 서버를 만들고 싶어.
프롬프트에 현재 가지고 있는 재료 정보를 함께 보낼거야.
그리고 그걸 바탕으로 어떤 음식을 만들 수 있는지 레시피를 만들어줘야 해.
레시피 정보는 firestore에 저장할거야.
재료 정보는 firestore에 이렇게 저장되어 있어.
Ingredients/{문서ID}/필드
필드는 다음처럼 구성되어 있어.

_SUID: "nJ2Zh4xP"
createDateTime: "2025-11-19T11:47:44.884Z"
ingredient_group: "발효 식품"
ingredient_name: "김치"
ingredient_number: "1"
updateDateTime: "2025-11-19T11:47:44.884Z"


레시피 정보는 다음처럼 저장할거야.

recipe_description: "설명"
recipe_difficulty: "●●●●○"
recipe_estimated_time: "67분"
recipe_ingredient_0: "재료0"
recipe_ingredient_1: "재료1"
recipe_ingredient_2: "재료2"
recipe_ingredient_count: 35
recipe_is_marked: true
recipe_name: "인생을 날로 먹기"
recipe_status: "0"

어떻게 하면 효율적으로 레시피를 뽑아줄지 프롬프트를 만들어줘. 그리고 firestore에 맞춰서 데이터를 불러오거나 저장하도록 코드를 짜봐. 