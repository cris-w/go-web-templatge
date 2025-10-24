package app

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/viper"
)

// Config 应用配置结构
type Config struct {
	Debug  bool   `mapstructure:"debug"`
	Addr   string `mapstructure:"addr"`
	DB     DBConfig
	JWT    JWTConfig
	Log    LogConfig
	Server ServerConfig
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

// ServerConfig HTTP 服务器配置
type ServerConfig struct {
	ReadTimeout  int `mapstructure:"read_timeout"`  // 秒
	WriteTimeout int `mapstructure:"write_timeout"` // 秒
	IdleTimeout  int `mapstructure:"idle_timeout"`  // 秒
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

// GetReadTimeout 获取读超时时间，默认 15 秒
func (s *ServerConfig) GetReadTimeout() time.Duration {
	if s.ReadTimeout <= 0 {
		return 15 * time.Second
	}
	return time.Duration(s.ReadTimeout) * time.Second
}

// GetWriteTimeout 获取写超时时间，默认 15 秒
func (s *ServerConfig) GetWriteTimeout() time.Duration {
	if s.WriteTimeout <= 0 {
		return 15 * time.Second
	}
	return time.Duration(s.WriteTimeout) * time.Second
}

// GetIdleTimeout 获取空闲超时时间，默认 60 秒
func (s *ServerConfig) GetIdleTimeout() time.Duration {
	if s.IdleTimeout <= 0 {
		return 60 * time.Second
	}
	return time.Duration(s.IdleTimeout) * time.Second
}
