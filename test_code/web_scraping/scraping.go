package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

const url = "https://ja.wikipedia.org/wiki/SCADA"

func main() {
	res, err := http.Get(url)
	if err != nil {
		log.Println(err)
	}
	defer res.Body.Close()

	/* 見出しの取得 */
	doc, _ := goquery.NewDocumentFromReader(res.Body)
	doc.Find(".mw-headline").Each(func(i int, s *goquery.Selection) {
		fmt.Println(s.Text())
	})

	fmt.Println("\n--------------------------------------------------")
	/* 目次の取得 */
	doc.Find(".tocnumber").Each(func(i int, s *goquery.Selection) {
		fmt.Println(s.Text(), " ", s.Next().Text())
	})
}
