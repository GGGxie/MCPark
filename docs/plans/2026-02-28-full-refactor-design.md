# MCPark 完整重构设计文档

**日期：** 2026-02-28  
**目标：** 参照 SkillsHub (https://skills.gggxie.com/) 代码风格，全面重构 MCPark (https://mcpark.gggxie.com/)

---

## 一、架构决策

### 后端
- **风格：** SkillsHub 风格，Handler 直接操作 SQLite DB，去掉 Service/Repository 中间层
- **技术栈：** Gin + `database/sql` + SQLite + JWT + Google OAuth2/ID Token
- **配置：** 纯环境变量（`.env` + godotenv），去掉 YAML 配置
- **API 前缀：** `/api`（去掉 `/v1`）

### 前端
- **版本：** React 18 + Vite 5 + TypeScript + Tailwind 3（对标 SkillsHub）
- **UI 风格：** 完全对标 SkillsHub，内容替换为 MCP 相关

---

## 二、数据库表设计

```sql
-- 用户表（新增 is_admin 字段）
CREATE TABLE users (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  google_id TEXT UNIQUE,
  email TEXT UNIQUE NOT NULL,
  name TEXT NOT NULL,
  avatar TEXT DEFAULT '',
  is_admin INTEGER DEFAULT 0,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP
)

-- MCP 服务表（核心）
CREATE TABLE services (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  name TEXT NOT NULL,                    -- 技术名称（slug）
  title TEXT NOT NULL,                   -- 显示名（英文）
  title_zh TEXT DEFAULT '',              -- 显示名（中文）
  description TEXT NOT NULL,             -- 简介（英文）
  desc_zh TEXT DEFAULT '',               -- 简介（中文）
  category TEXT NOT NULL,                -- 分类
  icon TEXT DEFAULT '🤖',
  image TEXT DEFAULT '',
  config_json TEXT DEFAULT '',           -- MCP 配置示例（JSON）
  content TEXT DEFAULT '',               -- 详细文档（英文 Markdown）
  content_zh TEXT DEFAULT '',            -- 详细文档（中文 Markdown）
  tags TEXT DEFAULT '',
  author_id INTEGER NOT NULL,
  author_name TEXT DEFAULT '',
  author_avatar TEXT DEFAULT '',
  status TEXT DEFAULT 'pending',         -- pending/approved/rejected
  featured INTEGER DEFAULT 0,
  likes INTEGER DEFAULT 0,
  views INTEGER DEFAULT 0,
  install_count INTEGER DEFAULT 0,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (author_id) REFERENCES users(id)
)

-- 评论表
CREATE TABLE comments (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  service_id INTEGER NOT NULL,
  user_id INTEGER NOT NULL,
  user_name TEXT DEFAULT '',
  avatar TEXT DEFAULT '',
  content TEXT NOT NULL,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (service_id) REFERENCES services(id),
  FOREIGN KEY (user_id) REFERENCES users(id)
)

-- 点赞表
CREATE TABLE likes (
  user_id INTEGER NOT NULL,
  service_id INTEGER NOT NULL,
  PRIMARY KEY (user_id, service_id)
)

-- 文章表（用于 Learn 页面）
CREATE TABLE articles (
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
)
```

---

## 三、API 端点设计

```
GET  /api/auth/google              → OAuth 跳转 URL
GET  /api/auth/google/callback     → OAuth 回调（重定向到前端）
POST /api/auth/google/token        → ID Token 登录
GET  /api/auth/me                  → 当前用户（需登录）

GET    /api/services               → 列表（分页、分类、搜索、状态默认 approved）
GET    /api/services/featured      → 精选服务
GET    /api/services/categories    → 分类列表
GET    /api/services/:id           → 详情（访问量+1）
POST   /api/services               → 创建（需登录，默认 pending）
POST   /api/services/:id/like      → 点赞/取消（需登录）
GET    /api/services/:id/comments  → 评论列表
POST   /api/services/:id/comments  → 发表评论（需登录）
POST   /api/services/:id/install   → 安装计数（无需登录）

# 管理员
GET    /api/admin/services         → 所有状态服务列表
POST   /api/admin/services/:id/approve → 审核通过
POST   /api/admin/services/:id/reject  → 审核拒绝
DELETE /api/admin/services/:id    → 删除服务

GET  /api/articles                 → 文章列表
GET  /api/articles/:id             → 文章详情
```

---

## 四、前端路由

```
/                  → HomePage（精选 + 列表 + 搜索）
/service/:id       → ServiceDetailPage（详情 + 评论 + 点赞）
/submit            → SubmitPage（提交 MCP，需登录）
/auth              → AuthPage（Google 登录）
/auth/callback     → AuthCallback（OAuth 回调处理）
/learn             → LearnPage（文章列表）
/learn/:id         → ArticleDetailPage（文章详情）
/help              → HelpPage（帮助 + FAQ）
/admin             → AdminPage（管理员审核，需管理员）
*                  → NotFoundPage
```

---

## 五、关键配置

**.env 文件：**
```
PORT=8086
DATABASE_PATH=./data/mcpark.db
JWT_SECRET=...
GOOGLE_CLIENT_ID=...
GOOGLE_CLIENT_SECRET=...
GOOGLE_REDIRECT_URL=https://mcpark.gggxie.com/api/auth/google/callback
FRONTEND_URL=https://mcpark.gggxie.com
ADMIN_EMAIL=gavinxie1219@gmail.com
```

**前端 .env：**
```
VITE_GOOGLE_CLIENT_ID=...
```

---

## 六、MCP 服务分类

```
Development Tools    → 开发工具
Productivity         → 效率工具
Data & Analytics     → 数据分析
Communication        → 通信协作
AI & ML              → AI 与机器学习
File Management      → 文件管理
Web & API            → Web 与 API
Database             → 数据库
```
