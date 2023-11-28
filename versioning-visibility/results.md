export PATH=$PATH:/home/temporal/shahab/go/go/bin
export $(cat /secrets/elasticsearch/elasticsearchcreds.yml  | sed 's/: /=/' | tr '\n' ' ')
curl -XGET --user "${username}:${password}" "${scheme}://${hostname}:${port}/${visindex}/_search" -H 'Content-Type: application/json' -d'{
"size": 0,
"aggs": {
"Namespaces": {
"terms": {
"size": 100,
"field": "NamespaceId"
}
}
}
}' | jq .

curl -XGET --user "${username}:${password}" "${scheme}://${hostname}:${port}/${visindex}/_search" -H 'Content-Type: application/json' -d'{
"size": 0,
"query" : {
"bool" : {
"filter": [
{ "term": { "NamespaceId": "C6139bd23-bc97-4e49-943d-3f05ac1d3e5f" }},
{ "term": { "ExecutionStatus": "Running" }}
]
}
},
"aggs": {
"Namespaces": {
"terms": {
"size": 10000,
"field": "BuildIds"
}
}
}
}' | jq .

curl -XGET --user "${username}:${password}" "${scheme}://${hostname}:${port}/${visindex}/_count" -H 'Content-Type: application/json' -d'{
"query" : {
"bool" : {
"filter": [
{ "term": { "NamespaceId": "C6139bd23-bc97-4e49-943d-3f05ac1d3e5f" }},
{ "term": { "ExecutionStatus": "Running" }}
]
}
}
}' | jq .

curl -X POST --user "${username}:${password}" "${scheme}://${hostname}:${port}/${visindex}/_doc/_delete_by_query" -H 'Content-Type: application/json' -d'{
"query" : {
"term" : {
"NamespaceId": "X6139bd23-bc97-4e49-943d-3f05ac1d3e5f"
}
}
}' | jq .



########### EXP 1

        2.6M docs
        NamespaceId = "6139bd23-bc97-4e49-943d-3f05ac1d3e5f"
        Retention   = 3*24*time.Hour + 5*time.Minute
        Rps         = 10
        WfLen       = 10 * time.Minute
        BuildLen    = 60 * time.Minute
        NTaskQueues = 5

Run 1:

|                          Query |      Rep |      Avg |      Min |      Max |
|                  countOneBuild |       10 |    26689 |    14757 |    57608 |
|           countOneBuildCurrent |       10 |    28445 |    13087 |    55019 |
|      countOneBuildCurrentField |       10 |    27584 |    13762 |    47915 |
|         countOneBuildCurrentNS |       10 |    32878 |    11619 |    67511 |
|    countOneBuildCurrentNSField |       10 |    32196 |    12766 |    68738 |
|              countOneBuildOpen |       10 |    31987 |    13727 |    55890 |
|                   groupByBuild |       10 |   207520 |    66540 |   330641 |
|            groupByBuildCurrent |       10 |   378017 |   153664 |   700643 |
|       groupByBuildCurrentField |       10 |   108491 |    28587 |   443456 |
|          groupByBuildCurrentNS |       10 |   962975 |   126871 |  1938599 |
|     groupByBuildCurrentNSField |       10 |   317141 |   188590 |   687840 |
|               groupByBuildOpen |       10 |    65857 |    23323 |   207201 |


Run 2: (add limit 1)

|                          Query |      Rep |      Avg |      Min |      Max |
|                  countOneBuild |       10 |    28150 |    20994 |    46667 |
|           countOneBuildCurrent |       10 |    29916 |    20065 |    37887 |
|      countOneBuildCurrentField |       10 |    32502 |    21289 |    60971 |
|         countOneBuildCurrentNS |       10 |    23448 |    18159 |    31643 |
|    countOneBuildCurrentNSField |       10 |    24333 |    18606 |    34855 |
|              countOneBuildOpen |       10 |    27505 |    19916 |    45324 |
|                   groupByBuild |       10 |   165726 |    29543 |   273554 |
|            groupByBuildCurrent |       10 |   290634 |   204910 |   385608 |
|       groupByBuildCurrentField |       10 |    90417 |    57563 |   128012 |
|          groupByBuildCurrentNS |       10 |   802559 |    39858 |  1094621 |
|     groupByBuildCurrentNSField |       10 |   234357 |   157170 |   439491 |
|               groupByBuildOpen |       10 |    40462 |    29963 |    56201 |
|                  limitOneBuild |       10 |    38325 |    23468 |    63352 |
|         limitOneBuildOneTqOpen |       10 |    32501 |    23456 |    53420 |

