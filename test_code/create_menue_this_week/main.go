package main

import (
  _ "fmt"
  "net/url"

  "this_week/src/ControlDB"
)

func UpdateRecipeNameDB()  {
  // mainデータベースを更新 or 作成
  // 月に1回更新

  // =============================
  serch_word := "おかず"
  home_url := "https://www.kurashiru.com/"
  serch_url := home_url + "search?query=" + url.QueryEscape(serch_word)
  var db_name string = "RecipeName.db"
  // =============================

  // データベースの作成
  ControlDB.Create_RecipeNameDB(db_name)

  // 量の型を入れるテーブルを作成
  // ControlDB.Create_FoodTypeDB()

  // データベースを更新
  ControlDB.Update_RecipeNameDB(home_url, serch_url, db_name)

}

// ランダムでレシピを取り出す
func CreateMenueThisWeekdDB()  {
  ControlDB.Choice_RandomMenue(2)
}

func main()  {
  // main関数
  // 月１に起動するRecipName.dbを
  UpdateRecipeNameDB()

  //週１で起動する、１っ週間のメニューを決めてくれる
  CreateMenueThisWeekdDB()

}
