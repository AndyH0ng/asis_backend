package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/sashabaranov/go-openai"
)

// OpenAIService OpenAI API 상호작용을 처리합니다
type OpenAIService struct {
	client *openai.Client
}

// NewOpenAIService 새로운 OpenAI 서비스를 생성합니다
func NewOpenAIService(apiKey string) *OpenAIService {
	client := openai.NewClient(apiKey)
	return &OpenAIService{client: client}
}

// GenerateRecipePrompt 레시피 생성을 위한 최적화된 프롬프트를 생성합니다
func GenerateRecipePrompt(ingredientsText string) string {
	return fmt.Sprintf(`당신은 전문 요리사입니다. 주어진 재료를 기반으로 실용적이고 맛있는 레시피를 만들어주세요.

## 현재 보유한 재료:
%s

## 레시피 생성 규칙:
1. 위 재료를 최대한 활용하되, 필요한 경우 기본 양념(소금, 후추, 식용유 등)은 추가 가능합니다.
2. 난이도는 ●와 ○로 5개 표시 (예: ●●●○○는 중간 난이도)
3. 예상 조리 시간을 분 단위로 정확하게 표시
4. 재료는 구체적인 양과 함께 나열
5. 조리 방법은 큰 단계와 작은 단계로 구성:
   - 큰 단계: "재료 준비", "조리하기", "마무리" 등의 주요 과정
   - 작은 단계: 각 큰 단계를 이루는 세부 동작들

## 응답 형식 (반드시 JSON 형식으로):
{
  "recipe_name": "레시피 이름",
  "recipe_description": "요리에 대한 간단한 설명 (2-3문장)",
  "recipe_difficulty": "●●●○○",
  "recipe_estimated_time": "45분",
  "ingredients": [
    "재료1 - 200g",
    "재료2 - 1개",
    "재료3 - 2큰술"
  ],
  "cooking_steps": [
    {
      "title": "재료 준비",
      "sub_steps": [
        "채소를 깨끗이 씻어 물기를 제거합니다",
        "고기는 먹기 좋은 크기로 자릅니다",
        "양념 재료를 계량합니다"
      ]
    },
    {
      "title": "조리하기",
      "sub_steps": [
        "팬에 식용유를 두르고 중불로 가열합니다",
        "고기를 넣고 겉면이 익을 때까지 볶습니다",
        "채소를 넣고 함께 볶습니다"
      ]
    },
    {
      "title": "마무리",
      "sub_steps": [
        "양념을 넣고 골고루 섞습니다",
        "약불로 줄이고 5분간 더 조리합니다",
        "불을 끄고 그릇에 담아냅니다"
      ]
    }
  ]
}

JSON 형식만 응답해주세요. 다른 설명은 불필요합니다.`, ingredientsText)
}

// AIRecipeResponse OpenAI의 구조화된 응답을 나타냅니다
type AIRecipeResponse struct {
	RecipeName          string        `json:"recipe_name"`
	RecipeDescription   string        `json:"recipe_description"`
	RecipeDifficulty    string        `json:"recipe_difficulty"`
	RecipeEstimatedTime string        `json:"recipe_estimated_time"`
	Ingredients         []string      `json:"ingredients"`
	CookingSteps        []CookingStep `json:"cooking_steps"`
}

// GenerateRecipe OpenAI API를 호출하여 레시피를 생성합니다
func (os *OpenAIService) GenerateRecipe(ctx context.Context, ingredientsText string) (*Recipe, error) {
	prompt := GenerateRecipePrompt(ingredientsText)

	resp, err := os.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: openai.GPT4,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: "당신은 전문 요리사입니다. 항상 JSON 형식으로만 응답하세요.",
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
			Temperature: 0.7,
			MaxTokens:   2000,
		},
	)

	if err != nil {
		return nil, fmt.Errorf("error calling OpenAI API: %v", err)
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("no response from OpenAI")
	}

	content := resp.Choices[0].Message.Content

	// 응답 정리 (마크다운 코드 블록 제거)
	content = strings.TrimSpace(content)
	content = strings.TrimPrefix(content, "```json")
	content = strings.TrimPrefix(content, "```")
	content = strings.TrimSuffix(content, "```")
	content = strings.TrimSpace(content)

	// JSON 응답 파싱
	var aiResponse AIRecipeResponse
	if err := json.Unmarshal([]byte(content), &aiResponse); err != nil {
		return nil, fmt.Errorf("error parsing OpenAI response: %v\nResponse: %s", err, content)
	}

	// AIRecipeResponse를 Recipe 모델로 변환
	recipe := &Recipe{
		RecipeName:            aiResponse.RecipeName,
		RecipeDescription:     aiResponse.RecipeDescription,
		RecipeDifficulty:      aiResponse.RecipeDifficulty,
		RecipeEstimatedTime:   aiResponse.RecipeEstimatedTime,
		RecipeIngredientCount: len(aiResponse.Ingredients),
		RecipeStepCount:       len(aiResponse.CookingSteps),
		RecipeIsMarked:        false,
		RecipeStatus:          "0",
		Ingredients:           make(map[string]string),
		CookingSteps:          make(map[string]string),
	}

	// 재료 배열을 맵 형식으로 변환 (recipe_ingredient_0, recipe_ingredient_1 등)
	for i, ingredient := range aiResponse.Ingredients {
		key := fmt.Sprintf("recipe_ingredient_%d", i)
		recipe.Ingredients[key] = ingredient
	}

	// 조리 단계를 맵 형식으로 변환 (recipe_step_0_title, recipe_step_0_substep_0 등)
	for i, step := range aiResponse.CookingSteps {
		// 큰 단계 제목 저장
		titleKey := fmt.Sprintf("recipe_step_%d_title", i)
		recipe.CookingSteps[titleKey] = step.Title

		// 작은 단계들 저장
		for j, subStep := range step.SubSteps {
			subStepKey := fmt.Sprintf("recipe_step_%d_substep_%d", i, j)
			recipe.CookingSteps[subStepKey] = subStep
		}

		// 작은 단계 개수 저장
		subStepCountKey := fmt.Sprintf("recipe_step_%d_substep_count", i)
		recipe.CookingSteps[subStepCountKey] = fmt.Sprintf("%d", len(step.SubSteps))
	}

	return recipe, nil
}

// ParseIngredientCount ingredient_number 필드에서 숫자를 추출합니다
func ParseIngredientCount(ingredientNumber string) int {
	count, err := strconv.Atoi(ingredientNumber)
	if err != nil {
		return 1 // 파싱 실패시 기본값 1
	}
	return count
}
