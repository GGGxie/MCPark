package db

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

func InitDB(dbPath string) *sql.DB {
	if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
		log.Fatal("Failed to create database directory:", err)
	}

	database, err := sql.Open("sqlite3", dbPath+"?_journal_mode=WAL")
	if err != nil {
		log.Fatal("Failed to open database:", err)
	}

	if err = database.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	createTables(database)
	seedData(database)
	return database
}

func createTables(database *sql.DB) {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			google_id TEXT UNIQUE,
			email TEXT UNIQUE NOT NULL,
			name TEXT NOT NULL,
			avatar TEXT DEFAULT '',
			is_admin INTEGER DEFAULT 0,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS services (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT DEFAULT '',
			title TEXT NOT NULL,
			title_zh TEXT DEFAULT '',
			description TEXT NOT NULL,
			desc_zh TEXT DEFAULT '',
			category TEXT NOT NULL,
			icon TEXT DEFAULT '🤖',
			image TEXT DEFAULT '',
			config_json TEXT DEFAULT '',
			content TEXT DEFAULT '',
			content_zh TEXT DEFAULT '',
			tags TEXT DEFAULT '',
			author_id INTEGER NOT NULL,
			author_name TEXT DEFAULT '',
			author_avatar TEXT DEFAULT '',
			status TEXT DEFAULT 'pending',
			featured INTEGER DEFAULT 0,
			likes INTEGER DEFAULT 0,
			views INTEGER DEFAULT 0,
			install_count INTEGER DEFAULT 0,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (author_id) REFERENCES users(id)
		)`,
		`CREATE TABLE IF NOT EXISTS comments (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			service_id INTEGER NOT NULL,
			user_id INTEGER NOT NULL,
			user_name TEXT DEFAULT '',
			avatar TEXT DEFAULT '',
			content TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (service_id) REFERENCES services(id),
			FOREIGN KEY (user_id) REFERENCES users(id)
		)`,
		`CREATE TABLE IF NOT EXISTS likes (
			user_id INTEGER NOT NULL,
			service_id INTEGER NOT NULL,
			PRIMARY KEY (user_id, service_id)
		)`,
		`CREATE TABLE IF NOT EXISTS articles (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			title_zh TEXT DEFAULT '',
			description TEXT NOT NULL,
			desc_zh TEXT DEFAULT '',
			category TEXT DEFAULT 'general',
			content TEXT DEFAULT '',
			content_zh TEXT DEFAULT '',
			author_id INTEGER NOT NULL,
			author_name TEXT DEFAULT '',
			views INTEGER DEFAULT 0,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
	}

	for _, q := range queries {
		if _, err := database.Exec(q); err != nil {
			log.Fatal("Failed to create table:", err)
		}
	}
}

