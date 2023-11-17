package main

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"
	cmn "temporaltests/versioning-visibility"
	"time"

	"github.com/google/uuid"
)

var (
	workers   = 1
	tasks     = make(chan string)
	batchSize = 100
)

func main() {
	for w := 0; w < workers; w++ {
		go work()
	}

	makeTasks()
}

func makeTasks() {
	end := time.Now().Truncate(time.Second).UnixNano()
	start := end - int64(cmn.Retention)

	tqs := []string{}
	for i := 0; i < cmn.NTaskQueues; i++ {
		tqs = append(tqs, uuid.New().String())
	}

	for sec := start; sec < end; sec += int64(time.Second) {
		status := cmn.StatusCompleted
		if sec+int64(cmn.WfLen) > end {
			status = cmn.StatusRunning
		}
		for i := 0; i < cmn.Rps; i++ {
			bld := cmn.BuildId(int((sec - start) / int64(cmn.BuildLen)))
			tq := cmn.TaskQueue(int((sec / int64(time.Second)) % int64(cmn.NTaskQueues)))
			putDoc(bld, status, tq, sec)
		}

		fmt.Println(int(100*(sec-start)/int64(cmn.Retention)), "%")
	}
	fmt.Println("Done!")
	//os.Exit(0)
}

func putDoc(buildId string, status string, tq string, start int64) {
	timestamp := time.Unix(0, start).UTC().Format(time.RFC3339Nano)
	closeTime := ""
	currentBuildId := ""
	currentBuildIdField := ""
	if status == cmn.StatusCompleted {
		closeTime = `"CloseTime": "` + time.Unix(0, start+int64(cmn.WfLen)).UTC().Format(time.RFC3339Nano) + `",`
	} else {
		currentBuildId = `"current:versioned:` + cmn.NamespaceId + `:` + buildId + `",`
		currentBuildIdField = `"CurrentBuildId": "` + cmn.NamespaceId + `:` + buildId + `",`
	}
	var jsonStr = `{
		  ` + closeTime + currentBuildIdField + `
          "BuildIds": [
            ` + currentBuildId + `
            "versioned:` + buildId + `"
          ],
          "ExecutionStatus": "` + status + `",
          "ExecutionTime": "` + timestamp + `",
          "NamespaceId": "` + cmn.NamespaceId + `",
          "RunId": "` + uuid.New().String() + `",
          "StartTime": "` + timestamp + `",
          "TaskQueue": "` + tq + `",
          "VisibilityTaskKey": "1104~218709045",
          "WorkflowId": "` + uuid.New().String() + `",
          "WorkflowType": "my_workflow"
        }`
	tasks <- jsonStr
}

func work() {
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
		data.Write([]byte(`{ "create" : { "_index" : "` + cmn.Index + `", "_id" : "` + uuid.New().String() + `" } }
`))
		data.Write(t)
		data.Write([]byte("\n"))
	}

	url := cmn.BaseUrl + "_bulk/"
	req, err := http.NewRequest("POST", url, data)
	req.Header.Set("Content-Type", "application/json")

	if cmn.Password != "" {
		req.SetBasicAuth("temporal", cmn.Password)
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
