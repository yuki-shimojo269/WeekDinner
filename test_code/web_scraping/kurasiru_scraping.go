package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

const url = "https://www.kurashiru.com/search?query=%E3%81%8A%E3%81%8B%E3%81%9A"

func main() {
	res, err := http.Get(url)
	if err != nil {
		log.Println(err)
	}
	defer res.Body.Close()

	/* 見出しの取得 */
	doc, _ := goquery.NewDocumentFromReader(res.Body)

	doc.Find("div.video-list-info").Each(func(i int, s *goquery.Selection) {
    s.Find("a").Each(func(j int, a *goquery.Selection){
      recipi_name := a.Text()
      fmt.Println(recipi_name)

      recipe_url, is_url := a.Attr("href")
      if is_url{
        fmt.Println(recipe_url)
      }
      fmt.Println("------------")
    })
	})
}
