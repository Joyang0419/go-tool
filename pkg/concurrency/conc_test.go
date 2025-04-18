package concurrency_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"go-tool/pkg/concurrency"
)

func TestTasksWithResult_Success(t *testing.T) {
	// 創建測試上下文
	ctx := context.Background()

	// TestResult 用於測試的結果結構體
	type TestResult struct {
		StringValue string
		IntValue    int
		FloatValue  float64
		SliceValue  []string
	}

	// 定義測試任務
	tasks := []concurrency.Task[TestResult]{
		// 獲取字符串任務
		{
			FetchFn: func() (*TestResult, error) {
				// 模擬處理過程
				time.Sleep(100 * time.Millisecond)
				return &TestResult{StringValue: "test string"}, nil
			},
			MergeFn: func(part, final *TestResult) {
				final.StringValue = part.StringValue
			},
		},

		// 獲取整數任務
		{
			FetchFn: func() (*TestResult, error) {
				// 模擬處理過程
				time.Sleep(50 * time.Millisecond)
				return &TestResult{IntValue: 42}, nil
			},
			MergeFn: func(part, final *TestResult) {
				final.IntValue = part.IntValue
			},
		},

		// 獲取浮點數任務
		{
			FetchFn: func() (*TestResult, error) {
				// 模擬處理過程
				time.Sleep(150 * time.Millisecond)
				return &TestResult{FloatValue: 3.14}, nil
			},
			MergeFn: func(part, final *TestResult) {
				final.FloatValue = part.FloatValue
			},
		},

		// 獲取切片任務
		{
			FetchFn: func() (*TestResult, error) {
				// 模擬處理過程
				time.Sleep(75 * time.Millisecond)
				return &TestResult{SliceValue: []string{"a", "b", "c"}}, nil
			},
			MergeFn: func(part, final *TestResult) {
				final.SliceValue = append(final.SliceValue, part.SliceValue...)
			},
		},
	}

	// 執行測試
	startTime := time.Now()
	result, err := concurrency.TasksWithResult(ctx, tasks)
	duration := time.Since(startTime)

	// 驗證結果
	assert.NoError(t, err, "應該沒有錯誤")
	assert.NotNil(t, result, "結果不應為空")

	// 驗證各個字段
	assert.Equal(t, "test string", result.StringValue)
	assert.Equal(t, 42, result.IntValue)
	assert.Equal(t, 3.14, result.FloatValue)
	assert.Equal(t, []string{"a", "b", "c"}, result.SliceValue)

	t.Logf("result: %+v", result)

	// 驗證並行執行節省了時間
	// 總耗時應該接近最長的任務耗時 (150ms)，而不是所有任務耗時之和 (375ms)
	assert.Less(t, duration, 200*time.Millisecond,
		"並行執行應該顯著節省時間，總耗時應接近最長任務的耗時")

	t.Logf("所有任務並行完成耗時: %v", duration)
}
