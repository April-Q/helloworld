package es

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	elastic "github.com/olivere/elastic/v7"
)

func RunESClient() {
	// Obtain a client and connect to the default Elasticsearch installation
	// on 127.0.0.1:9200. Of course you can configure your client to connect
	// to other hosts and configure it in various other ways.
	client, err := elastic.NewClient(
		elastic.SetSniff(false),
		elastic.SetBasicAuth("elastic", "jH9b50Rk2MvZsns4sgt5TYez"),
		elastic.SetURL("https://observability-deployment-d6bc06.es.eastus2.azure.elastic-cloud.com:9243"),
		elastic.SetErrorLog(log.New(os.Stderr, "ELASTIC ", log.LstdFlags)),
		elastic.SetInfoLog(log.New(os.Stdout, "", log.LstdFlags)),
	)
	if err != nil {
		// Handle error
		fmt.Println(err, client)
		return
	}
	// Open a Point in Time
	pit, err := client.OpenPointInTime("test").KeepAlive("2m").Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	res, err := client.Search().
		Query(
			// Return random results
			elastic.NewFunctionScoreQuery().AddScoreFunc(elastic.NewRandomFunction()),
		).
		Size(10).
		PointInTime(
			elastic.NewPointInTime(pit.Id, "2m"),
		).
		Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	for _, hit := range res.Hits.Hits {
		var doc map[string]interface{}
		if err := json.Unmarshal(hit.Source, &doc); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%v\n", doc)
	}
}

// populate will fill an example index.
func populate(client *elastic.Client, indexName string) error {
	bulk := client.Bulk().Index(indexName)
	for i := 0; i < 10000; i++ {
		doc := map[string]interface{}{
			"name": fmt.Sprintf("Product %d", i+1),
		}
		bulk = bulk.Add(elastic.NewBulkIndexRequest().
			Id(fmt.Sprint(i)).
			Doc(doc),
		)
		if bulk.NumberOfActions() >= 100 {
			_, err := bulk.Do(context.Background())
			if err != nil {
				return err
			}
			// bulk is reset after Do, so you can reuse it
			// We ignore indexing errors here though!
		}
	}
	if bulk.NumberOfActions() > 0 {
		_, err := bulk.Do(context.Background())
		if err != nil {
			return err
		}
	}

	return nil
}

