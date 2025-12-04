# 云剪切板后端

## 项目结构

```
backend/
├── main.go                 # 主入口文件
├── app/                    # 应用核心代码
│   ├── api/                # API层（控制器和路由）
│   ├── services/           # 业务逻辑层
│   ├── models/             # 数据模型
│   └── config/             # 配置管理
├── internal/               # 内部包（不对外暴露）
│   ├── clipboard/          # 字符串剪切板相关
│   ├── file/               # 文件管理相关
│   └── utils/              # 工具函数
├── pkg/                    # 可以对外暴露的包
├── data/                   # 数据存储目录
├── uploads/                # 文件上传目录
├── go.mod                  # Go模块文件
└── go.sum                  # Go模块校验和
```

## 分层架构

1. **API层**：处理HTTP请求和响应，定义路由和控制器
2. **业务逻辑层**：实现核心业务逻辑
3. **数据模型层**：定义数据结构和存储方式
4. **配置管理层**：处理应用配置

## 核心功能

### 字符串剪切板
- 支持内存总量限制
- 保证最近访问优先排序
- 支持上传、查看、复制、删除操作

### 文件管理
- 支持文件上传和下载
- 实现了下载次数限制
- 支持文件删除功能
- 支持定期删除过期文件
- 支持文件大小限制和总存储空间限制
- 实现了传输速度限制

## 运行方式

```bash
# 直接运行
go run main.go

# 编译后运行
go build -o cloud-clipboard
./cloud-clipboard
```

## 部署方案

### 环境要求

1. **操作系统**：Linux/macOS/Windows
2. **Go版本**：1.20+（用于后端）
3. **Node.js版本**：18+（用于前端构建）
4. **浏览器支持**：Chrome 80+, Firefox 75+, Safari 13+, Edge 80+

### 部署架构

```
┌─────────────────────────────────────────────────────────────────┐
│                         客户端浏览器                            │
└─────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
┌─────────────────────────────────────────────────────────────────┐
│                       反向代理服务器 (Nginx)                     │
├─────────────────────────────────────────────────────────────────┤
│  静态文件服务 (前端资源)                 API请求转发 (后端服务)   │
│  http://your-domain.com/              http://your-domain.com/api│
└─────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
┌─────────────────────────────────────────────────────────────────┐
│                         后端服务 (Go)                           │
└─────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
┌─────────────────────────────────────────────────────────────────┐
│                     本地文件系统 (数据存储)                      │
│  - 上传文件目录: ./uploads                                      │
│  - 元数据文件: ./data/files.json                                │
└─────────────────────────────────────────────────────────────────┘
```

### 1. 前端部署

#### 1.1 构建前端项目

```bash
# 进入前端目录
cd ../frontend

# 安装依赖
npm install

# 构建生产版本
npm run build

# 构建产物将生成在 ./dist 目录中
```

#### 1.2 部署静态资源

将构建后的 `dist` 目录中的所有文件部署到静态文件服务器（如Nginx、Apache等）。

### 2. 后端部署

#### 2.1 构建后端项目

```bash
# 进入后端目录
cd backend

# 安装依赖
go mod download

# 构建生产版本
go build -o cloud-clipboard main.go
```

#### 2.2 配置文件

后端服务使用默认配置，如需自定义配置，可以修改代码中的默认配置：

- **端口配置**：默认端口为3000
- **文件上传目录**：默认目录为 `./uploads`
- **元数据文件**：默认文件为 `./data/files.json`
- **最大文件大小**：默认10MB
- **总存储限制**：默认1GB
- **最大下载次数**：默认5次
- **文件过期时间**：默认7天

#### 2.3 启动后端服务

```bash
# 直接启动
./cloud-clipboard

# 后台启动（Linux/macOS）
nohup ./cloud-clipboard > cloud-clipboard.log 2>&1 &

# 使用systemd管理（Linux）
# 创建服务文件: /etc/systemd/system/cloud-clipboard.service
# 内容见下文
```

### 3. 反向代理配置

#### 3.1 Nginx配置示例

```nginx
server {
    listen 80;
    server_name your-domain.com;
    
    # 前端静态资源
    location / {
        root /path/to/frontend/dist;
        index index.html;
        try_files $uri $uri/ /index.html;
    }
    
    # 后端API
    location /api {
        proxy_pass http://localhost:3000;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
    
    # 文件上传目录（可选，用于直接访问上传的文件）
    location /uploads {
        alias /path/to/backend/uploads;
        expires 30d;
    }
}
```

#### 3.2 Systemd服务配置

创建 `/etc/systemd/system/cloud-clipboard.service` 文件：

```ini
[Unit]
Description=Cloud Clipboard Service
After=network.target

[Service]
Type=simple
User=www-data
WorkingDirectory=/path/to/backend
ExecStart=/path/to/backend/cloud-clipboard
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```

启用并启动服务：

```bash
# 重新加载systemd配置
sudo systemctl daemon-reload

# 启用服务
sudo systemctl enable cloud-clipboard

# 启动服务
sudo systemctl start cloud-clipboard

# 查看服务状态
sudo systemctl status cloud-clipboard

# 查看日志
sudo journalctl -u cloud-clipboard -f
```

### 4. 容器化部署（Docker）

#### 4.1 后端Dockerfile

在 `backend` 目录中创建 `Dockerfile`：

