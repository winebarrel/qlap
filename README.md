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
    	Histogram interval (default "0")
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
  "StartedAt": "2021-04-05T19:21:33.815362+09:00",
  "FinishedAt": "2021-04-05T19:21:43.831067+09:00",
  "ElapsedTime": 10,
  "MysqlConfig": {
    "User": "root",
    "Passwd": "",
    "Net": "tcp",
    "Addr": "127.0.0.1:3306",
    "DBName": "qlap",
    "Params": null,
    "Collation": "utf8mb4_general_ci",
    "Loc": {},
    "MaxAllowedPacket": 4194304,
    "ServerPubKey": "",
    "TLSConfig": "",
    "Timeout": 0,
    "ReadTimeout": 0,
    "WriteTimeout": 0,
    "AllowAllFiles": false,
    "AllowCleartextPasswords": false,
    "AllowNativePasswords": true,
    "AllowOldPasswords": false,
    "CheckConnLiveness": true,
    "ClientFoundRows": false,
    "ColumnsWithAlias": false,
    "InterpolateParams": false,
    "MultiStatements": false,
    "ParseTime": false,
    "RejectReadOnly": false
  },
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
  "Token": "20c792d5-c8f8-45b9-818f-2c151ff0eece",
  "Queries": 2929,
  "QPS": 292.44186718006193,
  "MaxQPS": 304,
  "MinQPS": 202,
  "MedianQPS": 303,
  "ExpectedQPS": 300,
  "Response": {
    "Time": {
      "Cumulative": "1.687329538s",
      "HMean": "475.205µs",
      "Avg": "576.077µs",
      "P50": "518.562µs",
      "P75": "736.735µs",
      "P95": "957.849µs",
      "P99": "1.28271ms",
      "P999": "2.803472ms",
      "Long5p": "1.23627ms",
      "Short5p": "209.029µs",
      "Max": "4.93772ms",
      "Min": "138.463µs",
      "Range": "4.799257ms",
      "StdDev": "270.592µs"
    },
    "Rate": {
      "Second": 1735.8790526904177
    },
    "Samples": 2929,
    "Count": 2929,
    "Histogram": [
      {
        "138µs - 618µs": 1769
      },
      {
        "618µs - 1.098ms": 1097
      },
      {
        "1.098ms - 1.578ms": 41
      },
      {
        "1.578ms - 2.058ms": 15
      },
      {
        "2.058ms - 2.538ms": 3
      },
      {
        "2.538ms - 3.018ms": 3
      },
      {
        "3.018ms - 3.497ms": 1
      },
      {
        "3.497ms - 3.977ms": 0
      },
      {
        "3.977ms - 4.457ms": 0
      }
    ]
  }
}
```

## DSN Examples

* https://github.com/go-sql-driver/mysql#examples

## Related Links

* https://github.com/winebarrel/qrnl
