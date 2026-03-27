# AGENTS.md - meaw-oa 项目指南

## 项目概述

Go 语言 Web 应用，使用 [Fiber v3](https://gofiber.io/) 框架。
模块路径：`github.com/jiehui555/meaw-oa`

---

## 构建与运行命令

| 命令 | 说明 |
|------|------|
| `go build ./...` | 编译整个项目 |
| `go run .` | 运行主程序 |
| `go fmt ./...` | 格式化代码 |
| `go vet ./...` | 静态分析检查 |
| `go mod tidy` | 整理依赖 |
| `go test ./...` | 运行所有测试 |
| `go test -v ./...` | 运行所有测试（详细输出） |
| `go test -run TestXxx ./...` | 运行单个测试函数 |
| `go test -run TestXxx ./path/to/pkg` | 运行指定包的单个测试 |
| `go test -count=1 ./...` | 禁用缓存运行测试 |

如安装了 `golangci-lint`：
- `golangci-lint run` — 运行完整静态检查

---

## 代码风格指南

### 文件结构

每个 `.go` 文件按以下顺序组织：
1. `package` 声明
2. `import` 块（分组：标准库 → 第三方 → 项目内部，空行分隔）
3. 类型定义
4. 常量/变量
5. 函数/方法

```go
package main

import (
    "fmt"
    "log/slog"

    "github.com/gofiber/fiber/v3"

    "github.com/jiehui555/meaw-oa/internal/config"
)
```

### 命名规范

- **包名**：小写单词，简短（如 `handler`、`model`、`config`），不用下划线
- **导出标识符**：`PascalCase`（如 `NewServer`、`UserModel`）
- **非导出标识符**：`camelCase`（如 `parseInput`、`dbConn`）
- **接口**：单方法接口用 `er` 后缀（如 `Reader`、`Handler`）
- **常量**：`PascalCase` 或分组 `const` 块
- **文件名**：`snake_case.go`
- **测试文件**：`xxx_test.go`，与被测文件同目录

### 错误处理

- 必须显式处理每个错误，**禁止** 使用 `_` 忽略
- 包装错误时使用 `fmt.Errorf("context: %w", err)`
- Fiber handler 返回 `error`，由框架统一处理
- 启动失败使用 `slog.Error` + `panic`
- 日志统一使用 `log/slog`，不要使用 `log` 标准包

```go
if err := app.Listen(":3000"); err != nil {
    slog.Error("server failed to start", "error", err)
    panic(err)
}
```

### 路由与 Handler

- 使用 Fiber v3 的 `app.Get()`、`app.Post()` 等方法注册路由
- Handler 签名：`func(c fiber.Ctx) error`
- 响应使用 `c.JSON()`、`c.SendString()` 等方法
- 参数获取：`c.Params("id")`、`c.Query("key")`、`c.Body()`

### 测试规范

- 测试函数命名：`TestFunctionName_Scenario`
- 使用标准 `testing` 包，子测试使用 `t.Run()`
- API 测试使用 Fiber 的 `app.Test()` + 内存 SQLite，不启动真实服务器
- 测试基础设施放在 `handler_test.go`，业务测试放在 `xxx_test.go`

```go
func TestLogin(t *testing.T) {
    db := setupTestDB(t)
    app := setupApp(t, db)

    t.Run("success", func(t *testing.T) {
        _, res := doRequest(t, app, "POST", "/api/login", map[string]string{
            "name": "admin", "password": "password",
        })
        if res.Code != 0 {
            t.Errorf("expected code 0, got %d", res.Code)
        }
    })
}
```

### 格式化

- 使用 `gofmt` 或 `goimports` 格式化，**缩进为 Tab**
- 行宽建议不超过 120 字符
- 函数之间空一行，逻辑块之间可空行增强可读性
- **不添加代码注释**，除非用户明确要求

### 并发

- 使用 `sync.Mutex` 或 `sync.RWMutex` 保护共享状态
- 避免 goroutine 泄漏，使用 `context.Context` 控制生命周期
- channel 操作注意死锁风险

### 目录结构建议（随着项目增长）

```
.
├── main.go
├── go.mod
├── internal/
│   ├── common/        # 通用工具（响应、JWT 等）
│   ├── config/        # 配置加载
│   ├── database/      # 数据库连接与迁移
│   ├── handler/       # HTTP 处理器
│   ├── logger/        # 日志初始化
│   ├── middleware/     # 中间件
│   └── model/         # 数据模型
└── docs/              # 文档（如需要）
```

---

## 注意事项

- **不要运行程序或测试**：生成代码后只运行 `go fmt ./...` 和 `go vet ./...` 做风格检查，禁止运行 `go run`、`go test` 等可能启动进程的命令
- Go 版本：1.25.0，使用新语法前确认兼容性
- 提交前建议用户手动运行完整测试：`go test ./...`
- 不提交 `go.sum` 以外的生成文件
