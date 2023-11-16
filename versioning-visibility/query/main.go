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
	cmn "temporaltests/versioning-visibility"
	"time"
)

const repetition = 10

// const reachabilityInterval = time.Minute
// const reachabilityFreq = int(cmn.WfLen/reachabilityInterval) + 1
// const lastStartInterval = 2 * time.Minute
// const lastStartFreq = int(lastStartInterval/reachabilityInterval) + 1
const reachabilityFreq = 1
const lastStartFreq = 1
const groupBySize = 1000
const maxBuildId = int(int64(cmn.Retention-time.Second) / int64(cmn.BuildLen))
const scavengerInterval = 12 * time.Hour

type statsPoint struct {
	timer int64
	took  int64
}

var queryStats = make(map[string][]statsPoint)
var strategyStats = make(map[string][]int64)

var tests = []func(){
	//ReachabilityCountOneBuildPerTqOpen,
	//ReachabilityCountOneBuildPerTq,
	ReachabilityCountOneBuildAllTqOpen,
	ReachabilityCountOneBuildAllTq,
	ReachabilityCountOneBuildAllTqCurrent,
	ReachabilityCountOneBuildAllTqCurrentNS,
	ReachabilityCountOneBuildAllTqCurrentField,
	ReachabilityCountOneBuildAllTqCurrentNSField,
	//ReachabilityLimitOneBuildPerTqOpen,
	//ReachabilityLimitOneBuildAllTqOpen,
	//ReachabilityLimitOneBuildPerTq,
	//ReachabilityLimitOneBuildAllTq,
	ReachabilityGroupByBuild,
	//ReachabilityGroupByBuildTq,
	ReachabilityGroupByBuildOpen,
	ReachabilityGroupByBuildCurrent,
	ReachabilityGroupByBuildCurrentNS,
	ReachabilityGroupByBuildCurrentField,
	ReachabilityGroupByBuildCurrentNSField,
	//ReachabilityGroupByBuildTqOpen,
	//ScavengerCountOneBuildPerTqOpen,
	//ScavengerCountOneBuildAllTqOpen,
	//ScavengerCountOneBuildPerTq,
	//ScavengerCountOneBuildAllTq,
	//ScavengerLimitOneBuildPerTqOpen,
	//ScavengerLimitOneBuildAllTqOpen,
	//ScavengerLimitOneBuildPerTq,
	//ScavengerLimitOneBuildAllTq,
	//ScavengerGroupByBuild,
	//ScavengerGroupByBuildTq,
	//ScavengerGroupByBuildOpen,
	//ScavengerGroupByBuildTqOpen,
	//ReachabilityLastStartOneBuildAllTq,
	//ReachabilityLastStartOneBuildPerTq,
	//ReachabilityGroupByTqOneBuild,
	//ReachabilityGroupByTqOneBuildOpen,
}

func main() {
	for i := 0; i < repetition; i++ {
		order := shuffleOrder(len(tests))
		for _, o := range order {
			tests[o]()
		}
	}

	reportStats()
}

func shuffleOrder(items int) []int {
	a := []int{}
	for i := 0; i < items; i++ {
		a = append(a, i)
	}
	rand.Shuffle(len(a), func(i, j int) { a[i], a[j] = a[j], a[i] })
	fmt.Println(a)
	return a
}

func ReachabilityCountOneBuildPerTqOpen() {
	start := time.Now()
	bld := maxBuildId - 1
	for i := 0; i < reachabilityFreq; i++ {
		for tq := 0; tq < cmn.NTaskQueues; tq++ {
			countOneBuildOneTqOpen(cmn.BuildId(bld), cmn.TaskQueue(tq))
		}
	}
	recordStrategyStats("ReachabilityCountOneBuildPerTqOpen", start)
}

func ReachabilityCountOneBuildAllTqOpen() {
	start := time.Now()
	bld := maxBuildId - 1
	for i := 0; i < reachabilityFreq; i++ {
		countOneBuildOpen(cmn.BuildId(bld))
	}
	recordStrategyStats("ReachabilityCountOneBuildAllTqOpen", start)
}

func ReachabilityCountOneBuildAllTqCurrent() {
	start := time.Now()
	bld := maxBuildId - 1
	for i := 0; i < reachabilityFreq; i++ {
		countOneBuildCurrent(cmn.BuildId(bld))
	}
	recordStrategyStats("ReachabilityCountOneBuildAllTqCurrent", start)
}

