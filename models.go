package main

// Ingredient Firestore에 저장된 재료 정보를 나타냅니다
type Ingredient struct {
	SUID             string `firestore:"_SUID" json:"_SUID"`
	CreateDateTime   string `firestore:"createDateTime" json:"createDateTime"`
	IngredientGroup  string `firestore:"ingredient_group" json:"ingredient_group"`
	IngredientName   string `firestore:"ingredient_name" json:"ingredient_name"`
	IngredientNumber string `firestore:"ingredient_number" json:"ingredient_number"`
	UpdateDateTime   string `firestore:"updateDateTime" json:"updateDateTime"`
}

// Recipe Firestore에 저장될 레시피 정보를 나타냅니다
type Recipe struct {
	RecipeDescription     string            `firestore:"recipe_description" json:"recipe_description"`
	RecipeDifficulty      string            `firestore:"recipe_difficulty" json:"recipe_difficulty"`
	RecipeEstimatedTime   string            `firestore:"recipe_estimated_time" json:"recipe_estimated_time"`
	Ingredients           map[string]string `firestore:",flatten" json:"ingredients"`    // recipe_ingredient_0, recipe_ingredient_1 등으로 저장됩니다
	RecipeIngredientCount int               `firestore:"recipe_ingredient_count" json:"recipe_ingredient_count"`
	CookingSteps          map[string]string `firestore:",flatten" json:"cooking_steps"`  // recipe_step_0_title, recipe_step_0_substep_0 등으로 저장됩니다
	RecipeStepCount       int               `firestore:"recipe_step_count" json:"recipe_step_count"`
	RecipeIsMarked        bool              `firestore:"recipe_is_marked" json:"recipe_is_marked"`
	RecipeName            string            `firestore:"recipe_name" json:"recipe_name"`
	RecipeStatus          string            `firestore:"recipe_status" json:"recipe_status"`
	CreateDateTime        string            `firestore:"createDateTime" json:"createDateTime"`
	UpdateDateTime        string            `firestore:"updateDateTime" json:"updateDateTime"`
}

// CookingStep 큰 단계를 나타냅니다
type CookingStep struct {
	Title    string   `json:"title"`     // 큰 단계 제목
	SubSteps []string `json:"sub_steps"` // 작은 단계들
}

// RecipeRequest 레시피 생성 API 요청을 나타냅니다
type RecipeRequest struct {
	UserID string `json:"user_id"` // 선택사항: 사용자별 재료 필터링에 사용
}

// RecipeResponse API 응답을 나타냅니다
type RecipeResponse struct {
	Success bool    `json:"success"`
	Message string  `json:"message"`
	Recipe  *Recipe `json:"recipe,omitempty"`
}