Run 3: keep the good ones only

|                          Query |      Rep |      Avg |      Min |      Max |
|                  countOneBuild |       10 |    26573 |    14470 |    44436 |
|           countOneBuildCurrent |       10 |    26115 |    19252 |    43959 |
|      countOneBuildCurrentField |       10 |    22960 |    19414 |    28961 |
|         countOneBuildCurrentNS |       10 |    25251 |    19702 |    44284 |
|    countOneBuildCurrentNSField |       10 |    25160 |    18154 |    44599 |
|              countOneBuildOpen |       10 |    23990 |    17983 |    37954 |
|                   groupByBuild |       10 |   137681 |    21783 |   294451 |
|            groupByBuildCurrent |       10 |   235062 |    19503 |   418799 |
|       groupByBuildCurrentField |       10 |    69217 |    23656 |    99610 |
|               groupByBuildOpen |       10 |    32959 |    20520 |    42774 |

Run 4: with start times

|                          Query |      Rep |      Avg |      Min |      Max |
|                  countOneBuild |       10 |    20268 |    11152 |    24381 |
|      countOneBuildCurrentField |       10 |    37598 |    16429 |   173567 |
|             countOneBuildOneTq |       50 |    24256 |    17038 |    49946 |
|         countOneBuildOneTqOpen |       50 |    22427 |    11204 |    47474 |
|              countOneBuildOpen |       10 |    26516 |    18047 |    48752 |
|                   groupByBuild |       10 |   165925 |    20455 |   369841 |
|            groupByBuildCurrent |       10 |   185146 |    18773 |   303423 |
|       groupByBuildCurrentField |       10 |    75568 |    40116 |   128994 |
|               groupByBuildOpen |       10 |    38167 |    21200 |    66715 |
|              groupByTqOneBuild |       10 |    35449 |    22971 |    71006 |
|          groupByTqOneBuildOpen |       10 |    30456 |    14329 |    56515 |
|              lastStartOneBuild |       10 |    34446 |    21556 |    55880 |
|         lastStartOneBuildOneTq |       50 |    31213 |    14418 |    67475 |


####### EXP 2:

        60M docs
        60K open wfs
        NamespaceId = "B6139bd23-bc97-4e49-943d-3f05ac1d3e5f"
        Retention   = 7*24*time.Hour + 5*time.Minute
        Rps         = 100
        WfLen       = 10 * time.Minute
        BuildLen    = 60 * time.Minute
        NTaskQueues = 5
        groupby size = 1000

Run 1: (one repetition)
|                          Query |      Rep |      Avg |      Min |      Max |
|                  countOneBuild |        1 |    32547 |    32547 |    32547 |
|      countOneBuildCurrentField |        1 |    31276 |    31276 |    31276 |
|             countOneBuildOneTq |        5 |    63030 |    41310 |    95404 |
|         countOneBuildOneTqOpen |        5 |    32794 |    22539 |    52682 |
|              countOneBuildOpen |        1 |    21161 |    21161 |    21161 |
|                   groupByBuild |        1 |  4532789 |  4532789 |  4532789 |
|            groupByBuildCurrent |        1 |  5072986 |  5072986 |  5072986 |
|       groupByBuildCurrentField |        1 |  1557189 |  1557189 |  1557189 |
|               groupByBuildOpen |        1 |   159545 |   159545 |   159545 |
|              groupByTqOneBuild |        1 |    44325 |    44325 |    44325 |
|          groupByTqOneBuildOpen |        1 |    66117 |    66117 |    66117 |
|              lastStartOneBuild |        1 |    84858 |    84858 |    84858 |
|                  limitOneBuild |        1 |    38177 |    38177 |    38177 |
|         limitOneBuildOneTqOpen |        6 |    35970 |    26475 |    48747 |

Run 2: (ten repetition)

