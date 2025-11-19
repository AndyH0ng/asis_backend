package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

var (
	firestoreService *FirestoreService
	openaiService    *OpenAIService
)

func main() {
	// 환경 변수 로드
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}

	// 서비스 초기화
	ctx := context.Background()

	// Firestore 초기화
	credentialsPath := os.Getenv("FIREBASE_CREDENTIALS_PATH")
	if credentialsPath == "" {
		log.Fatal("FIREBASE_CREDENTIALS_PATH environment variable is required")
	}

	var err error
	firestoreService, err = NewFirestoreService(ctx, credentialsPath)
	if err != nil {
		log.Fatalf("Failed to initialize Firestore service: %v", err)
	}
	defer firestoreService.Close()

	// OpenAI 초기화
	openaiAPIKey := os.Getenv("OPENAI_API_KEY")
	if openaiAPIKey == "" {
		log.Fatal("OPENAI_API_KEY environment variable is required")
	}
	openaiService = NewOpenAIService(openaiAPIKey)

	// HTTP 라우트 설정
	http.HandleFunc("/health", healthCheckHandler)
	http.HandleFunc("/api/generate-recipe", generateRecipeHandler)
	http.HandleFunc("/api/ingredients", getIngredientsHandler)

	// 서버 시작
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s...", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// healthCheckHandler 헬스 체크 요청을 처리합니다
func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "healthy",
	})
}

// getIngredientsHandler Firestore에서 모든 재료를 조회합니다
func getIngredientsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()
	ingredients, err := firestoreService.GetAllIngredients(ctx)
	if err != nil {
		log.Printf("Error getting ingredients: %v", err)
		http.Error(w, "Failed to get ingredients", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":     true,
		"ingredients": ingredients,
		"count":       len(ingredients),
	})
}

// generateRecipeHandler 보유한 재료를 기반으로 레시피를 생성합니다
func generateRecipeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()

	// 1단계: Firestore에서 모든 재료 가져오기
	log.Println("Fetching ingredients from Firestore...")
	ingredients, err := firestoreService.GetAllIngredients(ctx)
	if err != nil {
		log.Printf("Error getting ingredients: %v", err)
		respondWithError(w, "Failed to get ingredients", http.StatusInternalServerError)
		return
	}

	if len(ingredients) == 0 {
		respondWithError(w, "No ingredients found in Firestore", http.StatusBadRequest)
		return
	}

	log.Printf("Found %d ingredients", len(ingredients))

	// 2단계: 프롬프트용 재료 포맷팅
	ingredientsText := FormatIngredientsForPrompt(ingredients)
	log.Printf("Ingredients formatted for prompt:\n%s", ingredientsText)

	// 3단계: OpenAI를 사용하여 레시피 생성
	log.Println("Generating recipe with OpenAI...")
	recipe, err := openaiService.GenerateRecipe(ctx, ingredientsText)
	if err != nil {
		log.Printf("Error generating recipe: %v", err)
		respondWithError(w, fmt.Sprintf("Failed to generate recipe: %v", err), http.StatusInternalServerError)
		return
	}

	log.Printf("Recipe generated: %s", recipe.RecipeName)

	// 4단계: 레시피를 Firestore에 저장
	log.Println("Saving recipe to Firestore...")
	recipeID, err := firestoreService.SaveRecipe(ctx, recipe)
	if err != nil {
		log.Printf("Error saving recipe: %v", err)
		respondWithError(w, "Failed to save recipe", http.StatusInternalServerError)
		return
	}

	log.Printf("Recipe saved with ID: %s", recipeID)

	// 5단계: 성공 응답 반환
	response := RecipeResponse{
		Success: true,
		Message: fmt.Sprintf("Recipe successfully generated and saved with ID: %s", recipeID),
		Recipe:  recipe,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// respondWithError 에러 응답을 전송합니다
func respondWithError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(RecipeResponse{
		Success: false,
		Message: message,
	})
}
