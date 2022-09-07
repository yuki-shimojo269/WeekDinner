package ControlDB

import (
	"fmt"
	"log"
	"strconv"
  "path/filepath"
	"math/rand"
	"strings"
	"time"

  "gorm.io/gorm"
  "gorm.io/driver/sqlite"
)
// ==============================================================
//SQLデータベースの作成と初期化
func Create_RecipeNameDB(db_name string) {
   db_path := filepath.Join("data", db_name)
   db, err := gorm.Open(sqlite.Open(db_path), &gorm.Config{})
   if err != nil {
     panic("failed to connect database")
   }
   db.Table("RecipeNameList").AutoMigrate(&RecipeNameList{})
}

//ページを更新しながらデータベースに追加
func Update_RecipeNameDB(home_url string, serch_url string, db_name string){
  db_path := filepath.Join("data", db_name)
  db, err := gorm.Open(sqlite.Open(db_path), &gorm.Config{})
  if err != nil {
    panic("failed to connect database")
  }

  // ページを更新しながら頑張る
	for page:=1; page<=2; page++{
		page := strconv.Itoa(page)
		acsess_url := serch_url + "&page=" + page
		RecipeNameList_list := GetIdTitleUrl_list(acsess_url, home_url, db_name)
    for _, RecipeNameList_data := range RecipeNameList_list{
			recipe_id := RecipeNameList_data.Recipe_id
			var exists bool
			err := db.Model(&RecipeNameList{}).
				 			  Table("RecipeNameList").
				 				Select("count(*) > 0").
				 				Where("Recipe_id = ?", recipe_id).
				 				Find(&exists).
         				Error
			if err != nil{
				log.Println(err)
			}
			if exists == false{
      	db.Table("RecipeNameList").Create(&RecipeNameList_data)
			}
		}
	}
}
// ==============================================================

// --------------------------------------------------------------
//ランダムでレシピを取り出す
func Choice_RandomMenue(choice_num int)  {
	db_path := filepath.Join("data", "RecipeName.db")
  db, err := gorm.Open(sqlite.Open(db_path), &gorm.Config{})
  if err != nil {
    panic("failed to connect database")
  }

	//データベースの総数を計算
	var recode_count int64
	db.Model(&RecipeNameList{}).
		 Table("RecipeNameList").
		 Find(&RecipeNameList{}).
		 Count(&recode_count)

	// 取り出すレシピのIDをランダムで選択
	var random_ID_list []int64
	for i := 0; i < choice_num; i++{
		random_ID := int64(rand.Intn(int(recode_count)))
		random_ID_list = append(random_ID_list, random_ID)
	}

	// レシピIDから食材の情報を取ってくる
	for _, recipe_id := range random_ID_list{
		var is_exist bool
		db.Table("RecipeNameList").
			 Select("Have_recipe").
			 Where("Data_id = ?", recipe_id).
			 Find(&is_exist)
		if is_exist{
			 continue
		}else{
			var recipe_url string
			db.Table("RecipeNameList").
				 Select("Recipe_url").
				 Where("Data_id = ?", recipe_id).
				 Find(&recipe_url)

			fmt.Println(recipe_url)
			Buff_FoodData := GetFoodData_list(recipe_url)
			Create_RecipeDB(recipe_id, Buff_FoodData)
			// テーブルを更新
			db.Table("RecipeNameList").
				 Select("Have_recipe").
				 Where("Data_id = ?", recipe_id).
				 Updates(RecipeNameList{Have_recipe: true})
			}
		}

	// 一週間分のメニューを作成
	Create_week_menu(random_ID_list)
}

// 取ってきたデータをデータベースに突っ込む
func Create_RecipeDB(recipe_id int64, food_data [][]string)  {
	db_path := filepath.Join("data", "RecipeName.db")
  db, err := gorm.Open(sqlite.Open(db_path), &gorm.Config{})
	if err != nil{
		fmt.Println("error")
	}
	recipe_id_converted := strconv.Itoa(int(recipe_id))
	db.Table("recipe_"+recipe_id_converted).AutoMigrate(&RecipeData{})

	for _, data := range food_data{
		data[0] = strings.Replace(data[0], "(A)", "", 1)
		data[0] = strings.Replace(data[0], "(B)", "", 1)
		data[0] = strings.Replace(data[0], "(C)", "", 1)
		data := RecipeData{Food_name:data[0], Food_amount_type:data[1]}
		db.Table("recipe_"+recipe_id_converted).Create(&data)
	}
}


func Create_week_menu(recipe_id_list []int64){
	// 日付の初期化----------------------
	const DateFormat = "2006/01/02"
	now_time := time.Now().Format(DateFormat)
	now_time = strings.Replace(now_time, "/", "", -1)
	// ----------------------------------

	// 一週間分のデータベースの作成=========
	db_path := filepath.Join("data", now_time+".db")
  db, err := gorm.Open(sqlite.Open(db_path), &gorm.Config{})
  if err != nil {
    panic("failed to connect database")
  }
	db.Table("RecipeList").AutoMigrate(&WeekRecipeList{})
	db.Table("WeekFood").AutoMigrate(&WeekFood{})
	// ===================================

	// データベースにデータを追加 ----------
	db_path = filepath.Join("data", "RecipeName.db")
  home_db, err := gorm.Open(sqlite.Open(db_path), &gorm.Config{})
	for _, recipe_id := range recipe_id_list{
		var recipe_data RecipeNameList
		home_db.Table("RecipeNameList").
			 			Select("*").
			 			Where("Data_id = ?", recipe_id).
			 			Find(&recipe_data)

		add_data := WeekRecipeList{Data_id:     recipe_id,
															 Recipe_id:   recipe_data.Recipe_id,
															 Recipe_name: recipe_data.Recipe_name,
															 Recipe_url:  recipe_data.Recipe_url}
		db.Table("RecipeList").Create(&add_data)
	}
	// -----------------------------------

	// ===================================
	// 一週間分の食材の合計
	for _, recipe_id := range recipe_id_list{
		var FoodData []RecipeData
		recipe_id_converted := strconv.Itoa(int(recipe_id))
		home_db.Table("recipe_"+recipe_id_converted).
						Select("*").
						Find(&FoodData)
		for _, data := range FoodData{
			db.Table("WeekFood").Create(WeekFood{
					Food_name: data.Food_name,
					Food_amount_type: data.Food_amount_type})
		}
	}

	// ===================================
