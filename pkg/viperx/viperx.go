package viperx

import (
	"embed"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

type LoadConfigSO struct {
	FilePath string // 完整的配置文件路徑，例如 "./config.yaml"
	Config   any
}

// LoadConfig 載入配置文件
// path: 配置文件路徑
// name: 配置文件名稱（不包含副檔名）
// configType: 配置文件類型（如 yaml, json 等）
// config: 配置結構體的指針
func LoadConfig(so LoadConfigSO) {
	// 檢查配置文件是否存在
	if _, err := os.Stat(so.FilePath); os.IsNotExist(err) {
		slog.Info(fmt.Sprintf("viperx.LoadConfig config file not found: %s", so.FilePath))
		return
	}

	slog.Info("viperx.LoadConfig from filepath", slog.String("filePath", so.FilePath))

	// 提取文件路徑、名稱和類型
	dir := filepath.Dir(so.FilePath)
	filename := filepath.Base(so.FilePath)
	ext := filepath.Ext(filename)
	name := strings.TrimSuffix(filename, ext)
	configType := strings.TrimPrefix(ext, ".")

	v := viper.New()

	// 設置配置文件路徑
	v.AddConfigPath(dir)

	// 設置配置文件名稱
	v.SetConfigName(name)

	// 設置配置文件類型
	v.SetConfigType(configType)

	// 設置環境變數前綴
	v.SetEnvPrefix(strings.ToUpper(name))

	// 將 . 轉換為 _
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// 自動載入環境變數
	v.AutomaticEnv()

	// 讀取配置文件
	if err := v.ReadInConfig(); err != nil {
		panic(fmt.Sprintf("viperx.LoadConfig v.ReadInConfig error: %v", err))
	}

	// 解析到結構體
	if err := v.Unmarshal(so.Config); err != nil {
		panic(fmt.Sprintf("viperx.LoadConfig v.Unmarshal error: %v", err))
	}

	return
}

type LoadEmbeddedConfigSO struct {
	Embed    embed.FS
	Filename string
	Config   any
}

// LoadEmbeddedConfig 加载嵌入式配置文件到指定的结构体
// config: 配置结构体的指针
func LoadEmbeddedConfig(so LoadEmbeddedConfigSO) {
	// 读取嵌入式配置
	data, err := so.Embed.ReadFile(so.Filename)
	if err != nil {
		panic(fmt.Sprintf("viperx.LoadEmbeddedConfig embed.ReadFile error: %v", err))
	}

	// 解析YAML到结构体
	if err = yaml.Unmarshal(data, so.Config); err != nil {
		panic(fmt.Sprintf("viperx.LoadEmbeddedConfig yaml.Unmarshal error %v", err))
	}

	slog.Info("LoadEmbeddedConfig success", slog.String("filename", so.Filename))
}
