# qlap

qlap is a MySQL benchmarking tool like [mysqlslap](https://dev.mysql.com/doc/refman/8.0/en/mysqlslap.html).

## Usage

```
Usage of qlap:
  -auto-generate-sql
    	Automatically generate SQL to execute
  -auto-generate-sql-load-type string
    	Test load type: 'mixed', 'update', 'write', or 'key' (default "mixed")
  -auto-generate-sql-secondary-indexes int
    	Number of secondary indexes in the table to be created
  -auto-generate-sql-write-number int
    	Number of rows to be pre-populated (default 100)
  -char-cols-index
    	Create an index on the VARCHAR column if 'true'
  -commit-rate int
    	Commit every X queries
  -drop-existing-db
    	Forcibly delete the existing DB
  -dsn string
    	Data Source Name
  -engine string
    	Engine of the table to be created
  -hinterval string
    	Histogram interval, e.g. '100ms' (default "0")
  -int-cols-index
    	Create an index on the INT column if 'true'
  -nagents int
    	Number of agents (default 1)
  -number-char-cols int
    	Number of VARCHAR columns in the table to be created (default 1)
  -number-int-cols int
    	Number of INT columns in the table to be created (default 1)
  -pre-query value
    	Queries to be pre-executed for each agent
  -query string
    	SQL to execute
  -rate int
    	Rate limit for each agent (qps). Zero is unlimited
  -time int
    	Test run time (sec). Zero is infinity (default 60)
```

```
$ qlap -dsn root:@/ -nagents 3 -rate 100 -time 10 \
    -auto-generate-sql -auto-generate-sql-load-type mixed \
    -number-int-cols 3 -number-char-cols 3

00:10 | 3 agents / run 2727 queries (303 qps)

{
  "StartedAt": "2021-04-05T20:05:52.928409+09:00",
  "FinishedAt": "2021-04-05T20:06:02.944544+09:00",
  "ElapsedTime": 10,
  "NAgents": 3,
  "Time": 10000000000,
  "Rate": 100,
  "AutoGenerateSql": true,
  "NumberPrePopulatedData": 100,
  "DropExistingDatabase": false,
  "Engine": "",
  "LoadType": "mixed",
  "NumberSecondaryIndexes": 0,
  "CommitRate": 0,
  "NumberIntCols": 3,
  "IntColsIndex": false,
  "NumberCharCols": 3,
  "CharColsIndex": false,
  "Query": "",
  "PreQueries": null,
  "Token": "52a44dcf-6636-4a83-b627-00806f369e2e",
  "Queries": 2828,
  "QPS": 282.34726747872327,
  "MaxQPS": 304,
  "MinQPS": 101,
  "MedianQPS": 303,
  "ExpectedQPS": 300,
  "Response": {
    "Time": {
      "Cumulative": "1.72131028s",
      "HMean": "524.378µs",
      "Avg": "608.667µs",
      "P50": "563.136µs",
      "P75": "780.498µs",
      "P95": "956.474µs",
      "P99": "1.228807ms",
      "P999": "2.56122ms",
      "Long5p": "1.190527ms",
      "Short5p": "263.358µs",
      "Max": "7.461918ms",
      "Min": "148.608µs",
      "Range": "7.31331ms",
      "StdDev": "268.573µs"
    },
    "Rate": {
      "Second": 1642.934474312208
    },
    "Samples": 2828,
    "Count": 2828,
    "Histogram": [
      {
        "148µs - 879µs": 2470
      },
      {
        "879µs - 1.611ms": 350
      },
      {
        "1.611ms - 2.342ms": 3
      },
      {
        "2.342ms - 3.073ms": 4
      },
      {
        "3.073ms - 3.805ms": 1
      },
      {
        "3.805ms - 4.536ms": 0
      },
      {
        "4.536ms - 5.267ms": 0
      },
      {
        "5.267ms - 5.999ms": 0
      },
      {
        "5.999ms - 6.73ms": 0
      }
    ]
  }
}
```

## DSN Examples

* https://github.com/go-sql-driver/mysql#examples

## Related Links

* https://github.com/winebarrel/qrnl
