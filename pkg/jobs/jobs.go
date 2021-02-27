package jobs

import (
	"fmt"

	"github.com/sean830314/GoCrawler/pkg/nosql"
	"github.com/sean830314/GoCrawler/pkg/service/ptt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type SaveArticlesJob struct {
	Board   string `json:"board" form:"board" valid:"Required;MaxSize(100)"`
	NumPage int    `json:"num_page" form:"num_page" valid:"Range(1,100)"`
}

func (saj SaveArticlesJob) ExecSaveArtilcesJob() {
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
		logrus.Info("Done crawled %s", url)
	}
	logrus.Info("comments %d", count)
	logrus.Info("article %d", count_article)
}
