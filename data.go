package qlap

import (
	"fmt"
	"math/rand"
	"qlap/randstr"
	"strconv"
	"strings"
	"time"
)

type AutoGenerateSqlLoadType string

const (
	LoadTypeMixed  = AutoGenerateSqlLoadType("mixed")  // require prepopulated data
	LoadTypeUpdate = AutoGenerateSqlLoadType("update") // require prepopulated data
	LoadTypeWrite  = AutoGenerateSqlLoadType("write")
	LoadTypeKey    = AutoGenerateSqlLoadType("key") // require prepopulated data
	// NOTE: "read" is a full scan, so I will not implement it
	CharColLen = 128
)

type DataOpts struct {
	LoadType               AutoGenerateSqlLoadType
	NumberSecondaryIndexes int
	CommitRate             int
	NumberIntCols          int
	IntColsIndex           bool
	NumberCharCols         int
	CharColsIndex          bool
	Query                  string
	PreQueries             []string
}

type Data struct {
	*DataOpts
	randSrc   rand.Source
	idList    []int
	idIdx     int
	coin      bool
	commitCnt int
}

func newData(opts *DataOpts, idList []int) (data *Data) {
	data = &Data{
		DataOpts: opts,
		randSrc:  rand.NewSource(time.Now().UnixNano()),
		idList:   idList,
	}

	return
}

func (data *Data) initStmts() []string {
	stmts := []string{}

	if data.CommitRate > 0 {
		stmts = append(stmts, "SET autocommit = 0")
	}

	if len(data.PreQueries) > 0 {
		stmts = append(stmts, data.PreQueries...)
	}

	return stmts
}

func (data *Data) next() string {
	if data.CommitRate > 0 {
		if data.commitCnt == data.CommitRate {
			data.commitCnt = 0
			return "COMMIT"
		}

		data.commitCnt++
	}

	if data.Query != "" {
		return data.Query
	}

	switch data.LoadType {
	case LoadTypeMixed:
		coin := data.coin
		data.coin = !coin

		if coin {
			return data.buildInsertStmt()
		} else {
			return data.buildSelectByKeyStmt()
		}
	case LoadTypeUpdate:
		return data.buildUpdateStmt()
	case LoadTypeWrite:
		return data.buildInsertStmt()
	case LoadTypeKey:
		return data.buildSelectByKeyStmt()
	default:
		panic("Failed to generate SQL statement: invalid load type: " + data.LoadType)
	}
}

func (data *Data) buildCreateTableStmt() string {
	sb := strings.Builder{}
	sb.WriteString("CREATE TABLE t1 (")
	sb.WriteString("id SERIAL")

	for i := 1; i <= data.NumberSecondaryIndexes; i++ {
		fmt.Fprintf(&sb, ", id%d VARCHAR(36) UNIQUE KEY", i)
	}

	for i := 1; i <= data.NumberIntCols; i++ {
		fmt.Fprintf(&sb, ", intcol%d INT(32)", i)

		if data.IntColsIndex {
			fmt.Fprintf(&sb, ", INDEX(intcol%d)", i)
		}
	}

	for i := 1; i <= data.NumberCharCols; i++ {
		fmt.Fprintf(&sb, ", charcol%d VARCHAR(128)", i)

		if data.CharColsIndex {
			fmt.Fprintf(&sb, ", INDEX(charcol%d)", i)
		}
	}

	sb.WriteString(")")

	return sb.String()
}

func (data *Data) buildSelectByKeyStmt() string {
	sb := strings.Builder{}
	sb.WriteString("SELECT ")

	for i := 1; i <= data.NumberIntCols; i++ {
		if i >= 2 {
			sb.WriteString(",")
		}

		fmt.Fprintf(&sb, "intcol%d", i)
	}

	for i := 1; i <= data.NumberCharCols; i++ {
		if data.NumberIntCols >= 1 || i >= 2 {
			sb.WriteString(",")
		}

		fmt.Fprintf(&sb, "charcol%d", i)
	}

	sb.WriteString(" FROM t1")
	fmt.Fprintf(&sb, " WHERE id = %d", data.nextId())

	return sb.String()
}

func (data *Data) buildInsertStmt() string {
	sb := strings.Builder{}
	sb.WriteString("INSERT INTO t1 VALUES (")
	sb.WriteString("NULL") // id

	for i := 1; i <= data.NumberSecondaryIndexes; i++ {
		sb.WriteString(", UUID()")
	}

	for i := 1; i <= data.NumberIntCols; i++ {
		sb.WriteString(",")
		num := data.randSrc.Int63() >> 32
		sb.WriteString(strconv.FormatInt(num, 10))
	}

	for i := 1; i <= data.NumberCharCols; i++ {
		sb.WriteString(",'")
		sb.WriteString(randstr.String(data.randSrc, CharColLen))
		sb.WriteString("'")
	}

	sb.WriteString(")")

	return sb.String()
}

func (data *Data) buildUpdateStmt() string {
	sb := strings.Builder{}
	sb.WriteString("UPDATE t1 SET ")

	for i := 1; i <= data.NumberIntCols; i++ {
		if i >= 2 {
			sb.WriteString(",")
		}

		v := data.randSrc.Int63() >> 32
		fmt.Fprintf(&sb, "intcol%d = %d", i, v)
	}

	for i := 1; i <= data.NumberCharCols; i++ {
		if data.NumberIntCols >= 1 || i >= 2 {
			sb.WriteString(",")
		}

		fmt.Fprintf(&sb, "charcol%d = '%s'", i, randstr.String(data.randSrc, CharColLen))
	}

	fmt.Fprintf(&sb, " WHERE id = %d", data.nextId())

	return sb.String()
}

func (data *Data) nextId() int {
	if data.idIdx >= len(data.idList) {
		data.idIdx = 0
	}

	return data.idList[data.idIdx]
}
