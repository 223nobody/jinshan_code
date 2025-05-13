## 基于前后端分离架构的试题管理系统，前端采用 React+Vite 技术栈，后端使用 Gin 框架实现 RESTful API，支持试题的创建、编辑、批量操作及多维统计。

**api 请求接口集**

1. 健康监测接口
   GET http://localhost:8080/health
2. AI 出题接口
   POST http://localhost:8080/api/questions/CreateByAI
3. 自主出题接口
   POST http://localhost:8080/api/questions/CreateByHand
4. 编辑题目接口
   PUT http://localhost:8080/api/questions/update
5. 批量插入接口
   POST http://localhost:8080/api/questions/batch-insert
6. 获取题目信息接口
   GET http://localhost:8080/api/stats/byid/:id
7. 批量删除接口
   DELETE http://localhost:8080/api/questions/batch-delete
8. 题目分页列表展示接口
   GET http://localhost:8080/api/stats/summary
9. 单选题分页列表展示接口
   GET http://localhost:8080/api/stats/bytype1
10. 多选题分页列表展示接口
    GET http://localhost:8080/api/stats/bytype2
11. 编程题分页列表展示接口
    GET http://localhost:8080/api/stats/bytype3

**数据库建表语句(表名 questions)**

```sql
CREATE TABLE IF NOT EXISTS questions (
   id INTEGER PRIMARY KEY AUTOINCREMENT,
   title TEXT NOT NULL,
   type INTEGER NOT NULL,
   language TEXT NOT NULL,
   answers TEXT COMMENT NOT NULL,
   rights TEXT COMMENT NOT NULL
);
```

## 项目结构

**项目前端结构说明**

```bash
client/
├── dist/                 # 构建输出
├── public/               # 静态资源
│   ├── images/           # SVG/PNG素材
│   └── index.html        # 入口模板
├── src/
│   │── QuestionBank.jsx  # 题库组件
│   │── QuestionEdit.jsx  # 题目编辑组件
│   │── QuestionCreateByHand.jsx  # 自主出题组件
│   │── QuestionCreateByAI.jsx    # AI出题组件
│   │── App.css
│   │── index.css
│   └── main.jsx          # 应用入口
├── vite.config.js        # Vite配置
└── package.json          # 依赖管理
```

**项目后端结构说明**

```bash
server/                      # 服务根目录
├── api/                     # API通用组件
│   └── response.go          # 统一响应格式封装
├── config/                  # 配置管理
│   └── config.go            # 配置加载与解析
├── controllers/             # 业务控制器
│   ├── actions.go           # 通用操作处理
│   └── question.go          # 题目业务逻辑
├── services/                # 服务层组件
│   ├── client.go            # 基础服务客户端
│   ├── deepseek.go          # 深度求索AI服务集成
│   └── tongyi.go            # 通义千问服务集成
├── storage/                 # 数据存储层
│   ├── database.go          # 数据库连接管理
│   └── storage.go           # 文件存储操作
├── log/                     # 日志目录
├── .env                     # 环境变量文件
├── .gitignore               # Git忽略配置
├── go.mod                   # Go模块依赖
├── go.sum                   # 依赖校验文件
├── main.go                  # 服务入口文件
└── question_service.db      # SQLite数据库文件
```

**快速启动**

```bash
## 开发环境
# 克隆仓库
git clone https://git.wpsit.cn/wust-2026-server-task/fukun.git
# 启动前端
cd client && npm install
npm run dev -- --port 3000
# 启动后端
cd ../server
go run main.go -port 8080 -env dev

## 生产构建
# 前端构建
cd client && npm run build
# 后端编译
cd ../server
go build -ldflags "-s -w" -o server main.go
```

## 核心功能

### 试题管理

​​ **智能创建** ​​：支持 AI 生成（QuestionCreateByAI.jsx）与手动录入(QuestionCreateByHand.jsx)  
​​ **批量操作** ​​：提供批量删除/导入/导出功能  
​​ **版本控制** ​​：试题修改历史追踪（通过 question_code 字段管理）

### 统计展示

​​ **全局概览** ​​：/stats/summary 展示总量统计  
​​ **分类查询** ​​：/stats/bytype[1-3] 分题型统计  
​​ **详情查看** ​​：/stats/byid/:id 显示完整题目数据

## 技术亮点

1. **前后端分离架构**：前端采用 React + Vite 技术栈，后端使用 Gin 框架实现 RESTful API，这种架构使得前后端开发可以并行进行，提高开发效率，同时也方便维护和扩展。
2. **模块化设计**：代码遵循模块化设计思想，将不同功能模块分离到不同文件中。如 services 包负责业务逻辑，controllers 包处理 HTTP 请求，storage 包负责数据存储，提高了代码的可维护性和可扩展性。
   配置管理：通过 config.LoadConfig() 函数加载配置，将配置信息与代码分离，便于在不同环境中进行配置管理，增强了代码的灵活性。
