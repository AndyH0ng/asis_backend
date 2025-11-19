package main

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

// FirestoreService 모든 Firestore 작업을 처리합니다
type FirestoreService struct {
	client *firestore.Client
}

// NewFirestoreService 새로운 Firestore 서비스를 생성합니다
// credentialsPath: 로컬 개발시 파일 경로 또는 credentialsJSON: Railway 배포시 JSON 문자열
func NewFirestoreService(ctx context.Context, credentialsPath string, credentialsJSON string) (*FirestoreService, error) {
	var opt option.ClientOption

	// credentialsJSON이 제공되면 우선 사용 (Railway 등 클라우드 배포시)
	if credentialsJSON != "" {
		opt = option.WithCredentialsJSON([]byte(credentialsJSON))
	} else if credentialsPath != "" {
		// 파일 경로 사용 (로컬 개발시)
		opt = option.WithCredentialsFile(credentialsPath)
	} else {
		return nil, fmt.Errorf("firebase credentials not provided")
	}

	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		return nil, fmt.Errorf("error initializing firebase app: %v", err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		return nil, fmt.Errorf("error initializing firestore client: %v", err)
	}

	return &FirestoreService{client: client}, nil
}

// Close Firestore 클라이언트를 종료합니다
func (fs *FirestoreService) Close() error {
	return fs.client.Close()
}

// GetAllIngredients Firestore에서 모든 재료를 조회합니다
func (fs *FirestoreService) GetAllIngredients(ctx context.Context) ([]Ingredient, error) {
	var ingredients []Ingredient

	iter := fs.client.Collection("Ingredients").Documents(ctx)
	defer iter.Stop()

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error iterating ingredients: %v", err)
		}

		var ingredient Ingredient
		if err := doc.DataTo(&ingredient); err != nil {
			return nil, fmt.Errorf("error converting document to ingredient: %v", err)
		}

		ingredients = append(ingredients, ingredient)
	}

	return ingredients, nil
}

// SaveRecipe 레시피를 Firestore에 저장합니다
func (fs *FirestoreService) SaveRecipe(ctx context.Context, recipe *Recipe) (string, error) {
	// 타임스탬프 설정
	now := time.Now().Format(time.RFC3339)
	recipe.CreateDateTime = now
	recipe.UpdateDateTime = now

	// Recipe 구조체를 Firestore용 맵으로 변환
	recipeData := make(map[string]interface{})
	recipeData["recipe_description"] = recipe.RecipeDescription
	recipeData["recipe_difficulty"] = recipe.RecipeDifficulty
	recipeData["recipe_estimated_time"] = recipe.RecipeEstimatedTime
	recipeData["recipe_ingredient_count"] = recipe.RecipeIngredientCount
	recipeData["recipe_step_count"] = recipe.RecipeStepCount
	recipeData["recipe_is_marked"] = recipe.RecipeIsMarked
	recipeData["recipe_name"] = recipe.RecipeName
	recipeData["recipe_status"] = recipe.RecipeStatus
	recipeData["createDateTime"] = recipe.CreateDateTime
	recipeData["updateDateTime"] = recipe.UpdateDateTime

	// 개별 재료 필드 추가
	for key, value := range recipe.Ingredients {
		recipeData[key] = value
	}

	// 조리 단계 필드 추가
	for key, value := range recipe.CookingSteps {
		recipeData[key] = value
	}

	// Firestore에 저장
	docRef, _, err := fs.client.Collection("Recipes").Add(ctx, recipeData)
	if err != nil {
		return "", fmt.Errorf("error saving recipe to firestore: %v", err)
	}

	return docRef.ID, nil
}

// FormatIngredientsForPrompt 재료 정보를 AI 프롬프트용 문자열로 포맷팅합니다
func FormatIngredientsForPrompt(ingredients []Ingredient) string {
	if len(ingredients) == 0 {
		return "재료 없음"
	}

	result := ""
	for i, ing := range ingredients {
		result += fmt.Sprintf("%d. %s (%s) - 수량: %s개\n",
			i+1,
			ing.IngredientName,
			ing.IngredientGroup,
			ing.IngredientNumber)
	}

	return result
}
