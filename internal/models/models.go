package models

import "time"

// DagRun 表示DAG的一次运行
type DagRun struct {
	ID              int64     `db:"id"`
	DagID           string    `db:"dag_id"`
	ExecutionDate   time.Time `db:"execution_date"`
	State           string    `db:"state"`
	RunID           string    `db:"run_id"`
	ExternalTrigger bool      `db:"external_trigger"`
	StartDate       time.Time `db:"start_date"`
	EndDate         time.Time `db:"end_date"`
}

// TaskInstance 表示任务实例
type TaskInstance struct {
	TaskID        string    `db:"task_id"`
	DagID         string    `db:"dag_id"`
	ExecutionDate time.Time `db:"execution_date"`
	StartDate     time.Time `db:"start_date"`
	EndDate       time.Time `db:"end_date"`
	State         string    `db:"state"`
	TryNumber     int       `db:"try_number"`
	MaxTries      int       `db:"max_tries"`
	Hostname      string    `db:"hostname"`
	Unixname      string    `db:"unixname"`
	JobID         int64     `db:"job_id"`
	QueuedDttm    time.Time `db:"queued_dttm"`
	RunID         string    `db:"run_id"`
}

// XCom 表示任务间通信的数据
type XCom struct {
	ID            int64     `db:"id"`
	Key           string    `db:"key"`
	Value         string    `db:"value"`
	TaskID        string    `db:"task_id"`
	DagID         string    `db:"dag_id"`
	ExecutionDate time.Time `db:"execution_date"`
	RunID         string    `db:"run_id"`
}

// Log 表示日志记录
type Log struct {
	ID          int64     `db:"id"`
	DagID       string    `db:"dag_id"`
	TaskID      string    `db:"task_id"`
	ExecutionID int64     `db:"execution_id"`
	EventType   string    `db:"event_type"`
	LogData     string    `db:"log_data"`
	CreatedAt   time.Time `db:"created_at"`
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
	HostName        string    `db:"hostname"`
}

// TableConfig 存储表的配置信息
type TableConfig struct {
	TableName     string
	RetentionDays int
}

// Config 存储所有清理配置
type Config struct {
	RetentionDays map[string]int
	BatchSize     int
	DryRun        bool
	Verbose       bool
}
