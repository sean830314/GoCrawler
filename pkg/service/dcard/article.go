package dcard

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/sirupsen/logrus"
)

type Article struct {
	ID           string   `json:"id"`
	Title        string   `json:"title"`
	Content      string   `json:"content"`
	Excerpt      string   `json:"excerpt"`
	CreatedAt    string   `json:"createdAt"`
	UpdatedAt    string   `json:"updatedAt"`
	CommentCount int32    `json:"commentCount"`
	ForumName    string   `json:"forumName"`
	ForumAlias   string   `json:"forumAlias"`
	Gender       string   `json:"gender"`
	School       string   `json:"school"`
	LikeCount    int32    `json:"likeCount"`
	Topics       []string `json:"topics"`
	Tags         []string `json:"tags"`
}

type Comment struct {
	ID        string `json:"id"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
	Floor     int32  `json:"floor"`
	Content   string `json:"content"`
	LikeCount int32  `json:"likeCount"`
	Gender    string `json:"gender"`
	School    string `json:"school"`
}

func GetArticle(url string) (*Article, error) {
	logrus.Info(fmt.Sprintf("[Get] article url: %s", url))
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var article Article
	json.Unmarshal(body, &article)
	return &article, nil
}

func GetComment(url string) ([]Comment, error) {
	logrus.Info(fmt.Sprintf("[Get] comment url: %s", url))
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var comments []Comment
	json.Unmarshal(body, &comments)
	return comments, nil
}

func main() {
	articleId := 235492620
	articleURL := fmt.Sprintf("https://www.dcard.tw/_api/posts/%d", articleId)
	article, err := GetArticle(articleURL)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(article)
	commentURL := fmt.Sprintf("https://www.dcard.tw/_api/posts/%d/comments", articleId)
	comments, err := GetComment(commentURL)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(comments)
}
