package models

import "time"

// DagRun 表示DAG的一次运行
type DagRun struct {
	ID                     int64     `db:"id"`
	DagID                  string    `db:"dag_id"`
	ExecutionDate          time.Time `db:"execution_date"`
	State                  string    `db:"state"`
	RunID                  string    `db:"run_id"`
	ExternalTrigger        bool      `db:"external_trigger"`
	Conf                   []byte    `db:"conf"`
	StartDate              time.Time `db:"start_date"`
	EndDate                time.Time `db:"end_date"`
	DataIntervalStart      time.Time `db:"data_interval_start"`
	DataIntervalEnd        time.Time `db:"data_interval_end"`
	LastSchedulingDecision time.Time `db:"last_scheduling_decision"`
	RunType                string    `db:"run_type"`
	CreatingJobID          int64     `db:"creating_job_id"`
	DagHash                string    `db:"dag_hash"`
	UpdatedAt              time.Time `db:"updated_at"`
	QueuedAt               time.Time `db:"queued_at"`
	LogTemplateID          int64     `db:"log_template_id"`
}

// TaskInstance 表示任务实例
type TaskInstance struct {
	TaskID             string    `db:"task_id"`
	DagID              string    `db:"dag_id"`
	RunID              string    `db:"run_id"`
	MapIndex           int       `db:"map_index"`
	StartDate          time.Time `db:"start_date"`
	EndDate            time.Time `db:"end_date"`
	Duration           float64   `db:"duration"`
	State              string    `db:"state"`
	TryNumber          int       `db:"try_number"`
	MaxTries           int       `db:"max_tries"`
	Hostname           string    `db:"hostname"`
	Unixname           string    `db:"unixname"`
	JobID              int64     `db:"job_id"`
	Pool               string    `db:"pool"`
	PoolSlots          int       `db:"pool_slots"`
	Queue              string    `db:"queue"`
	PriorityWeight     int       `db:"priority_weight"`
	Operator           string    `db:"operator"`
	QueuedDttm         time.Time `db:"queued_dttm"`
	PID                int       `db:"pid"`
	ExecutorConfig     []byte    `db:"executor_config"`
	ExternalExecutorID string    `db:"external_executor_id"`
	TriggerID          int64     `db:"trigger_id"`
	TriggerTimeout     time.Time `db:"trigger_timeout"`
	NextMethod         string    `db:"next_method"`
	NextKwargs         []byte    `db:"next_kwargs"`
	QueuedByJobID      int64     `db:"queued_by_job_id"`
	CustomOperatorName string    `db:"custom_operator_name"`
	UpdatedAt          time.Time `db:"updated_at"`
}

// XCom 表示任务间通信的数据
type XCom struct {
	DagID     string    `db:"dag_id"`
	TaskID    string    `db:"task_id"`
	RunID     string    `db:"run_id"`
	MapIndex  int       `db:"map_index"`
	Key       string    `db:"key"`
	Value     []byte    `db:"value"`
	Timestamp time.Time `db:"timestamp"`
	DagRunID  int64     `db:"dag_run_id"`
}

// Log 表示日志记录
type Log struct {
	ID            int64     `db:"id"`
	DagID         string    `db:"dag_id"`
	TaskID        string    `db:"task_id"`
	ExecutionDate time.Time `db:"execution_date"`
	Dttm          time.Time `db:"dttm"`
	Event         string    `db:"event"`
	Owner         string    `db:"owner"`
	Extra         string    `db:"extra"`
	MapIndex      int       `db:"map_index"`
}

// Job 表示Airflow的作业记录
type Job struct {
	ID              int64     `db:"id"`
	DagID           string    `db:"dag_id"`
	State           string    `db:"state"`
	JobType         string    `db:"job_type"`
	StartDate       time.Time `db:"start_date"`
	EndDate         time.Time `db:"end_date"`
	LatestHeartbeat time.Time `db:"latest_heartbeat"`
	Hostname        string    `db:"hostname"`
	Unixname        string    `db:"unixname"`
	ExecutorClass   string    `db:"executor_class"`
}

// TableConfig 存储表的配置信息
type TableConfig struct {
	TableName     string
	RetentionDays int
	DateColumn    string
}

// Config 存储所有清理配置
type Config struct {
	RetentionDays map[string]int
	BatchSize     int
	DryRun        bool
	Verbose       bool
	SleepSeconds  int
}
