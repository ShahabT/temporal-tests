package experiment

import "time"

type Stage struct {
	StartWorkflowPerSec float32
	Duration            time.Duration
}

type Experiment struct {
	ActivitiesCanHandlePerSec float32
	ConcurrentPollers         int
	ConcurrentActivities      int
	Stages                    []Stage
}

var Experiments = map[string]Experiment{
	"large_worker": {
		ActivitiesCanHandlePerSec: 100000,
		ConcurrentPollers:         100,
		ConcurrentActivities:      100,
	},
	"plenty_workers_small": {
		ActivitiesCanHandlePerSec: 10,
		ConcurrentPollers:         2,
		ConcurrentActivities:      2,
		Stages: []Stage{
			{
				StartWorkflowPerSec: 6,
				Duration:            3 * time.Minute,
			},
		},
	},
	"plenty_workers_large": {
		ActivitiesCanHandlePerSec: 400,
		ConcurrentPollers:         100,
		ConcurrentActivities:      100,
		Stages: []Stage{
			{
				StartWorkflowPerSec: 150,
				Duration:            3 * time.Minute,
			},
		},
	},
	"just_enough_workers_small": {
		ActivitiesCanHandlePerSec: 22, // effectively 20
		ConcurrentPollers:         4,
		ConcurrentActivities:      4,
		Stages: []Stage{
			{
				StartWorkflowPerSec: 19,
				Duration:            3 * time.Minute,
			},
		},
	},
	"just_enough_workers_large": {
		ActivitiesCanHandlePerSec: 160,
		ConcurrentPollers:         40,
		ConcurrentActivities:      40,
		Stages: []Stage{
			{
				StartWorkflowPerSec: 130,
				Duration:            3 * time.Minute,
			},
		},
	},
	"not_enough_workers_small": {
		ActivitiesCanHandlePerSec: 20,
		ConcurrentPollers:         4,
		ConcurrentActivities:      4,
		Stages: []Stage{
			{
				StartWorkflowPerSec: 40,
				Duration:            3 * time.Minute,
			},
		},
	},
	"not_enough_workers_large": {
		ActivitiesCanHandlePerSec: 90,
		ConcurrentPollers:         10,
		ConcurrentActivities:      10,
		Stages: []Stage{
			{
				StartWorkflowPerSec: 180,
				Duration:            3 * time.Minute,
			},
		},
	},
	"few_workers_2": {
		ActivitiesCanHandlePerSec: 20,
		ConcurrentPollers:         2,
		ConcurrentActivities:      2,
		Stages: []Stage{
			{
				StartWorkflowPerSec: 40,
				Duration:            1 * time.Minute,
			},
		},
	},
	"few_workers_1": {
		ActivitiesCanHandlePerSec: 5,
		ConcurrentPollers:         1,
		ConcurrentActivities:      2,
		Stages: []Stage{
			{
				StartWorkflowPerSec: 10,
				Duration:            1 * time.Minute,
			},
		},
	},
	"few_workers_super_slow": {
		ActivitiesCanHandlePerSec: .1,
		ConcurrentPollers:         1,
		ConcurrentActivities:      1,
		Stages: []Stage{
			{
				StartWorkflowPerSec: 5,
				Duration:            5 * time.Second,
			},
		},
	},
}
