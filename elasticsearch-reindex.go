package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"regexp"
	"sort"
	"time"

	elastic "gopkg.in/olivere/elastic.v5"
)

var (
	// url cluster ES
	url string
	// single index
	singleIndex string
	// pattern index
	patternIndex string
)

func init() {
	flag.StringVar(&url, "u", url, "e.g. http://localhost:9200")
	flag.StringVar(&singleIndex, "s", singleIndex, "single index rendex")
	flag.StringVar(&patternIndex, "p", singleIndex, "rendex by pattern")
}

// ReindexMatched todo
func ReindexMatched(ReindexIndex string) string {

	client, err := elastic.NewClient(
		elastic.SetURL(url),
		// for es on docker disable setsniff
		elastic.SetSniff(true))
	if err != nil {
		// Handle error
		log.Fatal(err)
	}

	WriteDisable := `{"index.blocks.write": true}`

	settings, err := client.IndexPutSettings().
		Index(ReindexIndex).
		BodyString(WriteDisable).
		Do(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	if settings.Acknowledged {
		fmt.Printf("%s: successfully disabled write operations against the index\n", ReindexIndex)
	}

	t := time.Now()
	NewIndexName := ReindexIndex + t.Format("-20060102150405")

	reindex, err := client.Reindex().WaitForActiveShards("all").
		SourceIndex(ReindexIndex).
		DestinationIndex(NewIndexName).
		Conflicts("abort").
		WaitForCompletion(true).
		Do(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s: reindexed a total of %d documents\n", NewIndexName, reindex.Total)

	health, err := client.ClusterHealth().Index(NewIndexName).Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s: index status is %q\n", NewIndexName, health.Status)

	delete, err := client.DeleteIndex(ReindexIndex).Do(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	if delete.Acknowledged {
		fmt.Printf("%s: index deleted\n", ReindexIndex)
	}

	alias, err := client.Alias().Add(NewIndexName, ReindexIndex).Do(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	if alias.Acknowledged {
		fmt.Printf("alias created: old name - %s, new name %s\n", ReindexIndex, NewIndexName)
	}
	return ReindexIndex
}

func main() {

	flag.Parse()
	log.SetFlags(0)

	if url == "" {
		log.Fatal("missing url parameter")
	}

	if singleIndex == "" {
		if patternIndex == "" {
			log.Fatal("missing index name or index patter parameter")
		}
	}

	client, err := elastic.NewClient(
		elastic.SetURL(url),
		// for es on docker disable setsniff
		elastic.SetSniff(true))
	if err != nil {
		log.Fatal(err)
	}

	indices, err := client.IndexNames()
	if err != nil {
		log.Fatal(err)
	}

	// Sort by default
	sort.Strings(indices)

	var SliceIndex []string
	for _, index := range indices {
		if len(patternIndex) > 0 {
			matched, err := regexp.MatchString(patternIndex, index)
			if err != nil {
				log.Fatal("invalid index pattern")
			}
			if matched {
				SliceIndex = append(SliceIndex, index)
			}
		}
	}

	if singleIndex != "" {
		ReindexMatched(singleIndex)
	}

	if patternIndex != "" {
		for _, element := range SliceIndex {
			ReindexMatched(element)
		}
	}
}
