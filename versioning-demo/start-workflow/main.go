package main

import (
	"context"
	"fmt"
	common "temporaltests"
	"time"

	"go.temporal.io/sdk/client"
)

func main() {
	c, _ := common.NewCloudClient("shahab-test4")
	defer c.Close()

	timer := time.Tick(3 * time.Second)

	for {
		select {
		case <-timer:
			fmt.Println("Starting new execution of MyWorkflow")
			_, err := c.ExecuteWorkflow(
				context.Background(), client.StartWorkflowOptions{
					TaskQueue: "my-tq",
				}, "MyWorkflow",
			)
			if err != nil {
				panic(err)
			}
		}
	}
}
