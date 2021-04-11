# qlap

qlap is a MySQL load testing tool like [mysqlslap](https://dev.mysql.com/doc/refman/8.0/en/mysqlslap.html).

## Usage

```
qlap - MySQL load testing tool like mysqlslap.

  Flags:
       --version                               Displays the program version string.
    -h --help                                  Displays help with available flag, subcommand, and positional value parameters.
    -d --dsn                                   Data Source Name, see https://github.com/go-sql-driver/mysql#examples.
    -n --nagents                               Number of agents. (default: 1)
    -t --time                                  Test run time (sec). Zero is infinity. (default: 60)
       --number-queries                        Number of queries to execute per agent. Zero is infinity. (default: 0)
    -r --rate                                  Rate limit for each agent (qps). Zero is unlimited. (default: 0)
    -a --auto-generate-sql                     Automatically generate SQL to execute.
       --auto-generate-sql-guid-primary        Use GUID as the primary key of the table to be created.
    -q --query                                 SQL to execute.
       --auto-generate-sql-write-number        Number of rows to be pre-populated for each agent. (default: 100)
    -l --auto-generate-sql-load-type           Test load type: 'mixed', 'update', 'write', 'key', or 'read'. (default: mixed)
       --auto-generate-sql-secondary-indexes   Number of secondary indexes in the table to be created. (default: 0)
       --commit-rate                           Commit every X queries. (default: 0)
    -e --engine                                Engine of the table to be created.
    -x --number-char-cols                      Number of VARCHAR columns in the table to be created. (default: 1)
       --char-cols-index                       Create indexes on VARCHAR columns in the table to be created.
    -y --number-int-cols                       Number of INT columns in the table to be created. (default: 1)
       --int-cols-index                        Create indexes on INT columns in the table to be created.
       --pre-query                             Queries to be pre-executed for each agent.
       --create                                SQL for creating custom tables.
       --drop-db                               Forcibly delete the existing DB.
       --no-drop                               Do not drop database after testing.
       --hinterval                             Histogram interval, e.g. '100ms'. (default: 0)
    -F --delimiter                             SQL statements delimiter. (default: ;)
```

```
$ qlap -d root@/ -n 3 -r 100 -t 10 -a -l mixed -x 3 -y 3
00:10 | 3 agents / run 2727 queries (303 qps)

{
  "DSN": "root@/",
  "StartedAt": "2021-04-05T20:47:48.122543+09:00",
  "FinishedAt": "2021-04-05T20:47:58.140224+09:00",
  "ElapsedTime": 10,
  "NAgents": 3,
  "Rate": 100,
  "AutoGenerateSql": true,
  "NumberPrePopulatedData": 100,
  "DropExistingDatabase": false,
  "UseExistingDatabase": false,
  "NoDropDatabase": false,
  "Engine": "",
  "LoadType": "mixed",
  "NumberSecondaryIndexes": 0,
  "CommitRate": 0,
  "NumberIntCols": 3,
  "IntColsIndex": false,
  "NumberCharCols": 3,
  "CharColsIndex": false,
  "Queries": null,
  "PreQueries": null,
  "Token": "bf9716b4-c9c8-4539-9295-701cc46daa21",
  "QueryCount": 2930,
  "QPS": 292.4863396144265,
  "MaxQPS": 305,
  "MinQPS": 202,
  "MedianQPS": 303,
  "ExpectedQPS": 300,
  "Response": {
    "Time": {
      "Cumulative": "1.70378355s",
      "HMean": "499.554µs",
      "Avg": "581.496µs",
      "P50": "535.414µs",
      "P75": "703.491µs",
      "P95": "951.235µs",
      "P99": "1.224525ms",
      "P999": "3.198177ms",
      "Long5p": "1.242647ms",
      "Short5p": "249.39µs",
      "Max": "7.235566ms",
      "Min": "145.575µs",
      "Range": "7.089991ms",
      "StdDev": "274.566µs"
    },
    "Rate": {
      "Second": 1719.7020126177413
    },
    "Samples": 2930,
    "Count": 2930,
    "Histogram": [
      {
        "145µs - 854µs": 2610
      },
      {
        "854µs - 1.563ms": 312
      },
      {
        "1.563ms - 2.272ms": 1
      },
      {
        "2.272ms - 2.981ms": 1
      },
      {
        "2.981ms - 3.69ms": 5
      },
      {
        "3.69ms - 4.399ms": 1
      },
      {
        "4.399ms - 5.108ms": 0
      },
      {
        "5.108ms - 5.817ms": 0
      },
      {
        "5.817ms - 6.526ms": 0
      }
    ]
  }
}
```

## Use Custom Query

```
qlap -d root@/ \
  --create 'create table test (id int); insert into test values (1)' \
  -q 'select id from test; select count(id) from test'
```

## Related Links

* MySQL/PostgreSQL load testing tool using query log
    * https://github.com/winebarrel/qrn
