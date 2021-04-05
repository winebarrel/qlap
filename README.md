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
$ qlap -dsn root:@/ -nagents 3 -time 10 \
    -auto-generate-sql -auto-generate-sql-load-type mixed \
    -number-int-cols 3 -number-char-cols 3

00:10 | 3 agents / run 158675 queries (17704 qps)

{
  "StartedAt": "2021-04-05T19:11:45.309493+09:00",
  "FinishedAt": "2021-04-05T19:11:55.309901+09:00",
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
  "Rate": 0,
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
  "Token": "63b2933f-ac01-480a-9bf6-63047bcacd1d",
  "Queries": 158675,
  "QPS": 15866.91373657773,
  "MaxQPS": 17919,
  "MinQPS": 17347,
  "MedianQPS": 17716,
  "ExpectedQPS": 0,
  "Response": {
    "Time": {
      "Cumulative": "26.437103022s",
      "HMean": "139.583µs",
      "Avg": "166.611µs",
      "P50": "152.834µs",
      "P75": "216.269µs",
      "P95": "287.631µs",
      "P99": "336.278µs",
      "P999": "596.327µs",
      "Long5p": "340.025µs",
      "Short5p": "81.759µs",
      "Max": "7.417935ms",
      "Min": "68.841µs",
      "Range": "7.349094ms",
      "StdDev": "89.355µs"
    },
    "Rate": {
      "Second": 6001.98137700475
    },
    "Samples": 158675,
    "Count": 158675,
    "Histogram": [
      {
        "68µs - 803µs": 158564
      },
      {
        "803µs - 1.538ms": 74
      },
      {
        "1.538ms - 2.273ms": 11
      },
      {
        "2.273ms - 3.008ms": 14
      },
      {
        "3.008ms - 3.743ms": 1
      },
      {
        "3.743ms - 4.478ms": 2
      },
      {
        "4.478ms - 5.213ms": 5
      },
      {
        "5.213ms - 5.948ms": 1
      },
      {
        "5.948ms - 6.683ms": 1
      },
      {
        "6.683ms - 7.417ms": 2
      }
    ]
  }
}
```

## DSN Examples

* https://github.com/go-sql-driver/mysql#examples

## Related Links

* https://github.com/winebarrel/qrnl