|                          Query |      Rep |      Avg |      Min |      Max |
|                  countOneBuild |       10 |    40307 |    27348 |    85012 |
|      countOneBuildCurrentField |       10 |    26544 |    13026 |    55703 |
|             countOneBuildOneTq |       50 |    36610 |    17643 |   158908 |
|         countOneBuildOneTqOpen |       50 |    31367 |    15564 |   106721 |
|              countOneBuildOpen |       10 |    29820 |    21268 |    45482 |
|                   groupByBuild |       10 |  3309390 |  2788460 |  4260725 |
|            groupByBuildCurrent |       10 |  4972547 |  3830900 |  8948044 |
|       groupByBuildCurrentField |       10 |  1391120 |   946249 |  2413527 |
|               groupByBuildOpen |       10 |    45728 |    32792 |    60057 |
|              groupByTqOneBuild |       10 |    75608 |    37411 |   184497 |
|          groupByTqOneBuildOpen |       10 |    42774 |    23361 |    82285 |
|              lastStartOneBuild |       10 |    69935 |    31661 |   290371 |
|                  limitOneBuild |       10 |    36038 |    25785 |    68729 |
|         limitOneBuildOneTqOpen |       60 |    35348 |    16718 |    99906 |

Run 3: (groupby size = 100)

|                          Query |      Rep |      Avg |      Min |      Max |
|                  countOneBuild |       10 |    36961 |    26067 |    60463 |
|             countOneBuildOneTq |       50 |    28213 |    20358 |    70435 |
|         countOneBuildOneTqOpen |       50 |    25987 |    12784 |    46487 |
|              countOneBuildOpen |       10 |    26350 |    19694 |    38008 |
|                   groupByBuild |       10 |  3123903 |  2650499 |  3733405 |
|            groupByBuildCurrent |       10 |  4675072 |  3892202 |  6505079 |
|          groupByBuildCurrentNS |       10 |  5611272 |  4190853 |  6935420 |
|               groupByBuildOpen |       10 |    49478 |    32545 |    73599 |
|              groupByTqOneBuild |       10 |    52552 |    34464 |   124096 |
|          groupByTqOneBuildOpen |       10 |    27172 |    20477 |    31547 |
|              lastStartOneBuild |       10 |    66859 |    40037 |   167484 |
|         lastStartOneBuildOneTq |       50 |    39353 |    26130 |   116481 |


### EXP 3

         121M docs
         1.2M open wfs
         NamespaceId = "C6139bd23-bc97-4e49-943d-3f05ac1d3e5f"
         Retention   = 7*24*time.Hour + 5*time.Minute
         Rps         = 200
         WfLen       = 100 * time.Minute
         BuildLen    = 60 * time.Minute
         NTaskQueues = 5

|                          Query |      Rep |      Avg |      Min |      Max |
|                  countOneBuild |       10 |    57845 |    25896 |   161970 |
|             countOneBuildOneTq |       50 |    37279 |    19635 |    80035 |
|         countOneBuildOneTqOpen |       50 |    45453 |    22214 |   170477 |
|              countOneBuildOpen |       10 |    53644 |    41794 |    94568 |
|                   groupByBuild |       10 |  6476570 |  5514112 |  7309902 |
|            groupByBuildCurrent |       10 |  9800303 |  7347307 | 13555863 |
|          groupByBuildCurrentNS |       10 | 13276881 |  9999999 | 17506542 |
|               groupByBuildOpen |       10 |   145545 |    23486 |   288233 |
|              groupByTqOneBuild |       10 |    78726 |    58122 |   148520 |
|          groupByTqOneBuildOpen |       10 |    91672 |    51928 |   187276 |
|              lastStartOneBuild |       10 |    90064 |    61115 |   141875 |
|         lastStartOneBuildOneTq |       50 |    51979 |    27607 |   143714 |


Run 2: add "filling" to clear caches

