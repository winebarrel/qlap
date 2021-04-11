package qlap

import (
	"context"
	"database/sql"
	"fmt"
	"math/rand"
	"time"
)

const (
	RecordPeriod = 1 * time.Second
)

type Agent struct {
	id          int
	mysqlConfig *MysqlConfig
	db          *sql.DB
	taskOps     *TaskOpts
	dataOpts    *DataOpts
	data        *Data
}

func newAgent(id int, myCfg *MysqlConfig, taskOps *TaskOpts, dataOpts *DataOpts) (agent *Agent) {
	agent = &Agent{
		id:          id,
		mysqlConfig: myCfg,
		taskOps:     taskOps,
		dataOpts:    dataOpts,
	}

	return
}

func (agent *Agent) prepare(maxIdleConns int, idList []string) error {
	db, err := agent.mysqlConfig.openAndPing(maxIdleConns)

	if err != nil {
		dsn := agent.mysqlConfig.FormatDSN()
		return fmt.Errorf("Failed to open/ping DB (agent id=%d, dsn=%s): %w", agent.id, dsn, err)
	}

	agent.db = db

	newIdList := make([]string, len(idList))
	copy(newIdList, idList)
	rand.Shuffle(len(newIdList), func(i, j int) { newIdList[i], newIdList[j] = newIdList[j], newIdList[i] })
	agent.data = newData(agent.dataOpts, newIdList)

	inits := agent.data.initStmts()

	for _, stmt := range inits {
		_, err = db.Exec(stmt)

		if err != nil {
			return fmt.Errorf("Failed to execute initial query (agent id=%d, query=%s): %w", agent.id, stmt, err)
		}
	}

	return nil
}

func (agent *Agent) run(ctx context.Context, recorder *Recorder, token string) error {
	_, err := agent.db.Exec(fmt.Sprintf("SELECT 'agent(%d) start: token=%s'", agent.id, token))

	if err != nil {
		return fmt.Errorf("Failed to execute start query (agent id=%d): %w", agent.id, err)
	}

	recordTick := time.NewTicker(RecordPeriod)
	defer recordTick.Stop()
	recDps := []recorderDataPoint{}

	err = loopWithThrottle(agent.taskOps.Rate, func(i int) (bool, error) {
		if agent.taskOps.NumberQueriesToExecute > 0 && i >= agent.taskOps.NumberQueriesToExecute {
			return false, nil
		}

		select {
		case <-ctx.Done():
			return false, nil
		case <-recordTick.C:
			recorder.add(recDps)
			recDps = recDps[:0]
		default:
			// Nothing to do
		}

		q := agent.data.next()
		rt, err := agent.query(q)

		if err != nil {
			return false, fmt.Errorf("Execute query error (query=%s): %w", q, err)
		}

		recDps = append(recDps, recorderDataPoint{
			timestamp: time.Now(),
			resTime:   rt,
		})

		return true, nil
	})

	if err != nil {
		return fmt.Errorf("Failed to transact (agent id=%d): %w", agent.id, err)
	}

	_, err = agent.db.Exec(fmt.Sprintf("SELECT 'agent(%d) end: token=%s'", agent.id, token))

	if err != nil {
		return fmt.Errorf("Failed to execute exit query (agent id=%d): %w", agent.id, err)
	}

	return nil
}

func (agent *Agent) close() error {
	err := agent.db.Close()

	if err != nil {
		return fmt.Errorf("Failed to close DB (agent id=%d): %w", agent.id, err)
	}

	return nil
}

func (agent *Agent) query(q string) (time.Duration, error) {
	start := time.Now()
	_, err := agent.db.Exec(q)
	end := time.Now()

	if err != nil {
		return 0, err
	}

	return end.Sub(start), nil
}