func ReachabilityCountOneBuildAllTqCurrentNS() {
	start := time.Now()
	bld := maxBuildId - 1
	for i := 0; i < reachabilityFreq; i++ {
		countOneBuildCurrentNS(cmn.BuildId(bld))
	}
	recordStrategyStats("ReachabilityCountOneBuildAllTqCurrentNS", start)
}

func ReachabilityCountOneBuildAllTqCurrentField() {
	start := time.Now()
	bld := maxBuildId - 1
	for i := 0; i < reachabilityFreq; i++ {
		countOneBuildCurrentField(cmn.BuildId(bld))
	}
	recordStrategyStats("ReachabilityCountOneBuildAllTqCurrentField", start)
}

func ReachabilityCountOneBuildAllTqCurrentNSField() {
	start := time.Now()
	bld := maxBuildId - 1
	for i := 0; i < reachabilityFreq; i++ {
		countOneBuildCurrentNSField(cmn.BuildId(bld))
	}
	recordStrategyStats("ReachabilityCountOneBuildAllTqCurrentNSField", start)
}

func ReachabilityLastStartOneBuildPerTq() {
	start := time.Now()
	bld := maxBuildId - 1
	for i := 0; i < lastStartFreq; i++ {
		for tq := 0; tq < cmn.NTaskQueues; tq++ {
			lastStartOneBuildOneTq(cmn.BuildId(bld), cmn.TaskQueue(tq))
		}
	}
	recordStrategyStats("ReachabilityLastStartOneBuildPerTq", start)
}

func ReachabilityLastStartOneBuildAllTq() {
	start := time.Now()
	bld := maxBuildId - 1
	for i := 0; i < lastStartFreq; i++ {
		lastStartOneBuild(cmn.BuildId(bld))
	}
	recordStrategyStats("ReachabilityLastStartOneBuildAllTq", start)
}

func ScavengerLimitOneBuildPerTqOpen() {
	start := time.Now()
	minBld := maxBuildId - int(scavengerInterval/cmn.BuildLen) - 1
	for b := minBld; b <= maxBuildId; b++ {
		for tq := 0; tq < cmn.NTaskQueues; tq++ {
			limitOneBuildOneTqOpen(cmn.BuildId(b), cmn.TaskQueue(tq))
		}
	}
	recordStrategyStats("ScavengerLimitOneBuildPerTqOpen", start)
}

func ScavengerLimitOneBuildAllTqOpen() {
	start := time.Now()
	minBld := maxBuildId - int(scavengerInterval/cmn.BuildLen) - 1
	for b := minBld; b <= maxBuildId; b++ {
		countOneBuildOpen(cmn.BuildId(b))
	}
	recordStrategyStats("ScavengerLimitOneBuildAllTqOpen", start)
}

func ScavengerLimitOneBuildPerTq() {
	start := time.Now()
	minBld := maxBuildId - int(scavengerInterval/cmn.BuildLen) - 1
	for b := minBld; b <= maxBuildId; b++ {
		for tq := 0; tq < cmn.NTaskQueues; tq++ {
			limitOneBuildOneTq(cmn.BuildId(b), cmn.TaskQueue(tq))
		}
	}
	recordStrategyStats("ScavengerLimitOneBuildPerTq", start)
}

func ScavengerLimitOneBuildAllTq() {
	start := time.Now()
	minBld := maxBuildId - int(scavengerInterval/cmn.BuildLen) - 1
	for b := minBld; b <= maxBuildId; b++ {
		countOneBuild(cmn.BuildId(b))
	}
	recordStrategyStats("ScavengerLimitOneBuildAllTq", start)
}

func ScavengerCountOneBuildPerTqOpen() {
	start := time.Now()
	minBld := maxBuildId - int(scavengerInterval/cmn.BuildLen) - 1
	for b := minBld; b <= maxBuildId; b++ {
		for tq := 0; tq < cmn.NTaskQueues; tq++ {
			countOneBuildOneTqOpen(cmn.BuildId(b), cmn.TaskQueue(tq))
		}
	}
	recordStrategyStats("ScavengerCountOneBuildPerTqOpen", start)
}

