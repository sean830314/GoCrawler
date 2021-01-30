package ptt

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/sean830314/GoCrawler/pkg/utils"

	"github.com/PuerkitoBio/goquery"
)

type BoardArticle struct {
	Id        string
	BoardName string
	Title     string
	Url       string
}

func GetPagesFromBoard(boardName string) (int, error) {
	url := "https://www.ptt.cc/bbs/" + boardName + "/index.html"
	board := &BoardArticle{}
	dom, err := board.ParseBoard(url)
	if err != nil {
		return 0, err
	}
	found := board.CheckIsFound(dom)
	if found != true {
		return 0, errors.New(fmt.Sprintf("Error: url %s not found", url))
	}
	page_group := dom.Find("div.btn-group-paging")
	var pages int
	page_group.Find("a").Each(func(i int, s *goquery.Selection) {
		text := s.Text()
		if text == "‹ 上頁" {
			url, exists := s.Attr("href")
			if exists {
				url = strings.Split(strings.Split(url, "index")[1], ".")[0]
				pages, _ = strconv.Atoi(url)
			} else {
				pages = 0
			}
		}
	})
	return pages, nil
}

func GetArticlesFromBoard(url string) ([]*BoardArticle, error) {
	board := &BoardArticle{}
	dom, err := board.ParseBoard(url)
	if err != nil {
		return nil, err
	}
	found := board.CheckIsFound(dom)
	if found != true {
		return nil, errors.New(fmt.Sprintf("Error: url %s not found", url))
	}
	main_content := dom.Find("div#main-container")
	boardArticles := make([]*BoardArticle, 0)
	main_content.Find("div.r-ent").Each(func(i int, s *goquery.Selection) {
		ba := BoardArticle{}
		title := strings.Trim(s.Find("div.title").Text(), ": \t\n\r")
		url, exists := s.Find("div.title").Find("a").Attr("href")
		if !exists {
			fmt.Println("url href is not exists")
		} else {
			ba.Title = title
			ba.Url = fmt.Sprintf("https://www.ptt.cc%s", url)
			urlStr := strings.Split(url, "/")
			ba.Id = urlStr[3]
			ba.BoardName = urlStr[2]
			boardArticles = append(boardArticles, &ba)
		}
	})
	return boardArticles, nil
}

func (b *BoardArticle) ParseBoard(url string) (*goquery.Document, error) {
	resp, err := utils.GetResponseWithCookie(url)
	if err != nil {
		return nil, err
	}
	dom, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		return nil, err
	}
	return dom, err
}

func (b *BoardArticle) CheckIsFound(dom *goquery.Document) bool {
	noFound := dom.Find("div.bbs-content").Text()
	if noFound == "404 - Not Found." {
		return false
	}
	return true
}
