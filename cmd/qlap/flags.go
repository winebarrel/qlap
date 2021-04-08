package main

import (
	"fmt"
	"os"
	"qlap"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/integrii/flaggy"
)

var version string

const (
	DefaultTime                   = 60
	DefaultDBName                 = "qlap"
	DefaultNumberPrePopulatedData = 100
	DefaultLoadType               = string(qlap.LoadTypeMixed)
	DefaultNumberIntCols          = 1
	DefaultNumberCharCols         = 1
	DefaultDelimiter              = ";"
)

type Flags struct {
	qlap.TaskOpts
	qlap.DataOpts
	qlap.RecorderOpts
}

func parseFlags() (flags *Flags) {
	flaggy.SetVersion(version)
	flaggy.SetDescription("MySQL load testing tool like mysqlslap.")
	flags = &Flags{}
	var dsn string
	flaggy.String(&dsn, "d", "dsn", "Data Source Name, see https://github.com/go-sql-driver/mysql#examples.")
	flags.NAgents = 1
	flaggy.Int(&flags.NAgents, "n", "nagents", "Number of agents.")
	argTime := DefaultTime
	flaggy.Int(&argTime, "t", "time", "Test run time (sec). Zero is infinity.")
	flaggy.Int(&flags.Rate, "r", "rate", "Rate limit for each agent (qps). Zero is unlimited.")
	flaggy.Bool(&flags.AutoGenerateSql, "a", "auto-generate-sql", "Automatically generate SQL to execute.")
	var queries string
	flaggy.String(&queries, "q", "query", "SQL to execute.")
	flags.NumberPrePopulatedData = DefaultNumberPrePopulatedData
	flaggy.Int(&flags.NumberPrePopulatedData, "", "auto-generate-sql-write-number", "Number of rows to be pre-populated for each agent.")
	strLoadType := DefaultLoadType
	flaggy.String(&strLoadType, "", "auto-generate-sql-load-type", "Test load type: 'mixed', 'update', 'write', 'key', or 'read'.")
	flaggy.Int(&flags.NumberSecondaryIndexes, "", "auto-generate-sql-secondary-indexes", "Number of secondary indexes in the table to be created.")
	flaggy.Int(&flags.CommitRate, "", "commit-rate", "Commit every X queries.")
	flaggy.String(&flags.Engine, "e", "engine", "Engine of the table to be created.")
	flags.NumberCharCols = DefaultNumberCharCols
	flaggy.Int(&flags.NumberCharCols, "x", "number-char-cols", "Number of VARCHAR columns in the table to be created.")
	flaggy.Bool(&flags.CharColsIndex, "", "char-cols-index", "Create indexes on VARCHAR columns in the table to be created.")
	flags.NumberIntCols = DefaultNumberIntCols
	flaggy.Int(&flags.NumberIntCols, "y", "number-int-cols", "Number of INT columns in the table to be created.")
	flaggy.Bool(&flags.IntColsIndex, "", "int-cols-index", "Create indexes on INT columns in the table to be created.")
	var preqs string
	flaggy.String(&preqs, "", "pre-query", "Queries to be pre-executed for each agent.")
	var creates string
	flaggy.String(&creates, "", "create", "SQL for creating custom tables.")
	flaggy.Bool(&flags.DropExistingDatabase, "", "drop-db", "Forcibly delete the existing DB.")
	flaggy.Bool(&flags.NoDropDatabase, "", "no-drop", "Do not drop database after testing.")
	hinterval := "0"
	flaggy.String(&hinterval, "", "hinterval", "Histogram interval, e.g. '100ms'.")
	delimiter := DefaultDelimiter
	flaggy.String(&delimiter, "F", "delimiter", "SQL statements delimiter.")
	flaggy.Parse()

	if len(os.Args) <= 1 {
		flaggy.ShowHelpAndExit("")
	}

	// DSN
	if dsn == "" {
		printErrorAndExit("'--dsn(-d)' is required")
	}

	myCfg, err := mysql.ParseDSN(dsn)

	if err != nil {
		printErrorAndExit("DSN parsing error: " + err.Error())
	}

	flags.DSN = dsn

	if myCfg.DBName == "" {
		myCfg.DBName = DefaultDBName
	}

	flags.MysqlConfig = &qlap.MysqlConfig{Config: myCfg}

	// NAgents
	if flags.NAgents < 1 {
		printErrorAndExit("'--nagents(-n)' must be >= 1")
	}

	// Time
	if argTime < 0 {
		printErrorAndExit("'--time(-t)' must be >= 0")
	}

	flags.Time = time.Duration(argTime) * time.Second

	// Rate
	if flags.Rate < 0 {
		printErrorAndExit("'--rate(-r)' must be >= 0")
	}

	// Delimiter
	if delimiter == "" {
		printErrorAndExit("'--delimiter(-F)' must not be empty")
	}

	// AutoGenerateSql / Queries
	if !flags.AutoGenerateSql && queries == "" {
		printErrorAndExit("Either '--auto-generate-sql(-a)' or '--query(-q)' is required")
	} else if flags.AutoGenerateSql && queries != "" {
		printErrorAndExit("Cannot set both '--auto-generate-sql(-a)' and '--query(-q)'")
	}

	// Queries
	if queries != "" {
		flags.Queries = strings.Split(queries, delimiter)
	}

	// Creates
	if creates != "" {
		if queries == "" {
			printErrorAndExit("'--query(-q)' is required for '--create'")
		}

		flags.Creates = strings.Split(creates, delimiter)
	}

	// NumberPrePopulatedData
	if flags.NumberPrePopulatedData < 0 {
		printErrorAndExit("'--auto-generate-sql-write-number' must be >= 0")
	}

	// LoadType
	loadType := qlap.AutoGenerateSqlLoadType(strLoadType)

	if loadType != qlap.LoadTypeMixed &&
		loadType != qlap.LoadTypeUpdate &&
		loadType != qlap.LoadTypeWrite &&
		loadType != qlap.LoadTypeKey &&
		loadType != qlap.LoadTypeRead {
		printErrorAndExit("Invalid load type: " + strLoadType)
	}

	if flags.NumberPrePopulatedData == 0 && (loadType == qlap.LoadTypeMixed ||
		loadType == qlap.LoadTypeUpdate ||
		loadType == qlap.LoadTypeKey ||
		loadType == qlap.LoadTypeRead) {
		printErrorAndExit("Pre-populated data is required for 'mixed', 'update', 'key', and 'read'")
	}

	flags.LoadType = loadType

	// NumberSecondaryIndexes
	if flags.NumberSecondaryIndexes < 0 {
		printErrorAndExit("'--auto-generate-sql-secondary-indexes' must be >= 0")
	}

	// CommitRate
	if flags.CommitRate < 0 {
		printErrorAndExit("'--commit-rate' must be >= 0")
	}

	// NumberIntCols
	if flags.NumberIntCols < 1 {
		printErrorAndExit("'--number-int-cols(-y)' must be >= 1")
	}

	// NumberCharCols
	if flags.NumberCharCols < 1 {
		printErrorAndExit("'--number-char-cols(-x)' must be >= 1")
	}

	// PreQueries
	if preqs != "" {
		flags.PreQueries = strings.Split(preqs, delimiter)
	}

	// HInterval
	if hi, err := time.ParseDuration(hinterval); err != nil {
		printErrorAndExit("Failed to parse hinterval: " + err.Error())
	} else {
		flags.HInterval = hi
	}

	return
}

func printErrorAndExit(msg string) {
	fmt.Fprintln(os.Stderr, msg)
	os.Exit(1)
}