```dockerfile
# 使用官方Go镜像作为构建环境
FROM golang:1.22-alpine AS builder

# 设置工作目录
WORKDIR /app

# 复制go.mod和go.sum文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -o cloud-clipboard main.go

# 使用轻量级镜像作为运行环境
FROM alpine:latest

# 设置工作目录
WORKDIR /app

# 复制构建产物
COPY --from=builder /app/cloud-clipboard .

# 创建必要的目录
RUN mkdir -p data uploads

# 暴露端口
EXPOSE 3000

# 启动应用
CMD ["./cloud-clipboard"]
```

#### 4.2 前端Dockerfile

在 `frontend` 目录中创建 `Dockerfile`：

```dockerfile
# 使用官方Node.js镜像作为构建环境
FROM node:18-alpine AS builder

# 设置工作目录
WORKDIR /app

# 复制package.json和package-lock.json
COPY package*.json ./

# 安装依赖
RUN npm install

# 复制源代码
COPY . .

# 构建应用
RUN npm run build

# 使用Nginx作为静态文件服务器
FROM nginx:alpine

# 复制Nginx配置
COPY nginx.conf /etc/nginx/conf.d/default.conf

# 复制构建产物
COPY --from=builder /app/dist /usr/share/nginx/html

# 暴露端口
EXPOSE 80

# 启动Nginx
CMD ["nginx", "-g", "daemon off;"]
```

在 `frontend` 目录中创建 `nginx.conf`：

```nginx
server {
    listen 80;
    server_name localhost;
    
    location / {
        root /usr/share/nginx/html;
        index index.html;
        try_files $uri $uri/ /index.html;
    }
    
    location /api {
        proxy_pass http://backend:3000;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

#### 4.3 Docker Compose配置

在项目根目录创建 `docker-compose.yml`：

```yaml
version: '3.8'

services:
  backend:
    build:
      context: ./backend
    ports:
      - "3000:3000"
    volumes:
      - ./backend/data:/app/data
      - ./backend/uploads:/app/uploads
    restart: always
    environment:
      - GIN_MODE=release

  frontend:
    build:
      context: ./frontend
    ports:
      - "80:80"
    depends_on:
      - backend
    restart: always
```

使用Docker Compose启动服务：

```bash
docker-compose up -d
```

### 5. 安全性建议

1. **使用HTTPS**：配置SSL证书，使用HTTPS协议保护数据传输
2. **限制访问IP**：在防火墙中限制后端服务的访问IP
3. **定期备份数据**：定期备份 `data` 和 `uploads` 目录
4. **更新依赖**：定期更新前端和后端的依赖包，修复安全漏洞
5. **配置强密码**：如果后续添加认证功能，使用强密码策略

### 6. 监控与日志

1. **后端日志**：默认输出到标准输出，可以通过重定向保存到文件
2. **Nginx日志**：默认日志路径为 `/var/log/nginx/access.log` 和 `/var/log/nginx/error.log`
3. **系统监控**：可以使用Prometheus + Grafana监控系统资源使用情况
4. **应用监控**：可以在后端添加监控指标，如请求量、响应时间等

### 7. 常见问题及解决方案

1. **端口被占用**：
   - 检查端口使用情况：`lsof -i :3000` 或 `netstat -tlnp | grep 3000`
   - 修改后端端口：修改 `main.go` 中的配置

2. **上传文件失败**：
   - 检查文件大小是否超过限制
   - 检查总存储空间是否已满
   - 检查上传目录权限

3. **前端无法访问后端API**：
   - 检查Nginx配置是否正确
   - 检查后端服务是否正常运行
   - 检查防火墙设置

4. **下载次数不准确**：
   - 检查元数据文件权限
   - 检查后端服务是否正常重启

### 8. 升级步骤

1. **备份数据**：备份 `data` 和 `uploads` 目录
2. **更新代码**：拉取最新代码
3. **重新构建**：重新构建前端和后端
4. **重启服务**：重启后端服务和Nginx

### 9. 性能优化

1. **启用Gzip压缩**：在Nginx中配置Gzip压缩，减少传输大小
2. **启用缓存**：为静态资源配置适当的缓存策略
3. **调整并发数**：根据服务器配置调整后端服务的并发处理能力
4. **优化数据库**：如果后续添加数据库，可以优化数据库查询
5. **使用CDN**：对于静态资源，可以使用CDN加速访问

## 开发环境配置

### 前端开发

```bash
# 进入前端目录
cd ../frontend

# 安装依赖
npm install

# 启动开发服务器
npm run dev
```

### 后端开发

```bash
# 进入后端目录
cd backend

# 安装依赖
go mod download

# 启动开发服务器
go run main.go
```

## 测试

### API测试

可以使用Postman、curl等工具测试API端点：

```bash
# 测试健康检查
curl http://localhost:3000/health

# 测试获取文件列表
curl http://localhost:3000/api/files
```

### 前端测试

在浏览器中访问 `http://localhost:5173`（开发环境）或 `http://localhost`（生产环境），测试前端功能。

## 贡献指南

1. Fork仓库
2. 创建特性分支：`git checkout -b feature/your-feature`
3. 提交代码：`git commit -am 'Add some feature'`
4. 推送到分支：`git push origin feature/your-feature`
5. 提交Pull Request

## 许可证

MIT License