|                          Query |      Rep |      Avg |      Min |      Max |
|                  countOneBuild |       10 |    44027 |    33089 |    63447 |
|             countOneBuildOneTq |       50 |    38881 |    19833 |    92257 |
|         countOneBuildOneTqOpen |       50 |    40012 |    28995 |   117623 |
|              countOneBuildOpen |       10 |    56937 |    39797 |    86845 |
|                        filling |      240 |  3814496 |  2872884 |  8645842 |
|                   groupByBuild |       10 |  6284874 |  5563665 |  8007021 |
|            groupByBuildCurrent |       10 |  9073530 |  7287619 | 11679941 |
|          groupByBuildCurrentNS |       10 | 13475796 |  9999999 | 16614344 |
|               groupByBuildOpen |       10 |   144922 |   106479 |   270216 |
|              groupByTqOneBuild |       10 |    82574 |    57256 |   130255 |
|          groupByTqOneBuildOpen |       10 |    81291 |    63910 |   117078 |
|              lastStartOneBuild |       10 |    71869 |    59004 |    91644 |
|         lastStartOneBuildOneTq |       50 |    61351 |    28477 |   341928 |

Run 3: all queries

|                          Query |      Rep |      Avg |      Min |      Max |
|                  countOneBuild |       10 |    47280 |    32742 |    82741 |
|           countOneBuildCurrent |       10 |    47300 |    33468 |    89476 |
|      countOneBuildCurrentField |       10 |    42834 |    26811 |    76609 |
|         countOneBuildCurrentNS |       10 |    34924 |    25504 |    56959 |
|    countOneBuildCurrentNSField |       10 |    40605 |    24802 |   120772 |
|             countOneBuildOneTq |       50 |    33360 |    19965 |    76475 |
|         countOneBuildOneTqOpen |       50 |    38399 |    28480 |   106227 |
|              countOneBuildOpen |       10 |    52347 |    40741 |    71318 |
|                        filling |      420 |  3750399 |  2915269 |  9686493 |
|                   groupByBuild |       10 |  7006776 |  5507244 | 10391065 |
|            groupByBuildCurrent |       10 |  9509352 |  7798607 | 14647076 |
|       groupByBuildCurrentField |       10 |  2949867 |  2462337 |  5266524 |
|          groupByBuildCurrentNS |       10 | 14177939 |  9999999 | 17999098 |
|     groupByBuildCurrentNSField |       10 |  3850477 |  3161425 |  7178754 |
|               groupByBuildOpen |       10 |   173553 |   113094 |   323894 |
|              groupByTqOneBuild |       10 |    68745 |    55632 |    95116 |
|          groupByTqOneBuildOpen |       10 |    74726 |    65125 |    84574 |
|              lastStartOneBuild |       10 |    79514 |    58544 |   155640 |
|         lastStartOneBuildOneTq |       50 |    46323 |    26077 |   112189 |
|                  limitOneBuild |       10 |    37997 |    27202 |    64396 |
|             limitOneBuildOneTq |       50 |    45214 |    29367 |    99323 |
|         limitOneBuildOneTqOpen |       60 |    43770 |    28144 |   115931 |


Run 4: xl AWS instances, terminate_after=1 for count queries

|                          Query |      Rep |      Avg |      Min |      Max | Took Avg | Took Min | Took Max |
|                  countOneBuild |       10 |    22301 |    17843 |    30693 |       -1 |       -1 |       -1 |
|           countOneBuildCurrent |       10 |    23330 |    19465 |    38101 |       -1 |       -1 |       -1 |
|      countOneBuildCurrentField |       10 |    23090 |    19692 |    34980 |       -1 |       -1 |       -1 |
|         countOneBuildCurrentNS |       10 |    20427 |    18729 |    21894 |       -1 |       -1 |       -1 |
|    countOneBuildCurrentNSField |       10 |    18980 |    10243 |    23457 |       -1 |       -1 |       -1 |
|             countOneBuildOneTq |       50 |    22142 |    12275 |    41991 |       -1 |       -1 |       -1 |
|         countOneBuildOneTqOpen |       50 |    22797 |    10966 |    45981 |       -1 |       -1 |       -1 |
|              countOneBuildOpen |       10 |    22375 |    16824 |    34380 |       -1 |       -1 |       -1 |
|                        filling |      380 |  1455234 |    24101 |  2890727 |     1443 |        6 |     2883 |
|                   groupByBuild |       10 |  2745461 |  2456140 |  3386125 |     2729 |     2439 |     3370 |
|            groupByBuildCurrent |       10 |  3648920 |  2957844 |  4585832 |     3632 |     2939 |     4569 |
|       groupByBuildCurrentField |       10 |   998773 |   816808 |  1232678 |      982 |      799 |     1218 |
|          groupByBuildCurrentNS |       10 |  5091872 |  4585806 |  5500916 |     5075 |     4571 |     5485 |
|     groupByBuildCurrentNSField |       10 |  1052430 |   785132 |  1418091 |     1035 |      769 |     1401 |
|               groupByBuildOpen |       10 |    76906 |    55862 |   117635 |       59 |       39 |       93 |
|              groupByTqOneBuild |       10 |    39348 |    32199 |    62042 |       21 |       17 |       30 |
|          groupByTqOneBuildOpen |       10 |    45625 |    32347 |    55000 |       29 |       17 |       37 |
|              lastStartOneBuild |       10 |    43087 |    31030 |    77233 |       27 |       20 |       60 |
|  lastStartOneBuildGroupByPerTq |       10 |    96254 |    60903 |   236646 |       78 |       43 |      220 |
|                  limitOneBuild |       10 |    24106 |    20787 |    30154 |        8 |        6 |       14 |
|             limitOneBuildOneTq |       50 |    27468 |    21737 |    57488 |        9 |        6 |       35 |
|         limitOneBuildOneTqOpen |       60 |    26468 |    14717 |    61382 |        9 |        6 |       26 |


