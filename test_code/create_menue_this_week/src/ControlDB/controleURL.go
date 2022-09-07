package ControlDB

import (
  "fmt"
  "net/http"
  "strings"

  "github.com/PuerkitoBio/goquery"
)


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

// レシピURLから材料のリストを取ってくる
func GetFoodData_list(recipe_url string)[][]string{
  var Buff_FoodData [][]string
  //スクレイピング
  res, err := http.Get(recipe_url)
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()

	doc, _ := goquery.NewDocumentFromReader(res.Body)
  doc.Find(".ingredient-list-item").Each(func(i int, data *goquery.Selection){
    text_data := data.Find("a").Text()
    text_data = strings.Replace(text_data,"\n", "", -1)
    Food_name := strings.Replace(text_data," ", "", -1)
    How_many := data.Find("span").Text()
    var buff_data = []string{Food_name, How_many}
    Buff_FoodData = append(Buff_FoodData, buff_data)
  })
  return Buff_FoodData
}
