package main

import (
	"context"
	"fmt"
	"log"
	common "temporaltests"
	"time"

	"github.com/google/uuid"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

const (
	oldVersioningTq = "oldVersioned-tq"
	newVersioningTq = "newVersioned-tq"
	unversionedTq   = "unversioned-tq"
)

const bld0, bld1, bld2, bld3 = "v1.0", "v1.1", "v1.2", "v1.3"

func main() {
	//c, _ := experiment.NewLocalClient("default")
	c, _ := common.NewCloudClient("shahab-test3")
	defer c.Close()

	//runNewVersioningStuff(c)
	//startWf(c, unversionedTq)
	//startWf(c, oldVersioningTq)
	//go runWorker(c, unversionedTq, "BUILD-U.1", false)
	//go runWorker(c, oldVersioningTq, "BUILD-V.0", true)
	//go runWorker(c, oldVersioningTq, "", false)
	time.Sleep(time.Hour)
}

func runNewVersioningStuff(c client.Client) {
	runNewVersioningWorkers(c)
	cleanRules(c)
	addRules(c)
	//addWorkflows(c)
	describeWf(c, "main-wf")
	describeWf(c, "child-with-inherit")
	describeWf(c, "child-using-rules")
}

func describeWf(c client.Client, wfId string) {
	wf, err := c.DescribeWorkflowExecution(context.Background(), wfId, "")
	if err != nil {
		panic(err)
	}
	fmt.Println(wf.String())
}

func act(desc string) (string, error) {
	return "OK.", nil
}

func childWf(ctx workflow.Context, desc string) (string, error) {
	return "OK.", nil
}

func wf(ctx workflow.Context, desc string) (string, error) {
	switch desc {
	case "should be assigned to v1.1 and complete in that build":
		// run two activities, one with two activities with different versioning intent
		fut := workflow.ExecuteActivity(
			workflow.WithActivityOptions(
				ctx, workflow.ActivityOptions{
					StartToCloseTimeout: 10 * time.Second,
					VersioningIntent:    temporal.VersioningIntentCompatible,
				},
			), act, "should complete in v1.1",
		)
		var val string
		err := fut.Get(ctx, &val)
		if err != nil {
			panic(err)
		}
		fut = workflow.ExecuteActivity(
			workflow.WithActivityOptions(
				ctx, workflow.ActivityOptions{
					StartToCloseTimeout: 10 * time.Second,
					VersioningIntent:    temporal.VersioningIntentDefault,
				},
			), act, "should use rules and complete in v1.2",
		)
		err = fut.Get(ctx, &val)
		if err != nil {
			panic(err)
		}
		return "", workflow.NewContinueAsNewError(
			workflow.WithWorkflowVersioningIntent(
				ctx,
				temporal.VersioningIntentCompatible,
			), wf, "should inherit build ID v1.1 from previous run and stay there",
		)
	case "should inherit build ID v1.1 from previous run and stay there":
		var val string
		// run compatible child wf
		fut := workflow.ExecuteChildWorkflow(
			workflow.WithWorkflowVersioningIntent(
				workflow.WithWorkflowID(ctx, "child-with-inherit"),
				temporal.VersioningIntentCompatible,
			), childWf, "should inherit build ID v1.1 from parent wf and stay there",
		)
		err := fut.Get(ctx, &val)
		if err != nil {
			panic(err)
		}
		return "", workflow.NewContinueAsNewError(
			workflow.WithWorkflowVersioningIntent(
				ctx,
				temporal.VersioningIntentDefault,
			), wf, "should start on latest build ID v1.2 and end in v1.3 via redirect rules",
		)
	case "should start on latest build ID v1.2 and end in v1.3 via redirect rules":
		var val string
		// run latest child wf
		fut := workflow.ExecuteChildWorkflow(
			workflow.WithWorkflowVersioningIntent(
				workflow.WithWorkflowID(ctx, "child-using-rules"),
				temporal.VersioningIntentDefault,
			), childWf, "should start on latest build ID v1.2 and end in v1.3 via redirect rules",
		)
		err := fut.Get(ctx, &val)
		if err != nil {
			panic(err)
		}
	}
	return "OK.", nil
}

func addWorkflows(c client.Client) {
	_, err := c.ExecuteWorkflow(
		context.Background(), client.StartWorkflowOptions{
			TaskQueue: newVersioningTq,
			ID:        "main-wf",
		}, wf, "should be assigned to v1.1 and complete in that build",
	)
	if err != nil {
		panic(err)
	}

}

func runNewVersioningWorkers(c client.Client) {
	go runWorker(c, newVersioningTq, "", false)
	go runWorker(c, newVersioningTq, bld0, true)
	go runWorker(c, newVersioningTq, bld1, true)
	go runWorker(c, newVersioningTq, bld2, true)
	go runWorker(c, newVersioningTq, bld3, true)
}

func cleanRules(c client.Client) {
	ctx := context.Background()
	rules, err := c.GetWorkerVersioningRules(
		ctx, &client.GetWorkerVersioningOptions{
			TaskQueue: newVersioningTq,
		},
	)
	if err != nil {
		panic(err)
	}
	count := len(rules.AssignmentRules)
	for i := 0; i < count; i++ {
		rules, err = c.UpdateWorkerVersioningRules(
			ctx, &client.UpdateWorkerVersioningRulesOptions{
				TaskQueue:     newVersioningTq,
				ConflictToken: rules.ConflictToken,
				Operation: &client.VersioningOpDeleteAssignmentRule{
					RuleIndex: 0,
					Force:     true,
				},
			},
		)
		if err != nil {
			panic(err)
		}
	}
	redirects := rules.RedirectRules
	for _, r := range redirects {
		rules, err = c.UpdateWorkerVersioningRules(
			ctx, &client.UpdateWorkerVersioningRulesOptions{
				TaskQueue:     newVersioningTq,
				ConflictToken: rules.ConflictToken,
				Operation: &client.VersioningOpDeleteRedirectRule{
					SourceBuildID: r.Rule.SourceBuildID,
				},
			},
		)
		if err != nil {
			panic(err)
		}
	}
}

func addRules(c client.Client) {
	ctx := context.Background()
	rules, err := c.GetWorkerVersioningRules(
		ctx, &client.GetWorkerVersioningOptions{
			TaskQueue: newVersioningTq,
		},
	)
	if err != nil {
		panic(err)
	}
	rules, err = c.UpdateWorkerVersioningRules(
		ctx, &client.UpdateWorkerVersioningRulesOptions{
			TaskQueue:     newVersioningTq,
			ConflictToken: rules.ConflictToken,
			Operation: &client.VersioningOpInsertAssignmentRule{
				Rule: client.VersioningAssignmentRule{
					TargetBuildID: bld2,
				},
			},
		},
	)
	if err != nil {
		panic(err)
	}
	rules, err = c.UpdateWorkerVersioningRules(
		ctx, &client.UpdateWorkerVersioningRulesOptions{
			TaskQueue:     newVersioningTq,
			ConflictToken: rules.ConflictToken,
			Operation: &client.VersioningOpInsertAssignmentRule{
				Rule: client.VersioningAssignmentRule{
					TargetBuildID: bld0,
					Ramp: &client.VersioningRampByPercentage{
						Percentage: 0,
					},
				},
			},
		},
	)
	if err != nil {
		panic(err)
	}
	rules, err = c.UpdateWorkerVersioningRules(
		ctx, &client.UpdateWorkerVersioningRulesOptions{
			TaskQueue:     newVersioningTq,
			ConflictToken: rules.ConflictToken,
			Operation: &client.VersioningOpAddRedirectRule{
				Rule: client.VersioningRedirectRule{
					SourceBuildID: bld2,
					TargetBuildID: bld3,
				},
			},
		},
	)
	if err != nil {
		panic(err)
	}
}

func startWf(c client.Client, tq string) {
	options := client.StartWorkflowOptions{
		ID:        tq + "_" + uuid.New().String(),
		TaskQueue: tq,
	}
	we, err := c.ExecuteWorkflow(context.Background(), options, MyWorkflow)
	if err != nil {
		log.Fatalln("Unable to execute workflow", err)
	}
	log.Println("Started workflow", "WorkflowID", we.GetID())
}

func runWorker(c client.Client, tq string, buildId string, useVersioning bool) {
	options := worker.Options{UseBuildIDForVersioning: useVersioning}
	if buildId != "" {
		options.BuildID = buildId
	}

	w := worker.New(c, tq, options)
	w.RegisterWorkflow(wf)
	w.RegisterWorkflow(childWf)
	w.RegisterActivity(act)
	//w.RegisterWorkflow(MyWorkflow)
	//w.RegisterActivity(MyActivity)

	err := w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start worker for "+tq, err)
	}
}

func MyWorkflow(ctx workflow.Context) (string, error) {
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
	time.Sleep(time.Hour)
	return fmt.Sprintf("Returning after 1 sec"), nil
}
