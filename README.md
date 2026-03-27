# meaw-oa

一个基于 Go + Fiber v3 的轻量级 OA 后端服务，提供用户认证、验证码等功能。

## 技术栈

| 组件 | 技术 |
|------|------|
| Web 框架 | [Fiber v3](https://gofiber.io/) |
| ORM | [GORM](https://gorm.io/) |
| 数据库 | SQLite (纯 Go 驱动) |
| 认证 | JWT (golang-jwt/jwt/v5) |
| 密码加密 | bcrypt |
| 验证码 | base64Captcha |

## 快速开始

### 本地运行

```bash
# 克隆项目
git clone https://github.com/jiehui555/meaw-oa.git
cd meaw-oa

# 安装依赖
go mod tidy

# 运行
go run .
```

服务默认监听 `:3000`，日志写入 `app.log`，数据库文件为 `app.db`。

### 使用 Docker 运行

```bash
docker run -d \
  --name meaw-oa \
  -p 3000:3000 \
  -v ./data:/app/data \
  ghcr.io/jiehui555/meaw-oa:latest
```

挂载 `/app/data` 目录持久化数据：

```bash
docker run -d \
  --name meaw-oa \
  -p 3000:3000 \
  -e DB_PATH=/app/data/app.db \
  -e LOG_PATH=/app/data/app.log \
  -v ./data:/app/data \
  ghcr.io/jiehui555/meaw-oa:latest
```

### 使用 Docker Compose

```yaml
version: "3.8"
services:
  meaw-oa:
    image: ghcr.io/jiehui555/meaw-oa:latest
    container_name: meaw-oa
    ports:
      - "3000:3000"
    environment:
      - PORT=3000
      - DB_PATH=/app/data/app.db
      - LOG_PATH=/app/data/app.log
      - TZ=Asia/Shanghai
    volumes:
      - ./data:/app/data
    restart: unless-stopped
```

## 配置项

通过环境变量配置：

| 变量 | 默认值 | 说明 |
|------|--------|------|
| `PORT` | `3000` | 服务监听端口 |
| `DB_PATH` | `app.db` | SQLite 数据库文件路径 |
| `LOG_PATH` | `app.log` | 日志文件路径 |

## API 接口

| 方法 | 路径 | 描述 | 认证 |
|------|------|------|------|
| `GET` | `/api/captcha` | 获取验证码图片 | - |
| `POST` | `/api/login` | 用户登录 | 需验证码 |
| `POST` | `/api/refresh` | 刷新令牌 | - |
| `GET` | `/api/admin/dashboard` | 管理员仪表板 | JWT + Admin |

### 获取验证码

```bash
curl http://localhost:3000/api/captcha
```

```json
{
  "code": 0,
  "data": {
    "captcha_id": "...",
    "captcha_image": "base64..."
  }
}
```

### 用户登录

```bash
curl -X POST http://localhost:3000/api/login \
  -H "Content-Type: application/json" \
  -d '{
    "name": "admin",
    "password": "password",
    "captcha_id": "...",
    "captcha_answer": "..."
  }'
```

响应：

```json
{
  "code": 0,
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIs..."
  }
}
```

默认管理员账号：`admin` / `password`

### 刷新令牌

```bash
curl -X POST http://localhost:3000/api/refresh \
  -H "Content-Type: application/json" \
  -d '{"refresh_token": "eyJhbGciOiJIUzI1NiIs..."}'
```

### 访问受保护接口

```bash
curl http://localhost:3000/api/admin/dashboard \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIs..."
```

## Docker 镜像

镜像托管在 GitHub Container Registry：

```bash
ghcr.io/jiehui555/meaw-oa
```

可用标签：

| 标签 | 说明 |
|------|------|
| `latest` | 最新版本 |
| `v1.2.3` | 精确版本 |
| `1.2` | 主版本.次版本 |

支持平台：`linux/amd64`、`linux/arm64`

## Release

项目使用 [GoReleaser](https://goreleaser.com/) 自动构建跨平台二进制文件：

- Linux: `amd64` / `arm64`
- Windows: `amd64` / `arm64`
- macOS: `amd64` / `arm64`

发布时推送 `v*` 标签即可触发自动构建：

```bash
git tag v1.0.0
git push origin v1.0.0
```

## 开发

```bash
# 编译
go build ./...

# 格式化
go fmt ./...

# 静态检查
go vet ./...

# 运行测试
go test -v ./...
```

## License

MIT
