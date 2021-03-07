package dcard

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/sirupsen/logrus"
)

type Board struct {
	ID                string   `json:"id"`
	Alias             string   `json:"alias"`
	Name              string   `json:"name"`
	Description       string   `json:"description"`
	SubscriptionCount string   `json:"subscriptionCount"`
	CreatedAt         string   `json:"createdAt"`
	UpdatedAt         string   `json:"updatedAt"`
	Topics            []string `json:"topics"`
	Subcategories     []string `json:"subcategories"`
	PostCount         int32    `json:"postCount"`
}

func ListBoards() (*[]Board, error) {
	url := "http://dcard.tw/_api/forums"
	logrus.Info(fmt.Sprintf("[Get] all boards url: %s", url))
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var boards []Board
	json.Unmarshal(body, &boards)
	return &boards, nil
}

type BoardArticle struct {
	ID int32 `json:"id"`
}

func GetArticlesFromBoard(url string) ([]BoardArticle, error) {
	logrus.Info(fmt.Sprintf("[Get] articles of board url: %s", url))
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var articles []BoardArticle
	json.Unmarshal(body, &articles)
	return articles, nil
}
