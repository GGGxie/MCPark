package database

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"MCPark/internal/config"
	"MCPark/internal/model"

	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// InitDatabase 初始化数据库连接
func InitDatabase() {
	var err error

	// 检查配置是否已初始化
	if config.GlobalConfig == nil {
		log.Fatal("Config not initialized. Please call config.InitConfig() first")
	}

	// 配置GORM日志级别
	logLevel := getLogLevel()
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	}

	// 根据配置选择数据库驱动
	DB, err = openDatabase(gormConfig)
	if err != nil {
		log.Fatalf("Failed to connect to %s database: %v", config.GlobalConfig.Database.Type, err)
	}

	// 自动迁移数据库结构
	err = DB.AutoMigrate(
		&model.MCPService{},
		&model.Category{},
	)
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// 初始化默认分类
	initDefaultCategories()

	log.Printf("✅ Database (%s) initialized successfully", config.GlobalConfig.Database.Type)
}

// openDatabase 根据配置打开数据库连接
func openDatabase(gormConfig *gorm.Config) (*gorm.DB, error) {
	cfg := config.GlobalConfig

	switch cfg.Database.Type {
	case "mysql":
		return openMySQL(gormConfig)
	case "sqlite":
		fallthrough
	default:
		return openSQLite(gormConfig)
	}
}

// openSQLite 打开SQLite数据库连接
func openSQLite(gormConfig *gorm.Config) (*gorm.DB, error) {
	cfg := config.GlobalConfig
	dbPath := cfg.Database.SQLite.Path

	// 确保数据目录存在
	dataDir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dataDir, 0o755); err != nil {
		return nil, fmt.Errorf("failed to create data directory: %w", err)
	}

	log.Printf("📁 Using SQLite database: %s", dbPath)
	return gorm.Open(sqlite.Open(dbPath), gormConfig)
}

// openMySQL 打开MySQL数据库连接
func openMySQL(gormConfig *gorm.Config) (*gorm.DB, error) {
	cfg := config.GlobalConfig
	dsn := cfg.GetDatabaseDSN()

	log.Printf("🐬 Connecting to MySQL: %s@%s:%d/%s",
		cfg.Database.MySQL.Username,
		cfg.Database.MySQL.Host,
		cfg.Database.MySQL.Port,
		cfg.Database.MySQL.Database,
	)

	return gorm.Open(mysql.Open(dsn), gormConfig)
}

// getLogLevel 根据配置获取GORM日志级别
func getLogLevel() logger.LogLevel {
	if config.GlobalConfig == nil {
		return logger.Info
	}

	switch config.GlobalConfig.Log.Level {
	case "debug":
		return logger.Info
	case "info":
		return logger.Warn
	case "warn":
		return logger.Error
	case "error":
		return logger.Silent
	default:
		return logger.Info
	}
}

// initDefaultCategories 初始化默认分类
func initDefaultCategories() {
	defaultCategories := []model.Category{
		{Name: "database", DisplayName: "数据库", Description: "数据库相关的MCP服务", Icon: "🗃️", SortOrder: 1},
		{Name: "api", DisplayName: "API服务", Description: "API接口相关的MCP服务", Icon: "🌐", SortOrder: 2},
		{Name: "file", DisplayName: "文件系统", Description: "文件操作相关的MCP服务", Icon: "📁", SortOrder: 3},
		{Name: "tool", DisplayName: "开发工具", Description: "开发辅助工具", Icon: "🔧", SortOrder: 4},
		{Name: "other", DisplayName: "其他", Description: "其他类型的MCP服务", Icon: "📦", SortOrder: 5},
	}

	for _, category := range defaultCategories {
		var existingCategory model.Category
		result := DB.Where("name = ?", category.Name).First(&existingCategory)
		if result.Error != nil {
			// 分类不存在，创建新的
			if err := DB.Create(&category).Error; err != nil {
				log.Printf("Failed to create category %s: %v", category.Name, err)
			} else {
				log.Printf("✅ Created default category: %s", category.DisplayName)
			}
		}
	}
}

// GetDB 获取数据库实例
func GetDB() *gorm.DB {
	return DB
}
