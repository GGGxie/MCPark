> ⚠️ **注意：本文档内容与当前代码不一致（例如启动脚本、端口、API 路径、配置方式均可能已过期）。请以 `docs/` 目录与代码为准。**

# MCP 服务收集管理系统

## 项目概述

MCP 服务收集管理系统是一个 B/S 架构的 Web 应用，用于收集、管理和分享 Model Context Protocol (MCP) 服务配置。用户可以添加自己开发的或收集到的 MCP 服务，系统会自动生成 Cursor IDE 的 deep link，方便其他用户一键安装。

## 快速开始

### 一键启动

1. **克隆项目**

   ```bash
   git clone <your-repo-url>
   cd MCPark
   ```

2. **启动系统**
   ```bash
   ./start.sh
   ```

脚本会自动：

- 构建前端应用
- 启动后端服务 (端口 8080)
- 打开浏览器访问系统

### 停止服务

```bash
./stop.sh
```

### 系统要求

- **Go**: 1.21+
- **Node.js**: 18+
- **npm**: 包管理器

### 数据库配置

系统支持 **SQLite** 和 **MySQL** 两种数据库：

#### SQLite (默认，推荐个人使用)

- ✅ **零配置**：开箱即用，无需安装额外服务
- ✅ **便于备份**：数据存储在 `backend/data/MCPark.db` 文件
- ✅ **轻量级**：适合个人和小团队使用

#### MySQL (适合生产环境)

- ✅ **高性能**：支持大并发和大数据量
- ✅ **团队协作**：支持多用户同时访问
- ✅ **成熟稳定**：企业级数据库

**切换数据库**：

```bash
# 使用 MySQL (需先安装 MySQL 并创建数据库)
cp backend/conf/config.mysql.yaml backend/conf/config.yaml
# 编辑 backend/conf/config.yaml 设置连接信息
./start.sh
```

详细配置说明请参考：📖 [数据库配置指南](docs/DATABASE_CONFIG.md)

## 功能需求

### 核心功能

1. **MCP 服务管理**

   - 添加新的 MCP 服务配置
   - 编辑已有的服务配置
   - 删除服务配置
   - 查看服务配置详情

2. **配置展示**

   - 展示 JSON 格式的 MCP 配置
   - 语法高亮显示
   - 一键复制配置

3. **Deep Link 生成**

   - 自动生成 Cursor IDE 的安装链接
   - 支持一键复制链接
   - 提供二维码分享

4. **服务分类**

   - 按功能分类（数据库、API、工具等）
   - 按语言分类
   - 标签系统

5. **搜索和筛选**
   - 按名称搜索
   - 按分类筛选
   - 按标签筛选

### 附加功能

1. **统计功能**
   - 服务使用次数
   - 热门服务排行
   - 最新添加服务

## 技术栈

### 后端

- **语言**: Go 1.21+
- **框架**: Gin
- **数据库**: SQLite3
- **ORM**: GORM
- **配置管理**: Viper
- **日志**: Logrus
- **API 文档**: Swagger

### 前端

- **框架**: React 18
- **UI 库**: Ant Design 5.x
- **HTTP 客户端**: Axios
- **构建工具**: Vite

### 部署

- **容器化**: Docker
- **反向代理**: Nginx
- **部署**: Docker Compose

## 系统架构

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   前端 (React)  │────│  后端 (Gin)    │────│ 数据库 (SQLite) │
│                 │    │                 │    │                 │
│ - 服务列表      │    │ - REST API      │    │ - 服务配置      │
│ - 配置展示      │    │ - 配置验证      │    │ - 分类信息      │
│ - Deep Link     │    │ - Link生成      │    │ - 用户数据      │
│ - 搜索筛选      │    │ - 文件上传      │    │                 │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

## 数据库设计

### 表结构

#### 1. mcp_services (MCP 服务表)

