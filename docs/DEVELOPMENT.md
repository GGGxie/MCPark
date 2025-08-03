# 开发指南

## 开发环境搭建

### 后端开发

1. **安装Go 1.21+**
   ```bash
   # macOS
   brew install go
   
   # 或下载官方安装包
   # https://golang.org/dl/
   ```

2. **初始化项目**
   ```bash
   cd backend
   go mod tidy
   ```

3. **运行开发服务器**
   ```bash
   go run cmd/main.go
   ```

4. **数据库迁移**
   ```bash
   go run cmd/migrate.go
   ```

### 前端开发

1. **安装Node.js 18+**
   ```bash
   # macOS
   brew install node
   
   # 或使用nvm
   nvm install 18
   nvm use 18
   ```

2. **安装依赖**
   ```bash
   cd frontend
   npm install
   ```

3. **运行开发服务器**
   ```bash
   npm run dev
   ```

## 项目结构详解

### 后端结构
```
backend/
├── cmd/
│   ├── main.go              # 主程序入口
│   └── migrate.go           # 数据库迁移工具
├── internal/
│   ├── api/                 # API处理器
│   │   ├── handler/         # HTTP处理器
│   │   ├── middleware/      # 中间件
│   │   └── router/          # 路由配置
│   ├── model/               # 数据模型
│   ├── service/             # 业务逻辑层
│   └── repository/          # 数据访问层
├── pkg/                     # 公共包
│   ├── config/             # 配置管理
│   ├── database/           # 数据库连接
│   ├── logger/             # 日志管理
│   └── utils/              # 工具函数
├── configs/                # 配置文件
└── docs/                   # API文档
```

### 前端结构
```
frontend/
├── src/
│   ├── components/         # 可复用组件
│   │   ├── ServiceCard/    # 服务卡片
│   │   ├── ConfigViewer/   # 配置查看器
│   │   ├── DeepLinkGenerator/ # Deep Link生成器
│   │   └── Layout/         # 布局组件
│   ├── pages/              # 页面组件
│   │   ├── Home/           # 首页
│   │   ├── ServiceDetail/  # 服务详情
│   │   ├── AddService/     # 添加服务
│   │   ├── Category/       # 分类管理
│   │   └── Stats/          # 统计页面
│   ├── services/           # API服务
│   │   ├── api.ts          # API客户端
│   │   ├── serviceApi.ts   # 服务相关API
│   │   └── categoryApi.ts  # 分类相关API
│   ├── store/              # 状态管理
│   ├── utils/              # 工具函数
│   ├── types/              # TypeScript类型定义
│   └── styles/             # 样式文件
├── public/                 # 静态资源
└── package.json
```

## 开发规范

### 后端开发规范

1. **代码结构**
   - 使用Clean Architecture架构
   - 分层清晰：Handler -> Service -> Repository
   - 依赖注入，便于测试

2. **命名规范**
   - 文件名：snake_case
   - 包名：小写单词
   - 函数名：PascalCase（公开）或camelCase（私有）
   - 常量：UPPER_CASE

3. **错误处理**
   - 使用自定义错误类型
   - 统一错误响应格式
   - 记录错误日志

4. **API设计**
   - RESTful风格
   - 统一响应格式
   - 适当的HTTP状态码
   - API版本控制

### 前端开发规范

1. **组件设计**
   - 函数式组件 + Hooks
   - 组件职责单一
   - Props类型严格定义

2. **状态管理**
   - 使用Zustand进行全局状态管理
   - 本地状态优先使用useState
   - 复杂逻辑使用useReducer

3. **样式规范**
   - 使用Ant Design组件库
   - CSS Modules或styled-components
   - 响应式设计

4. **TypeScript规范**
   - 严格类型检查
   - 避免any类型
   - 接口定义清晰

## API文档

### 认证
目前版本不需要认证，后续版本可能添加JWT认证。

### 请求格式
```json
{
  "Content-Type": "application/json"
}
```

### 响应格式
```json
{
  "code": 0,
  "message": "success",
  "data": {},
  "timestamp": "2024-01-01T00:00:00Z"
}
```

### 错误码
- 0: 成功
- 400: 请求参数错误
- 404: 资源不存在
- 500: 服务器内部错误

## 测试策略

### 后端测试
```bash
# 运行所有测试
go test ./...

# 运行特定包的测试
go test ./internal/service

# 生成测试覆盖率报告
go test -cover ./...
```

### 前端测试
```bash
# 运行单元测试
npm test

# 运行e2e测试
npm run test:e2e
```

## 数据库管理

### 迁移
```bash
# 创建迁移文件
go run cmd/migrate.go create add_new_table

# 执行迁移
go run cmd/migrate.go up

# 回滚迁移
go run cmd/migrate.go down
```

### 种子数据
```bash
# 生成测试数据
go run cmd/seed.go
```

## 部署

### 开发环境
```bash
# 启动后端
cd backend && go run cmd/main.go

# 启动前端
cd frontend && npm run dev
```

### 生产环境
```bash
# 使用Docker Compose
cd deployment
docker-compose up -d
```

## 常见问题

### Q: 如何添加新的API端点？
A: 
1. 在`internal/model`中定义请求/响应结构
2. 在`internal/repository`中添加数据访问方法
3. 在`internal/service`中添加业务逻辑
4. 在`internal/api/handler`中添加HTTP处理器
5. 在`internal/api/router`中注册路由

### Q: 为什么没有用户系统？
A: 这是一个私人使用的系统，不需要复杂的用户管理和权限控制，简化了系统设计和维护成本。

### Q: 如何添加新的前端页面？
A:
1. 在`src/pages`中创建页面组件
2. 在`src/services`中添加相关API调用
3. 在路由配置中添加新路由
4. 更新导航菜单

### Q: 如何处理CORS问题？
A: 后端已配置CORS中间件，前端开发时确保API地址正确。

## 贡献指南

1. Fork项目
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 创建Pull Request

## 联系方式

如有问题，请提交Issue或联系项目维护者。