func ScavengerCountOneBuildAllTqOpen() {
	start := time.Now()
	minBld := maxBuildId - int(scavengerInterval/cmn.BuildLen) - 1
	for b := minBld; b <= maxBuildId; b++ {
		countOneBuildOpen(cmn.BuildId(b))
	}
	recordStrategyStats("ScavengerCountOneBuildAllTqOpen", start)
}

func ScavengerCountOneBuildPerTq() {
	start := time.Now()
	maxBld := int(scavengerInterval/cmn.BuildLen) + 1
	for b := 0; b <= maxBld; b++ {
		for tq := 0; tq < cmn.NTaskQueues; tq++ {
			countOneBuildOneTq(cmn.BuildId(b), cmn.TaskQueue(tq))
		}
	}
	recordStrategyStats("ScavengerCountOneBuildPerTq", start)
}

func ScavengerCountOneBuildAllTq() {
	start := time.Now()
	maxBld := int(scavengerInterval/cmn.BuildLen) + 1
	for b := 0; b <= maxBld; b++ {
		countOneBuild(cmn.BuildId(b))
	}
	recordStrategyStats("ScavengerCountOneBuildAllTq", start)
}

func ReachabilityCountOneBuildPerTq() {
	start := time.Now()
	bld := maxBuildId - 1
	for i := 0; i < reachabilityFreq; i++ {
		for tq := 0; tq < cmn.NTaskQueues; tq++ {
			countOneBuildOneTq(cmn.BuildId(bld), cmn.TaskQueue(tq))
		}
	}
	recordStrategyStats("ReachabilityCountOneBuildPerTq", start)
}

func ReachabilityCountOneBuildAllTq() {
	start := time.Now()
	bld := maxBuildId - 1
	for i := 0; i < reachabilityFreq; i++ {
		countOneBuild(cmn.BuildId(bld))
	}
	recordStrategyStats("ReachabilityCountOneBuildAllTq", start)
}

func ReachabilityLimitOneBuildPerTq() {
	start := time.Now()
	bld := maxBuildId - 1
	for i := 0; i < reachabilityFreq; i++ {
		for tq := 0; tq < cmn.NTaskQueues; tq++ {
			limitOneBuildOneTq(cmn.BuildId(bld), cmn.TaskQueue(tq))
		}
	}
	recordStrategyStats("ReachabilityLimitOneBuildPerTq", start)
}

func ReachabilityLimitOneBuildAllTq() {
	start := time.Now()
	bld := maxBuildId - 1
	for i := 0; i < reachabilityFreq; i++ {
		limitOneBuild(cmn.BuildId(bld))
	}
	recordStrategyStats("ReachabilityLimitOneBuildAllTq", start)
}

func ReachabilityLimitOneBuildPerTqOpen() {
	start := time.Now()
	bld := maxBuildId - 1
	for i := 0; i < reachabilityFreq; i++ {
		for tq := 0; tq < cmn.NTaskQueues; tq++ {
			limitOneBuildOneTqOpen(cmn.BuildId(bld), cmn.TaskQueue(tq))
		}
	}
	recordStrategyStats("ReachabilityLimitOneBuildPerTqOpen", start)
}

func ReachabilityLimitOneBuildAllTqOpen() {
	start := time.Now()
	bld := maxBuildId - 1
	for i := 0; i < reachabilityFreq; i++ {
		limitOneBuildOpen(cmn.BuildId(bld))
	}
	recordStrategyStats("ReachabilityLimitOneBuildAllTqOpen", start)
}

func ReachabilityGroupByBuildTq() {
	start := time.Now()
	for i := 0; i < reachabilityFreq; i++ {
		for bld := 0; bld <= maxBuildId; bld += groupBySize / 20 {
			filter := []string{}
			for b := bld; b < bld+groupBySize/20; b++ {
				filter = append(filter, cmn.BuildId(b))
			}
			groupByBuildTqFiltered(filter)
		}
	}
	recordStrategyStats("ReachabilityGroupByBuildTq", start)
}

func ReachabilityGroupByBuild() {
	start := time.Now()
	for i := 0; i < reachabilityFreq; i++ {
		groupByBuild()
		//for bld := 0; bld <= maxBuildId; bld += groupBySize {
		//	filter := []string{}
		//	for b := bld; b < bld+groupBySize; b++ {
		//		filter = append(filter, cmn.BuildId(b))
		//	}
		//	groupByBuildFiltered(filter)
		//}
	}
	recordStrategyStats("ReachabilityGroupByBuild", start)
}