func RunESClient1() {
	// Starting with elastic.v5, you must pass a context to execute each service
	ctx := context.Background()
	// Obtain a client and connect to the default Elasticsearch installation
	// on 127.0.0.1:9200. Of course you can configure your client to connect
	// to other hosts and configure it in various other ways.
	client, err := elastic.NewClient(
		elastic.SetSniff(false),
		elastic.SetBasicAuth("elastic", "i6u2oejwnbkIQPWaFNFu"),
		elastic.SetURL("https://10.108.189.145:9200"),
		elastic.SetErrorLog(log.New(os.Stderr, "ELASTIC ", log.LstdFlags)),
		elastic.SetInfoLog(log.New(os.Stdout, "", log.LstdFlags)),
	)
	if err != nil {
		// Handle error
		fmt.Println(err, client)
		return
	}

	// Ping the Elasticsearch server to get e.g. the version number
	info, code, err := client.Ping("https://10.108.189.145:9200").Do(ctx)
	if err != nil {
		// Handle error
		fmt.Println(err, client)
		return
	}
	fmt.Printf("Elasticsearch returned with code %d and version %s\n", code, info.Version.Number)

	// Use the IndexExists service to check if a specified index exists.
	exists, err := client.IndexExists("test").Do(ctx)
	if err != nil {
		// Handle error
		panic(err)
	}
	if !exists {
		// Create a new index.
		createIndex, err := client.CreateIndex("test").BodyString("").Do(ctx)
		if err != nil {
			// Handle error
			panic(err)
		}
		if !createIndex.Acknowledged {
			// Not acknowledged
			fmt.Println(createIndex.Acknowledged, createIndex.Index)
		}
		fmt.Println("create index:", createIndex.Index)
	}
	// Index a tweet (using JSON serialization)
	tweet1 := Tweet{User: "olivere", Message: "Take Five", Retweets: 0}
	put1, err := client.Index().
		Index("test").
		Id("1").
		BodyJson(tweet1).
		Do(ctx)
	if err != nil {
		// Handle error
		panic(err)
	}
	fmt.Printf("Indexed tweet %s to index %s, type %s,version %d\n", put1.Id, put1.Index, put1.Type, put1.Version)

	// Index a second tweet (by string)
	tweet2 := `{"user" : "olivere", "message" : "It's a Raggy Waltz"}`
	put2, err := client.Index().
		Index("test").
		Id("2").
		BodyString(tweet2).
		Do(ctx)
	if err != nil {
		// Handle error
		panic(err)
	}
	fmt.Printf("Indexed tweet %s to index %s, type %s, version %d\n", put2.Id, put2.Index, put2.Type, put2.Version)

	// Get tweet with specified ID
	get1, err := client.Get().
		Index("test").
		Id("1").
		Do(ctx)
	if err != nil {
		// Handle error
		panic(err)
	}
	if get1.Found {
		fmt.Printf("Got document %s in version %d from index %s, type %s\n", get1.Id, get1.Version, get1.Index, get1.Type)
	}

	// Flush to make sure the documents got written.
	_, err = client.Flush().Index("test").Do(ctx)
	if err != nil {
		panic(err)
	}

	// Search with a term query
	ExistQuery := elastic.NewExistsQuery("user")
	searchResult, err := client.Search().
		Index("test").     // search in index "twitter"
		Query(ExistQuery). // specify the query
		// Sort("user", true). // sort by "user" field, ascending
		From(0).Size(10). // take documents 0-9
		Pretty(true).     // pretty print request and response JSON
		Do(ctx)           // execute
	if err != nil {
		// Handle error
		panic(err)
	}
	// searchResult is of type SearchResult and returns hits, suggestions,
	// and all kinds of other information from Elasticsearch.
	fmt.Printf("Query took %d milliseconds\n", searchResult.TookInMillis)

	// Here's how you iterate through results with full control over each step.
	if searchResult.TotalHits() > 0 {
		fmt.Printf("Found a total of %d tweets\n", searchResult.TotalHits())

		// Iterate through results
		for _, hit := range searchResult.Hits.Hits {
			// hit.Index contains the name of the index

			// Deserialize hit.Source into a Tweet (could also be just a map[string]interface{}).
			var t Tweet
			err := json.Unmarshal(hit.Source, &t)
			if err != nil {
				// Deserialization failed
			}

			// Work with tweet
			fmt.Printf("Tweet by %s: %s\n", t.User, t.Message)
		}
	} else {
		// No hits
		fmt.Print("Found no doc\n")
	}

	// Update a tweet by the update API of Elasticsearch.
	// We just increment the number of retweets.
	// update, err := client.Update().Index("twitter").Type("tweet").Id("1").
	// 	Script(elastic.NewScriptInline("ctx._source.retweets += params.num").Lang("painless").Param("num", 1)).
	// 	Upsert(map[string]interface{}{"retweets": 0}).
	// 	Do(ctx)
	// if err != nil {
	// 	// Handle error
	// 	panic(err)
	// }
	// fmt.Printf("New version of tweet %q is now %d\n", update.Id, update.Version)

	// ...

	// Delete an index.
	// deleteIndex, err := client.DeleteIndex("twitter").Do(ctx)
	// if err != nil {
	// 	// Handle error
	// 	panic(err)
	// }
	// if !deleteIndex.Acknowledged {
	// 	fmt.Println(deleteIndex.Acknowledged)
	// 	// Not acknowledged
	// }
}

// Tweet is a structure used for serializing/deserializing data in Elasticsearch.
type Tweet struct {
	User     string                `json:"user"`
	Message  string                `json:"message"`
	Retweets int                   `json:"retweets"`
	Image    string                `json:"image,omitempty"`
	Created  time.Time             `json:"created,omitempty"`
	Tags     []string              `json:"tags,omitempty"`
	Location string                `json:"location,omitempty"`
	Suggest  *elastic.SuggestField `json:"suggest_field,omitempty"`
}

const mapping = `
{
	"settings":{
		"number_of_shards": 1,
		"number_of_replicas": 0
	},
	"mappings":{
		"tweet":{
			"properties":{
				"user":{
					"type":"keyword"
				},
				"message":{
					"type":"text",
					"store": true,
					"fielddata": true
				},
				"image":{
					"type":"keyword"
				},
				"created":{
					"type":"date"
				},
				"tags":{
					"type":"keyword"
				},
				"location":{
					"type":"geo_point"
				},
				"suggest_field":{
					"type":"completion"
				}
			}
		}
	}
}`
