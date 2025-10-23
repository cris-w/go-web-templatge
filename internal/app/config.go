package app

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/viper"
)

// Config 应用配置结构
type Config struct {
	Debug bool   `mapstructure:"debug"`
	Addr  string `mapstructure:"addr"`
	DB    DBConfig
	JWT   JWTConfig
	Log   LogConfig
}

// DBConfig 数据库配置
type DBConfig struct {
	DSN             string `mapstructure:"dsn"`
	MaxIdleConns    int    `mapstructure:"max_idle_conns"`
	MaxOpenConns    int    `mapstructure:"max_open_conns"`
	ConnMaxLifetime int    `mapstructure:"conn_max_lifetime"`
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret      string `mapstructure:"secret"`
	ExpireHours int    `mapstructure:"expire_hours"`
}

// LogConfig 日志配置
type LogConfig struct {
	Level    string `mapstructure:"level"`
	FilePath string `mapstructure:"file_path"`
}

var GlobalConfig *Config

// LoadConfig 加载配置文件
func LoadConfig() (*Config, error) {
	// 获取环境变量，默认为 dev
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "dev"
	}

	viper.SetConfigName(fmt.Sprintf("config_%s", env))
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")
	viper.AddConfigPath("../configs")
	viper.AddConfigPath("../../configs")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %w", err)
	}

	GlobalConfig = &config
	return &config, nil
}

// GetConnMaxLifetime 获取连接最大生命周期
func (c *DBConfig) GetConnMaxLifetime() time.Duration {
	return time.Duration(c.ConnMaxLifetime) * time.Second
}
