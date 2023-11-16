Versioning Visibility Queries
=============================

Useful commands
---------------

check schema
```shell
curl localhost:9200/temporal_visibility_v1_dev/_mapping | jq
```

find some docs
```shell
curl localhost:9200/temporal_visibility_v1_dev/_search -H 'Content-Type: application/json' -d '{                                     
    "size": 100}' | jq    
```

count docs
```shell
curl localhost:9200/temporal_visibility_v1_dev/_count -H 'Content-Type: application/json' -d '{}' | jq    
```

delete all docs
```shell
curl -X POST localhost:9200/temporal_visibility_v1_dev/_doc/_delete_by_query -H 'Content-Type: application/json' -d '{    "query": { 
        "match_all": {}   
    }}' | jq                    
```

delete index!
```shell
curl -X DELETE "localhost:9200/temporal_visibility_v1_dev" | jq                    
```

count by build id
```shell
curl localhost:9200/temporal_visibility_v1_dev/_search -H 'Content-Type: application/json' -d '{
   "size": 0,
   "aggs": {
        "group_by": {
          "terms": {
            "size": 10000,
            "field": "BuildIds"
          }
        }
    }
}' | jq .
```

count by build id
```shell
curl localhost:9200/temporal_visibility_v1_dev/_search -H 'Content-Type: application/json' -d '{
   "size": 10,
   "query" : {
        "bool" : {
            "filter": [
              { "term": { "NamespaceId": "7139bd23-bc97-4e49-943d-3f05ac1d3e5f" }},
              { "term": { "ExecutionStatus": "Running" }}
            ]
        }
    },
   "aggs": {
        "group_by": {
          "terms": {
            "size": 10000,
            "field": "BuildIds"
          }
        }
    }
}' | jq .
```

count by build id
```shell
curl localhost:9200/temporal_visibility_v1_dev/_search -H 'Content-Type: application/json' -d '{
   "size": 0,
   "aggs": {
        "group_by": {
          "terms": {
            "size": 10000,
            "include": "current:.*",
            "field": "BuildIds"
          }
        }
    }
}' | jq .
```

count by task queue
```shell
curl localhost:9200/temporal_visibility_v1_dev/_search -H 'Content-Type: application/json' -d '{
   "size": 0,
   "aggs": {
        "group_by": {
          "terms": {
          "size": 100,
            "field": "TaskQueue"
          }
        }
    }
}' | jq .
```

count by build id and task queue
```shell
curl localhost:9200/temporal_visibility_v1_dev/_search -H 'Content-Type: application/json' -d '{
   "size": 0,
   "aggs": {
        "group_by_BuildIdsTQ": {
            "multi_terms": {
                "size": 100,
                "terms": [{
                  "field": "BuildIds" 
                }, {
                  "field": "TaskQueue"
                }]
              }
        }
    }
}' | jq .
```

count by build id and task queue
```shell
curl localhost:9200/temporal_visibility_v1_dev/_search -H 'Content-Type: application/json' -d '{
  "size": 0,
  "aggs": {
    "group_by_BuildIdsTQ": {
      "composite": {
        "sources": [
          { "BuildIds": { "terms": { "field": "BuildIds" } } },
          { "TaskQueue": { "terms": { "field": "TaskQueue" } } }
        ]
      }
    }
  }
}' | jq .
```

count by status
```shell
curl localhost:9200/temporal_visibility_v1_dev/_search -H 'Content-Type: application/json' -d '{
   "size": 0,
   "aggs": {
        "group_by": {
          "terms": {
          "size": 100,
            "field": "ExecutionStatus"
          }
        }
    }
}' | jq .
```

min and max StartTime
```shell
curl localhost:9200/temporal_visibility_v1_dev/_search -H 'Content-Type: application/json' -d '{
   "size": 0,
   "aggs": {
        "minStart": {
          "min": {
                "field": "StartTime"
              }
        },
        "maxStart": {
          "max": {
                "field": "StartTime"
              }
        }
    }
}' | jq .
```

put a doc
```shell
curl -X PUT localhost:9200/temporal_visibility_v1_dev/_doc/15  -H 'Content-Type: application/json' -d '{                             
          "BuildIds": [
            "unversioned",
            "unversioned:@temporalio/worker@1.8.4+ca6af558bed79f78bb9fe99565c5233d8af106fd254a55bfec118c929bc28a83"
          ],
          "ExecutionStatus": "Running",
          "ExecutionTime": "2023-09-20T15:58:02.215974181Z",
          "NamespaceId": "6139bd23-bc97-4e49-943d-3f05ac1d3e5f",
          "RunId": "34e7b9d3-74f8-47a7-9d85-2c0fb9bac015",
          "StartTime": "2023-09-20T14:57:02.215974181Z",
          "TaskQueue": "provolone-pooled-default-native",
          "VisibilityTaskKey": "1104~218709045",
          "WorkflowId": "f2d6d6f2-a821-4ec1-bbd5-fbfa0150ada7",
          "WorkflowType": "change_request_fetcher"
        }'
```

put bulk
```shell
curl -X POST localhost:9200/_bulk  -H 'Content-Type: application/json' -d '{"create" : { "_index" : "temporal_visibility_v1_dev", "_id" : "8e7e4fdb-7cdc-42ad-bf6f-1abba05afa52" } }
{          "BuildIds": [            "versioned:B-0-b63b42e0-b3a0-4c01-9784-22cf04071d27"          ],          "ExecutionStatus": "Completed",          "ExecutionTime": "2023-10-30T07:07:22Z",          "NamespaceId": "6139bd23-bc97-4e49-943d-3f05ac1d3e5f",          "RunId": "7808e445-e35e-4c33-85bc-bf343fddbcd8",          "StartTime": "2023-10-30T07:07:22Z",          "TaskQueue": "TQ-2-241dceac-d599-4b85-93e4-0c3217ecb4f0",          "VisibilityTaskKey": "1104~218709045",          "WorkflowId": "eacb7415-34d5-4000-86d4-c29508a6fef1",          "WorkflowType": "my_workflow"        }
{ "create" : { "_index" : "temporal_visibility_v1_dev", "_id" : "6802272e-d894-464e-8853-2c80944ae4c5" } }
{          "BuildIds": [            "versioned:B-0-b63b42e0-b3a0-4c01-9784-22cf04071d27"          ],          "ExecutionStatus": "Completed",          "ExecutionTime": "2023-10-30T07:07:33Z",          "NamespaceId": "6139bd23-bc97-4e49-943d-3f05ac1d3e5f",          "RunId": "94390e5c-e693-4c6d-8302-7d1aef275ae2",          "StartTime": "2023-10-30T07:07:33Z",          "TaskQueue": "TQ-3-241dceac-d599-4b85-93e4-0c3217ecb4f0",          "VisibilityTaskKey": "1104~218709045",          "WorkflowId": "cf60fed2-42f0-43e8-9685-316b3c58c7c7",          "WorkflowType": "my_workflow"        }
'
```

cluster health
```shell
curl "localhost:9200/_cluster/health/?level=shards" | jq
curl "localhost:9200/_cluster/allocation/explain" | jq
```

data size
```shell
curl "localhost:9200/_cat/shards?v=true"
```


Running elasticsearch locally
-----------------------------
export your java home
```shell
export ES_JAVA_HOME=/opt/homebrew/Cellar/openjdk/21.0.1/libexec/openjdk.jdk/Contents/Home
```
export java options
```shell
export ES_JAVA_OPTS="$ES_JAVA_OPTS -XX:-MaxFDLimit"
```
maybe increase open file limit
```shell
ulimit -S -n 400000
```