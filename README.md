# manu_choice
自動で一週間の献立を決めてくれるwebアプリ

# app install
## install golang
```
$ sudo add-apt-repository ppa:longsleep/golang-backports  
$ sudo apt update  
$ sudo apt install golang  
```
## install sqlite3
```
$ sudo apt install sqlite3  
```
## install go package
```
go get github.com/PuerkitoBio/goquery

go get gorm.io/driver/sqlite
go get gorm.io/gorm

go get github.com/gin-gonic/gin

go get github.com/bamzi/jobrunner

go get -u google.golang.org/api/calendar/v3
go get -u golang.org/x/oauth2/google
```

# 機能一覧
- Kurasiruからレシピをスクレイピング
- 日曜日に一週間の献立とその材料を一覧でリストアップ
- Googleカレンダーに記録する

# GoogleカレンダーAPI
https://console.cloud.google.com/home/dashboard?project=grand-ability-350113

# スケジュール確認用コマンド
```
curl -s http://localhost:8080/jobrunner/status | jq .
```
# スケジュール変更用URL
```
http://localhost:8080/
```
