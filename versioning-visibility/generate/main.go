package main

import (
	"fmt"
	cmn "temporaltests/versioning-visibility"
	"time"

	"github.com/google/uuid"
)

var (
	workers   = 1
	batchSize = 100
)

func main() {
	for w := 0; w < workers; w++ {
		go cmn.Work(batchSize)
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
			cmn.PutDoc(cmn.NamespaceId, bld, status, tq, sec)
		}

		fmt.Println(int(100*(sec-start)/int64(cmn.Retention)), "%")
	}
	fmt.Println("Done!")
	//os.Exit(0)
}
