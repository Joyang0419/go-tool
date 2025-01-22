package cronjobx

import (
	"time"

	"github.com/robfig/cron/v3"
)

// IJob 定義一個可執行的任務介面
type IJob interface {
	Name() string
	// Run 執行任務邏輯，返回錯誤表示執行失敗
	Run() error
	// Spec 返回 Cron 表達式
	Spec() string
}

// JobInfo 儲存任務的資訊
type JobInfo struct {
	Name      string
	EntryID   cron.EntryID
	LastRun   time.Time
	LastError error
}
