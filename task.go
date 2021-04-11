package qlap

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
	"golang.org/x/term"
)

const (
	ProgressReportPeriod = 1
)

type TaskOpts struct {
	MysqlConfig            *MysqlConfig `json:"-"`
	NAgents                int
	Time                   time.Duration `json:"-"`
	Rate                   int
	AutoGenerateSql        bool
	NumberPrePopulatedData int
	NumberQueriesToExecute int
	DropExistingDatabase   bool
	UseExistingDatabase    bool
	NoDropDatabase         bool
	Engine                 string
	Creates                []string
	OnlyPrint              bool
}

type Task struct {
	*TaskOpts
	agents   []*Agent
	dataOpts *DataOpts
	recOpts  *RecorderOpts
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func NewTask(taskOpts *TaskOpts, dataOpts *DataOpts, recOpts *RecorderOpts) (task *Task) {
	agents := make([]*Agent, taskOpts.NAgents)

	for i := 0; i < taskOpts.NAgents; i++ {
		agents[i] = newAgent(i, taskOpts.MysqlConfig, taskOpts, dataOpts)
	}

	task = &Task{
		TaskOpts: taskOpts,
		agents:   agents,
		dataOpts: dataOpts,
		recOpts:  recOpts,
	}

	return
}

func (task *Task) Prepare() error {
	idList, err := task.setupDB()

	if err != nil {
		return fmt.Errorf("Failed to setup DB: %w", err)
	}

	for _, agent := range task.agents {
		if err := agent.prepare(task.NAgents, idList); err != nil {
			return fmt.Errorf("Failed to prepare Agent: %w", err)
		}
	}

	return nil
}

func (task *Task) setupDB() ([]string, error) {
	// Temporarily empty the DB name
	orgDBName := task.MysqlConfig.DBName
	task.MysqlConfig.DBName = ""

	db, err := task.MysqlConfig.openAndPing(1)

	if err != nil {
		return nil, fmt.Errorf("Connection error: %w", err)
	}

	if task.Engine != "" {
		_, err = db.Exec(fmt.Sprintf("SET default_storage_engine = %s", task.Engine))

		if err != nil {
			return nil, fmt.Errorf("Set default_storage_engine error: %w", err)
		}
	}

	defer db.Close()
	task.MysqlConfig.DBName = orgDBName

	if task.DropExistingDatabase {
		_, err = db.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS `%s`", task.MysqlConfig.DBName))

		if err != nil {
			return nil, fmt.Errorf("Drop database error: %w", err)
		}
	}

	row := db.QueryRow(fmt.Sprintf("SELECT COUNT(1) FROM information_schema.SCHEMATA WHERE SCHEMA_NAME = '%s'", task.MysqlConfig.DBName))
	var dbCnt int

	if _, ok := db.(*NullDB); ok {
		err = nil
		dbCnt = 0
	} else {
		err = row.Scan(&dbCnt)
	}

	if err != nil {
		return nil, fmt.Errorf("Database existence check error: %w", err)
	}

	if dbCnt < 1 {
		_, err = db.Exec(fmt.Sprintf("CREATE DATABASE `%s`", task.MysqlConfig.DBName))

		if err != nil {
			return nil, fmt.Errorf("Create database error: %w", err)
		}
	} else {
		task.UseExistingDatabase = true
	}

	_, err = db.Exec(fmt.Sprintf("USE `%s`", task.MysqlConfig.DBName))

	if err != nil {
		return nil, fmt.Errorf("Use database error: %w", err)
	}

	if len(task.Creates) > 0 {
		for _, stmt := range task.Creates {
			_, err = db.Exec(stmt)

			if err != nil {
				return nil, fmt.Errorf("Create table error (query=%s): %w", stmt, err)
			}
		}

		return []string{}, nil
	}

	_, err = db.Exec("DROP TABLE IF EXISTS " + AutoGenerateTableName)

	if err != nil {
		return nil, fmt.Errorf("Drop table error: %w", err)
	}

	tblStmt := newData(task.dataOpts, nil).buildCreateTableStmt()
	_, err = db.Exec(tblStmt)

	if err != nil {
		return nil, fmt.Errorf("Create table error (query=%s): %w", tblStmt, err)
	}

	ctxWithoutCancel := context.Background()
	ctx, cancel := context.WithCancel(ctxWithoutCancel)
	eg := task.prePopulateData(ctx)
	task.trapSigint(ctx, cancel, eg)
	err = eg.Wait()
	cancel()

	if err != nil {
		return nil, fmt.Errorf("Pre-populate data error: %w", err)
	}

	idList := make([]string, task.NumberPrePopulatedData*task.NAgents)
	rs, err := db.Query("SELECT id FROM t1")

	if _, ok := db.(*NullDB); ok {
		return idList, nil
	}

	if err != nil {
		return nil, fmt.Errorf("Ftech id error: %w", err)
	}

	for i := 0; rs.Next(); i++ {
		err = rs.Scan(&idList[i])

		if err != nil {
			return nil, fmt.Errorf("Scan id error: %w", err)
		}
	}

	return idList, nil
}

func (task *Task) prePopulateData(ctx context.Context) *errgroup.Group {
	eg, ctx := errgroup.WithContext(ctx)

	for i := 0; i < task.NAgents; i++ {
		eg.Go(func() error {
			data := newData(task.dataOpts, nil)
			db, err := task.MysqlConfig.openAndPing(1)

			if err != nil {
				return fmt.Errorf("Connection error: %w", err)
			}

			defer db.Close()

			for i := 0; i < task.NumberPrePopulatedData; i++ {
				select {
				case <-ctx.Done():
					return nil
				default:
					insStmt := data.buildInsertStmt()
					_, err = db.Exec(insStmt)

					if err != nil {
						return fmt.Errorf("Insert error (query=%s): %w", insStmt, err)
					}
				}
			}

			return nil
		})
	}

	return eg
}

func (task *Task) Run() (*Recorder, error) {
	uuid, _ := uuid.NewRandom()
	token := uuid.String()
	rec := newRecorder(task.recOpts, task.TaskOpts, task.dataOpts, token)

	defer func() {
		rec.close()

		for _, agent := range task.agents {
			err := agent.close()

			if err != nil {
				fmt.Fprintf(os.Stderr, "[WARN] Failed to cloge Agent: %s", err)
			}
		}
	}()

	eg, ctxWithoutCancel := errgroup.WithContext(context.Background())
	ctx, cancel := context.WithCancel(ctxWithoutCancel)
	progressTick := time.NewTicker(ProgressReportPeriod * time.Second)
	rec.start(task.NAgents * 3)
	var numTermAgents int32

	// Variables for progress line
	taskStart := time.Now()
	prevExecCnt := 0

	// Run agents
	for _, v := range task.agents {
		agent := v
		eg.Go(func() error {
			err := agent.run(ctx, rec, token)
			atomic.AddInt32(&numTermAgents, 1)
			return err
		})
	}

	// Periodic report progress
	go func() {
	LOOP:
		for {
			select {
			case <-ctx.Done():
				progressTick.Stop()
				break LOOP
			case <-progressTick.C:
				if !task.OnlyPrint {
					execCnt := rec.Count()
					termAgentCnt := int(atomic.LoadInt32(&numTermAgents))
					task.printProgress(execCnt, prevExecCnt, taskStart, termAgentCnt)
					prevExecCnt = execCnt
				}
			}
		}
	}()

	// Time-out processing
	// NOTE: If it is zero, it will not time out
	if task.Time > 0 {
		go func() {
			select {
			case <-ctx.Done():
				// Nothing to do
			case <-time.After(task.Time):
				cancel()
			}
		}()
	}

	task.trapSigint(ctx, cancel, eg)
	err := eg.Wait()
	cancel()

	// Clear progress line
	if !rec.OnlyPrint {
		fmt.Fprintf(os.Stderr, "\r\n\n")
	}

	if err != nil {
		return nil, fmt.Errorf("Error during agent running: %w", err)
	}

	return rec, nil
}

func (task *Task) Close() error {
	err := task.teardownDB()

	if err != nil {
		return fmt.Errorf("Failed to teardown DB: %w", err)
	}

	return nil
}

func (task *Task) teardownDB() error {
	if !task.NoDropDatabase && !task.UseExistingDatabase {
		db, err := task.MysqlConfig.openAndPing(1)

		if err != nil {
			return fmt.Errorf("Connection error: %w", err)
		}

		defer db.Close()
		_, err = db.Exec(fmt.Sprintf("DROP DATABASE `%s`", task.MysqlConfig.DBName))

		if err != nil {
			return fmt.Errorf("Drop database error: %w", err)
		}
	}

	return nil
}

func (task *Task) printProgress(execCnt int, prevExecCnt int, taskStart time.Time, numTermAgents int) {
	qps := float64(execCnt-prevExecCnt) / ProgressReportPeriod
	elapsedTime := time.Since(taskStart)
	numRunAgents := task.NAgents - int(numTermAgents)
	termWidth, _, err := term.GetSize(0)

	if err != nil {
		panic("Failed to get terminal width: " + err.Error())
	}

	elapsedTimeSec := elapsedTime.Round(time.Second)
	min := elapsedTimeSec / time.Minute
	sec := (elapsedTimeSec - min*time.Minute) / time.Second
	progressLine := fmt.Sprintf("%02d:%02d | %d agents / run %d queries (%.0f qps)", min, sec, numRunAgents, execCnt, qps)
	fmt.Fprintf(os.Stderr, "\r%-*s", termWidth, progressLine)
}

func (task *Task) trapSigint(ctx context.Context, cancel context.CancelFunc, eg *errgroup.Group) {
	// SIGINT
	sgnlCh := make(chan os.Signal, 1)
	signal.Notify(sgnlCh, os.Interrupt)

	go func() {
		select {
		case <-ctx.Done():
			// Nothing to do
		case <-sgnlCh:
			cancel()
			_ = eg.Wait()
			_ = task.teardownDB()
			os.Exit(130)
		}
	}()
}