func ReachabilityGroupByTqOneBuild() {
	start := time.Now()
	bld := maxBuildId - 1
	for i := 0; i < reachabilityFreq; i++ {
		groupByTqOneBuild(cmn.BuildId(bld))
	}
	recordStrategyStats("ReachabilityGroupByTqOneBuild", start)
}

func ReachabilityGroupByTqOneBuildOpen() {
	start := time.Now()
	bld := maxBuildId - 1
	for i := 0; i < reachabilityFreq; i++ {
		groupByTqOneBuildOpen(cmn.BuildId(bld))
	}
	recordStrategyStats("ReachabilityGroupByTqOneBuildOpen", start)
}

func ReachabilityGroupByBuildTqOpen() {
	start := time.Now()
	for i := 0; i < reachabilityFreq; i++ {
		groupByBuildTqOpen()
	}
	recordStrategyStats("ReachabilityGroupByBuildTqOpen", start)
}

func ReachabilityGroupByBuildOpen() {
	start := time.Now()
	for i := 0; i < reachabilityFreq; i++ {
		groupByBuildOpen()
	}
	recordStrategyStats("ReachabilityGroupByBuildOpen", start)
}

func ReachabilityGroupByBuildCurrent() {
	start := time.Now()
	for i := 0; i < reachabilityFreq; i++ {
		groupByBuildCurrent()
	}
	recordStrategyStats("ReachabilityGroupByBuildCurrent", start)
}

func ReachabilityGroupByBuildCurrentNS() {
	start := time.Now()
	for i := 0; i < reachabilityFreq; i++ {
		groupByBuildCurrentNS()
	}
	recordStrategyStats("ReachabilityGroupByBuildCurrentNS", start)
}

func ReachabilityGroupByBuildCurrentField() {
	start := time.Now()
	for i := 0; i < reachabilityFreq; i++ {
		groupByBuildCurrentField()
	}
	recordStrategyStats("ReachabilityGroupByBuildCurrentField", start)
}

func ReachabilityGroupByBuildCurrentNSField() {
	start := time.Now()
	for i := 0; i < reachabilityFreq; i++ {
		groupByBuildCurrentNSField()
	}
	recordStrategyStats("ReachabilityGroupByBuildCurrentNSField", start)
}

func ScavengerGroupByBuildTq() {
	start := time.Now()
	maxBld := int(scavengerInterval/cmn.BuildLen) + 1
	for bld := 0; bld <= maxBld; bld += groupBySize / 20 {
		filter := []string{}
		for b := bld; b < bld+groupBySize/20; b++ {
			filter = append(filter, cmn.BuildId(b))
		}
		groupByBuildTqFiltered(filter)
	}
	recordStrategyStats("ScavengerGroupByBuildTq", start)
}

func ScavengerGroupByBuild() {
	start := time.Now()
	maxBld := int(scavengerInterval/cmn.BuildLen) + 1
	for bld := 0; bld <= maxBld; bld += groupBySize {
		filter := []string{}
		for b := bld; b < bld+groupBySize; b++ {
			filter = append(filter, cmn.BuildId(b))
		}
		groupByBuildFiltered(filter)
	}
	recordStrategyStats("ScavengerGroupByBuild", start)
}

func ScavengerGroupByBuildTqOpen() {
	start := time.Now()
	groupByBuildTqOpen()
	recordStrategyStats("ScavengerGroupByBuildTqOpen", start)
}

func ScavengerGroupByBuildOpen() {
	start := time.Now()
	groupByBuildOpen()
	recordStrategyStats("ScavengerGroupByBuildOpen", start)
}