Run 5: without terminate_after param


|                          Query |      Rep |      Avg |      Min |      Max | Took Avg | Took Min | Took Max |
|                  countOneBuild |       10 |    25599 |    21464 |    29997 |       -1 |       -1 |       -1 |
|           countOneBuildCurrent |       10 |    27715 |    21760 |    39091 |       -1 |       -1 |       -1 |
|      countOneBuildCurrentField |       10 |    24045 |    22372 |    29060 |       -1 |       -1 |       -1 |
|         countOneBuildCurrentNS |       10 |    22828 |    19384 |    31875 |       -1 |       -1 |       -1 |
|    countOneBuildCurrentNSField |       10 |    22445 |    17181 |    25833 |       -1 |       -1 |       -1 |
|             countOneBuildOneTq |       50 |    23294 |    13625 |    35665 |       -1 |       -1 |       -1 |
|         countOneBuildOneTqOpen |       50 |    23936 |    18713 |    42039 |       -1 |       -1 |       -1 |
|              countOneBuildOpen |       10 |    25458 |    22102 |    30394 |       -1 |       -1 |       -1 |
|                        filling |      380 |  1410174 |    19216 |  3059614 |     1398 |        5 |     3050 |
|                   groupByBuild |       10 |  2650789 |  2555916 |  2800737 |     2635 |     2541 |     2781 |
|            groupByBuildCurrent |       10 |  3617230 |  3242976 |  4954548 |     3602 |     3226 |     4939 |
|       groupByBuildCurrentField |       10 |   979321 |   850874 |  1113454 |      962 |      837 |     1099 |
|          groupByBuildCurrentNS |       10 |  4871279 |  4673621 |  5131932 |     4856 |     4660 |     5125 |
|     groupByBuildCurrentNSField |       10 |   908126 |   746368 |  1388230 |      892 |      731 |     1372 |
|               groupByBuildOpen |       10 |    68944 |    58673 |    80235 |       53 |       44 |       63 |
|              groupByTqOneBuild |       10 |    36783 |    29913 |    40953 |       21 |       16 |       27 |
|          groupByTqOneBuildOpen |       10 |    36557 |    28119 |    56686 |       20 |       17 |       28 |
|              lastStartOneBuild |       10 |    41084 |    33440 |    53915 |       25 |       19 |       38 |
|  lastStartOneBuildGroupByPerTq |       10 |    69725 |    54318 |    81807 |       53 |       38 |       68 |
|                  limitOneBuild |       10 |    24751 |    20433 |    34740 |        9 |        6 |       20 |
|             limitOneBuildOneTq |       50 |    25967 |    13347 |    56912 |       10 |        6 |       40 |
|         limitOneBuildOneTqOpen |       60 |    24905 |    19646 |    45895 |        9 |        6 |       30 |



EXP 4: Load test

Baseline cpu: ~22%


Run 1:
RPS = 100
OpenFilterRatio         = .6
TQFilterRatio           = .1
BuildIdNonexistentRatio = .2
terminate_after=1
start: 11-28 2:17 UTC
end: 11-28 2:33 UTC

    Avg Latency: 21ms
    CPU: stabilized at ~45%
    impact on other requests: insignificant at p99 and/or avg latency


