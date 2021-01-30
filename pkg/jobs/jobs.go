package jobs

import (
	"fmt"

	"github.com/sean830314/GoCrawler/pkg/service/ptt"
	"github.com/sirupsen/logrus"
)

type SaveArticlesJob struct {
	Board   string `json:"board" form:"board" valid:"Required;MaxSize(100)"`
	NumPage int    `json:"num_page" form:"num_page" valid:"Range(1,100)"`
}

type PttArticleStruct struct {
	meta ptt.BoardArticle
	data ptt.Article
}

func (saj SaveArticlesJob) ExecSaveArtilcesJob() {
	pages, err := ptt.GetPagesFromBoard(saj.Board)
	if err != nil {
		logrus.Error(err)
	}
	for i := 0; i < saj.NumPage; i++ {
		url := fmt.Sprintf("https://www.ptt.cc/bbs/%s/%s", saj.Board, fmt.Sprintf("index%d.html", pages-i))
		articlesMeta, err := ptt.GetArticlesFromBoard(url)
		if err != nil {
			logrus.Error(err)
		}
		for i := 0; i < len(articlesMeta); i++ {
			articleData, err := ptt.GetArticle(articlesMeta[i].Url)
			if err != nil {
				logrus.Error(err)
			}
			article := PttArticleStruct{
				meta: *articlesMeta[i],
				data: *articleData,
			}
			fmt.Printf("%v\n", article.meta)
		}
	}
}
