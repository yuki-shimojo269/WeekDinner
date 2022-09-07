package main

import (
	"fmt"
	"log"
  "os"
	"net/http"
	//"net/url"
  "strings"
	//"strconv"

	"github.com/PuerkitoBio/goquery"
  "gorm.io/gorm"
  "gorm.io/driver/sqlite"

	"serch_menue/pkg"
)


type resipe_data_form struct{
  Food_name string
  How_Meny string
}

//SQLデータベースの作成と初期化
func Create_db(db_name string){
  if err := os.Mkdir("menu_folder", 0777); err != nil {
       log.Println(err)
   }
   var db_path string = "menu_folder" + "/" + db_name
   db, err := gorm.Open(sqlite.Open(db_path), &gorm.Config{})
   if err != nil {
     panic("failed to connect database")
   }
   db.Table("samples").AutoMigrate(&resipe_data_form{})
}


// 食材、個数を取得しデータベースに突っ込む
func add_data_to_table(search_url string, recipe_id string, db_name string){
  // テーブルの作成と初期化
  recipe_id = "T"+strings.Replace(recipe_id, "-", "_", -1)
  var db_path string = "menu_folder" + "/" + db_name
  db, err := gorm.Open(sqlite.Open(db_path), &gorm.Config{})
  if err != nil {
    panic("failed to connect database")
  }
  db.Table(recipe_id).AutoMigrate(&resipe_data_form{})

  //スクレイピング
  res, err := http.Get(search_url)
	if err != nil {
		log.Println(err)
	}
	defer res.Body.Close()

	doc, _ := goquery.NewDocumentFromReader(res.Body)
  doc.Find(".ingredient-list-item").Each(func(i int, data *goquery.Selection){
    text_data := data.Find("a").Text()
    text_data = strings.Replace(text_data,"\n", "", -1)
    Food_name := strings.Replace(text_data," ", "", -1)
    How_many := data.Find("span").Text()
    db.Table(recipe_id).Create(&resipe_data_form{Food_name: Food_name, How_Meny: How_many})
  })
}

func main() {
  // home_url := "https://www.kurashiru.com/"

	create_menue_db.Create_database()
	db_name := "menue.db"

	/*
	var db_path string = "menu_folder" + "/" + db_name
	db, err := gorm.Open(sqlite.Open(db_path), &gorm.Config{})

	url := db.Select("Menu_url").First(&Menu_tile_url_format, "id=?", 1)
	fmt.Println(url)
	url_list := strings.Split(url, '/')
	recipe_id = url_list[len(url_list) - 1]
	*/

  Create_db(db_name)
	db_name = "recipe_data.db"
  url := "https://www.kurashiru.com/recipes/1ac5fbb3-9437-4869-a4b7-f53dd7912270"
  recipe_id := "1ac5fbb3-9437-4869-a4b7-f53dd7912270"
  fmt.Println(url)
  add_data_to_table(url, recipe_id, db_name)

}