Run 2:
RPS = 300
terminate_after=1
start: 11-28 2:46 UTC
end: 11-28 2:59 UTC

    Avg Latency: 500ms, peaked at 2 sec
    CPU: stabilized at ~80%
    impact on other requests: significant slowdown, spikes of latency of 10s of seconds.


Run 3:
RPS = 100
fallback to full query: 20% (queries with expected count = 0)
preference=_shard:#num
terminate_after=1
start: 11-28 3:53 UTC
end: 11-28 4:08 UTC

    Avg Latency: 22ms for full queries, 15ms for single-shard
    CPU: stabilized at ~38%
    impact on other requests: insignificant at p99 and/or avg latency


Run 4:
RPS = 300
fallback to full query: 20% (queries with expected count = 0)
preference=_shard:#num
terminate_after=1
start: 11-28 4:24 UTC
end: 11-28 4:41 UTC

    Avg Latency: 32ms for full queries, 21ms for single-shard
    CPU: stabilized at ~63%
    impact on other requests: avg latency almost doubled but still acceptable. p99 spiked to 500ms at times



Run 5:
RPS = 500
fallback to full query: 20% (queries with expected count = 0)
preference=_shard:#num|_only_local
terminate_after=1
start: 11-28 4:54 UTC
end: 11-28 5:14 UTC

    Avg Latency: 115ms for full queries, 73ms for single-shard
    CPU: stabilized at ~85%
    impact on other requests: avg latency increased from ~30ms to 90ms, 250ms, 350ms depending on query. p99 increased to .5, 2.7, 4.3 secs depending on query


Run 6:

        RPS = 300

        OpenFilterRatio         = .6
        BuildIdNonexistentRatio = .2

        TQFilterRatio  = .3
        TQExcludeRatio = .5

        NumBuildIdsInFilter = 50
        NumTQsInFilter      = 4

|                                              Query |      Rep |      Avg |      Min |      Max |
|                             Full--All-AllTQs-Empty |      837 |  3406033 |    26523 |  9074936 |
|                         Full--All-ExcludeTQs-Empty |      171 |  3581425 |   107667 |  8594875 |
|                         Full--All-IncludeTQs-Empty |      174 |  3056710 |    77695 |  8780554 |
|                            Full--Open-AllTQs-Empty |     1230 |  3464274 |    38967 |  8892628 |
|                        Full--Open-ExcludeTQs-Empty |      276 |  3258410 |    56934 |  8758734 |
|                        Full--Open-IncludeTQs-Empty |      264 |  3688729 |    35994 |  8539955 |
|                            Shard--All-AllTQs-Empty |      970 |  1323709 |    16516 |  3826441 |
|                         Shard--All-AllTQs-NonEmpty |     3783 |  1287178 |    14793 |  3857002 |
|                        Shard--All-ExcludeTQs-Empty |      201 |  1399524 |    24339 |  3767401 |
|                     Shard--All-ExcludeTQs-NonEmpty |      852 |  1306890 |    26642 |  3838466 |
|                        Shard--All-IncludeTQs-Empty |      201 |  1256764 |    33857 |  3797404 |
|                     Shard--All-IncludeTQs-NonEmpty |      801 |  1310098 |    20609 |  3827186 |
|                           Shard--Open-AllTQs-Empty |     1423 |  1313349 |    15169 |  3889038 |
|                        Shard--Open-AllTQs-NonEmpty |     5734 |  1294483 |    14968 |  3896933 |
|                       Shard--Open-ExcludeTQs-Empty |      317 |  1252869 |    22434 |  3737847 |
|                    Shard--Open-ExcludeTQs-NonEmpty |     1217 |  1310194 |    10183 |  3757573 |
|                       Shard--Open-IncludeTQs-Empty |      311 |  1411086 |    30651 |  3818055 |
|                    Shard--Open-IncludeTQs-NonEmpty |     1244 |  1289982 |    19155 |  3860240 |

Run 7: (removed NS filter)

