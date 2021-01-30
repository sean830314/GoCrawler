package nosql

import (
	"fmt"
	"sync"

	"github.com/gocql/gocql"
	"github.com/sirupsen/logrus"
)

type CassandraClient struct {
	Host string
	Port int
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
				continue // Keep trying on i/o error.
			} else {
				fmt.Println("=====================")
				logrus.Info("Batch executed.")
			}
			break
		}
	}
	wg.Done()
}

func (c *CassandraClient) InsertArticles(articles []PttArticle) {
	goroutines := 8
	BatchSizeMaximum := 10
	cluster := gocql.NewCluster(c.Host)
	cluster.Port = c.Port
	cluster.Keyspace = "ptt_keyspace"
	cluster.Consistency = gocql.Any
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
	for i := 0; i < len(articles); i++ {
		if i == (len(articles) - 1) { // Send in the last batch.
			in <- b
			break
		}
		counter++
		b.Query("INSERT INTO ptt_keyspace.article (url, board_name, title, author, date, content, ip, all, count, b, p, n) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", articles[i].Url, articles[i].BoardName, articles[i].Title, articles[i].Author, articles[i].Date, articles[i].Content, articles[i].Ip, articles[i].All, articles[i].Count, articles[i].B, articles[i].P, articles[i].N)
		if counter == BatchSizeMaximum { // Send in the batch and start a new one.
			in <- b
			b = gocql.NewBatch(gocql.LoggedBatch)
			counter = 0
		}
	}
	close(in)
	wg.Wait()
}

func (c *CassandraClient) InsertArticleComments(comments []PttComment) {
	goroutines := 8
	BatchSizeMaximum := 10
	cluster := gocql.NewCluster(c.Host)
	cluster.Port = c.Port
	cluster.Keyspace = "ptt_keyspace"
	cluster.Consistency = gocql.Any
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
	for i := 0; i < len(comments); i++ {
		if i == (len(comments) - 1) { // Send in the last batch.
			in <- b
			break
		}
		counter++
		b.Query("INSERT INTO ptt_keyspace.comment (url, push_content, push_date, push_tag, push_uid) VALUES (?, ?, ?, ?, ?)", comments[i].Url, comments[i].PushContent, comments[i].PushIpdatetime, comments[i].PushTag, comments[i].PushUserID)
		if counter == BatchSizeMaximum { // Send in the batch and start a new one.
			in <- b
			b = gocql.NewBatch(gocql.LoggedBatch)
			counter = 0
		}
	}
	close(in)
	wg.Wait()
}

func (c *CassandraClient) InitCassandra() {
	cluster := gocql.NewCluster(c.Host)
	cluster.Port = c.Port
	c.createKeyspace(cluster, "ptt_keyspace")
	c.createTable(cluster, "ptt_keyspace")
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
	err = session.Query(fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s.article(url text, board_name text, title text, author text, date text, content text, ip text, all int, count int, b int, p int, n int, PRIMARY KEY(url));`, keyspace)).Exec()
	if err != nil {
		logrus.Error("CreateTable", err)
	}
	err = session.Query(fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s.comment(url text, push_tag text, push_uid text, push_content text, push_date text, PRIMARY KEY(push_content));`, keyspace)).Exec()
	if err != nil {
		logrus.Error("CreateTable", err)
	}
}
