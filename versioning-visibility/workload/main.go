package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"sync"
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

const (
	QueryRPS        = 1
	InsertRPS       = 1
	InsertWorkers   = 1
	InsertBatchSize = 1

	OpenFilterRatio         = .6
	BuildIdNonexistentRatio = .2

	TQFilterRatio  = .3
	TQExcludeRatio = .5

	NumBuildIdsInFilter = 10
	NumTQsInFilter      = 4

	NumShards = 20
)

var workloads = []NamespaceConfig{
	//{
	//	"C6139bd23-bc97-4e49-943d-3f05ac1d3e5f",
	//	168,
	//	166,
	//	4,
	//	10,
	//},
	//{
	//	"B6139bd23-bc97-4e49-943d-3f05ac1d3e5f",
	//	168,
	//	167,
	//	4,
	//	5,
	//},
	//{
	//	"7139bd23-bc97-4e49-943d-3f05ac1d3e5f",
	//	2,
	//	1,
	//	4,
	//	5,
	//},
	{
		"6139bd23-bc97-4e49-943d-3f05ac1d3e5f",
		720,
		719,
		4,
		5,
	},
}

const terminateAfter = "terminate_after=1&"

var insertCounter = 0

var qCounter = 0
var qTime int64 = 0

var sqCounter = 0
var sqTime int64 = 0

var queryStats = make(map[string][]int64)
var mutex sync.Mutex

func main() {
	reportRPS()

	for w := 0; w < InsertWorkers; w++ {
		go cmn.Work(InsertBatchSize)
	}
	go insertDocs()

	runQueries()
}

func insertDocs() {
	ticker := time.NewTicker(time.Second / InsertRPS)
	for range ticker.C {
		ns := getNamespace()
		bld := cmn.BuildId(rand.Intn(ns.maxBuildId))
		tq := cmn.TaskQueue(rand.Intn(ns.maxTq))
		cmn.PutDoc(ns.id, bld, cmn.StatusCompleted, tq, time.Now().UnixNano())
		insertCounter++
	}
}

func runQueries() {
	ticker := time.NewTicker(time.Second / QueryRPS)
	for range ticker.C {
		go runOneQuery()
	}
}

func reportRPS() {
	ticker := time.NewTicker(time.Second)
	go func() {
		var sec = 1
		var lastCount = 0
		var lastSCount = 0
		var lastICount = 0
		for range ticker.C {
			fullQAvg := int64(-1)
			if qCounter > 0 {
				fullQAvg = qTime / int64(qCounter)
			}
			shardQAvg := int64(-1)
			if sqCounter > 0 {
				shardQAvg = sqTime / int64(sqCounter)
			}
			fmt.Printf(
				"FullQ RPS: %d\t Full Q Avg microSec: %d\t ShardQ RPS: %d\t ShardQ Avg microSec: %d\t Doc Insert RPS: %d \n",
				qCounter-lastCount,
				fullQAvg,
				sqCounter-lastSCount,
				shardQAvg,
				insertCounter-lastICount,
			)
			lastCount = qCounter
			lastSCount = sqCounter
			lastICount = insertCounter

			if sec%60 == 0 {
				reportStats()
			}

			sec++
		}
	}()
}

func getNamespace() NamespaceConfig {
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

	return ns
}

func runOneQuery() {
	ns := getNamespace()

	tqs := ""
	excludeTQs := rand.Float32() < TQExcludeRatio
	if rand.Float32() < TQFilterRatio {
		tqs = getTQs(ns)
	}

	nonexistent := rand.Float32() < BuildIdNonexistentRatio
	openOnly := rand.Float32() < OpenFilterRatio
	buildIds := getBuildIds(ns, openOnly, nonexistent)

	query(ns, buildIds, tqs, openOnly, nonexistent, excludeTQs)
}

func getBuildIds(ns NamespaceConfig, openOnly bool, nonexistent bool) string {
	res := []string{}

	for i := 0; i < NumBuildIdsInFilter; i++ {
		if openOnly {
			build := rand.Intn(ns.maxBuildId-ns.openMinBuildId+1) + ns.openMinBuildId
			if nonexistent {
				build = rand.Intn(ns.openMinBuildId) - 1
			}
			res = append(res, "versioned:"+cmn.BuildId(build))
		} else {
			build := rand.Intn(ns.maxBuildId + 1)
			if nonexistent {
				build += ns.maxBuildId + 1
			}
			res = append(res, "versioned:"+cmn.BuildId(build))
		}

	}

	return `"` + strings.Join(res, `", "`) + `"`
}

