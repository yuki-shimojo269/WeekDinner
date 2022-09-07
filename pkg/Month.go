package week_menue

import (
	"fmt"
	"log"
	"strconv"
  "path/filepath"
  "net/http"
	"strings"

  "gorm.io/gorm"
  "gorm.io/driver/sqlite"
  "github.com/PuerkitoBio/goquery"
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
// ==============================================================

// ==============================================================


// タイトル、レシピID、URLを取得しデータベースに突っ込む
func GetIdTitleUrl_list(search_url string, home_url string,db_name string) []RecipeNameList{
  //input  : url
  //output : url,title

  res, err := http.Get(search_url)
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()


  RecipeNameList_list := []RecipeNameList{}

	doc, _ := goquery.NewDocumentFromReader(res.Body)
	doc.Find("div.video-list-info").Each(func(i int, s *goquery.Selection){
    s.Find("a").Each(func(j int, a *goquery.Selection){

      // Recipe_name を取得 ===
      recipi_name := strings.Replace(a.Text(), "\n", "", -1)
      recipi_name = strings.Replace(recipi_name, " ", "", -1)
      // ======================

      // URLを取得 ============
      recipe_url, is_url := a.Attr("href")
      if is_url == false{
        fmt.Println("%q error", "scraping")
      }
      // ======================

      // Menu_idを取得 ========
      recipe_id := strings.Split(recipe_url, "/")[2]
      // ======================

      // urlを再構築 -----------
      recipe_url = home_url[:len([]rune(home_url))-1] + recipe_url
      // ----------------------

      // dataをリストに追加
      data := RecipeNameList{Recipe_id: recipe_id,
                             Recipe_name: recipi_name,
                             Recipe_url: recipe_url,
                             Have_recipe: false}
      RecipeNameList_list = append(RecipeNameList_list, data)
    })
	})
	return RecipeNameList_list
}

//ページを更新しながらデータベースに追加
func Update_RecipeNameDB(home_url string, serch_url string, db_name string){
  db_path := filepath.Join("data", db_name)
  db, err := gorm.Open(sqlite.Open(db_path), &gorm.Config{})
  if err != nil {
    panic("failed to connect database")
  }

  // ページを更新しながら頑張る
	for page:=1; page<=5; page++{
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