|                                              Query |      Rep |      Avg |      Min |      Max |
|                             Full--All-AllTQs-Empty |      740 |    48701 |    17910 |   720340 |
|                         Full--All-ExcludeTQs-Empty |      150 |    50333 |    23041 |   609157 |
|                         Full--All-IncludeTQs-Empty |      154 |    56291 |    20092 |   664619 |
|                            Full--Open-AllTQs-Empty |      665 |    98613 |    15321 |   947125 |
|                        Full--Open-ExcludeTQs-Empty |      164 |    98546 |    21010 |   663669 |
|                        Full--Open-IncludeTQs-Empty |      149 |   104119 |    19219 |   626893 |
|                            Shard--All-AllTQs-Empty |      998 |    27486 |     8861 |   368281 |
|                         Shard--All-AllTQs-NonEmpty |     3868 |    35422 |     6231 |   457989 |
|                        Shard--All-ExcludeTQs-Empty |      196 |    28801 |    12805 |   319129 |
|                     Shard--All-ExcludeTQs-NonEmpty |      836 |    36016 |    12784 |   844318 |
|                        Shard--All-IncludeTQs-Empty |      216 |    27969 |     8030 |   273989 |
|                     Shard--All-IncludeTQs-NonEmpty |      819 |    34795 |    12855 |   276240 |
|                           Shard--Open-AllTQs-Empty |     1391 |    34606 |     7044 |   433516 |
|                        Shard--Open-AllTQs-NonEmpty |     5764 |    27239 |     5055 |   371347 |
|                       Shard--Open-ExcludeTQs-Empty |      344 |    33706 |    10153 |   202477 |
|                    Shard--Open-ExcludeTQs-NonEmpty |     1207 |    29281 |     5250 |   399784 |
|                       Shard--Open-IncludeTQs-Empty |      321 |    40703 |    13089 |   367859 |
|                    Shard--Open-IncludeTQs-NonEmpty |     1265 |    26888 |     5655 |   307041 |

Run 8:  increased       NumBuildIdsInFilter = 100

|                                              Query |      Rep |      Avg |      Min |      Max |
|                             Full--All-AllTQs-Empty |      716 |    48347 |    19411 |   370856 |
|                         Full--All-ExcludeTQs-Empty |      154 |    58851 |    27538 |   278508 |
|                         Full--All-IncludeTQs-Empty |      152 |    51697 |    21347 |   413524 |
|                            Full--Open-AllTQs-Empty |      447 |    83972 |    15550 |   918048 |
|                        Full--Open-ExcludeTQs-Empty |       94 |   106765 |    20282 |   890073 |
|                        Full--Open-IncludeTQs-Empty |       92 |    83551 |    18600 |   381151 |
|                            Shard--All-AllTQs-Empty |      956 |    26118 |    12427 |   258078 |
|                         Shard--All-AllTQs-NonEmpty |     3975 |    37280 |     7697 |   321900 |
|                        Shard--All-ExcludeTQs-Empty |      215 |    27116 |    12467 |   216766 |
|                     Shard--All-ExcludeTQs-NonEmpty |      795 |    39964 |    12854 |   328221 |
|                        Shard--All-IncludeTQs-Empty |      200 |    24482 |     5600 |   109491 |
|                     Shard--All-IncludeTQs-NonEmpty |      834 |    40411 |     6992 |   374372 |
|                           Shard--Open-AllTQs-Empty |     1430 |    38070 |    12276 |   335373 |
|                        Shard--Open-AllTQs-NonEmpty |     5851 |    25342 |     8867 |   289165 |
|                       Shard--Open-ExcludeTQs-Empty |      321 |    39938 |    13213 |   274022 |
|                    Shard--Open-ExcludeTQs-NonEmpty |     1229 |    26329 |     7101 |   223636 |
|                       Shard--Open-IncludeTQs-Empty |      314 |    41620 |    13219 |   317279 |
|                    Shard--Open-IncludeTQs-NonEmpty |     1310 |    25822 |    10312 |   219713 |


Run 8:  increased       NumBuildIdsInFilter = 150

