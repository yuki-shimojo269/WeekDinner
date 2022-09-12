package main

import (
  "fmt"
  "net/url"
  "net/http"
  "path/filepath"
  "os"
  "strings"
  "io/ioutil"
  _ "time"

  "github.com/bamzi/jobrunner"
  "github.com/gin-gonic/gin"

  "TodayMenu/pkg"
)

const(
  json_file = "..."
  calendar_id = "..."
)


var WeeksToInt map[string]int = map[string]int{"Sunday"    : 0,
                                               "Monday"    : 1,
                                               "Tuesday"   : 2,
                                               "Wednesday" : 3,
                                               "Thursday"  : 4,
                                               "Friday"    : 5,
                                               "Saturday"  : 6}

var WeeksEnJp map[string]string = map[string]string{"Sunday"    : "日曜日",
                                                    "Monday"    : "月曜日",
                                                    "Tuesday"   : "火曜日",
                                                    "Wednesday" : "水曜日",
                                                    "Thursday"  : "木曜日",
                                                    "Friday"    : "金曜日",
                                                    "Saturday"  : "土曜日"}

// 月に1回動くやつ=====================
type MonthJob struct {
  sample string
}
//月に1回動かすやつ
func (e MonthJob) Run() {
    fmt.Println("Run Month Job!")
    UpdateRecipeNameDB("おかず")
    fmt.Println("Fin Month Job!")
}
func UpdateRecipeNameDB(serch_word string){
  // mainデータベースを更新 or 作成
  // 月に1回更新

  // =============================
  // serch_word = "おかず"
  home_url := "https://www.kurashiru.com/"
  serch_url := home_url + "search?query=" + url.QueryEscape(serch_word)
  var db_name string = "Recipe.db"
  // =============================

  // データベースの作成
  week_menue.Create_RecipeNameDB(db_name)

  // 量の型を入れるテーブルを作成
  // week_menue.Create_FoodTypeDB()

  // データベースを更新
  week_menue.Update_RecipeNameDB(home_url, serch_url, db_name)
}
// ==============================


// 週に１回動くやつ---------------
type WeekJob struct{
  Sample string
}
func (e WeekJob) Run() {
    fmt.Println("Run Week Job!")
    CreateMenueThisWeekdDB()
    fmt.Println("Fin Week Job!")
}

func CreateMenueThisWeekdDB() {
  weeks, err := ioutil.ReadFile(filepath.Join("data", "NextWeeks.txt"))
  if err != nil{
    fmt.Println(err)
  }
  week_list := strings.Split(string(weeks), ",")

  date_list := []int{}
  for _, week := range week_list[:len(week_list)]{
    date_list = append(date_list, WeeksToInt[week])
  }
  date_length := len(date_list)


  week_menue.Choice_RandomMenue(date_length)
  /*
  Monday,Wednesday
  に入力する
  */
  week_menue.AddEvent(date_length, date_list, json_file, calendar_id)
}
// ------------------------------

// ==============================
func InitData()WebDataForm{
  if err := os.Mkdir("data", 0777); err != nil {
        fmt.Println(err)
  }

  fp, err := os.Create(filepath.Join("data", "NextWeeks.txt"))
  if err != nil {
      panic(err)
  }
  defer fp.Close()


  var web_data WebDataForm
  web_data.File_path = filepath.Join("data", "NextWeeks.txt")
  web_data.Weeks = append(web_data.Weeks, "Sunday","Monday","Tuesday","Wednesday","Thursday","Friday","Saturday")
  web_data.Write()
  web_data.Read()

  fmt.Println("Run Month Job!")
  UpdateRecipeNameDB("おかず")
  fmt.Println("Fin Month Job!")
  fmt.Println("Run Week Job!")
  CreateMenueThisWeekdDB()
  fmt.Println("Fin Week Job!")

  return web_data
}
// ==============================

// 必要なデータを格納するやつ===========
type WebDataForm struct{
  File_path     string
  Weeks       []string  `form:"Weeks[]"`
  Weeks_jp    []string
}
func (e *WebDataForm) Read(){
  weeks, err := ioutil.ReadFile(e.File_path)
  if err != nil{
    fmt.Println(err)
  }
  week_list := strings.Split(string(weeks), ",")

  e.Weeks_jp = nil
  e.Weeks = nil
  for _, week := range week_list[:len(week_list)]{
    e.Weeks = append(e.Weeks, week)
    e.Weeks_jp = append(e.Weeks_jp, WeeksEnJp[week])
  }
}
func (e *WebDataForm) Write(){
  write_text := strings.Join(e.Weeks, ",")
  ioutil.WriteFile(e.File_path,  []byte(write_text), os.ModePerm)
  e.Read()
}
// ==================================

func main() {
    web_data := InitData()

    jobrunner.Start()
    // 毎月の１日の３時にレシピ一覧を取得
    jobrunner.Schedule("0 3 1 * *", MonthJob{})
    // "0 3 1 * *"

    // 毎週土曜日の2時にその週の夕飯を決める
    jobrunner.Schedule("0 2 * * SAT", WeekJob{})
    // "0 3 * * SAT"



    gin.SetMode(gin.ReleaseMode)
    r := gin.Default()
    r.GET("/jobrunner/status", JobJSON)

    r.LoadHTMLGlob("temp/*.tmpl")

    home := r.Group("/")
    {
      home.GET("/", func(c *gin.Context) {
        c.HTML(http.StatusOK, "index.tmpl", gin.H{"eventtime": web_data.Weeks_jp})
      })
      home.POST("/", func(c *gin.Context) {
        c.ShouldBind(&web_data)
        c.HTML(http.StatusOK, "index.tmpl", gin.H{"weks": web_data.Weeks, "eventtime": web_data.Weeks_jp})
        fmt.Println(web_data.Weeks_jp)
        web_data.Write()
        web_data.Read()
        c.Redirect(http.StatusMovedPermanently, "/")
      })
    }

    r.Run(":8080")
}

// JobJSON ...
func JobJSON(c *gin.Context) {
    c.JSON(http.StatusOK, jobrunner.StatusJson())
}
