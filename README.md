# 自動で一週間の献立を決めてくれるwebアプリ

# 機能一覧
- Kurasiruからレシピをスクレイピング
- 日曜日に一週間の献立とその材料を一覧でリストアップ
- Googleカレンダーに記録する

# 開発OS
Ubuntu-20.04 LTS (windows10 WSL)

# 手順
1. golang と golangのパッケージ、sqlite3 をインストール
2. GCP(Google Cloud Platform)の設定
3. main.goを編集し実行

# 1 各種ソフトをインストール
## golangをインストール
```
$ sudo add-apt-repository ppa:longsleep/golang-backports  
$ sudo apt update  
$ sudo apt install golang  
```
## sqlite3をインストール
```
$ sudo apt install sqlite3  
```
## golang のパッケージをインストール
```
go get github.com/PuerkitoBio/goquery

go get gorm.io/driver/sqlite
go get gorm.io/gorm

go get github.com/gin-gonic/gin

go get github.com/bamzi/jobrunner

go get -u google.golang.org/api/calendar/v3
go get -u golang.org/x/oauth2/google
```

# 2.GoogleカレンダーAPI
google cloud platformのURL
https://console.cloud.google.com/
参考になるサイト
https://www.coppla-note.net/posts/tutorial/google-calendar-api/

1. プロジェクトを作成する [ここ](https://www.coppla-note.net/posts/tutorial/google-calendar-api/#%E3%83%97%E3%83%AD%E3%82%B8%E3%82%A7%E3%82%AF%E3%83%88%E3%81%AE%E4%BD%9C%E6%88%90)
2. Calender APIを有効化 [ここ](https://www.coppla-note.net/posts/tutorial/google-calendar-api/#calendar-api-%E3%81%AE%E6%9C%89%E5%8A%B9%E5%8C%96)
3. サービスアカウントの設定とjsonのダウンロード [ここ](https://www.coppla-note.net/posts/tutorial/google-calendar-api/#%E8%A3%9C%E8%B6%B3%E3%82%B5%E3%83%BC%E3%83%93%E3%82%B9%E3%82%A2%E3%82%AB%E3%82%A6%E3%83%B3%E3%83%88%E3%82%92%E4%BD%BF%E3%81%A3%E3%81%9F%E3%82%84%E3%82%8A%E6%96%B9)
4. ダウンロードしたjsonを直下に置く(別に直下である必要はないが、pathを考えるとめんどい)
5. Googleカレンダーにサービスアカウントを使えるようにする。この時に「カレンダーID」をメモしておく[ここ](https://www.coppla-note.net/posts/tutorial/google-calendar-api/#%E3%82%B5%E3%83%BC%E3%83%93%E3%82%B9%E3%82%A2%E3%82%AB%E3%82%A6%E3%83%B3%E3%83%88%E3%81%AB%E3%82%AB%E3%83%AC%E3%83%B3%E3%83%80%E3%83%BC%E3%82%92%E5%85%B1%E6%9C%89%E3%81%99%E3%82%8B)

# 3 main.goを編集し実行
1. main.goの18~21行目を編集する
```
const(
  json_file = "..."   //ここにjsonファイルの名前（hogehoge.json）
  calendar_id = "..." //ここにカレンダーID (hoge@group.calendar.google.com)
)
```
2. go run main.go

# 確認用
## ジョブスケジュールの確認用
```
curl -s http://localhost:8080/jobrunner/status | jq .
```
## スケジュール変更用URL
```
http://localhost:8080/
```
