package jobs

import (
	"fmt"

	"github.com/sean830314/GoCrawler/pkg/nosql"
	"github.com/sean830314/GoCrawler/pkg/service/dcard"
	"github.com/sean830314/GoCrawler/pkg/service/ptt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type SavePttArticlesJob struct {
	Board   string `json:"board" form:"board" valid:"Required;MaxSize(100)"`
	NumPage int    `json:"num_page" form:"num_page" valid:"Range(1,100)"`
}

func (saj SavePttArticlesJob) ExecSaveArtilcesJob() {
	c := nosql.CassandraClient{
		Host: viper.GetString("cassandra.host"),
		Port: viper.GetInt("cassandra.port"),
	}
	c.InitCassandra()
	pages, err := ptt.GetPagesFromBoard(saj.Board)
	if err != nil {
		logrus.Error(err)
	}
	count := 0
	count_article := 0
	for i := 0; i < saj.NumPage; i++ {
		url := fmt.Sprintf("https://www.ptt.cc/bbs/%s/%s", saj.Board, fmt.Sprintf("index%d.html", pages-i))
		logrus.Info("start crawling %s", url)
		articlesMeta, err := ptt.GetArticlesFromBoard(url)
		if err != nil {
			logrus.Error(err)
		}
		articles := []nosql.PttArticle{}
		for j := 0; j < len(articlesMeta); j++ {
			count_article = count_article + 1
			articleData, err := ptt.GetArticle(articlesMeta[j].Url)
			if err != nil {
				logrus.Error(err)
			}
			article := nosql.PttArticle{
				Url:       articlesMeta[j].Url,
				BoardName: articlesMeta[j].BoardName,
				Title:     articleData.Title,
				Author:    articleData.Author,
				Date:      articleData.Date,
				Content:   articleData.Content,
				Ip:        articleData.Ip,
				All:       articleData.All,
				Count:     articleData.Count,
				P:         articleData.P,
				B:         articleData.B,
				N:         articleData.N,
			}
			comments := []nosql.PttComment{}
			for k := 0; k < len(articleData.Comments); k++ {
				comment := nosql.PttComment{
					Url:            articlesMeta[j].Url,
					PushTag:        articleData.Comments[k].PushTag,
					PushUserID:     articleData.Comments[k].PushUserID,
					PushContent:    articleData.Comments[k].PushContent,
					PushIpdatetime: articleData.Comments[k].PushIpdatetime,
				}
				comments = append(comments, comment)
			}
			c.InsertArticleComments(comments)
			count = count + len(comments)
			articles = append(articles, article)
		}
		c.InsertArticles(articles)
		logrus.Info("Done crawled", url)
	}
	logrus.Info("num of ptt comments", count)
	logrus.Info("num of ptt article", count_article)
}

type SaveDcardArticlesJob struct {
	BoardID    string `json:"board_id" form:"board_id" valid:"Required;MaxSize(100)"`
	NumArticle int    `json:"num_article" form:"num_article" valid:"Range(1,100)"`
}

func (saj SaveDcardArticlesJob) ExecSaveArtilcesJob() {
	// c := nosql.CassandraClient{
	// 	Host: viper.GetString("cassandra.host"),
	// 	Port: viper.GetInt("cassandra.port"),
	// }
	//c.InitCassandra()
	boardURL := fmt.Sprintf("http://dcard.tw/_api/forums/%s/posts", saj.BoardID)
	articles, err := dcard.GetArticlesFromBoard(boardURL)
	if err != nil {
		logrus.Error(err)
	}
	if articles != nil && len(articles) >= saj.NumArticle {
		var urls []string
		for i := 0; i < saj.NumArticle; i++ {
			urls = append(urls, fmt.Sprintf("https://www.dcard.tw/_api/posts/%d", articles[i].ID))
		}
		logrus.Info("Crawling urls: ", urls)
		for i := 0; i < len(urls); i++ {
			logrus.Info("start crawling ", urls[i])
			article, err := dcard.GetArticle(urls[i])
			if err != nil {
				logrus.Error(err)
			}
			commentURL := fmt.Sprintf("%s/comments", urls[i])
			comments, err := dcard.GetComment(commentURL)
			if err != nil {
				logrus.Error(err)
			}
			logrus.Info(fmt.Sprintf("num of dcard article(%s) comment %d", article.ID, len(comments)))
		}
		logrus.Info("num of dcard article ", len(urls))
	}
}
