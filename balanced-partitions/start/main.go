package main

import (
	"context"
	"fmt"
	"log"
	"os"
	experiment "temporaltests/balanced-partitions"
	"time"

	"github.com/google/uuid"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

var expName string
var exp experiment.Experiment
var activityWait time.Duration

func main() {
	//c, _ := experiment.NewLocalClient("backlog_exp")
	c, _ := experiment.NewCloudClient("shahab-test2")
	defer c.Close()

	expName = os.Args[1]
	exp = experiment.Experiments[expName]
	activityWait = time.Duration(int64(float32(time.Second) * float32(exp.ConcurrentActivities) / exp.ActivitiesCanHandlePerSec))

	if exp.ConcurrentPollers > 0 {
		go runClient(c)
		runWorker(c)
	} else {
		runClient(c)
	}
}

func runClient(c client.Client) {
	time.Sleep(time.Second * 2)
	for _, s := range exp.Stages {
		start := time.Now()
		sleep := time.Duration(int64(float32(time.Second) / s.StartWorkflowPerSec))
		for {
			if time.Since(start) > s.Duration {
				break
			}
			go startWf(c)
			time.Sleep(sleep)
		}
	}
}

func startWf(c client.Client) {
	options := client.StartWorkflowOptions{
		ID:        uuid.New().String(),
		TaskQueue: "tq",
	}
	we, err := c.ExecuteWorkflow(context.Background(), options, MyWorkflow, expName)
	if err != nil {
		log.Fatalln("Unable to execute workflow", err)
	}
	log.Println("Started workflow", "WorkflowID", we.GetID())
}

func runWorker(c client.Client) {
	w := worker.New(
		c, "tq", worker.Options{
			MaxConcurrentActivityExecutionSize:     exp.ConcurrentActivities,
			MaxConcurrentActivityTaskPollers:       exp.ConcurrentPollers,
			MaxConcurrentWorkflowTaskPollers:       200,
			MaxConcurrentWorkflowTaskExecutionSize: 2000,
		},
	)

	w.RegisterWorkflow(MyWorkflow)
	w.RegisterActivity(MyActivity)

	err := w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start worker", err)
	}
}

func MyWorkflow(ctx workflow.Context, expName string) (string, error) {
	options := workflow.ActivityOptions{
		StartToCloseTimeout: time.Second * 20,
	}
	ctx = workflow.WithActivityOptions(ctx, options)

	var activityOutput string
	err := workflow.ExecuteActivity(ctx, MyActivity).Get(ctx, &activityOutput)
	if err != nil {
		return "", err
	}

	return "Activity: " + activityOutput, nil
}

func MyActivity(ctx context.Context) (string, error) {
	time.Sleep(activityWait)
	return fmt.Sprintf("Returning after %s", activityWait), nil
}