func getTQs(ns NamespaceConfig) string {
	res := []string{}

	for i := 0; i < NumTQsInFilter; i++ {
		res = append(res, cmn.TaskQueue(rand.Intn(ns.maxTq)))
	}

	return `"` + strings.Join(res, `", "`) + `"`
}

func query(ns NamespaceConfig, buildIds string, taskQueues string, openOnly bool, nonexistent bool, excludeTqs bool) {
	running := ""
	if openOnly {
		running = `{ "term": { "ExecutionStatus": "Running" }},`
	}

	tqFilter := ""
	excludeTq := ""
	if taskQueues != "" {
		if excludeTqs {
			excludeTq = `"must_not" : { "terms" : { "TaskQueue" : [` + taskQueues + `] }},`
		} else {
			tqFilter = `{ "terms" : { "TaskQueue" : [` + taskQueues + `] }},`
		}
	}

	query := []byte(`{
		"query" : {
        "bool" : {
			` + excludeTq + `
            "filter": [
			  ` + running + `
 			  ` + tqFilter + `
              { "term": { "NamespaceId": "` + ns.id + `" }},
              { "terms" : { "BuildIds" : [` + buildIds + `] }}
            ]
        }
    }
        }`)

	name := ""

	if openOnly {
		name += "-Open"
	} else {
		name += "-All"
	}
	if taskQueues != "" {
		if excludeTqs {
			name += "-ExcludeTQs"
		} else {
			name += "-IncludeTQs"
		}
	} else {
		name += "-AllTQs"
	}
	if nonexistent {
		name += "-Empty"
	} else {
		name += "-NonEmpty"
	}

	if NumShards > -1 {
		shard := rand.Intn(NumShards)
		count, err := runQuery(query, "Shard-"+name, false, nonexistent, shard)
		if err == nil && count > 0 {
			return
		}
	}

	runQuery(query, "Full-"+name, false, nonexistent, -1)
}

func runQuery(q []byte, name string, log bool, expectingZero bool, shard int) (int, error) {
	url := cmn.BaseUrl + cmn.Index + "/_count?" + terminateAfter

	if shard > -1 {
		url += "preference=_shards:" + strconv.Itoa(shard) + "|_only_local"
	}

	req, err := http.NewRequest("GET", url, bytes.NewBuffer(q))
	req.Header.Set("Content-Type", "application/json")

	start := time.Now()
	if cmn.Password != "" {
		req.SetBasicAuth("temporal", cmn.Password)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("error while querying ES", err)
		return -1, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode > 299 {
		if resp.StatusCode == 429 {
			fmt.Println("Was throttled")
		} else {
			fmt.Println("####### QUERY:")
			fmt.Println(string(q))
			fmt.Println("response Status:", resp.Status)
			fmt.Println("Body:", body)
		}
		return -1, fmt.Errorf(resp.Status)
	}

	timer := int64(time.Since(start) / time.Microsecond)
	if shard > -1 {
		sqCounter++
		sqTime += timer
	} else {
		qCounter++
		qTime += timer
	}

	mutex.Lock()
	s := queryStats[name]
	if s == nil {
		s = []int64{timer}
		queryStats[name] = s
	} else {
		queryStats[name] = append(s, timer)
	}
	mutex.Unlock()

	var data map[string]interface{}
	json.Unmarshal(body, &data)

	count := int(data["count"].(float64))

	countOK := shard > -1 || (count > 0 && !expectingZero) || (expectingZero && count == 0)

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

	return count, nil
}

func reportStats() {
	mutex.Lock()
	stats := queryStats
	queryStats = make(map[string][]int64)
	mutex.Unlock()
	fmt.Printf("\n| %50s | %8s | %8s | %8s | %8s |\n", "Query", "Rep", "Avg", "Min", "Max")

	keys := make([]string, 0, len(stats))
	for k := range stats {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		count, avg, min, max := summarizeQuery(stats[k])
		fmt.Printf("| %50s | %8d | %8d | %8d | %8d |\n", k, count, avg, min, max)
	}
}

func summarizeQuery(v []int64) (int64, int64, int64, int64) {
	var sumTimer int64
	var minTimer = int64(999999999)
	var maxTimer = int64(-1)
	for _, p := range v {
		sumTimer += p
		minTimer = min(minTimer, p)
		maxTimer = max(maxTimer, p)
	}

	count := int64(len(v))

	return count, sumTimer / count, minTimer, maxTimer
}
