# MCP服务配置示例

本文档提供了一些常见的MCP服务配置示例，可以作为添加新服务时的参考。

## 数据库服务

### PostgreSQL
```json
{
  "command": "npx",
  "args": [
    "-y",
    "@modelcontextprotocol/server-postgres",
    "postgresql://localhost/mydb"
  ]
}
```

### SQLite
```json
{
  "command": "npx",
  "args": [
    "-y", 
    "@modelcontextprotocol/server-sqlite",
    "/path/to/database.db"
  ]
}
```

## 文件系统服务

### 文件系统访问
```json
{
  "command": "npx",
  "args": [
    "-y",
    "@modelcontextprotocol/server-filesystem",
    "/path/to/allowed/directory"
  ]
}
```

### Git仓库
```json
{
  "command": "npx",
  "args": [
    "-y",
    "@modelcontextprotocol/server-git",
    "/path/to/git/repo"
  ]
}
```

## Web服务

### Web爬虫
```json
{
  "command": "npx",
  "args": [
    "-y",
    "@modelcontextprotocol/server-fetch"
  ]
}
```

### Puppeteer浏览器自动化
```json
{
  "command": "npx",
  "args": [
    "-y",
    "@modelcontextprotocol/server-puppeteer"
  ]
}
```

## 开发工具

### 内存KV存储
```json
{
  "command": "npx",
  "args": [
    "-y",
    "@modelcontextprotocol/server-memory"
  ]
}
```

### 时间和日期
```json
{
  "command": "npx",
  "args": [
    "-y",
    "@modelcontextprotocol/server-time"
  ]
}
```

## 自定义服务

### Python脚本服务
```json
{
  "command": "python",
  "args": [
    "/path/to/your/mcp_server.py"
  ]
}
```

### Node.js服务
```json
{
  "command": "node",
  "args": [
    "/path/to/your/mcp-server.js"
  ]
}
```

### Go服务
```json
{
  "command": "/path/to/your/mcp-server",
  "args": []
}
```

## 配置最佳实践

### 1. 使用绝对路径
对于本地文件和可执行文件，建议使用绝对路径以避免路径问题：
```json
{
  "command": "/usr/local/bin/python3",
  "args": ["/home/user/scripts/my_mcp_server.py"]
}
```

### 2. 环境变量
对于需要配置的服务，可以通过环境变量传递配置：
```json
{
  "command": "node",
  "args": ["./server.js"],
  "env": {
    "DATABASE_URL": "postgresql://localhost/mydb",
    "API_KEY": "your-api-key"
  }
}
```

### 3. 参数验证
确保所有必要的参数都已提供，并且格式正确：
```json
{
  "command": "npx",
  "args": [
    "-y",
    "@modelcontextprotocol/server-postgres",
    "postgresql://user:password@host:port/database"
  ]
}
```

## 常见分类

建议按以下分类组织MCP服务：

- **数据库** - PostgreSQL, MySQL, SQLite, MongoDB等
- **文件系统** - 文件访问, Git, 文档处理等
- **Web服务** - HTTP客户端, 网页抓取, API接口等
- **开发工具** - 代码分析, 测试工具, 构建工具等
- **AI/ML** - 机器学习模型, 数据处理等
- **系统工具** - 系统监控, 日志分析等
- **自定义** - 用户自己开发的服务

## 配置验证

系统会对JSON配置进行基本验证：

1. **必要字段**: `command`字段是必需的
2. **数组格式**: `args`应该是字符串数组
3. **JSON格式**: 配置必须是有效的JSON格式
4. **安全检查**: 避免危险的命令和路径

## 故障排除

### 常见问题

1. **命令找不到**: 确保命令在PATH中或使用绝对路径
2. **权限问题**: 确保有执行权限
3. **端口冲突**: 如果服务使用固定端口，避免冲突
4. **依赖缺失**: 确保所需的npm包或Python包已安装

### 调试技巧

1. 在终端中手动运行命令确认可用性
2. 检查服务的日志输出
3. 验证配置参数的正确性
4. 确认网络连接和防火墙设置