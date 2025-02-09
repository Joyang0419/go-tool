package cronjobx

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
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
}

// ManagerParams 定義 Manager 的依賴參數
type ManagerParams struct {
	fx.In
	Jobs   []IJob `group:"jobs"`
	Config Config
}

func (receiver *Manager) Start() {
	receiver.cron.Start()
}

func (receiver *Manager) Stop() {
	receiver.cron.Stop()
}

func NewManager(lc fx.Lifecycle, params ManagerParams) {
	m := &Manager{
		cron: cron.New(
			cron.WithSeconds(),
			cron.WithChain(
				cron.Recover(cron.DefaultLogger),
				cron.SkipIfStillRunning(cron.DefaultLogger),
			),
		),
	}

	// 註冊所有任務
	for _, job := range params.Jobs {
		slog.Info(fmt.Sprintf("[NewManager]Registering job: %s, spec: %s", job.Name(), job.Spec()))
		if _, err := m.cron.AddFunc(job.Spec(), func() {
			slog.Info(fmt.Sprintf("[NewManager]Running job: %s", job.Name()))
			if errJob := job.Run(); errJob != nil {
				slog.Error(fmt.Sprintf("[NewManager]Job.Run() error: %v, jobName: %s", errJob, job.Name()))
			}
			slog.Info(fmt.Sprintf("[NewManager]Job finished: %s", job.Name()))
		}); err != nil {
			slog.Error("[NewManager]m.cron.AddFunc error: %v", err)
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
