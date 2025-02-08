package cronjobx

// IJob 定義一個可執行的任務介面
type IJob interface {
	Name() string
	// Run 執行任務邏輯，返回錯誤表示執行失敗
	Run() error
	// Spec 返回 Cron 表達式
	Spec() string
}
