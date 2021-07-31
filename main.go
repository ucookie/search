package main

import (
	"fmt"
	"log"
	"os"

	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/document"
	"github.com/blevesearch/bleve/v2/mapping"

	"github.com/blevesearch/bleve/v2/analysis/analyzer/custom"
	"github.com/blevesearch/bleve/v2/analysis/char/html"
	"github.com/blevesearch/bleve/v2/analysis/token/lowercase"
	"github.com/blevesearch/bleve/v2/analysis/tokenizer/unicode"
)

const (
	dataFile = "my.bleve"
)

func createIndex() {
	// DEBUG: 删除旧
	if _, err := os.Stat(dataFile); err == nil {
		log.Printf("index %s is already exists", dataFile)
		return
	}

	m := bleve.NewIndexMapping()
	err := m.AddCustomAnalyzer("my", map[string]interface{}{
		"type": custom.Name,
		"char_filters": []string{
			html.Name,
		},
		"tokenizer": unicode.Name,
		"token_filters": []string{
			lowercase.Name,
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	index, err := bleve.New(dataFile, m)
	if err != nil {
		log.Fatal(err)
	}

	mm := index.Mapping().(*mapping.IndexMappingImpl)
	fmt.Printf(">>>>%+v\n", mm.CustomAnalysis)
	log.Println("创建索引完成:", index.Name())

}

type doc struct {
	Name string `json:"name"`
	Text string `json:"text"`
}

func insertDoc(id string, d *doc) {
	index, err := bleve.Open(dataFile)
	if err != nil {
		log.Fatal(err)
	}
	defer index.Close()

	if err = index.Index(id, d); err != nil {
		log.Fatal(err)
	}
}

func queryDocs(q string) {
	index, err := bleve.Open(dataFile)
	if err != nil {
		log.Fatal(err)
	}
	defer index.Close()

	query := bleve.NewMatchQuery(q)
	search := bleve.NewSearchRequest(query)
	ret, err := index.Search(search)
	if err != nil {
		log.Fatal(err)
	}

	for _, hit := range ret.Hits {

		f, _ := index.Document(hit.ID)
		doc := f.(*document.Document)

		fmt.Printf("%s, %.2f, ", hit.ID, hit.Score)
		for _, v := range doc.Fields {
			// fmt.Printf("%s:%+v, ", v.Name(), string(v.Value()))
			fmt.Printf("%s, ", v.Name())
		}
		fmt.Printf("\n")
	}
}

func main() {
	createIndex()

	// d1 := &doc{Name: "d1", Text: "There are many apple trees in a garden. They’re good friends. One day an old tree is ill. There are many pests in the tree. Leaves of the tree turn yellow. The old tree feels very sad and unwell. Another tree sends for a doctor for him. At first, they send for a pigeon, but she has no idea about it. Then they send for an oriole, and she can’t treat the old tree well. Then they send for a woodpecker. She is a good doctor. She pecks a hole in the tree and eats lots of pests. At last the old tree becomes better and better. Leaves turn green and green."}
	// d2 := &doc{Name: "d2", Text: "Today is Sunday! On Sundays, I usually play the flute.My father usually reads the newspaper. My motherusuallycleansthe house. Buttoday my mother is in bed. She is ill. My father has to do the housework. Now, he is cleaning the house. “Sam, can you help me?” “Yes, Dad!” Now, we’re washing the car. Where’s my sister, Amy? She is playing my flute. What a lucky girl!"}
	// d3 := &doc{Name: "d3", Text: "Today is Susan’s birthday. She is nine years old. Her friends are in her home now. There is a birthday party in the evening. Look! Mary is listening to the music. And Tom is drinking orange juice. Jack and Sam are playing cards on the floor. Lily and Amy are watching TV. Someone is knocking at the door. It’s Henry. He brings a big teddy bear for Susan. The teddy bear is yellow. Susan is very happy. All the children are happy. They sing a birthday song for Susan"}
	// d4 := &doc{Name: "d4", Text: "It was a cold winter day.A farmer found a snake on the ground. It was nearly dead by cold. The Farmer was a kind man. Hepicked up thesnake carefully and put it under the coat. Soon the snake Began to move and it raised its mouth and bit the farmer. “Oh, My god!” said the farmer, “I save your life, but you thank me in that way. You must die.” Then he killed the snake with a stick. At last he died, too"}
	// d5 := &doc{Name: "d5", Text: "Hong Kong is a nice place, especially in summer. JulyisahotmonthinHongKong.Butit’san excellent time for swimming. There is a beautiful beach at Repulse Bay (浅水湾). To get there, you can take a bus from Central. Lots of people go to the beach on Sundays and Saturdays. But if you go on a weekday, it is will be not so crowded.Visitors to Hongkong need passports. But people from many countries do not need visas. Hongkong is a nice place for holiday. There are many shops"}
	// insertDoc("001", d1)
	// insertDoc("002", d2)
	// insertDoc("003", d3)
	// insertDoc("004", d4)
	// insertDoc("005", d5)

	// queryDocs("summer")

}
