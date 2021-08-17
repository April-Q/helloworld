package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jmoiron/jsonq"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/printers"
)

type example struct {
	Id        int    `json:"id"`
	CreatedAt string `json:"created_at"`
	Subobj    Subobj `json:"subobj"`
	Array     []int  `json:"array"`
}
type Subobj struct {
	Foo int `json:"foo"`
}

func prettyjson() {

	// c := `{
	// 	"foo": 1,
	// 	"bar": 2,
	// 	"test": "Hello, world!",
	// 	"baz": 123.1,
	// 	"array": [
	// 		{"foo": 1},
	// 		{"bar": 2},
	// 		{"baz": 3}
	// 	],
	// 	"subobj": {
	// 		"foo": 1,
	// 		"subarray": [1,2,3],
	// 		"subsubobj": {
	// 			"bar": 2,
	// 			"baz": 3,
	// 			"array": ["hello", "world"]
	// 		}
	// 	},
	// 	"bool": true
	// 	}`

	aa := example{
		Id:        55555,
		CreatedAt: "4444",
		Subobj: Subobj{
			Foo: 444,
		},
		Array: []int{5, 6, 7},
	}
	bb, err := json.Marshal(aa)
	data := map[string]interface{}{}
	dec := json.NewDecoder(strings.NewReader(string(bb)))
	dec.Decode(&data)

	jq := jsonq.NewQuery(data)
	// data["foo"] -> 1
	obj, err := jq.Interface("id")
	fmt.Println(obj, err)

	path := []string{"subobj", "foo"}
	// data["subobj"]["subarray"][1] -> 2
	obj1, err := jq.Int(path...)
	fmt.Println(obj1)
	obj2, err := jq.Array("array")
	fmt.Println(obj2)

}

func printTable(columns []metav1.TableColumnDefinition, rows []metav1.TableRow) {
	table := &metav1.Table{
		ColumnDefinitions: columns,
		Rows:              rows,
	}
	out := bytes.NewBuffer([]byte{})
	printer := printers.NewTablePrinter(printers.PrintOptions{})
	printer.PrintObj(table, out)
	fmt.Println(out.String())

}

func printable2() {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)

	header := []string{"Name", "Type", "Age", "Des"}
	var tableRow table.Row

	for i := range header {
		tableRow = append(tableRow, header[i])
	}

	t.AppendHeader(tableRow)
	t.AppendRows([]table.Row{
		{"Arya", "Stark", 3000},
		{"Jon", "Snow", 2000, "You know nothing, Jon Snow!"},
	})
	t.AppendSeparator()
	t.AppendRow([]interface{}{"Tyrion", "Lannister", 5000})
	// t.AppendFooter(table.Row{"", "Total", 10000})
	// t.SetStyle(table.StyleColoredBright)
	// t.SetPageSize(1)
	t.Render()

}
