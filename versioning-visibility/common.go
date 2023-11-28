package versioning_visibility

import (
	"bytes"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	BaseUrl  = "http://localhost:9200/"
	Index    = "temporal_visibility_v1_dev"
	Password = ""

	//NamespaceId = "6139bd23-bc97-4e49-943d-3f05ac1d3e5f"
	//Retention   = 30*24*time.Hour + 5*time.Minute
	//Rps         = 10
	//WfLen       = 10 * time.Minute
	//BuildLen    = 60 * time.Minute
	//NTaskQueues = 5
	NamespaceId = "9139bd23-bc97-4e49-943d-3f05ac1d3e5f"
	Retention   = 1*time.Hour + 5*time.Minute
	Rps         = 1
	WfLen       = 10 * time.Minute
	BuildLen    = 30 * time.Minute
	NTaskQueues = 5

	FillingNamespaceId = ""
)

const (
	StatusRunning   = "Running"
	StatusCompleted = "Completed"
)

func BuildId(index int) string {
	return "B-" + strconv.Itoa(index) + "-b63b42e0-b3a0-4c01-9784-22cf04071d27"
}

func TaskQueue(index int) string {
	return "TQ-" + strconv.Itoa(index) + "-241dceac-d599-4b85-93e4-0c3217ecb4f0"
}

var (
	tasks = make(chan string)
)

func PutDoc(nsId string, buildId string, status string, tq string, start int64) {
	timestamp := time.Unix(0, start).UTC().Format(time.RFC3339Nano)
	closeTime := ""
	currentBuildId := ""
	currentBuildIdField := ""
	if status == StatusCompleted {
		closeTime = `"CloseTime": "` + time.Unix(0, start+int64(WfLen)).UTC().Format(time.RFC3339Nano) + `",`
	} else {
		currentBuildId = `"current:versioned:` + nsId + `:` + buildId + `",`
		currentBuildIdField = `"CurrentBuildId": "` + nsId + `:` + buildId + `",`
	}
	var jsonStr = `{
		  ` + closeTime + currentBuildIdField + `
          "BuildIds": [
            ` + currentBuildId + `
            "versioned:` + buildId + `"
          ],
          "ExecutionStatus": "` + status + `",
          "ExecutionTime": "` + timestamp + `",
          "NamespaceId": "` + nsId + `",
          "RunId": "` + uuid.New().String() + `",
          "StartTime": "` + timestamp + `",
          "TaskQueue": "` + tq + `",
          "VisibilityTaskKey": "1104~218709045",
          "WorkflowId": "` + uuid.New().String() + `",
          "WorkflowType": "my_workflow"
        }`
	tasks <- jsonStr
}

func Work(batchSize int) {
	client := &http.Client{}
	buf := [][]byte{}
	for t := range tasks {
		buf = append(buf, []byte(strings.Replace(t, "\n", "", -1)))
		if len(buf) >= batchSize {
			submit(client, buf)
			buf = [][]byte{}
		}
	}
	if len(buf) > 0 {
		submit(client, buf)
	}
}

func submit(client *http.Client, buf [][]byte) {
	data := bytes.NewBuffer([]byte{})
	for _, t := range buf {
		data.Write([]byte(`{ "create" : { "_index" : "` + Index + `", "_id" : "` + uuid.New().String() + `" } }
`))
		data.Write(t)
		data.Write([]byte("\n"))
	}

	url := BaseUrl + "_bulk/"
	req, err := http.NewRequest("POST", url, data)
	req.Header.Set("Content-Type", "application/json")

	if Password != "" {
		req.SetBasicAuth("temporal", Password)
	}

	var resp *http.Response
	for i := 100; ; i = min(i*2, 10000) {
		resp, err = client.Do(req)
		if err != nil {
			fmt.Println("sleeping Error:", err)
			time.Sleep(time.Duration(int64(time.Millisecond) * int64(i)))
			continue
		}
		if resp.Status[0] != '2' {
			fmt.Println("response Status:", resp.Status)
			fmt.Println("sleeping Error:", err)
			time.Sleep(time.Duration(int64(time.Millisecond) * int64(i)))
			continue
		}
		break
	}

	defer resp.Body.Close()
	fmt.Println("Batch Created")
	//fmt.Println("response Headers:", resp.Header)
	//body, _ := io.ReadAll(resp.Body)
	//fmt.Println("response Body:", string(body))
}
