package es

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	elasticsearch "github.com/elastic/go-elasticsearch/v7"
	"github.com/go-logr/logr"
	"github.com/robfig/cron"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func Start() {
	url := "https://observability-deployment-d6bc06.es.eastus2.azure.elastic-cloud.com:9243"
	username := "elastic"
	password := "jH9b50Rk2MvZsns4sgt5TYez"
	esconfig, err := NewEsConfig(url, username, password)
	if err != nil {
		fmt.Println(err)
	}

	// start to run es alert
	ctx := context.Background()
	esconfig.run(ctx)
}

func TryChan() {
	c1 := make(chan string)
	c2 := make(chan string)
	go func() {
		time.Sleep(time.Second * 1)
		c1 <- "one"
	}()
	go func() {
		time.Sleep(time.Second * 2)
		c2 <- "two"
	}()
	for i := 0; i < 2; i++ {
		select {
		// case msg1 := <-c1:
		// 	fmt.Println("received", msg1)
		case msg2 := <-c2:
			fmt.Println("received", msg2)
		default:
			fmt.Println("Default")
		}
	}
}

func TryChan2() {
	triggerChan := make(chan int, 1)
	// go uu(triggerChan)

	go func() {
		for i := 0; i < 10; i++ {
			QueueElasticAlertTrigger(context.Background(), triggerChan, i)
		}
	}()

	time.Sleep(3 * time.Second)
}

func uu(channel chan int) {
	for {
		select {
		case a := <-channel:
			fmt.Println(a)
			break
		}
	}
}

func QueueElasticAlertTrigger(ctx context.Context, channel chan int, trigger int) error {
	select {
	case <-ctx.Done():
		return nil
	case channel <- trigger:
		fmt.Println("send data", trigger)
		return nil
	default:
		return fmt.Errorf("channel is blocked")
	}
}

const (
	defaultTimestampFormat string = time.RFC3339
	elasticsearchPrefix    string = "elasticSearchPrefix"
)

type EsConfig struct {
	esclient   elasticsearch.Client
	ruleConfig RuleConfig
	schedule   cron.Schedule
	log        logr.Logger
	// cache knows how to load Kubernetes objects.
	cache cache.Cache
	// Context carries values across API boundaries.
	context.Context
	// client knows how to perform CRUD operations on Kubernetes objects.
	client client.Client
}

// RuleConfig represents a rule configuration file.
type RuleConfig struct {
	// Name is the name of the rule. This value should come
	// from the 'name' field of the rule configuration file
	Name string `json:"name"`

	// ElasticsearchIndex is the index that this rule should
	// query. This value should come from the 'index' field
	// of the rule configuration file
	ElasticsearchIndex string `json:"index"`

	// CronSchedule is the interval at which the
	// *github.com/morningconsult/go-elasticsearch-alerts/command/query.QueryHandler
	// will execute the query. This value should come from
	// the 'schedule' field of the rule configuration file
	CronSchedule string `json:"schedule"`

	// ElasticsearchBodyRaw is the untyped query that this
	// alert should send when querying Elasticsearch. This
	// value should come from the 'body' field of the
	// rule configuration file
	ElasticsearchBodyRaw interface{} `json:"body"`
}

func NewEsConfig(url string, username, password string) (*EsConfig, error) {
	cfg := elasticsearch.Config{
		Username: username,
		Password: password,
		Addresses: []string{
			url,
		},
	}
	es, _ := elasticsearch.NewClient(cfg)
	res, err := es.Info()
	if err != nil {
		fmt.Println(err)
	}
	var r map[string]interface{}
	defer res.Body.Close()
	// Check response status
	if res.IsError() {
		return nil, fmt.Errorf("Error: %s", res.String())
	}
	// Deserialize the response into a map.
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, fmt.Errorf("Error parsing the response body: %s", err)
	}
	schedule, err := cron.Parse("@every 2m")
	if err != nil {
		return nil, fmt.Errorf("error parsing cron schedule: %v", err)
	}

	return &EsConfig{
		esclient: *es,
		schedule: schedule,
		log:      logr.FromContext(context.Background()),
	}, nil
}

func (es *EsConfig) run(ctx context.Context) {
	var (
		now  = time.Now()
		next = now
	)

	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(next.Sub(now)):
			data, err := es.query()
			if err != nil {
				es.log.Error(err, "error querying Elasticsearch")
				break
			}
			// send alert
			if data["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"].(float64) > 0 {
				es.process(data)
			}
		}
		now = time.Now()
		next = es.schedule.Next(now)
	}
}

func (es *EsConfig) process(respData map[string]interface{}) error {
	// Print the ID and document source for each hit.
	for _, hit := range respData["hits"].(map[string]interface{})["hits"].([]interface{}) {
		alert := hit.(map[string]interface{})["_source"].(map[string]interface{})["title"]
		// match trigger : assume get trigger name from hit
		fmt.Println("alert will create diagnosis", alert)

	}

	return nil
}

func (es *EsConfig) query() (map[string]interface{}, error) {

	// Build the request body.
	var buf bytes.Buffer
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"title": "alert",
			},
		},
	}

	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, fmt.Errorf("Error encoding query: %s", err)
	}

	// Perform the search request.
	res, err := es.esclient.Search(
		es.esclient.Search.WithContext(context.Background()),
		es.esclient.Search.WithIndex("test"),
		es.esclient.Search.WithBody(&buf),
		es.esclient.Search.WithTrackTotalHits(true),
		es.esclient.Search.WithPretty(),
	)
	if err != nil {
		return nil, fmt.Errorf("Error getting response: %s", err)
	}

	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			return nil, fmt.Errorf("Error parsing the response body: %s", err)
		} else {
			// Print the response status and error information.
			return nil, fmt.Errorf("[%s] %s: %s",
				res.Status(),
				e["error"].(map[string]interface{})["type"],
				e["error"].(map[string]interface{})["reason"],
			)
		}
	}

	var data map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("Error parsing the response body: %s", err)
	}
	// Print the response status, number of results, and request duration.
	log.Printf(
		"[%s] %d hits; took: %dms",
		res.Status(),
		int(data["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"].(float64)),
		int(data["took"].(float64)),
	)

	return data, nil
}
