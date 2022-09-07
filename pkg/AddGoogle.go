package week_menue

import (
  "fmt"
  "time"
  "strings"
  "path/filepath"
  "log"
  "context"
  "unicode/utf8"
  "strconv"

  "google.golang.org/api/calendar/v3"
  "google.golang.org/api/option"
  "gorm.io/gorm"
  "gorm.io/driver/sqlite"
)


func createEvent(data GoolgeCalendar) *calendar.Event {

    //start_datatime := "2022" + "-" + "8" + "-" + "18" + "T" + "11:00:00" + ":00+09:00"
    //end_datatime :=  "2022" + "-" + "8" + "-" + "18" + "T" + "12:00:00" + ":00+09:00"


    event := &calendar.Event{
      Summary: "夕飯："+data.Recipe_name,
      Location: data.Recipe_url,
      Description: data.Ingredient,
      Start: &calendar.EventDateTime{
        DateTime: data.Event_day+"T17:00:00",
        TimeZone: "Asia/Tokyo",
      },
      End: &calendar.EventDateTime{
        DateTime: data.Event_day+"T18:00:00",
        TimeZone: "Asia/Tokyo",
      },
    }

    return event
}

func createEventData_fromDb(event_num int, date_list []int)[]GoolgeCalendar{
  // 日付の初期化----------------------
	// 稼働した週の月曜日を返す
	const DateFormat = "2006-01-02"
	now_time := time.Now()
	Sunday := now_time.AddDate(0, 0, -1*int(now_time.Weekday()))
	Sunday_day := strings.Replace(Sunday.Format(DateFormat), "-", "", -1)
	// ----------------------------------

  // 一週間分のデータベースの作成=========
  db_path  := filepath.Join("data", "LogTable.db")
  db, err := gorm.Open(sqlite.Open(db_path), &gorm.Config{})
  if err != nil {
    panic("failed to connect database")
  }
	// ===================================

  // データを取り出す ===================
  var data_list []WeekRecipeList
  db.Model(&WeekRecipeList{}).
     Table("RecipeList").
     Select("*").
     Where("Week = ?", Sunday_day).
     Find(&data_list)
  // ===================================

  // 材料のデータセットを持ってくる=======
  db_path = filepath.Join("data", "Recipe.db")
  home_db, err := gorm.Open(sqlite.Open(db_path), &gorm.Config{})
  // ===================================

  // データを入れる----------------------
  var eventdata_list []GoolgeCalendar
  var eventdata GoolgeCalendar
  for i:=0; i<event_num; i++{
    var FoodData []RecipeData
    recipe_id := strconv.Itoa(int(data_list[i].Data_id))
    home_db.Table("recipe_"+recipe_id).
            Select("*").
            Find(&FoodData)
    var Ingredient string

    for i:=0; i<len(FoodData); i++{
      Ingredient = Ingredient + FoodData[i].Food_name
      for j:=0; j<12-utf8.RuneCountInString(FoodData[i].Food_name); j++{
        Ingredient = Ingredient+"　"
      }
      Ingredient = Ingredient + FoodData[i].Food_amount_type + "\n"
    }
    fmt.Println(Ingredient)

    eventdata.Recipe_id = data_list[i].Data_id
    eventdata.Recipe_name = data_list[i].Recipe_name
    eventdata.Recipe_url = data_list[i].Recipe_url
    eventdata.Ingredient = Ingredient
    eventdata.Event_day = Sunday.AddDate(0, 0, date_list[i]).Format(DateFormat)
    eventdata_list = append(eventdata_list, eventdata )
  }
  // -----------------------------------
  return eventdata_list
}

func AddEvent(event_num int, date_list []int, json_file string, calendar_id string) (err error){
  // イベントを作成^^^^^^^^^^^^^^^^^
  eventdata_list := createEventData_fromDb(event_num, date_list)
  // ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

  // 初期化=========================
  ctx := context.Background()
	calendarService, err := calendar.NewService(ctx, option.WithCredentialsFile(json_file))
	if err != nil {
		return err
	}
  // ==============================

  // イベント追加-------------------
  for i := 0; i < event_num; i++{
    event, err := calendarService.Events.Insert(calendar_id, createEvent(eventdata_list[i])).Do()
    if err != nil {
      log.Fatalf("Unable to create event. %v\n", err)
    }
    fmt.Printf("Event created: %s\n", event.HtmlLink)
  }
  // ------------------------------
  return nil
}
