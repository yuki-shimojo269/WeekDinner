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
3. 実行

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

1. プロジェクトを作成する 参考サイト(https://www.coppla-note.net/posts/tutorial/google-calendar-api/#%E3%83%97%E3%83%AD%E3%82%B8%E3%82%A7%E3%82%AF%E3%83%88%E3%81%AE%E4%BD%9C%E6%88%90)
2. Calender APIを有効化 参考サイト(https://www.coppla-note.net/posts/tutorial/google-calendar-api/#calendar-api-%E3%81%AE%E6%9C%89%E5%8A%B9%E5%8C%96)

# 確認用
## ジョブスケジュールの確認用
```
curl -s http://localhost:8080/jobrunner/status | jq .
```
## スケジュール変更用URL
```
http://localhost:8080/
```
