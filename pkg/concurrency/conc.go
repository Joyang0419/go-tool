package concurrency

import (
	"context"
	"fmt"
	"reflect"

	"github.com/samber/lo"
	"github.com/sourcegraph/conc/pool"
)

type Task[TResult any] struct {
	FetchFn func() (*TResult, error)
	MergeFn func(partResult, finalResult *TResult)
}

// TasksWithResult 執行多個並發的函式，並將其部分結果整合為最終結果。
// 注意事項:
// - TResult 必須是結構體類型 (struct)，若傳入的泛型型別不是結構體，函式將返回錯誤。
// - TResult 的每個欄位應對應至一個執行緒的結果，請避免執行緒間的資料競態條件 (race condition)。
// - fns 切片不得為空，若 fns 為空，函式將返回錯誤。
// - 每個 fns 函式需提供 convertFn，該函數負責將部分結果整合至最終結果中。請確保整合過程中邏輯正確無誤且資料不會誤覆蓋。
func TasksWithResult[TResult any](ctx context.Context, tasks []Task[TResult]) (*TResult, error) {
	// 確保 TResult 是結構體類型
	var zero TResult
	resultType := reflect.TypeOf(zero)
	if resultType.Kind() != reflect.Struct {
		return nil, fmt.Errorf("TasksWithResult TResult must be a struct")
	}

	if len(tasks) == 0 {
		return nil, fmt.Errorf("DifferentTasks tasks is empty")
	}

	var result = new(TResult)
	p := pool.New().WithContext(ctx)

	lo.ForEach(tasks, func(task Task[TResult], index int) {
		_ = index
		pGoFn := func(ctx context.Context) error {
			part, err := task.FetchFn()
			if !lo.IsNil(err) {
				return err
			}
			task.MergeFn(part, result)
			return nil
		}

		p.Go(pGoFn)

	})

	if err := p.Wait(); !lo.IsNil(err) {
		return nil, err
	}

	return result, nil
}