|                                              Query |      Rep |      Avg |      Min |      Max |
|                             Full--All-AllTQs-Empty |      764 |    44221 |    21476 |   339934 |
|                         Full--All-ExcludeTQs-Empty |      161 |    54910 |    24727 |   491988 |
|                         Full--All-IncludeTQs-Empty |      156 |    49220 |    23362 |   434002 |
|                            Full--Open-AllTQs-Empty |      425 |    46127 |     9141 |   429241 |
|                        Full--Open-ExcludeTQs-Empty |       81 |    70412 |    19629 |   851960 |
|                        Full--Open-IncludeTQs-Empty |       86 |    48944 |    20044 |   343368 |
|                            Shard--All-AllTQs-Empty |     1004 |    22917 |     9523 |   304770 |
|                         Shard--All-AllTQs-NonEmpty |     3959 |    35423 |     7095 |   397515 |
|                        Shard--All-ExcludeTQs-Empty |      226 |    23781 |     8947 |   150917 |
|                     Shard--All-ExcludeTQs-NonEmpty |      838 |    36254 |     9584 |   385316 |
|                        Shard--All-IncludeTQs-Empty |      203 |    24144 |    12559 |   220411 |
|                     Shard--All-IncludeTQs-NonEmpty |      859 |    35885 |    12654 |   354407 |
|                           Shard--Open-AllTQs-Empty |     1501 |    35691 |     6751 |   303194 |
|                        Shard--Open-AllTQs-NonEmpty |     6036 |    21601 |     5200 |   250254 |
|                       Shard--Open-ExcludeTQs-Empty |      316 |    38993 |    13852 |   457195 |
|                    Shard--Open-ExcludeTQs-NonEmpty |     1258 |    21580 |     5842 |   183703 |
|                       Shard--Open-IncludeTQs-Empty |      327 |    34490 |    13392 |   318844 |
|                    Shard--Open-IncludeTQs-NonEmpty |     1334 |    22205 |     5405 |   171432 |

Run 9: increased         NumTQsInFilter      = 20

|                                              Query |      Rep |      Avg |      Min |      Max |
|                             Full--All-AllTQs-Empty |      691 |   386269 |    21428 |  3256576 |
|                         Full--All-ExcludeTQs-Empty |      175 |   429575 |    35659 |  3317291 |
|                      Full--All-ExcludeTQs-NonEmpty |       32 |  2937267 |    83651 |  5262110 |
|                         Full--All-IncludeTQs-Empty |      135 |   489468 |    29585 |  3207733 |
|                      Full--All-IncludeTQs-NonEmpty |        8 |   208640 |    31883 |   762942 |
|                            Full--Open-AllTQs-Empty |      379 |   414784 |    16222 |  3207745 |
|                        Full--Open-ExcludeTQs-Empty |      103 |   458253 |    24608 |  2250656 |
|                     Full--Open-ExcludeTQs-NonEmpty |       56 |   389130 |    38076 |  3401559 |
|                        Full--Open-IncludeTQs-Empty |       82 |   498304 |    24712 |  4130777 |
|                     Full--Open-IncludeTQs-NonEmpty |       18 |   261268 |    28718 |  2010371 |
|                            Shard--All-AllTQs-Empty |      933 |   189182 |     9679 |  1576039 |
|                         Shard--All-AllTQs-NonEmpty |     3873 |   209886 |    10857 |  1687793 |
|                        Shard--All-ExcludeTQs-Empty |      219 |   169113 |    10409 |  1363252 |
|                     Shard--All-ExcludeTQs-NonEmpty |      863 |   219686 |    12622 |  1601487 |
|                        Shard--All-IncludeTQs-Empty |      197 |   210760 |    13261 |  1478189 |
|                     Shard--All-IncludeTQs-NonEmpty |      834 |   216822 |    13144 |  1592212 |
|                           Shard--Open-AllTQs-Empty |     1492 |   211907 |    12249 |  1617276 |
|                        Shard--Open-AllTQs-NonEmpty |     5964 |   187221 |     5235 |  1632468 |
|                       Shard--Open-ExcludeTQs-Empty |      311 |   191787 |    14225 |  1581659 |
|                    Shard--Open-ExcludeTQs-NonEmpty |     1293 |   184568 |    12013 |  1608568 |
|                       Shard--Open-IncludeTQs-Empty |      314 |   220812 |    13939 |  1611962 |
|                    Shard--Open-IncludeTQs-NonEmpty |     1242 |   187607 |     8702 |  1532317 |


Run 10: added the NS filter back. Could not handle for more than 10 sec, with bellow latency:

FullQ RPS: 6     Full Q Avg microSec: 3997407    ShardQ RPS: 166         ShardQ Avg microSec: 2183277
