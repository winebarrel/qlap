package main

import (
	"flag"
	"fmt"
	"os"
	"qlap"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
)

const (
	DefaultTime                   = 60
	DefaultDBName                 = "qlap"
	DefaultNumberPrePopulatedData = 100
	DefaultLoadType               = string(qlap.LoadTypeMixed)
	DefaultNumberIntCols          = 1
	DefaultNumberCharCols         = 1
)

type Flags struct {
	qlap.TaskOpts
	qlap.DataOpts
	qlap.RecorderOpts
}

func parseFlags() (flags *Flags) {
	flags = &Flags{}
	dsn := flag.String("dsn", "", "Data Source Name, see https://github.com/go-sql-driver/mysql#examples")
	flag.IntVar(&flags.NAgents, "nagents", 1, "Number of agents")
	argTime := flag.Int("time", DefaultTime, "Test run time (sec). Zero is infinity")
	flag.IntVar(&flags.Rate, "rate", 0, "Rate limit for each agent (qps). Zero is unlimited")
	flag.BoolVar(&flags.AutoGenerateSql, "auto-generate-sql", false, "Automatically generate SQL to execute")
	queries := flag.String("query", "", "SQL to execute")
	flag.IntVar(&flags.NumberPrePopulatedData, "auto-generate-sql-write-number", DefaultNumberPrePopulatedData, "Number of rows to be pre-populated for each agent")
	strLoadType := flag.String("auto-generate-sql-load-type", DefaultLoadType, "Test load type: 'mixed', 'update', 'write', or 'key'")
	flag.IntVar(&flags.NumberSecondaryIndexes, "auto-generate-sql-secondary-indexes", 0, "Number of secondary indexes in the table to be created")
	flag.IntVar(&flags.CommitRate, "commit-rate", 0, "Commit every X queries")
	flag.StringVar(&flags.Engine, "engine", "", "Engine of the table to be created")
	flag.IntVar(&flags.NumberIntCols, "number-int-cols", DefaultNumberIntCols, "Number of INT columns in the table to be created")
	flag.BoolVar(&flags.IntColsIndex, "int-cols-index", false, "Create indexes on INT columns in the table to be created")
	flag.IntVar(&flags.NumberCharCols, "number-char-cols", DefaultNumberCharCols, "Number of VARCHAR columns in the table to be created")
	flag.BoolVar(&flags.CharColsIndex, "char-cols-index", false, "Create indexes on VARCHAR columns in the table to be created")
	preqs := flag.String("pre-query", "", "Queries to be pre-executed for each agent")
	creates := flag.String("create", "", "SQL for creating custom tables")
	flag.BoolVar(&flags.DropExistingDatabase, "drop-db", false, "Forcibly delete the existing DB")
	hinterval := flag.String("hinterval", "0", "Histogram interval, e.g. '100ms'")
	delimiter := flag.String("delimiter", ";", "SQL statements delimiter")
	flag.Parse()

	if flag.NFlag() == 0 {
		printUsageAndExit()
	}

	// DSN
	if *dsn == "" {
		printErrorAndExit("'-dsn' is required")
	}

	myCfg, err := mysql.ParseDSN(*dsn)

	if err != nil {
		printErrorAndExit("DSN parsing error: " + err.Error())
	}

	flags.DSN = *dsn

	if myCfg.DBName == "" {
		myCfg.DBName = DefaultDBName
	}

	flags.MysqlConfig = &qlap.MysqlConfig{Config: myCfg}

	// NAgents
	if flags.NAgents < 1 {
		printErrorAndExit("'-nagents' must be >= 1")
	}

	// Time
	if *argTime < 0 {
		printErrorAndExit("'-time' must be >= 0")
	}

	flags.Time = time.Duration(*argTime) * time.Second

	// Rate
	if flags.Rate < 0 {
		printErrorAndExit("'-rate' must be >= 0")
	}

	// Delimiter
	if *delimiter == "" {
		printErrorAndExit("'-delimiter' must not be empty")
	}

	// AutoGenerateSql / Queries
	if !flags.AutoGenerateSql && *queries == "" {
		printErrorAndExit("Either '-auto-generate-sql' or '-query' is required")
	} else if flags.AutoGenerateSql && *queries != "" {
		printErrorAndExit("Cannot set both '-auto-generate-sql' and '-query'")
	}

	// Queries
	if *queries != "" {
		flags.Queries = strings.Split(*queries, *delimiter)
	}

	// Creates
	if *creates != "" {
		if *queries == "" {
			printErrorAndExit("'-query' is required for '-create'")
		}

		flags.Creates = strings.Split(*creates, *delimiter)
	}

	// NumberPrePopulatedData
	if flags.NumberPrePopulatedData < 0 {
		printErrorAndExit("'-auto-generate-sql-write-number' must be >= 0")
	}

	// LoadType
	loadType := qlap.AutoGenerateSqlLoadType(*strLoadType)

	if loadType != qlap.LoadTypeMixed &&
		loadType != qlap.LoadTypeUpdate &&
		loadType != qlap.LoadTypeWrite &&
		loadType != qlap.LoadTypeKey {
		printErrorAndExit("Invalid load type: " + *strLoadType)
	}

	if flags.NumberPrePopulatedData == 0 && (loadType == qlap.LoadTypeMixed ||
		loadType == qlap.LoadTypeUpdate ||
		loadType == qlap.LoadTypeKey) {
		printErrorAndExit("Pre-populated data is required for 'mixed', 'update', and 'key'")
	}

	flags.LoadType = loadType

	// NumberSecondaryIndexes
	if flags.NumberSecondaryIndexes < 0 {
		printErrorAndExit("'-auto-generate-sql-secondary-indexes' must be >= 0")
	}

	// CommitRate
	if flags.CommitRate < 0 {
		printErrorAndExit("'-commit-rate' must be >= 0")
	}

	// NumberIntCols
	if flags.NumberIntCols < 1 {
		printErrorAndExit("'-number-int-cols' must be >= 1")
	}

	// NumberCharCols
	if flags.NumberCharCols < 1 {
		printErrorAndExit("'-number-char-cols' must be >= 1")
	}

	// PreQueries
	if *preqs != "" {
		flags.PreQueries = strings.Split(*preqs, *delimiter)
	}

	// HInterval
	if hi, err := time.ParseDuration(*hinterval); err != nil {
		printErrorAndExit("Failed to parse hinterval: " + err.Error())
	} else {
		flags.HInterval = hi
	}

	return
}

func printUsageAndExit() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	flag.PrintDefaults()
	os.Exit(2)
}

func printErrorAndExit(msg string) {
	fmt.Fprintln(os.Stderr, msg)
	os.Exit(1)
}
