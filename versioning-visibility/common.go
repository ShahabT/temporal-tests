package versioning_visibility

import (
	"strconv"
	"time"
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