func reportStats() {
	fmt.Printf("\n| %30s | %8s | %8s | %8s | %8s |\n", "Query", "Rep", "Avg", "Min", "Max")

	keys := make([]string, 0, len(queryStats))
	for k := range queryStats {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		count, avg, min, max, _, _, _ := summarizeQuery(queryStats[k])
		fmt.Printf("| %30s | %8d | %8d | %8d | %8d |\n", k, count, avg, min, max)
	}

	fmt.Printf("\n| %50s | %8s | %8s | %8s | %8s |\n", "Strategy", "Rep", "Avg", "Min", "Max")

	keys = make([]string, 0, len(strategyStats))
	for k := range strategyStats {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		count, avg, min, max := summarizeStrategy(strategyStats[k])
		fmt.Printf("| %50s | %8d | %8d | %8d | %8d |\n", k, count, avg, min, max)
	}
}

func summarizeQuery(v []statsPoint) (int64, int64, int64, int64, int64, int64, int64) {
	var sumTimer int64
	var minTimer = int64(9999999)
	var maxTimer = int64(-1)
	for _, p := range v {
		sumTimer += p.timer
		minTimer = min(minTimer, p.timer)
		maxTimer = max(maxTimer, p.timer)
	}

	var sumTook int64
	var minTook = int64(9999999)
	var maxTook = int64(-1)
	for _, p := range v {
		sumTook += p.took
		minTook = min(minTook, p.took)
		maxTook = max(maxTook, p.took)
	}

	count := int64(len(v))

	return count, sumTimer / count, minTimer, maxTimer, sumTook / count, minTook, maxTook
}

func summarizeStrategy(v []int64) (int64, int64, int64, int64) {
	var sumTimer int64
	var minTimer = int64(9999999)
	var maxTimer = int64(-1)
	for _, p := range v {
		sumTimer += p
		minTimer = min(minTimer, p)
		maxTimer = max(maxTimer, p)
	}

	count := int64(len(v))

	return count, sumTimer / count, minTimer, maxTimer
}

func countOneBuildOneTq(build string, tq string) {
	runQuery([]byte(`{
		"query" : {
        "bool" : {
            "filter": [
              { "term": { "NamespaceId": "`+cmn.NamespaceId+`" }},
              { "term" : { "TaskQueue" : "`+tq+`" }},
              { "term" : { "BuildIds" : "versioned:`+build+`" }}
            ]
        }
    }
        }`),
		"count",
		"countOneBuildOneTq",
		false,
	)
}

func countOneBuild(build string) {
	runQuery([]byte(`{
		"query" : {
        "bool" : {
            "filter": [
              { "term": { "NamespaceId": "`+cmn.NamespaceId+`" }},
              { "term" : { "BuildIds" : "versioned:`+build+`" }}
            ]
        }
    }
        }`),
		"count",
		"countOneBuild",
		false,
	)
}

func countOneBuildOneTqOpen(build string, tq string) {
	runQuery([]byte(`{
		"query" : {
        "bool" : {
            "filter": [
              { "term": { "NamespaceId": "`+cmn.NamespaceId+`" }},
              { "term" : { "TaskQueue" : "`+tq+`" }},
              { "term" : { "BuildIds" : "versioned:`+build+`" }},
			  { "term": { "ExecutionStatus": "Running" }}
            ]
        }
    }
        }`),
		"count",
		"countOneBuildOneTqOpen",
		false,
	)
}

func countOneBuildOpen(build string) {
	runQuery([]byte(`{
		"query" : {
        "bool" : {
            "filter": [
              { "term": { "NamespaceId": "`+cmn.NamespaceId+`" }},
              { "term" : { "BuildIds" : "versioned:`+build+`" }},
			  { "term": { "ExecutionStatus": "Running" }}
            ]
        }
    }
        }`),
		"count",
		"countOneBuildOpen",
		false,
	)
}

func countOneBuildCurrent(build string) {
	runQuery([]byte(`{
		"query" : {
        "bool" : {
            "filter": [
              { "term": { "NamespaceId": "`+cmn.NamespaceId+`" }},
              { "term" : { "BuildIds" : "versioned:`+build+`" }}
            ]
        }
    }
        }`),
		"count",
		"countOneBuildCurrent",
		false,
	)
}

func countOneBuildCurrentField(build string) {
	runQuery([]byte(`{
		"query" : {
        "bool" : {
            "filter": [
              { "term": { "NamespaceId": "`+cmn.NamespaceId+`" }},
              { "term" : { "CurrentBuildId" : "versioned:`+build+`" }}
            ]
        }
    }
        }`),
		"count",
		"countOneBuildCurrentField",
		false,
	)
}

