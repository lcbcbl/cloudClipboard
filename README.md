# 云剪切板应用

一个基于Node.js和React的云剪切板应用，支持跨设备共享字符串和文件。

## 功能特性

### 字符串剪切板
- ✅ 上传字符串
- ✅ 查看所有字符串（按最近访问排序）
- ✅ 复制字符串到本地剪贴板
- ✅ 删除字符串
- ✅ 内存存储（LRU缓存）
- ✅ 内存总量限制
- ✅ 单条字符串大小限制

### 文件管理
- ✅ 上传文件
- ✅ 下载文件（带速度限制）
- ✅ 删除文件
- ✅ 定期删除过期文件
- ✅ 文件大小限制
- ✅ 总存储空间限制
- ✅ 下载次数限制

## 技术栈

### 后端
- **语言**: Go 1.21
- **框架**: Gin
- **字符串存储**: 内存存储（LRU缓存）
- **文件存储**: 本地文件系统
- **元数据存储**: JSON文件

### 前端
- **框架**: React 18
- **构建工具**: Vite
- **UI组件库**: Ant Design
- **HTTP客户端**: Axios

## 项目结构

```
├── backend/                # 后端代码
│   ├── cmd/                # 命令行入口
│   │   └── api/            # API服务入口
│   │       └── main.go     # 主程序入口
│   ├── internal/           # 内部包
│   │   ├── config/         # 配置文件
│   │   ├── controllers/    # 控制器
│   │   ├── middleware/     # 中间件
│   │   ├── services/       # 业务逻辑
│   │   └── utils/          # 工具函数
│   ├── uploads/            # 文件存储目录
│   ├── data/               # 元数据存储
│   └── go.mod              # Go模块配置
├── frontend/               # 前端代码
│   ├── src/                # 源代码
│   │   ├── components/     # React组件
│   │   ├── hooks/          # 自定义Hooks
│   │   ├── services/       # API服务
│   │   ├── utils/          # 工具函数
│   │   ├── App.jsx         # 应用入口
│   │   ├── main.jsx        # 渲染入口
│   │   └── index.css       # 全局样式
│   ├── index.html          # HTML模板
│   ├── vite.config.js      # Vite配置
│   └── package.json        # 依赖配置
├── README.md               # 项目说明
└── TECHNICAL_SPEC.md       # 技术方案
```

## 安装和运行

### 前提条件
- Node.js 16+ 和 npm 8+（前端）
- Go 1.21+（后端）

### 后端安装和运行

1. 进入后端目录
```bash
cd backend
```

2. 初始化Go模块（如果尚未初始化）
```bash
go mod tidy
```

3. 运行后端服务
```bash
# 开发模式
go run cmd/api/main.go

# 编译并运行
go build -o cloud-clipboard cmd/api/main.go
./cloud-clipboard
```

后端服务将在 `http://localhost:3000` 运行

### 前端安装和运行

1. 进入前端目录
```bash
cd frontend
```

2. 安装依赖
```bash
npm install
```

3. 运行前端服务
```bash
# 开发模式
npm run dev

# 生产构建
npm run build

# 预览生产构建
npm run preview
```

前端服务将在 `http://localhost:5173` 运行

## API文档

### 字符串剪切板API

| 方法 | 路径 | 功能 |
|------|------|------|
| POST | /api/clipboard/text | 上传字符串 |
| GET | /api/clipboard/text | 获取所有字符串（按最近访问排序） |
| GET | /api/clipboard/text/:id | 获取指定字符串 |
| DELETE | /api/clipboard/text/:id | 删除指定字符串 |

### 文件API

| 方法 | 路径 | 功能 |
|------|------|------|
| POST | /api/files | 上传文件 |
| GET | /api/files | 获取所有文件列表 |
| GET | /api/files/:id | 获取文件信息 |
| GET | /api/files/:id/download | 下载文件（带速度限制） |
| DELETE | /api/files/:id | 删除文件 |

### 健康检查API

| 方法 | 路径 | 功能 |
|------|------|------|
| GET | /health | 检查服务健康状态 |

## 配置说明

### 后端配置

配置文件位于 `backend/src/config/index.js`，可根据需要修改以下配置：

- **服务器配置**: 端口、主机
- **字符串剪切板配置**: 最大内存、最大项数、单项最大大小
- **文件配置**: 上传目录、元数据文件、最大文件大小、总存储限制、最大下载次数、速度限制、清理间隔、最大文件年龄

### 前端配置

配置文件位于 `frontend/vite.config.js`，主要配置代理服务器地址。

## 性能优化建议

1. **后端优化**
   - 使用PM2进行进程管理和负载均衡
   - 考虑使用Redis替代内存存储，支持分布式部署
   - 实现文件分片上传，支持大文件上传
   - 添加缓存机制，减少重复请求

2. **前端优化**
   - 使用CDN加速静态资源
   - 实现组件懒加载
   - 添加图片压缩和优化
   - 实现请求防抖和节流

3. **部署优化**
   - 使用Nginx作为反向代理和静态资源服务器
   - 配置HTTPS
   - 实现自动化部署
   - 添加监控和日志系统

## 测试

### 后端测试

```bash
cd backend
go test ./...
```

### 前端测试

```bash
cd frontend
npm run lint
```

## 许可证

MIT
