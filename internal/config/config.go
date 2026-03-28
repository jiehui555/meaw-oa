package config

import "os"

// Config 应用配置结构
type Config struct {
	Port    string // HTTP 服务端口
	DBPath  string // 数据库文件路径
	LogPath string // 日志文件路径
}

// Load 从环境变量加载配置
// 默认值：端口 3000，数据库 app.db，日志 app.log
func Load() *Config {
	return &Config{
		Port:    getEnv("PORT", "3000"),
		DBPath:  getEnv("DB_PATH", "app.db"),
		LogPath: getEnv("LOG_PATH", "app.log"),
	}
}

// getEnv 获取环境变量值，如果未设置则返回默认值
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
