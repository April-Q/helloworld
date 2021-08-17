package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/prometheus/alertmanager/template"
)

func main() {
	// es.TryChan()

	// es.Start()
	// es.EsProcessor()
	// es.RunESClient1()
	// prettyjson()
	// printable2()
	// es.RunESClient()

	// fileupload()
	httpServer()

}

// Test Command:
// curl -X POST  -H "Content-Type: application/json" -d '{"ation":"ssss"}'  http://127.0.0.1:8080/clusters/?cluster=2
func httpServer() {
	r := mux.NewRouter()
	r.HandleFunc("/clusters/", NotificationHandler)
	r.HandleFunc("/sendEmail", NotificationHandler)
	r.HandleFunc("/sendSMS", NotificationHandler)
	http.ListenAndServe(":8080", r)
}

func NotificationHandler(w http.ResponseWriter, r *http.Request) {
	cache := make(map[string]interface{})
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var parameters template.Data

	err = json.Unmarshal(body, &parameters)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cache["url"] = r.RequestURI
	// Update cache with parameters.
	cache["notification"] = parameters
	// data, err := json.Marshal(cache)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	hits, err := json.MarshalIndent(cache, "", "    ")
	if err != nil {
		return
	}
	fmt.Println(string(hits))
	// Response with cache.
	w.Header().Set("Content-Type", "application/json")
	w.Write(hits)
}

type GrafanaNotification struct {
	DashboardId int `json:"dashboardId"`
	EvalMatches []struct {
		Value  float64     `json:"value"`
		Metric string      `json:"metric"`
		Tags   interface{} `json:"tags"`
	} `json:"evalMatches"`
	ImageUrl string      `json:"imageUrl"`
	Message  string      `json:"message"`
	OrgId    int         `json:"orgId"`
	PanelId  int         `json:"panelId"`
	RuleId   int         `json:"ruleId"`
	RuleName string      `json:"ruleName"`
	RuleUrl  string      `json:"ruleUrl"`
	State    string      `json:"state"`
	Tags     interface{} `json:"tags"`
	Title    string      `json:"title"`
}

// Notification is the message which sends to the notification center.
type Notification struct {
	Type string `json:"type"`
	Data *Data  `json:"data"`
}

// Data is the payload of the request body.
type Data struct {
	To        string `json:"to"`
	Content   string `json:"content"`
	IsSync    bool   `json:"isSync"`
	AppName   string `json:"ac_appName"`
	Timestamp string `json:"ac_timestamp"`
	Signature string `json:"ac_signature"`
	Sender    string `json:"sender,omitempty"`
	Title     string `json:"title,omitempty"`
}

func Fib(n int) int {
	if n < 2 {
		return n
	}
	return Fib(n-1) + Fib(n-2)
}

func printMap() {
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": []map[string]interface{}{
					{
						"match": map[string]interface{}{
							"message": "match",
						},
					},
				},
			},
		},
	}
	timeFrom := "sss"
	timeTo := "sss"
	if timeFrom != "" && timeTo != "" {
		queryRange := map[string]interface{}{
			"range": map[string]interface{}{
				"@timestamp": map[string]interface{}{
					"gte": timeFrom,
					"lte": timeTo,
				},
			},
		}
		mustMatch := query["query"].(map[string]interface{})["bool"].(map[string]interface{})["must"].([]map[string]interface{})
		mustMatch = append(mustMatch, queryRange)
		query["query"].(map[string]interface{})["bool"].(map[string]interface{})["must"] = mustMatch
	}
	fmt.Println(query)

	data, err := json.MarshalIndent(query, "", "   ")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(data))

}
