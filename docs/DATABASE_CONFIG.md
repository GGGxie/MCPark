# 数据库配置指南

MCP Info 系统现在支持多种数据库类型，通过配置文件灵活切换。

## 🗂️ 配置文件位置

```
backend/
├── conf/
│   ├── config.yaml          # 默认配置 (SQLite)
│   └── config.mysql.yaml    # MySQL配置示例
└── data/
    └── MCPark.db          # SQLite数据库文件
```

## 🔧 支持的数据库类型

### 1. SQLite (默认推荐)

**优点：**

- ✅ 零配置，开箱即用
- ✅ 无需安装额外服务
- ✅ 数据文件便于备份和迁移
- ✅ 适合个人使用和小团队

**配置示例：**

```yaml
database:
  type: "sqlite"
  sqlite:
    path: "./data/MCPark.db"
```

### 2. MySQL

**优点：**

- ✅ 适合大规模生产环境
- ✅ 支持并发访问
- ✅ 丰富的管理工具
- ✅ 适合团队协作

**配置示例：**

```yaml
database:
  type: "mysql"
  mysql:
    host: "localhost"
    port: 3306
    username: "root"
    password: "your_password"
    database: "MCPark"
    charset: "utf8mb4"
    parseTime: true
    loc: "Local"
```

## 🚀 使用方法

### 使用 SQLite (默认)

1. **无需额外配置**，直接启动：

   ```bash
   ./start.sh
   ```

2. **数据自动存储在**：`backend/data/MCPark.db`

### 切换到 MySQL

1. **安装 MySQL 服务器**

2. **创建数据库**：

   ```sql
   CREATE DATABASE MCPark CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
   ```

3. **复制配置文件**：

   ```bash
   cp backend/conf/config.mysql.yaml backend/conf/config.yaml
   ```

4. **修改配置**：
   编辑 `backend/conf/config.yaml`，设置正确的：

   - 主机地址 (`host`)
   - 端口 (`port`)
   - 用户名 (`username`)
   - 密码 (`password`)
   - 数据库名 (`database`)

5. **启动服务**：
   ```bash
   ./start.sh
   ```

## 🔄 数据迁移

### 从 SQLite 迁移到 MySQL

1. **导出 SQLite 数据**：

   ```bash
   sqlite3 backend/data/MCPark.db .dump > backup.sql
   ```

2. **转换并导入 MySQL**：
   ```bash
   # 需要手动转换SQL语法或使用专门的迁移工具
   ```

### 从 MySQL 迁移到 SQLite

1. **导出 MySQL 数据**：

   ```bash
   mysqldump -u username -p MCPark > backup.sql
   ```

2. **转换并导入 SQLite**：
   ```bash
   # 需要手动转换SQL语法
   ```

## ⚙️ 环境变量配置

也可以通过环境变量覆盖配置：

```bash
# 设置数据库类型
export MCP_DATABASE_TYPE=mysql

# 设置MySQL连接信息
export MCP_DATABASE_MYSQL_HOST=localhost
export MCP_DATABASE_MYSQL_PORT=3306
export MCP_DATABASE_MYSQL_USERNAME=root
export MCP_DATABASE_MYSQL_PASSWORD=your_password
export MCP_DATABASE_MYSQL_DATABASE=MCPark

# 启动服务
./start.sh
```

## 🛠️ 故障排除

### SQLite 问题

**问题**：`database is locked`
**解决**：检查是否有其他进程正在使用数据库文件

**问题**：`no such table`
**解决**：删除数据库文件，重新启动服务自动重建

### MySQL 问题

**问题**：`connection refused`
**解决**：检查 MySQL 服务是否启动，防火墙设置

**问题**：`access denied`
**解决**：检查用户名、密码和权限设置

**问题**：`unknown database`
**解决**：手动创建数据库，或检查数据库名拼写

## 📊 性能对比

| 特性       | SQLite     | MySQL      |
| ---------- | ---------- | ---------- |
| 安装复杂度 | ⭐         | ⭐⭐⭐     |
| 并发支持   | ⭐⭐       | ⭐⭐⭐⭐⭐ |
| 数据量支持 | ⭐⭐⭐     | ⭐⭐⭐⭐⭐ |
| 备份便利性 | ⭐⭐⭐⭐⭐ | ⭐⭐⭐     |
| 管理工具   | ⭐⭐       | ⭐⭐⭐⭐⭐ |
| 资源占用   | ⭐⭐⭐⭐⭐ | ⭐⭐⭐     |

## 💡 建议

- **个人使用**：推荐 SQLite
- **团队协作**：推荐 MySQL
- **生产环境**：推荐 MySQL
- **开发测试**：推荐 SQLite

## 🔗 相关链接

- [SQLite 官方文档](https://www.sqlite.org/docs.html)
- [MySQL 官方文档](https://dev.mysql.com/doc/)
- [GORM 文档](https://gorm.io/docs/)