3. **AI 服务集成**：支持多种 AI 模型（如 deepseek 和 tongyi），通过 AIService 接口进行抽象，方便后续添加更多 AI 模型，增强了系统的扩展性。
4. **路由管理**：使用 Gin 框架进行路由管理，提供清晰的 API 接口，同时支持 CORS 配置，方便前后端分离开发。
5. **数据库支持**：使用 modernc.org/sqlite 作为数据库，支持数据的持久化存储，同时使用 github.com/jmoiron/sqlx 进行数据库操 5.作，简化了数据库开发。
6. **功能丰富**：支持试题的智能创建（AI 生成与手动录入）、批量操作（删除、导入、导出）以及版本控制（试题修改历史追踪），同时提供多维度的统计展示功能。

## 主要不足

1. **错误处理与日志记录**：从日志文件中可以看到存在 API 请求失败的情况，但整体项目的错误处理和日志记录可能不够完善。例如，没有详细记录失败的具体原因和上下文信息，不利于问题的排查和修复。
2. **AI 服务依赖**：系统依赖外部 AI 服务（如 DashScope）生成试题，当这些服务出现问题（如网络超时）时，会影响系统的正常使用，且缺乏有效的容错机制。
3. **缺乏测试覆盖**：从提供的代码和文档来看，没有明显的测试代码，缺乏对各个功能模块的单元测试、集成测试，可能导致代码在后续维护和扩展过程中出现潜在的问题。
4. **性能优化不足**：随着试题数据量的增加，数据库查询性能可能会受到影响。目前代码中没有明显的数据库索引优化、缓存机制等性能优化措施。
5. **安全性问题**：文档中未提及系统的安全机制，如用户认证、授权、数据加密等，可能存在安全隐患。

## 优化建议

1. **完善错误处理与日志记录**
   详细错误信息返回：在控制器中，当出现错误时，返回更详细的错误信息。

```go
// 示例修改 GenerateQuestion 函数中的错误处理
func (h *StatsHandler) GenerateQuestion(c *gin.Context) {
    var req config.QuestionRequest1
    if err := c.ShouldBindJSON(&req); err != nil {
        api.Error(c, http.StatusBadRequest, fmt.Sprintf("参数绑定失败: %v", err))
        return
    }
    id, err := h.db.CreateQuestion(&req)
    if err != nil {
        api.Error(c, http.StatusBadRequest, fmt.Sprintf("数据库创建失败: %v", err))
        return
    }
    // 其他逻辑...
}
```

2. 增强 AI 服务容错机制
   重试机制：在调用 AI 服务时，实现重试逻辑。

```go
package services
import (
    "time"
    "github.com/sirupsen/logrus"
)
const maxRetries = 3
func callAIService(req config.QuestionRequest) (interface{}, error) {
    var result interface{}
    var err error
    for i := 0; i < maxRetries; i++ {
        result, err = actualAICall(req)
        if err == nil {
            return result, nil
        }
        logrus.Errorf("AI 服务调用失败，第 %d 次重试: %v", i+1, err)
        time.Sleep(2 * time.Second) // 重试间隔 2 秒
   }
    return nil, fmt.Errorf("AI 服务调用失败，达到最大重试次数: %v", err)
}
func actualAICall(req config.QuestionRequest) (interface{}, error) {
    // 实际调用 AI 服务的代码
    return nil, nil
}
```

3. **性能优化**
   数据库索引优化：在数据库表中添加合适的索引，提高查询性能。

```sql
-- 在 questions 表中对 type 和 keyword 字段添加索引
CREATE INDEX idx_questions_type ON questions (type);
CREATE INDEX idx_questions_keyword ON questions (keyword);
```

4. **加强系统安全性**
   用户认证和授权：使用 JWT（JSON Web Token）实现用户认证和授权。

```go
package middleware
import (
    "errors"
    "net/http"
    "strings"
    "time"

    "github.com/dgrijalva/jwt-go"
    "github.com/gin-gonic/gin"
)
var jwtKey = []byte("your_secret_key")
type Claims struct {
    UserID string `json:"user_id"`
    jwt.StandardClaims
}
func GenerateToken(userID string) (string, error) {
    expirationTime := time.Now().Add(24 * time.Hour)
    claims := &Claims{
        UserID: userID,
        StandardClaims: jwt.StandardClaims{
            ExpiresAt: expirationTime.Unix(),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, err := token.SignedString(jwtKey)
    if err != nil {
        return "", err
    }
    return tokenString, nil
}
func Authenticate() gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "未提供授权信息"})
            c.Abort()
            return
        }

        tokenString := strings.Split(authHeader, " ")[1]
        claims := &Claims{}

        token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
            return jwtKey, nil
        })

        if err != nil ||!token.Valid {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的授权信息"})
            c.Abort()
            return
        }
        c.Set("user_id", claims.UserID)
        c.Next()
    }
}
```