```sql
CREATE TABLE mcp_services (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(100) NOT NULL UNIQUE,           -- 服务名称
    display_name VARCHAR(200) NOT NULL,          -- 显示名称
    description TEXT,                            -- 服务描述
    config_json TEXT NOT NULL,                   -- JSON配置内容
    category_id INTEGER,                         -- 分类ID
    tags VARCHAR(500),                           -- 标签(逗号分隔)
    author VARCHAR(100),                         -- 作者
    repository_url VARCHAR(500),                 -- 仓库地址
    documentation_url VARCHAR(500),              -- 文档地址
    version VARCHAR(50),                         -- 版本号
    install_count INTEGER DEFAULT 0,            -- 安装次数
    is_verified BOOLEAN DEFAULT FALSE,          -- 是否验证
    is_active BOOLEAN DEFAULT TRUE,             -- 是否启用
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

#### 2. categories (分类表)

```sql
CREATE TABLE categories (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(100) NOT NULL UNIQUE,          -- 分类名称
    display_name VARCHAR(200) NOT NULL,         -- 显示名称
    description TEXT,                           -- 分类描述
    icon VARCHAR(100),                          -- 图标
    sort_order INTEGER DEFAULT 0,              -- 排序
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

<!-- 移除用户相关表，私人系统不需要用户管理 -->

## API 设计

### RESTful API 规范

#### 服务管理

```
GET    /api/v1/services              # 获取服务列表
POST   /api/v1/services              # 创建新服务
GET    /api/v1/services/{id}         # 获取服务详情
PUT    /api/v1/services/{id}         # 更新服务
DELETE /api/v1/services/{id}         # 删除服务
GET    /api/v1/services/{id}/config  # 获取配置JSON
GET    /api/v1/services/{id}/deeplink # 获取Deep Link
POST   /api/v1/services/{id}/install # 记录安装（增加计数）
```

#### 分类管理

```
GET    /api/v1/categories            # 获取分类列表
POST   /api/v1/categories            # 创建分类
PUT    /api/v1/categories/{id}       # 更新分类
DELETE /api/v1/categories/{id}       # 删除分类
```

#### 搜索和筛选

```
GET    /api/v1/search?q=keyword      # 搜索服务
GET    /api/v1/services?category={id} # 按分类筛选
GET    /api/v1/services?tags={tags}   # 按标签筛选
```

#### 系统信息

```
GET    /api/v1/health                # 健康检查
GET    /api/v1/stats                 # 系统统计信息
```

### API 响应格式

```json
{
  "code": 0,
  "message": "success",
  "data": {},
  "timestamp": "2024-01-01T00:00:00Z"
}
```

## Deep Link 生成规则

根据 Cursor 官方文档，Deep Link 格式为：

```
cursor://anysphere.cursor-deeplink/mcp/install?name=$NAME&config=$BASE64_ENCODED_CONFIG
```

### 生成步骤

1. 获取 MCP 服务的 JSON 配置
2. 将 JSON 字符串进行 Base64 编码
3. 构造完整的 Deep Link URL
4. 提供复制和二维码功能

### 示例

```go
func GenerateDeepLink(serviceName string, config map[string]interface{}) string {
    jsonBytes, _ := json.Marshal(config)
    encodedConfig := base64.StdEncoding.EncodeToString(jsonBytes)
    return fmt.Sprintf("cursor://anysphere.cursor-deeplink/mcp/install?name=%s&config=%s",
        url.QueryEscape(serviceName), encodedConfig)
}
```

## 前端页面设计

### 页面结构

1. **首页** - 服务列表和搜索
2. **服务详情页** - 配置展示和 Deep Link
3. **添加/编辑服务页** - 服务配置表单
4. **分类管理页** - 分类列表和管理
5. **统计页面** - 服务使用统计
6. **关于页面** - 项目介绍

### 主要组件

- **ServiceCard** - 服务卡片组件
- **ConfigViewer** - 配置查看器（JSON 高亮）
- **DeepLinkGenerator** - Deep Link 生成器
- **ServiceForm** - 服务编辑表单
- **SearchFilter** - 搜索筛选组件

## 项目结构

```
MCPark/
├── backend/                    # 后端代码
│   ├── cmd/                   # 主程序入口
│   ├── internal/              # 内部代码
│   │   ├── api/              # API处理器
│   │   ├── model/            # 数据模型
│   │   ├── service/          # 业务逻辑
│   │   └── repository/       # 数据访问
│   ├── pkg/                  # 公共包
│   ├── configs/              # 配置文件
│   └── docs/                 # API文档
├── frontend/                  # 前端代码
│   ├── src/
│   │   ├── components/       # React组件
│   │   ├── pages/           # 页面组件
│   │   ├── services/        # API服务
│   │   ├── store/           # 状态管理
│   │   ├── utils/           # 工具函数
│   │   └── styles/          # 样式文件
│   ├── public/              # 静态资源
│   └── package.json
├── deployment/               # 部署配置
│   ├── docker-compose.yml
│   └── nginx.conf
├── docs/                    # 项目文档
└── README.md
```

## 技术栈

### 后端技术

- **Go 1.21+** - 高性能后端开发语言
- **Gin** - 轻量级 Web 框架，提供路由和中间件
- **GORM** - Go ORM 库，支持多种数据库
- **Viper** - 配置管理，支持多种配置格式
- **数据库**：
  - **SQLite** - 默认嵌入式数据库，零配置
  - **MySQL** - 可选关系型数据库，适合生产环境

### 前端技术

- **React 18** - 现代化前端框架
- **TypeScript** - 类型安全的 JavaScript
- **Vite** - 快速构建工具和开发服务器
- **shadcn/ui** - 基于 Radix UI 的现代组件库
- **Tailwind CSS** - 实用优先的 CSS 框架

### 开发工具

- **VS Code / Cursor** - 推荐 IDE，完整配置支持
- **Go Language Server** - Go 语言智能提示和调试
- **ESLint + Prettier** - 代码质量和格式化

## 开发计划

### Phase 1: 基础功能 (2-3 周)

1. 搭建基础项目结构
2. 实现后端基础 API
3. 实现前端基础页面
4. 完成服务的 CRUD 操作

### Phase 2: 核心功能 (2-3 周)

1. 实现 Deep Link 生成
2. 完善搜索和筛选功能
3. 添加分类管理
4. 优化 UI/UX

### Phase 3: 增强功能 (1-2 周)

1. 实现统计功能
2. 添加数据导入/导出
3. 添加配置验证和格式化
4. 性能优化

### Phase 4: 部署优化 (1 周)

1. Docker 容器化
2. CI/CD 配置
3. 文档完善
4. 测试和部署

## 技术细节和注意事项

### 安全考虑

1. JSON 配置验证，防止恶意内容
2. 输入 sanitization 和 XSS 防护
3. 文件上传安全检查
4. 基础的访问控制（可选 IP 白名单）

### 性能优化

1. 数据库索引优化
2. 前端懒加载
3. API 响应缓存
4. 静态资源 CDN

### 兼容性

1. 支持不同版本的 MCP 协议
2. 向后兼容的 API 设计
3. 响应式前端设计

## 结论

这个 MCP 服务收集管理系统将为 MCP 生态系统提供一个中心化的服务发现和分享平台。通过友好的 Web 界面和一键安装功能，降低用户使用 MCP 服务的门槛，促进社区发展。

项目采用现代化的技术栈，具有良好的可扩展性和维护性。通过分阶段开发，可以快速上线基础功能，然后逐步完善和优化。
