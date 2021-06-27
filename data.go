package qlap

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/winebarrel/randstr"
)

type AutoGenerateSqlLoadType string

const (
	LoadTypeMixed         = AutoGenerateSqlLoadType("mixed")  // require pre-populated data
	LoadTypeUpdate        = AutoGenerateSqlLoadType("update") // require pre-populated data
	LoadTypeWrite         = AutoGenerateSqlLoadType("write")
	LoadTypeKey           = AutoGenerateSqlLoadType("key")  // require pre-populated data
	LoadTypeRead          = AutoGenerateSqlLoadType("read") // require pre-populated data
	AutoGenerateTableName = "t1"
)

type DataOpts struct {
	LoadType               AutoGenerateSqlLoadType
	GuidPrimary            bool
	NumberSecondaryIndexes int
	CommitRate             int
	MixedSelRatio          int
	MixedInsRatio          int
	NumberIntCols          int
	IntColsIndex           bool
	NumberCharCols         int
	CharColsIndex          bool
	Queries                []string `json:"-"`
	PreQueries             []string
}

type Data struct {
	*DataOpts
	randSrc   rand.Source
	idList    []string
	idIdx     int
	mixedIdx  int
	commitCnt int
	queryIdx  int
}

func newData(opts *DataOpts, idList []string) (data *Data) {
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

	if len(data.Queries) > 0 {
		q := data.Queries[data.queryIdx]
		data.queryIdx++

		if data.queryIdx == len(data.Queries) {
			data.queryIdx = 0
		}

		return q
	}

	switch data.LoadType {
	case LoadTypeMixed:
		var stmt string
		if data.mixedIdx < data.MixedSelRatio {
			stmt = data.buildSelectStmt(true)
		} else {
			stmt = data.buildInsertStmt()
		}

		data.mixedIdx++

		if data.mixedIdx >= data.MixedSelRatio+data.MixedInsRatio {
			data.mixedIdx = 0
		}

		return stmt
	case LoadTypeUpdate:
		return data.buildUpdateStmt()
	case LoadTypeWrite:
		return data.buildInsertStmt()
	case LoadTypeKey:
		return data.buildSelectStmt(true)
	case LoadTypeRead:
		return data.buildSelectStmt(false)
	default:
		panic("Failed to generate SQL statement: invalid load type: " + data.LoadType)
	}
}

func (data *Data) buildCreateTableStmt() string {
	sb := strings.Builder{}
	sb.WriteString("CREATE TABLE " + AutoGenerateTableName + " (id ")

	if data.GuidPrimary {
		sb.WriteString("VARCHAR(36) PRIMARY KEY")
	} else {
		sb.WriteString("SERIAL")
	}

	for i := 1; i <= data.NumberSecondaryIndexes; i++ {
		fmt.Fprintf(&sb, ",id%d VARCHAR(36) UNIQUE KEY", i)
	}

	for i := 1; i <= data.NumberIntCols; i++ {
		fmt.Fprintf(&sb, ",intcol%d INT(32)", i)

		if data.IntColsIndex {
			fmt.Fprintf(&sb, ",INDEX(intcol%d)", i)
		}
	}

	for i := 1; i <= data.NumberCharCols; i++ {
		fmt.Fprintf(&sb, ",charcol%d VARCHAR(128)", i)

		if data.CharColsIndex {
			fmt.Fprintf(&sb, ",INDEX(charcol%d)", i)
		}
	}

	sb.WriteString(")")

	return sb.String()
}

func (data *Data) buildSelectStmt(key bool) string {
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

	sb.WriteString(" FROM " + AutoGenerateTableName)

	if key {
		fmt.Fprintf(&sb, " WHERE id = '%s'", data.nextId())
	}

	return sb.String()
}

func (data *Data) buildInsertStmt() string {
	sb := strings.Builder{}
	sb.WriteString("INSERT INTO " + AutoGenerateTableName + " VALUES (")

	if data.GuidPrimary {
		sb.WriteString("UUID()")
	} else {
		sb.WriteString("NULL")
	}

	for i := 1; i <= data.NumberSecondaryIndexes; i++ {
		sb.WriteString(",UUID()")
	}

	for i := 1; i <= data.NumberIntCols; i++ {
		sb.WriteString(",")
		num := data.randSrc.Int63() >> 32
		sb.WriteString(strconv.FormatInt(num, 10))
	}

	for i := 1; i <= data.NumberCharCols; i++ {
		sb.WriteString(",'")
		sb.WriteString(randstr.String(data.randSrc, 128))
		sb.WriteString("'")
	}

	sb.WriteString(")")

	return sb.String()
}

func (data *Data) buildUpdateStmt() string {
	sb := strings.Builder{}
	sb.WriteString("UPDATE " + AutoGenerateTableName + " SET ")

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

		fmt.Fprintf(&sb, "charcol%d = '%s'", i, randstr.String(data.randSrc, 128))
	}

	fmt.Fprintf(&sb, " WHERE id = '%s'", data.nextId())

	return sb.String()
}

func (data *Data) nextId() string {
	if data.idIdx >= len(data.idList) {
		data.idIdx = 0
	}

	id := data.idList[data.idIdx]
	data.idIdx++

	return id
}
