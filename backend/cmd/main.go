package main

import (
	"fmt"
	"log"

	"MCPark/internal/config"
	"MCPark/internal/database"
	"MCPark/internal/handler"
	"MCPark/internal/repository"
	"MCPark/internal/service"

	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化配置
	if err := config.InitConfig(); err != nil {
		log.Fatalf("Failed to initialize config: %v", err)
	}
	log.Printf("✅ Config initialized: %s v%s", config.GlobalConfig.App.Name, config.GlobalConfig.App.Version)

	// 设置Gin模式
	gin.SetMode(config.GlobalConfig.App.Mode)

	// 初始化数据库
	database.InitDatabase()

	// 初始化仓库层
	serviceRepo := repository.NewServiceRepository(database.GetDB())
	categoryRepo := repository.NewCategoryRepository(database.GetDB())

	// 初始化服务层
	mcpService := service.NewMCPService(serviceRepo, categoryRepo)

	// 初始化处理器层
	h := handler.NewHandler(mcpService)

	// 设置路由
	r := h.SetupRoutes()

	// 启动信息
	addr := config.GlobalConfig.GetServerAddr()
	port := config.GlobalConfig.App.Port

	fmt.Printf("🚀 %s v%s 服务启动成功！\n", config.GlobalConfig.App.Name, config.GlobalConfig.App.Version)
	fmt.Printf("📊 数据库类型: %s\n", config.GlobalConfig.Database.Type)
	fmt.Printf("📱 访问地址: http://localhost:%d\n", port)
	fmt.Printf("🔧 API地址: http://localhost:%d/api/v1\n", port)
	fmt.Printf("💚 健康检查: http://localhost:%d/api/v1/health\n", port)

	// 启动服务器
	log.Fatal(r.Run(addr))
}
