package week_menue
// 変数を定義するときは戦闘が大文字じゃなきゃダメらしい

// 初期値------------------------------------------------
type FuncInitArg struct{
  IsInit  *bool
}
// ------------------------------------------------------

// RecipeName.db ----------------

//レシピ一覧
type RecipeNameList struct {
  Data_id      int64    `gorm:"primaryKey"`
  Recipe_id    string
  Recipe_name  string
  Recipe_url   string
  Have_recipe  bool
}

//レシピの材料(名前はrecipeID)
type RecipeData struct{
  Food_name         string
  Food_amount_type  string
  Food_amount_num   int
}

// 大さじやこさじを判断するためのデータベース
type FoodType struct{
  TypeID    int     `gorm:"primaryKey"`
  Type_name string
  Sample    string
}


// 今週のレシピ
type WeekRecipeList struct{
  Data_id      int64
  Recipe_id    string
  Recipe_name  string
  Recipe_url   string
  Week         string
}
// 今週のレシピ一覧
type WeekFood struct{
  Food_name        string
  Food_amount_type string
  Food_amount_num  int
}

// ===============================================

//Google カレンダーに入れるデータセット
type GoolgeCalendar struct{
  Recipe_id    int64
  Recipe_name  string
  Recipe_url   string
  Event_day    string
  Ingredient   string
}
