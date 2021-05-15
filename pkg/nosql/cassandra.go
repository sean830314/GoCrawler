package nosql

import (
	"fmt"
	"sync"
	"time"

	"github.com/gocql/gocql"
	"github.com/sirupsen/logrus"
)

type CassandraClient struct {
	Host string
	Port int
}

type DcardArticle struct {
	ID           int      `json:"id"`
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

type DcardComment struct {
	ID         string `json:"id"`
	CommentURL string `json:"commentURL"`
	CreatedAt  string `json:"createdAt"`
	UpdatedAt  string `json:"updatedAt"`
	Floor      int32  `json:"floor"`
	Content    string `json:"content"`
	LikeCount  int32  `json:"likeCount"`
	Gender     string `json:"gender"`
	School     string `json:"school"`
}

type PttArticle struct {
	Url                 string
	BoardName           string
	Title               string
	Author              string
	Date                string
	Content             string
	Ip                  string
	All, Count, P, B, N int
}
type PttComment struct {
	Url            string
	PushTag        string
	PushUserID     string
	PushContent    string
	PushIpdatetime string
}

func processBatches(s *gocql.Session, in chan *gocql.Batch, wg *sync.WaitGroup) {
	wg.Add(1)
	for batch := range in {
		for {
			if err := s.ExecuteBatch(batch); err != nil {
				logrus.Info(fmt.Sprintf("Couldn't execute batch: %s", err))
			} else {
				logrus.Info("Batch executed.")
			}
			break
		}
	}
	wg.Done()
}

func (c *CassandraClient) InsertData(obj interface{}) {
	goroutines := 4
	BatchSizeMaximum := 10
	cluster := gocql.NewCluster(c.Host)
	cluster.Port = c.Port
	cluster.Keyspace = "social_data"
	cluster.Consistency = gocql.Any
	cluster.ProtoVersion = 4
	cluster.ConnectTimeout = time.Second * 10
	session, err := cluster.CreateSession()
	if err != nil {
		logrus.Error(err)
	}
	defer session.Close()
	in := make(chan *gocql.Batch, 0)
	var wg sync.WaitGroup
	for i := 0; i < goroutines; i++ {
		go processBatches(session, in, &wg)
	}
	//Accumulate batches, and let goroutines execute them.
	counter := 0
	b := session.NewBatch(gocql.LoggedBatch)
	switch v := obj.(type) {
	case []PttArticle:
		for i := 0; i < len(v); i++ {
			if i == (len(v) - 1) { // Send in the last batch.
				in <- b
				break
			}
			counter++
			b.Query("INSERT INTO social_data.ptt_article (url, board_name, title, author, date, content, ip, all, count, b, p, n) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", v[i].Url, v[i].BoardName, v[i].Title, v[i].Author, v[i].Date, v[i].Content, v[i].Ip, v[i].All, v[i].Count, v[i].B, v[i].P, v[i].N)
			if counter == BatchSizeMaximum { // Send in the batch and start a new one.
				in <- b
				b = gocql.NewBatch(gocql.LoggedBatch)
				counter = 0
			}
		}
	case []PttComment:
		for i := 0; i < len(v); i++ {
			if i == (len(v) - 1) { // Send in the last batch.
				in <- b
				break
			}
			counter++
			b.Query("INSERT INTO social_data.ptt_comment (url, push_content, push_date, push_tag, push_uid) VALUES (?, ?, ?, ?, ?)", v[i].Url, v[i].PushContent, v[i].PushIpdatetime, v[i].PushTag, v[i].PushUserID)
			if counter == BatchSizeMaximum { // Send in the batch and start a new one.
				in <- b
				b = gocql.NewBatch(gocql.LoggedBatch)
				counter = 0
			}
		}
	case []DcardArticle:
		for i := 0; i < len(v); i++ {
			if i == (len(v) - 1) { // Send in the last batch.
				in <- b
				break
			}
			counter++
			b.Query("INSERT INTO social_data.dcard_article (id, title, content, excerpt, createdAt, updatedAt, commentCount, forumName, forumAlias, gender, school, likeCount, topics, tags) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", v[i].ID, v[i].Title, v[i].Content, v[i].Excerpt, v[i].CreatedAt, v[i].UpdatedAt, v[i].CommentCount, v[i].ForumName, v[i].ForumAlias, v[i].Gender, v[i].School, v[i].LikeCount, v[i].Topics, v[i].Tags)
			if counter == BatchSizeMaximum { // Send in the batch and start a new one.
				in <- b
				b = gocql.NewBatch(gocql.LoggedBatch)
				counter = 0
			}
		}
	case []DcardComment:
		for i := 0; i < len(v); i++ {
			if i == (len(v) - 1) { // Send in the last batch.
				in <- b
				break
			}
			counter++
			b.Query("INSERT INTO social_data.dcard_comment (id, commentURL, createdAt, updatedAt, floor, content, likeCount, gender, school) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)", v[i].ID, v[i].CommentURL, v[i].CreatedAt, v[i].UpdatedAt, v[i].Floor, v[i].Content, v[i].LikeCount, v[i].Gender, v[i].School)
			if counter == BatchSizeMaximum { // Send in the batch and start a new one.
				in <- b
				b = gocql.NewBatch(gocql.LoggedBatch)
				counter = 0
			}
		}
	default:
		fmt.Printf("I don't know about type %T!\n", v)
	}
	close(in)
	wg.Wait()
}

func (c *CassandraClient) InitCassandra() {
	cluster := gocql.NewCluster(c.Host)
	cluster.Port = c.Port
	c.createKeyspace(cluster, "social_data")
	c.createTable(cluster, "social_data")
}

func (cc *CassandraClient) createKeyspace(cluster *gocql.ClusterConfig, keyspace string) {
	c := cluster
	session, err := c.CreateSession()
	if err != nil {
		logrus.Error("CreateSession Error: ", err)
	}
	defer session.Close()
	err = session.Query(fmt.Sprintf(`CREATE KEYSPACE IF NOT EXISTS %s with replication = { 'class' : 'SimpleStrategy', 'replication_factor' : 1 };`, keyspace)).Exec()
	if err != nil {
		logrus.Error("CreateKeyspace Error: ", err)
	}
}

func (cc *CassandraClient) createTable(cluster *gocql.ClusterConfig, keyspace string) {
	c := cluster
	session, err := c.CreateSession()
	if err != nil {
		logrus.Error("CreateSession Error: ", err)
	}
	defer session.Close()
	err = session.Query(fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s.ptt_article(url text, board_name text, title text, author text, date text, content text, ip text, all int, count int, b int, p int, n int, PRIMARY KEY(url));`, keyspace)).Exec()
	if err != nil {
		logrus.Error("CreateTable", err)
	}
	err = session.Query(fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s.ptt_comment(url text, push_tag text, push_uid text, push_content text, push_date text, PRIMARY KEY(push_content));`, keyspace)).Exec()
	if err != nil {
		logrus.Error("CreateTable", err)
	}
	err = session.Query(fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s.dcard_article(id int, title text, content text, excerpt text, createdAt text, updatedAt text, commentCount int, forumName text, forumAlias text, gender text, school text, likeCount int, topics list<text>, tags list<text>, PRIMARY KEY(id));`, keyspace)).Exec()
	if err != nil {
		logrus.Error("CreateTable", err)
	}
	err = session.Query(fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s.dcard_comment(id text, commentURL text, createdAt text, updatedAt text, floor int, content text, likeCount int, gender text, school text, PRIMARY KEY(content));`, keyspace)).Exec()
	if err != nil {
		logrus.Error("CreateTable", err)
	}

}
