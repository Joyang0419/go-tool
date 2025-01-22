package cronjobx

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	pkgerrors "github.com/pkg/errors"
	"github.com/robfig/cron/v3"
	"go.uber.org/fx"
)

// Config 定義 CronJob Manager 的配置
type Config struct {
	// 關閉時的超時時間
	ShutdownTimeout time.Duration `mapstructure:"shutdown_timeout"`
}

type Manager struct {
	cron *cron.Cron
	jobs map[string]*JobInfo
	mu   sync.RWMutex
}

// ManagerParams 定義 Manager 的依賴參數
type ManagerParams struct {
	fx.In
	Jobs   []IJob `group:"jobs"`
	Config Config
}

// RegisterJob 註冊一個新的任務
func (receiver *Manager) RegisterJob(job IJob) error {
	receiver.mu.Lock()
	defer receiver.mu.Unlock()

	name := job.Name()
	if _, exists := receiver.jobs[name]; exists {
		return pkgerrors.New("[Manager][RegisterJob]Job already exists")
	}

	wrapper := func() {
		info := receiver.jobs[name]
		info.LastRun = time.Now()

		if err := job.Run(); err != nil {
			slog.Error(fmt.Sprintf("[CronJob:%s]Failed to run job: %v", name, err))
			info.LastError = err
		} else {
			info.LastError = nil
		}
	}

	entryID, err := receiver.cron.AddFunc(job.Spec(), wrapper)
	if err != nil {
		return pkgerrors.Wrap(err, "[Manager][RegisterJob]receiver.cron.AddFunc(job.Spec(), wrapper) error")
	}

	receiver.jobs[name] = &JobInfo{
		Name:    name,
		EntryID: entryID,
	}

	return nil
}

// GetJobStatus 獲取任務狀態
func (receiver *Manager) GetJobStatus(name string) (*JobInfo, error) {
	receiver.mu.RLock()
	defer receiver.mu.RUnlock()

	info, exists := receiver.jobs[name]
	if !exists {
		return nil, fmt.Errorf("job %s not found", name)
	}
	return info, nil
}

// ListJobs 列出所有任務
func (receiver *Manager) ListJobs() []*JobInfo {
	receiver.mu.RLock()
	defer receiver.mu.RUnlock()

	var result []*JobInfo
	for _, info := range receiver.jobs {
		result = append(result, info)
	}
	return result
}

func (receiver *Manager) Start() {
	receiver.cron.Start()
}

func (receiver *Manager) Stop() {
	receiver.cron.Stop()
}

func NewManager(lc fx.Lifecycle, params ManagerParams) {
	m := &Manager{
		cron: cron.New(cron.WithSeconds()),
		jobs: make(map[string]*JobInfo),
	}

	// 註冊所有任務
	for _, job := range params.Jobs {
		if err := m.RegisterJob(job); err != nil {
			slog.Error(fmt.Sprintf("[NewManager]m.RegisterJob(job) error %s: jobName: %v", job.Name(), err))
			continue
		}
	}

	// 註冊生命週期鉤子
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			slog.Info("[NewManager]Starting cron manager")
			m.Start()

			// 啟動信號監聽，用於優雅停止
			go func() {
				signalChan := make(chan os.Signal, 1)
				signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
				<-signalChan

				slog.Info("[NewManager]Signal received, stopping cron manager")
				stopChan := make(chan struct{})
				go func() {
					m.Stop()
					close(stopChan)
				}()

				select {
				case <-time.After(params.Config.ShutdownTimeout):
					slog.Error("[NewManager]Stop cron manager timeout")
				case <-stopChan:
					slog.Info("[NewManager]Cron manager stopped gracefully")
				}
				os.Exit(0)
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			slog.Info("[NewManager]Stopping cron manager via lifecycle hook")
			stopChan := make(chan struct{})
			go func() {
				m.Stop()
				close(stopChan)
			}()

			select {
			case <-time.After(params.Config.ShutdownTimeout):
				return pkgerrors.New("[NewManager]Stop cron manager timeout")
			case <-stopChan:
				return nil
			}
		},
	})
}
