package main

import (
	"context"
	"fmt"
	common "temporaltests"
	"time"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

func main() {
	c, _ := common.NewCloudClient("shahab-test4")
	defer c.Close()

	w := worker.New(
		c, "my-tq", worker.Options{
			UseBuildIDForVersioning: true,
			BuildID:                 "B",
		},
	)
	w.RegisterWorkflowWithOptions(MyWorkflow, workflow.RegisterOptions{Name: "MyWorkflow"})
	w.RegisterActivity(MyActivity)

	err := w.Run(worker.InterruptCh())
	if err != nil {
		panic(err)
	}

	time.Sleep(time.Hour)
}

func MyActivity(input string) (string, error) {
	fmt.Printf("Working on %s\n", input)
	return "OK.", nil
}

func MyWorkflow(ctx workflow.Context, desc string) (string, error) {
	err := workflow.Sleep(ctx, 2*time.Second)
	if err != nil {
		panic(err)
	}

	fut := workflow.ExecuteActivity(
		workflow.WithActivityOptions(
			ctx, workflow.ActivityOptions{
				StartToCloseTimeout: 10 * time.Second,
			},
		), MyActivity, "step 1",
	)
	var val string
	err = fut.Get(ctx, &val)
	if err != nil {
		panic(err)
	}

	err = workflow.Sleep(ctx, 2*time.Second)
	if err != nil {
		panic(err)
	}

	fut = workflow.ExecuteActivity(
		workflow.WithActivityOptions(
			ctx, workflow.ActivityOptions{
				StartToCloseTimeout: 10 * time.Second,
			},
		), MyActivity, "step 2",
	)
	err = fut.Get(ctx, &val)
	if err != nil {
		panic(err)
	}
	return "OK.", nil
}

func addWorkflows(c client.Client) {
	_, err := c.ExecuteWorkflow(
		context.Background(), client.StartWorkflowOptions{
			TaskQueue: "my-tq",
		}, MyWorkflow,
	)
	if err != nil {
		panic(err)
	}
}