func countOneBuildCurrentNSField(build string) {
	runQuery([]byte(`{
		"query" : {
        "bool" : {
            "filter": [
              { "term" : { "CurrentBuildId" : "versioned:`+build+`" }}
            ]
        }
    }
        }`),
		"count",
		"countOneBuildCurrentNSField",
		false,
	)
}

func countOneBuildCurrentNS(build string) {
	runQuery([]byte(`{
		"query" : {
        "bool" : {
            "filter": [
              { "term" : { "BuildIds" : "versioned:`+build+`" }}
            ]
        }
    }
        }`),
		"count",
		"countOneBuildCurrentNS",
		false,
	)
}

func limitOneBuildOneTq(build string, tq string) {
	runQuery([]byte(`{
		"size": 1,
		"query" : {
        "bool" : {
            "filter": [
              { "term": { "NamespaceId": "`+cmn.NamespaceId+`" }},
              { "term" : { "TaskQueue" : "`+tq+`" }},
              { "term" : { "BuildIds" : "versioned:`+build+`" }}
            ]
        }
    }
        }`),
		"search",
		"limitOneBuildOneTq",
		false,
	)
}

func limitOneBuild(build string) {
	runQuery([]byte(`{
		"size": 1,
		"query" : {
        "bool" : {
            "filter": [
              { "term": { "NamespaceId": "`+cmn.NamespaceId+`" }},
              { "term" : { "BuildIds" : "versioned:`+build+`" }}
            ]
        }
    }
        }`),
		"search",
		"limitOneBuild",
		false,
	)
}

func limitOneBuildOneTqOpen(build string, tq string) {
	runQuery([]byte(`{
		"size": 1,
		"query" : {
        "bool" : {
            "filter": [
              { "term": { "NamespaceId": "`+cmn.NamespaceId+`" }},
              { "term" : { "TaskQueue" : "`+tq+`" }},
              { "term" : { "BuildIds" : "versioned:`+build+`" }},
			  { "term": { "ExecutionStatus": "Running" }}
            ]
        }
    }
        }`),
		"search",
		"limitOneBuildOneTqOpen",
		false,
	)
}

func limitOneBuildOpen(build string) {
	runQuery([]byte(`{
		"size": 1,
		"query" : {
        "bool" : {
            "filter": [
              { "term": { "NamespaceId": "`+cmn.NamespaceId+`" }},
              { "term" : { "BuildIds" : "versioned:`+build+`" }},
			  { "term": { "ExecutionStatus": "Running" }}
            ]
        }
    }
        }`),
		"search",
		"limitOneBuildOneTqOpen",
		false,
	)
}

func groupByBuild() {
	runQuery([]byte(`{
		"size": 0,
		"query" : {
        	"bool" : {
            	"filter": [
              		{ "term": { "NamespaceId": "`+cmn.NamespaceId+`" }}
            	]
        	}
    	},
		"aggs": {
			"group_by_BuildIds": {
          		"terms": {
          			"size": `+strconv.Itoa(groupBySize)+`,
            		"field": "BuildIds"
          		}
			}
		}
		}`),
		"search",
		"groupByBuild",
		false,
	)
}

func groupByBuildTq() {
	runQuery([]byte(`{
		"size": 0,
		"query" : {
        	"bool" : {
            	"filter": [
              		{ "term": { "NamespaceId": "`+cmn.NamespaceId+`" }}
            	]
        	}
    	},
		"aggs": {
			"group_by_BuildIdsTQ": {
				"multi_terms": {
          			"size": `+strconv.Itoa(groupBySize)+`,
					"terms": [{
					  "field": "BuildIds" 
					}, {
					  "field": "TaskQueue"
					}]
				  }
			}
		}
		}`),
		"search",
		"groupByBuildTq",
		false,
	)
}

func groupByTqOneBuild(build string) {
	runQuery([]byte(`{
		"size": 0,
		"query" : {
        	"bool" : {
            	"filter": [
              		{ "term": { "NamespaceId": "`+cmn.NamespaceId+`" }},
              		{ "term" : { "BuildIds" : "versioned:`+build+`" }}
            	]
        	}
    	},
		"aggs": {
			"group_by": {
				"terms": {
          			"size": `+strconv.Itoa(groupBySize)+`,
            		"field": "TaskQueue"
          		}
			}
		}
		}`),
		"search",
		"groupByTqOneBuild",
		false,
	)
}

