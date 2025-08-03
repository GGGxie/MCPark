package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Config 应用配置结构
type Config struct {
	App      AppConfig      `mapstructure:"app"`
	Database DatabaseConfig `mapstructure:"database"`
	Log      LogConfig      `mapstructure:"log"`
	Server   ServerConfig   `mapstructure:"server"`
	CORS     CORSConfig     `mapstructure:"cors"`
}

// AppConfig 应用基础配置
type AppConfig struct {
	Name    string `mapstructure:"name"`
	Version string `mapstructure:"version"`
	Mode    string `mapstructure:"mode"`
	Port    int    `mapstructure:"port"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Type   string       `mapstructure:"type"`
	SQLite SQLiteConfig `mapstructure:"sqlite"`
	MySQL  MySQLConfig  `mapstructure:"mysql"`
}

// SQLiteConfig SQLite配置
type SQLiteConfig struct {
	Path string `mapstructure:"path"`
}

// MySQLConfig MySQL配置
type MySQLConfig struct {
	Host      string `mapstructure:"host"`
	Port      int    `mapstructure:"port"`
	Username  string `mapstructure:"username"`
	Password  string `mapstructure:"password"`
	Database  string `mapstructure:"database"`
	Charset   string `mapstructure:"charset"`
	ParseTime bool   `mapstructure:"parseTime"`
	Loc       string `mapstructure:"loc"`
}

// LogConfig 日志配置
type LogConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
	Output string `mapstructure:"output"`
	File   string `mapstructure:"file"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	ReadTimeout  int    `mapstructure:"readTimeout"`
	WriteTimeout int    `mapstructure:"writeTimeout"`
	StaticPath   string `mapstructure:"staticPath"`
}

// CORSConfig 跨域配置
type CORSConfig struct {
	AllowOrigins []string `mapstructure:"allowOrigins"`
	AllowMethods []string `mapstructure:"allowMethods"`
	AllowHeaders []string `mapstructure:"allowHeaders"`
}

// 全局配置实例
var GlobalConfig *Config

// InitConfig 初始化配置
func InitConfig() error {
	v := viper.New()

	// 设置配置文件名和路径
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath("./conf")

	// 设置环境变量前缀
	v.SetEnvPrefix("MCP")
	v.AutomaticEnv()

	// 读取配置文件
	if err := v.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	// 解析配置到结构体
	GlobalConfig = &Config{}
	if err := v.Unmarshal(GlobalConfig); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// 处理相对路径
	if err := processRelativePaths(); err != nil {
		return fmt.Errorf("failed to process relative paths: %w", err)
	}

	return nil
}

// processRelativePaths 处理配置中的相对路径
func processRelativePaths() error {
	// 获取当前工作目录
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	// 处理 SQLite 路径
	if !filepath.IsAbs(GlobalConfig.Database.SQLite.Path) {
		GlobalConfig.Database.SQLite.Path = filepath.Join(wd, GlobalConfig.Database.SQLite.Path)
	}

	// 处理静态文件路径
	if !filepath.IsAbs(GlobalConfig.Server.StaticPath) {
		GlobalConfig.Server.StaticPath = filepath.Join(wd, GlobalConfig.Server.StaticPath)
	}

	// 处理日志文件路径
	if GlobalConfig.Log.Output == "file" && !filepath.IsAbs(GlobalConfig.Log.File) {
		GlobalConfig.Log.File = filepath.Join(wd, GlobalConfig.Log.File)
	}

	return nil
}

// GetDatabaseDSN 获取数据库连接字符串
func (c *Config) GetDatabaseDSN() string {
	switch c.Database.Type {
	case "mysql":
		return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%t&loc=%s",
			c.Database.MySQL.Username,
			c.Database.MySQL.Password,
			c.Database.MySQL.Host,
			c.Database.MySQL.Port,
			c.Database.MySQL.Database,
			c.Database.MySQL.Charset,
			c.Database.MySQL.ParseTime,
			c.Database.MySQL.Loc,
		)
	case "sqlite":
		fallthrough
	default:
		return c.Database.SQLite.Path
	}
}

// IsSQLite 判断是否使用SQLite
func (c *Config) IsSQLite() bool {
	return c.Database.Type == "sqlite" || c.Database.Type == ""
}

// IsMySQL 判断是否使用MySQL
func (c *Config) IsMySQL() bool {
	return c.Database.Type == "mysql"
}

// GetServerAddr 获取服务器监听地址
func (c *Config) GetServerAddr() string {
	return fmt.Sprintf(":%d", c.App.Port)
}
