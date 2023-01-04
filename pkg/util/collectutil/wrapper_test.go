package collectutil

import (
	"encoding/json"
	"io/ioutil"
	"math"
	"testing"
	"time"
)

type testArticle struct {
	ArticleId     string   `json:"article_id"`
	Title         string   `json:"title"`
	Content       string   `json:"content"`
	PublishedDate string   `json:"published_date"`
	URLtoMedia    []string `json:"url_to_media"`
	URltoArticle  string   `json:"url_to_article"`
}

type withScore struct {
	article testArticle
	score   float64
}

func getArticleArr() []testArticle {

	var (
		mapArt   map[string][]testArticle
		articles []testArticle
	)

	file, err := ioutil.ReadFile("./report.json")
	if err != nil {
		panic(err)
	}

	json.Unmarshal([]byte(file), &mapArt)
	for _, val := range mapArt {
		articles = append(articles, val...)
	}
	return articles
}

func TestRelevancyUsingCollection(t *testing.T) {
	articles := getArticleArr()
	urlVariety := map[string]int{}
	result := *Map(articles, func(v testArticle) withScore {
		var (
			prevVariety int = 1
			url             = v.URltoArticle
			pubTime, _      = time.Parse("2006-01-02 15:04:05", v.PublishedDate)
			relevancy       = time.Now().UTC().Sub(pubTime).Seconds()
		)

		val, ok := urlVariety[url]
		if ok {
			prevVariety = val
		}
		relevancy = math.Log(relevancy)
		if math.IsNaN(relevancy) {
			relevancy = 0.1
		}
		score := relevancy * float64(prevVariety) * 2
		urlVariety[url] = prevVariety * 2
		return withScore{
			score:   score,
			article: testArticle{},
		}
	})

	Sort(result, func(i, j int) bool {
		return result[i].score < result[j].score
	})

}

func TestChunking(t *testing.T) {
	articles := getArticleArr()
	result := Chunk(articles, 10)

	reqChunkLen := int(len(articles) / 10)

	if len(articles)%10 != 0 {
		reqChunkLen++
	}

	if reqChunkLen != len(result) {
		t.Fatal("Chunking not working as expected")
	}
}