func groupByTqOneBuildOpen(build string) {
	runQuery([]byte(`{
		"size": 0,
		"query" : {
        	"bool" : {
            	"filter": [
              		{ "term": { "NamespaceId": "`+cmn.NamespaceId+`" }},
              		{ "term" : { "BuildIds" : "versioned:`+build+`" }},
			  		{ "term": { "ExecutionStatus": "Running" }}
            	]
        	}
    	},
		"aggs": {
			"group_by": {
				"terms": {
          			"size": `+strconv.Itoa(groupBySize)+`,
            		"field": "TaskQueue"
          		}
			}
		}
		}`),
		"search",
		"groupByTqOneBuildOpen",
		false,
	)
}

func lastStartOneBuild(build string) {
	runQuery([]byte(`{
		"size": 0,
		"query" : {
        	"bool" : {
            	"filter": [
              		{ "term": { "NamespaceId": "`+cmn.NamespaceId+`" }},
              		{ "term" : { "BuildIds" : "versioned:`+build+`" }}
            	]
        	}
    	},
		"aggs": {
			"maxStart": {
			  "max": {
					"field": "StartTime"
				  }
			}
		}
		}`),
		"search",
		"lastStartOneBuild",
		false,
	)
}

func lastStartOneBuildOneTq(build string, tq string) {
	runQuery([]byte(`{
		"size": 0,
		"query" : {
        	"bool" : {
            	"filter": [
              		{ "term": { "NamespaceId": "`+cmn.NamespaceId+`" }},
              		{ "term" : { "TaskQueue" : "`+tq+`" }},
              		{ "term" : { "BuildIds" : "versioned:`+build+`" }}
            	]
        	}
    	},
		"aggs": {
			"maxStart": {
			  "max": {
					"field": "StartTime"
				  }
			}
		}
		}`),
		"search",
		"lastStartOneBuildOneTq",
		false,
	)
}

func groupByBuildOpen() {
	runQuery([]byte(`{
		"size": 0,
		"query" : {
        	"bool" : {
            	"filter": [
              		{ "term": { "NamespaceId": "`+cmn.NamespaceId+`" }},
              		{ "term": { "ExecutionStatus": "Running" }}
            	]
        	}
    	},
		"aggs": {
			"group_by_BuildIds": {
          		"terms": {
          			"size": `+strconv.Itoa(groupBySize)+`,
            		"field": "BuildIds"
          		}
			}
		}
		}`),
		"search",
		"groupByBuildOpen",
		false,
	)
}

func groupByBuildCurrent() {
	runQuery([]byte(`{
		"size": 0,
		"query" : {
        	"bool" : {
            	"filter": [
              		{ "term": { "NamespaceId": "`+cmn.NamespaceId+`" }}
            	]
        	}
    	},
		"aggs": {
			"group_by_BuildIds": {
          		"terms": {
          			"size": `+strconv.Itoa(groupBySize)+`,
            		"field": "BuildIds",`+
		`"include": "versioned:B-22.*"`+
		//`"include": "current:versioned:.*"`+
		`}
			}
		}
		}`),
		"search",
		"groupByBuildCurrent",
		false,
	)
}

func groupByBuildCurrentNS() {
	runQuery([]byte(`{
		"size": 0,
		"aggs": {
			"group_by_BuildIds": {
          		"terms": {
          			"size": `+strconv.Itoa(groupBySize)+`,
            		"field": "BuildIds",`+
		`"include": "versioned:B-22.*"`+
		//`"include": "current:versioned:.*"`+
		`}
			}
		}
		}`),
		"search",
		"groupByBuildCurrentNS",
		false,
	)
}

func groupByBuildCurrentField() {
	runQuery([]byte(`{
		"size": 0,
		"query" : {
        	"bool" : {
            	"filter": [
              		{ "term": { "NamespaceId": "`+cmn.NamespaceId+`" }}
            	]
        	}
    	},
		"aggs": {
			"group_by_BuildIds": {
          		"terms": {
          			"size": `+strconv.Itoa(groupBySize)+`,
            		"field": "CurrentBuildId",`+
		`"include": "versioned:B-22.*"`+
		//`"include": "current:versioned:.*"`+
		`}
			}
		}
		}`),
		"search",
		"groupByBuildCurrentField",
		false,
	)
}