func seedData(database *sql.DB) {
	var count int
	database.QueryRow("SELECT COUNT(*) FROM services").Scan(&count)
	if count > 0 {
		return
	}

	database.Exec(`INSERT INTO users (id, google_id, email, name, avatar, is_admin) VALUES (1, 'system', 'system@mcpark.cc', 'MCPark', '', 0)`)

	services := []struct {
		name, title, titleZH, desc, descZH, category, icon, configJSON, content, contentZH, tags string
		featured                                                                                   bool
		likes, views, installCount                                                                 int
	}{
		{
			name:    "filesystem",
			title:   "Filesystem MCP",
			titleZH: "文件系统 MCP",
			desc:    "Read and write files on your local filesystem with full path control",
			descZH:  "完全控制本地文件系统的读写操作",
			category: "File Management",
			icon:    "📁",
			configJSON: `{
  "mcpServers": {
    "filesystem": {
      "command": "npx",
      "args": ["-y", "@modelcontextprotocol/server-filesystem", "/path/to/allowed/files"]
    }
  }
}`,
			content: "# Filesystem MCP\n\nProvides file system access to Claude, enabling reading and writing of files within specified directories.\n\n## Features\n\n- Read file contents\n- Write and create files\n- List directory contents\n- Move and copy files\n- Create directories\n\n## Installation\n\n```json\n{\n  \"mcpServers\": {\n    \"filesystem\": {\n      \"command\": \"npx\",\n      \"args\": [\"-y\", \"@modelcontextprotocol/server-filesystem\", \"/path/to/files\"]\n    }\n  }\n}\n```\n\n## Usage\n\nAsk Claude to read, write, or manage files in your specified directory.",
			contentZH: "# 文件系统 MCP\n\n为 Claude 提供文件系统访问能力，允许在指定目录内读写文件。\n\n## 功能\n\n- 读取文件内容\n- 写入和创建文件\n- 列出目录内容\n- 移动和复制文件\n- 创建目录\n\n## 安装\n\n```json\n{\n  \"mcpServers\": {\n    \"filesystem\": {\n      \"command\": \"npx\",\n      \"args\": [\"-y\", \"@modelcontextprotocol/server-filesystem\", \"/path/to/files\"]\n    }\n  }\n}\n```",
			tags:         "filesystem,files,read,write",
			featured:     true,
			likes:        128,
			views:        3200,
			installCount: 890,
		},
		{
			name:    "github",
			title:   "GitHub MCP",
			titleZH: "GitHub MCP",
			desc:    "Interact with GitHub repositories, issues, PRs, and more via the GitHub API",
			descZH:  "通过 GitHub API 管理仓库、Issue、PR 等",
			category: "Development Tools",
			icon:    "🐙",
			configJSON: `{
  "mcpServers": {
    "github": {
      "command": "npx",
      "args": ["-y", "@modelcontextprotocol/server-github"],
      "env": {
        "GITHUB_PERSONAL_ACCESS_TOKEN": "<YOUR_TOKEN>"
      }
    }
  }
}`,
			content: "# GitHub MCP\n\nConnect Claude to GitHub for repository management, code review, and issue tracking.\n\n## Features\n\n- List and search repositories\n- Read file contents from repos\n- Create and manage issues\n- Review pull requests\n- Create branches and commits\n\n## Setup\n\n1. Generate a GitHub Personal Access Token\n2. Add to your MCP config\n\n```json\n{\n  \"mcpServers\": {\n    \"github\": {\n      \"command\": \"npx\",\n      \"args\": [\"-y\", \"@modelcontextprotocol/server-github\"],\n      \"env\": { \"GITHUB_PERSONAL_ACCESS_TOKEN\": \"your-token\" }\n    }\n  }\n}\n```",
			contentZH: "# GitHub MCP\n\n将 Claude 连接到 GitHub，实现仓库管理、代码审查和 Issue 跟踪。\n\n## 功能\n\n- 列出和搜索仓库\n- 读取仓库文件内容\n- 创建和管理 Issue\n- 审查 Pull Request\n- 创建分支和提交\n\n## 配置\n\n1. 生成 GitHub Personal Access Token\n2. 添加到 MCP 配置",
			tags:         "github,git,repository,issues,pr",
			featured:     true,
			likes:        96,
			views:        2800,
			installCount: 720,
		},
		{
			name:    "brave-search",
			title:   "Brave Search MCP",
			titleZH: "Brave 搜索 MCP",
			desc:    "Web and local search using Brave Search API with privacy focus",
			descZH:  "使用 Brave Search API 进行隐私优先的网络和本地搜索",
			category: "Web & API",
			icon:    "🔍",
			configJSON: `{
  "mcpServers": {
    "brave-search": {
      "command": "npx",
      "args": ["-y", "@modelcontextprotocol/server-brave-search"],
      "env": {
        "BRAVE_API_KEY": "<YOUR_BRAVE_API_KEY>"
      }
    }
  }
}`,
			content: "# Brave Search MCP\n\nGive Claude web search capabilities using the Brave Search API.\n\n## Features\n\n- Web search with privacy\n- Local business search\n- News search\n- Image search\n\n## Setup\n\n1. Get a Brave Search API key from https://brave.com/search/api/\n2. Configure your MCP settings\n\n```json\n{\n  \"mcpServers\": {\n    \"brave-search\": {\n      \"command\": \"npx\",\n      \"args\": [\"-y\", \"@modelcontextprotocol/server-brave-search\"],\n      \"env\": { \"BRAVE_API_KEY\": \"your-api-key\" }\n    }\n  }\n}\n```",
			contentZH: "# Brave Search MCP\n\n使用 Brave Search API 为 Claude 提供网络搜索能力。\n\n## 功能\n\n- 隐私保护的网络搜索\n- 本地商业搜索\n- 新闻搜索\n- 图片搜索",
			tags:         "search,web,brave,api",
			featured:     true,
			likes:        84,
			views:        2100,
			installCount: 560,
		},
		{
			name:    "sqlite",
			title:   "SQLite MCP",
			titleZH: "SQLite MCP",
			desc:    "Query and manage SQLite databases directly from Claude",
			descZH:  "直接从 Claude 查询和管理 SQLite 数据库",
			category: "Database",
			icon:    "🗄️",
			configJSON: `{
  "mcpServers": {
    "sqlite": {
      "command": "npx",
      "args": ["-y", "@modelcontextprotocol/server-sqlite", "--db-path", "/path/to/db.sqlite"]
    }
  }
}`,
			content: "# SQLite MCP\n\nEnable Claude to query and manage SQLite databases with full SQL support.\n\n## Features\n\n- Execute SQL queries\n- List tables and schema\n- Insert, update, delete records\n- Create and drop tables\n\n## Usage\n\nAsk Claude to \"query my database for all users created this week\" or similar natural language queries.",
			contentZH: "# SQLite MCP\n\n让 Claude 能够查询和管理 SQLite 数据库。\n\n## 功能\n\n- 执行 SQL 查询\n- 列出表和结构\n- 增删改记录\n- 创建和删除表",
			tags:         "sqlite,database,sql,query",
			featured:     false,
			likes:        67,
			views:        1800,
			installCount: 430,
		},
		{
			name:    "fetch",
			title:   "Fetch MCP",
			titleZH: "网络请求 MCP",
			desc:    "Fetch web content, convert HTML to Markdown, and extract data from URLs",
			descZH:  "抓取网页内容，将 HTML 转换为 Markdown，从 URL 提取数据",
			category: "Web & API",
			icon:    "🌐",
			configJSON: `{
  "mcpServers": {
    "fetch": {
      "command": "uvx",
      "args": ["mcp-server-fetch"]
    }
  }
}`,
			content: "# Fetch MCP\n\nAllow Claude to fetch web content and convert it to readable formats.\n\n## Features\n\n- Fetch HTML pages and convert to Markdown\n- Extract text from web pages\n- Follow redirects\n- Handle authentication headers\n\n## Requirements\n\n- Python with `uvx` installed\n- Internet access",
			contentZH: "# 网络请求 MCP\n\n允许 Claude 抓取网页内容并转换为可读格式。\n\n## 功能\n\n- 获取 HTML 页面并转换为 Markdown\n- 从网页中提取文本\n- 跟随重定向",
			tags:         "fetch,http,web,html,markdown",
			featured:     false,
			likes:        52,
			views:        1400,
			installCount: 380,
		},
		{
			name:    "postgres",
			title:   "PostgreSQL MCP",
			titleZH: "PostgreSQL MCP",
			desc:    "Connect to PostgreSQL databases for read-only queries and schema inspection",
			descZH:  "连接 PostgreSQL 数据库进行只读查询和结构检查",
			category: "Database",
			icon:    "🐘",
			configJSON: `{
  "mcpServers": {
    "postgres": {
      "command": "npx",
      "args": ["-y", "@modelcontextprotocol/server-postgres", "postgresql://user:password@localhost/dbname"]
    }
  }
}`,
			content: "# PostgreSQL MCP\n\nConnect Claude to PostgreSQL databases for intelligent query assistance.\n\n## Features\n\n- Read-only SQL queries\n- Table and schema discovery\n- Relationship mapping\n- Query explanation\n\n## Security Note\n\nThis server provides read-only access by default. Ensure your database user has minimal permissions.",
			contentZH: "# PostgreSQL MCP\n\n将 Claude 连接到 PostgreSQL 数据库，提供智能查询辅助。\n\n## 功能\n\n- 只读 SQL 查询\n- 表和结构发现\n- 关系映射\n\n## 安全说明\n\n默认提供只读访问，确保数据库用户权限最小化。",
			tags:         "postgres,database,sql,query",
			featured:     false,
			likes:        45,
			views:        1200,
			installCount: 290,
		},
	}

	for _, s := range services {
		featured := 0
		if s.featured {
			featured = 1
		}
		database.Exec(`INSERT INTO services (name, title, title_zh, description, desc_zh, category, icon, config_json, content, content_zh, tags, author_id, author_name, status, featured, likes, views, install_count)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 1, 'MCPark', 'approved', ?, ?, ?, ?)`,
			s.name, s.title, s.titleZH, s.desc, s.descZH, s.category, s.icon, s.configJSON,
			s.content, s.contentZH, s.tags, featured, s.likes, s.views, s.installCount)
	}

	articles := []struct {
		title, titleZH, desc, descZH, category, content, contentZH string
		views                                                        int
	}{
		{
			title:   "Getting Started with MCP",
			titleZH: "MCP 入门指南",
			desc:    "Learn what MCP is and how to configure your first MCP server",
			descZH:  "了解什么是 MCP 以及如何配置你的第一个 MCP 服务",
			category: "Getting Started",
			content: "# Getting Started with MCP\n\nModel Context Protocol (MCP) is an open standard that enables AI assistants to securely connect to external data sources and tools.\n\n## What is MCP?\n\nMCP provides a standardized way for AI models to:\n\n- Access local files and databases\n- Connect to external APIs\n- Execute code and commands\n- Integrate with services like GitHub, Slack, and more\n\n## How to Install an MCP Server\n\n### Step 1: Choose an MCP Server\n\nBrowse MCPark to find an MCP server that fits your needs.\n\n### Step 2: Configure Your Client\n\nAdd the MCP server configuration to your Claude Desktop or Cursor settings:\n\n```json\n{\n  \"mcpServers\": {\n    \"server-name\": {\n      \"command\": \"npx\",\n      \"args\": [\"-y\", \"@package/server-name\"]\n    }\n  }\n}\n```\n\n### Step 3: Restart Your Client\n\nRestart Claude Desktop or your IDE to apply the changes.\n\n### Step 4: Start Using\n\nYour AI assistant now has access to the new capabilities provided by the MCP server.",
			contentZH: "# MCP 入门指南\n\n模型上下文协议（MCP）是一个开放标准，使 AI 助手能够安全地连接到外部数据源和工具。\n\n## 什么是 MCP？\n\nMCP 为 AI 模型提供了标准化的方式来：\n\n- 访问本地文件和数据库\n- 连接到外部 API\n- 执行代码和命令\n- 与 GitHub、Slack 等服务集成\n\n## 如何安装 MCP 服务\n\n### 步骤 1：选择 MCP 服务\n\n在 MCPark 浏览，找到适合你需求的 MCP 服务。\n\n### 步骤 2：配置客户端\n\n将 MCP 服务配置添加到你的 Claude Desktop 或 Cursor 设置：\n\n```json\n{\n  \"mcpServers\": {\n    \"server-name\": {\n      \"command\": \"npx\",\n      \"args\": [\"-y\", \"@package/server-name\"]\n    }\n  }\n}\n```\n\n### 步骤 3：重启客户端\n\n重启 Claude Desktop 或 IDE 以应用更改。",
			views: 4200,
		},
		{
			title:   "Building Your Own MCP Server",
			titleZH: "构建你自己的 MCP 服务",
			desc:    "Step-by-step guide to creating and publishing your own MCP server",
			descZH:  "创建和发布自己 MCP 服务的分步指南",
			category: "Development",
			content: "# Building Your Own MCP Server\n\nCreating an MCP server allows you to extend AI capabilities with your own tools and data sources.\n\n## Prerequisites\n\n- Node.js 18+ or Python 3.10+\n- Basic understanding of APIs\n- An idea for what capability to expose\n\n## Quick Start with TypeScript\n\n```typescript\nimport { Server } from '@modelcontextprotocol/sdk/server/index.js';\nimport { StdioServerTransport } from '@modelcontextprotocol/sdk/server/stdio.js';\n\nconst server = new Server(\n  { name: 'my-server', version: '1.0.0' },\n  { capabilities: { tools: {} } }\n);\n\nserver.setRequestHandler('tools/call', async (request) => {\n  // Handle tool calls\n});\n\nconst transport = new StdioServerTransport();\nawait server.connect(transport);\n```\n\n## Best Practices\n\n1. **Security**: Never expose sensitive credentials\n2. **Error handling**: Return meaningful error messages\n3. **Documentation**: Write clear usage instructions\n4. **Testing**: Test with Claude Desktop before publishing",
			contentZH: "# 构建你自己的 MCP 服务\n\n创建 MCP 服务允许你用自己的工具和数据源扩展 AI 能力。\n\n## 前提条件\n\n- Node.js 18+ 或 Python 3.10+\n- 基本的 API 了解\n- 一个想要暴露的能力想法\n\n## TypeScript 快速开始\n\n```typescript\nimport { Server } from '@modelcontextprotocol/sdk/server/index.js';\n// ... 完整示例见英文版\n```\n\n## 最佳实践\n\n1. **安全性**：不要暴露敏感凭证\n2. **错误处理**：返回有意义的错误信息\n3. **文档**：编写清晰的使用说明",
			views: 3100,
		},
		{
			title:   "MCP vs Traditional APIs",
			titleZH: "MCP 与传统 API 对比",
			desc:    "Understanding the difference between MCP servers and traditional REST APIs",
			descZH:  "理解 MCP 服务与传统 REST API 之间的区别",
			category: "Concepts",
			content: "# MCP vs Traditional APIs\n\nBoth MCP servers and REST APIs expose functionality, but they serve different purposes.\n\n## Quick Comparison\n\n| Aspect | MCP Server | REST API |\n|--------|-----------|----------|\n| Consumer | AI Models | Applications |\n| Protocol | JSON-RPC over stdio | HTTP |\n| Discovery | Schema-based | OpenAPI/Swagger |\n| Auth | Client-managed | Token/OAuth |\n| Use case | AI tool use | App integration |\n\n## When to Use MCP\n\n- The consumer is an AI model (Claude, GPT, etc.)\n- You want AI to discover and use tools dynamically\n- Local tool execution is needed\n- Low-latency local access is important\n\n## When to Use REST API\n\n- The consumer is a web or mobile app\n- You need wide client compatibility\n- Public API exposure is required\n- Standard HTTP caching and security matter",
			contentZH: "# MCP 与传统 API 对比\n\nMCP 服务和 REST API 都暴露功能，但服务于不同目的。\n\n## 快速对比\n\n| 方面 | MCP 服务 | REST API |\n|------|---------|----------|\n| 消费者 | AI 模型 | 应用程序 |\n| 协议 | JSON-RPC over stdio | HTTP |\n| 发现方式 | 基于 Schema | OpenAPI/Swagger |\n| 认证 | 客户端管理 | Token/OAuth |\n\n## 何时使用 MCP\n\n- 消费者是 AI 模型\n- 需要 AI 动态发现和使用工具\n- 需要本地工具执行",
			views: 2600,
		},
	}

	for _, a := range articles {
		database.Exec(`INSERT INTO articles (title, title_zh, description, desc_zh, category, content, content_zh, author_id, author_name, views)
			VALUES (?, ?, ?, ?, ?, ?, ?, 1, 'MCPark', ?)`,
			a.title, a.titleZH, a.desc, a.descZH, a.category, a.content, a.contentZH, a.views)
	}
}
