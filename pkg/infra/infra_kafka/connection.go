package infra_kafka

import (
	"fmt"
	"time"

	"github.com/IBM/sarama"
	"go.uber.org/fx"
)

type ConnectionConfig struct {
	// Kafka 基本配置
	Brokers []string `mapstructure:"brokers"`
	Group   string   `mapstructure:"group"`
	Topics  []string `mapstructure:"topics"`

	// 生產者配置
	Producer struct {
		// 重試次數
		RetryMax int `mapstructure:"retry_max"`
		// 重試間隔
		RetryBackoff time.Duration `mapstructure:"retry_backoff"`
		// 超時時間
		Timeout time.Duration `mapstructure:"timeout"`
		// 是否需要回應
		RequiredAcks sarama.RequiredAcks `mapstructure:"required_acks"`
	} `mapstructure:"producer"`

	// 消費者配置
	Consumer struct {
		// 消費起始位置：newest, oldest
		Offset string `mapstructure:"offset"`
		// 提交間隔
		CommitInterval time.Duration `mapstructure:"commit_interval"`
		// 每次獲取訊息的最大數量
		MaxFetchSize int32 `mapstructure:"max_fetch_size"`
	} `mapstructure:"consumer"`

	// 認證配置（如果需要）
	Auth struct {
		Username string `mapstructure:"username"`
		Password string `mapstructure:"password"`
		// SASL機制：PLAIN, SCRAM-SHA-256, SCRAM-SHA-512
		Mechanism string `mapstructure:"mechanism"`
	} `mapstructure:"auth"`
}

func NewKafkaProducer(config ConnectionConfig) (sarama.SyncProducer, error) {
	kafkaConfig := sarama.NewConfig()

	// 生產者配置
	// Retry.Max: 生產者重試的最大次數，超過次數後會放棄該消息
	kafkaConfig.Producer.Retry.Max = config.Producer.RetryMax
	// Retry.Backoff: 重試之間的等待時間
	kafkaConfig.Producer.Retry.Backoff = config.Producer.RetryBackoff
	// Timeout: 等待回應的超時時間
	kafkaConfig.Producer.Timeout = config.Producer.Timeout
	// RequiredAcks: 需要多少個分區副本確認才算成功
	// 0 = 不需要等待確認（可能丟失）
	// 1 = 等待 leader 確認（預設）
	// -1 = 等待所有副本確認（最安全，但最慢）
	kafkaConfig.Producer.RequiredAcks = config.Producer.RequiredAcks
	// Return.Successes: 是否等待並返回成功的響應，同步生產者必須設為 true
	kafkaConfig.Producer.Return.Successes = true

	// 認證配置
	if config.Auth.Username != "" {
		// Enable: 開啟 SASL 認證
		kafkaConfig.Net.SASL.Enable = true
		// User: SASL 認證使用者名稱
		kafkaConfig.Net.SASL.User = config.Auth.Username
		// Password: SASL 認證密碼
		kafkaConfig.Net.SASL.Password = config.Auth.Password
		// Mechanism: SASL 認證機制
		// 支援 PLAIN, SCRAM-SHA-256, SCRAM-SHA-512 等
		kafkaConfig.Net.SASL.Mechanism = sarama.SASLMechanism(config.Auth.Mechanism)
	}

	producer, err := sarama.NewSyncProducer(config.Brokers, kafkaConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create infra_kafka producer: %v", err)
	}

	return producer, nil
}

func NewKafkaConsumer(config ConnectionConfig) (sarama.ConsumerGroup, error) {
	kafkaConfig := sarama.NewConfig()

	// 消費者配置
	// 設置消費者的起始偏移量
	switch config.Consumer.Offset {
	case "newest":
		// OffsetNewest: 從最新的消息開始消費（只消費新消息）
		kafkaConfig.Consumer.Offsets.Initial = sarama.OffsetNewest
	case "oldest":
		// OffsetOldest: 從最舊的消息開始消費（包含歷史消息）
		kafkaConfig.Consumer.Offsets.Initial = sarama.OffsetOldest
	}

	// 設置消費者組的分區分配策略為輪詢
	// RoundRobin: 輪詢方式分配分區給消費者
	kafkaConfig.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRoundRobin()}
	// AutoCommit.Enable: 是否自動提交 offset
	kafkaConfig.Consumer.Offsets.AutoCommit.Enable = true
	// AutoCommit.Interval: 自動提交的間隔時間
	kafkaConfig.Consumer.Offsets.AutoCommit.Interval = config.Consumer.CommitInterval
	// Fetch.Max: 單次獲取消息的最大大小（bytes）
	kafkaConfig.Consumer.Fetch.Max = config.Consumer.MaxFetchSize

	// 認證配置
	if config.Auth.Username != "" {
		// Enable: 開啟 SASL 認證
		kafkaConfig.Net.SASL.Enable = true
		// User: SASL 認證使用者名稱
		kafkaConfig.Net.SASL.User = config.Auth.Username
		// Password: SASL 認證密碼
		kafkaConfig.Net.SASL.Password = config.Auth.Password
		// Mechanism: SASL 認證機制（PLAIN, SCRAM-SHA-256, SCRAM-SHA-512）
		kafkaConfig.Net.SASL.Mechanism = sarama.SASLMechanism(config.Auth.Mechanism)
	}

	consumer, err := sarama.NewConsumerGroup(config.Brokers, config.Group, kafkaConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create infra_kafka consumer: %v", err)
	}

	return consumer, nil
}

func InjectKafkaProducer(config ConnectionConfig) fx.Option {
	return fx.Options(
		fx.Supply(config),
		fx.Provide(NewKafkaProducer),
	)
}

func InjectKafkaConsumer(config ConnectionConfig) fx.Option {
	return fx.Options(
		fx.Supply(config),
		fx.Provide(NewKafkaConsumer),
	)
}