func groupByBuildCurrentNSField() {
	runQuery([]byte(`{
		"size": 0,
		"aggs": {
			"group_by_BuildIds": {
          		"terms": {
          			"size": `+strconv.Itoa(groupBySize)+`,
            		"field": "CurrentBuildId",`+
		`"include": "versioned:B-22.*"`+
		//`"include": "current:versioned:.*"`+
		`}
			}
		}
		}`),
		"search",
		"groupByBuildCurrentNSField",
		false,
	)
}

func groupByBuildTqOpen() {
	runQuery([]byte(`{
		"size": 0,
		"query" : {
        	"bool" : {
            	"filter": [
              		{ "term": { "NamespaceId": "`+cmn.NamespaceId+`" }},
              		{ "term": { "ExecutionStatus": "Running" }}
            	]
        	}
    	},
		"aggs": {
			"group_by_BuildIdsTQ": {
				"multi_terms": {
          			"size": `+strconv.Itoa(groupBySize)+`,
					"terms": [{
					  "field": "BuildIds" 
					}, {
					  "field": "TaskQueue"
					}]
				  }
			}
		}
		}`),
		"search",
		"groupByBuildTqOpen",
		false,
	)
}

func groupByBuildFiltered(buildIds []string) {
	filter := `["versioned:` + strings.Join(buildIds, `", "versioned:`) + `"]`

	runQuery([]byte(`{
		"size": 0,
		"query" : {
        	"bool" : {
            	"filter": [
              		{ "term": { "NamespaceId": "`+cmn.NamespaceId+`" }},
              		{ "terms" : { "BuildIds" : `+filter+` }}
            	]
        	}
    	},
		"aggs": {
			"group_by_BuildIds": {
          		"terms": {
          			"size": `+strconv.Itoa(groupBySize)+`,
            		"field": "BuildIds"
          		}
			}
		}
		}`),
		"search",
		"groupByBuildFiltered",
		false,
	)
}

func groupByBuildTqFiltered(buildIds []string) {
	filter := `["versioned:` + strings.Join(buildIds, `", "versioned:`) + `"]`

	runQuery([]byte(`{
		"size": 0,
		"query" : {
        	"bool" : {
            	"filter": [
              		{ "term": { "NamespaceId": "`+cmn.NamespaceId+`" }},
              		{ "terms" : { "BuildIds" : `+filter+` }}
            	]
        	}
    	},
		"aggs": {
			"group_by_BuildIdsTQ": {
				"multi_terms": {
          			"size": `+strconv.Itoa(groupBySize)+`,
					"terms": [{
					  "field": "BuildIds" 
					}, {
					  "field": "TaskQueue"
					}]
				  }
			}
		}
		}`),
		"search",
		"groupByBuildTqFiltered",
		false,
	)
}

func runQuery(q []byte, qtype string, name string, log bool) {
	url := cmn.BaseUrl + cmn.Index + "/_" + qtype
	req, err := http.NewRequest("GET", url, bytes.NewBuffer(q))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	start := time.Now()
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	timer := time.Since(start)
	defer resp.Body.Close()

	if log {
		fmt.Println("####### QUERY:", name)
		fmt.Println(string(q))
		fmt.Println("response Status:", resp.Status)
	}
	body, _ := io.ReadAll(resp.Body)
	var data map[string]interface{}
	json.Unmarshal(body, &data)

	statspoint := statsPoint{timer: timer.Microseconds()}
	took := data["took"]
	if took == nil {
		statspoint.took = -1
	} else {
		statspoint.took = int64(took.(float64))
	}

	if log {
		fmt.Printf("resp.took: %d timer: %d \n", statspoint.took, statspoint.timer)
		fmt.Println("response Body:", string(body))
	}
	s := queryStats[name]
	if s == nil {
		s = []statsPoint{statspoint}
		queryStats[name] = s
	} else {
		queryStats[name] = append(s, statspoint)
	}
}

func recordStrategyStats(name string, start time.Time) {
	took := time.Since(start).Microseconds()

	s := strategyStats[name]
	if s == nil {
		s = []int64{took}
		strategyStats[name] = s
	} else {
		strategyStats[name] = append(s, took)
	}
}
