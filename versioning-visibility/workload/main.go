package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	cmn "temporaltests/versioning-visibility"
	"time"
)

type NamespaceConfig struct {
	id             string
	maxBuildId     int
	openMinBuildId int
	maxTq          int
	loadShare      int
}

var workloads = []NamespaceConfig{
	{
		"C6139bd23-bc97-4e49-943d-3f05ac1d3e5f",
		168,
		166,
		4,
		10,
	},
	{
		"B6139bd23-bc97-4e49-943d-3f05ac1d3e5f",
		168,
		167,
		4,
		5,
	},
	{
		"7139bd23-bc97-4e49-943d-3f05ac1d3e5f",
		2,
		1,
		4,
		5,
	},
	//{
	//	"6139bd23-bc97-4e49-943d-3f05ac1d3e5f",
	//	720,
	//	719,
	//	4,
	//	5,
	//},
}

const (
	RPS = 100

	OpenFilterRatio         = .6
	TQFilterRatio           = .1
	BuildIdNonexistentRatio = .2
)

const terminateAfter = "?terminate_after=1"

var qCounter = 0

func main() {
	reportRPS()
	ticker := time.NewTicker(time.Second / RPS)
	for range ticker.C {
		go runOneQuery()
		qCounter++
	}
}

func reportRPS() {
	ticker := time.NewTicker(1 * time.Second)
	go func() {
		var lastCount = 0
		for range ticker.C {
			fmt.Printf("RPS: %d \n", qCounter-lastCount)
			lastCount = qCounter
		}
	}()
}

func runOneQuery() {
	loads := 0
	for _, wl := range workloads {
		loads += wl.loadShare
	}

	r := rand.Intn(loads)
	var ns = workloads[0]

	loads = 0
	for _, wl := range workloads {
		loads += wl.loadShare
		if loads > r {
			ns = wl
			break
		}
	}

	tq := -1
	if rand.Float32() < TQFilterRatio {
		tq = rand.Intn(ns.maxTq)
	}

	nonexistent := rand.Float32() < BuildIdNonexistentRatio
	if rand.Float32() < OpenFilterRatio {
		build := rand.Intn(ns.maxBuildId-ns.openMinBuildId+1) + ns.openMinBuildId
		if nonexistent {
			build = rand.Intn(ns.openMinBuildId)
		}
		query(ns, cmn.BuildId(build), tq, true, nonexistent)
	} else {
		build := rand.Intn(ns.maxBuildId + 1)
		if nonexistent {
			build += ns.maxBuildId
		}
		query(ns, cmn.BuildId(build), tq, false, nonexistent)
	}

}

func query(ns NamespaceConfig, buildId string, taskQueue int, openOnly bool, nonexistent bool) {
	running := ""
	if openOnly {
		running = `{ "term": { "ExecutionStatus": "Running" }},`
	}

	tq := ""
	if taskQueue > -1 {
		tq = `{ "term" : { "TaskQueue" : "` + cmn.TaskQueue(taskQueue) + `" }},`
	}

	runQuery([]byte(`{
		"query" : {
        "bool" : {
            "filter": [
			  `+running+`
 			  `+tq+`
              { "term": { "NamespaceId": "`+ns.id+`" }},
              { "term" : { "BuildIds" : "versioned:`+buildId+`" }}
            ]
        }
    }
        }`),
		"count",
		false,
		nonexistent,
	)
}

func runQuery(q []byte, qtype string, log bool, expectingZero bool) {
	if qtype == "count" {
		qtype += terminateAfter
	} else if qtype == "search" {
		qtype += "?track_total_hits=false"
	}
	url := cmn.BaseUrl + cmn.Index + "/_" + qtype
	req, err := http.NewRequest("GET", url, bytes.NewBuffer(q))
	req.Header.Set("Content-Type", "application/json")

	if cmn.Password != "" {
		req.SetBasicAuth("temporal", cmn.Password)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode > 299 {
		fmt.Println("####### QUERY:")
		fmt.Println(string(q))
		fmt.Println("response Status:", resp.Status)
		fmt.Println("Body:", body)
		return
	}

	var data map[string]interface{}
	json.Unmarshal(body, &data)

	count := data["count"].(float64)
	countOK := (count > 0 && !expectingZero) || (expectingZero && count == 0)

	if !countOK {
		fmt.Println("Unexpected Count, expectedZero=", expectingZero)
	}

	if log || !countOK {
		fmt.Println("####### QUERY:")
		fmt.Println(string(q))
		fmt.Println("response Status:", resp.Status)

		pretty, _ := json.MarshalIndent(data, "", "    ")
		fmt.Println("response Body:", string(pretty))
	}
}
