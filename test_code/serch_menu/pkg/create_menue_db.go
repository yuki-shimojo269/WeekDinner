package create_menue_db

import (
	"fmt"
	"log"
  "os"
	"net/http"
	"net/url"
  "strings"
	"strconv"

	"github.com/PuerkitoBio/goquery"
  "gorm.io/gorm"
  "gorm.io/driver/sqlite"
)

type Menu_tile_url_format struct {
  gorm.Model
  Menu_id string
  Menu_name  string
  Menu_url string
}

//SQLデータベースの作成と初期化
func Create_db(db_name string) bool{
  if err := os.Mkdir("menu_folder", 0777); err != nil {
       log.Println(err)
   }
   var db_path string = "menu_folder" + "/" + db_name
   db, err := gorm.Open(sqlite.Open(db_path), &gorm.Config{})
   if err != nil {
     panic("failed to connect database")
     return false
   }
   db.AutoMigrate(&Menu_tile_url_format{})
   return true
}

//　データを入力
func Add_data(db_name string, Menu_id string, Menu_name  string, Menu_url string) bool{
  var db_path string = "menu_folder" + "/" + db_name
  db, err := gorm.Open(sqlite.Open(db_path), &gorm.Config{})
  if err != nil {
    panic("failed to connect database")
    return false
  }
  db.Create(&Menu_tile_url_format{Menu_id: Menu_id, Menu_name: Menu_name, Menu_url: Menu_url})
  return true
}

// タイトル、レシピID、URLを取得しデータベースに突っ込む
func get_title_and_url(search_url string, home_url string,db_name string)bool{
	is_error := false
  //input  : url
  //output : url,title
  res, err := http.Get(search_url)
	if err != nil {
		log.Println(err)
		is_error = true
	}
	defer res.Body.Close()

	doc, _ := goquery.NewDocumentFromReader(res.Body)
	doc.Find("div.video-list-info").Each(func(i int, s *goquery.Selection){
    s.Find("a").Each(func(j int, a *goquery.Selection){
      recipi_name := strings.Replace(a.Text(), "\n", "", -1)
      recipe_url, is_url := a.Attr("href")
      if is_url == false{
        log.Println("%q error", "scraping")
				is_error = true
      }

      Menu_id := strings.Split(recipe_url, "/")[2]
      err := Add_data(db_name, Menu_id, recipi_name, home_url + recipe_url[1:])
      if err == false{
        log.Println("%q error", "data add")
				is_error = true
      }
    })
	})
	return is_error
}


func Create_database() {
  home_url := "https://www.kurashiru.com/"
	serch_url := home_url + "search?query=" + url.QueryEscape("おかず")
  var db_name string = "menue.db"
	var page_max int = 2

	// データベースの作成
  err := Create_db(db_name)
  if err == false {
    fmt.Println("false cereate db")
  }

	// ページを更新しながら頑張る
	for i:=1; i<=page_max; i++{
		page := strconv.Itoa(i)
		acsess_url := serch_url + "&page=" + page
		log.Println(acsess_url)
  	is_err := get_title_and_url(acsess_url, home_url, db_name)
		if is_err == true{
			break
		}
	}
}
