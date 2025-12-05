# 云剪贴板项目部署指南

## 一、项目概述

云剪贴板是一个基于React + Gin的全栈项目，提供文本和文件的跨设备同步功能。

## 二、前端部署

### 1. 构建前端项目

确保前端项目已安装依赖并构建完成：

```bash
cd frontend
npm install  # 如果尚未安装依赖
npm run build
```

构建完成后，会在`frontend/dist`目录下生成静态文件。

### 2. 部署方式

#### 方式一：与后端一起部署

本项目的后端已经配置了静态文件服务，可以直接托管前端文件。

**步骤：**

1. **创建web目录**
   ```bash
   mkdir -p backend/web
   ```

2. **复制前端静态文件**
   ```bash
   cp -r frontend/dist/* backend/web/
   ```

3. **启动后端服务**
   ```bash
   cd backend
   go run main.go
   ```

4. **访问应用**
   浏览器访问 `http://服务器IP:3000`

#### 方式二：使用Nginx单独部署

**步骤：**

1. **安装Nginx**
   ```bash
   # Ubuntu/Debian
   sudo apt update
   sudo apt install nginx
   ```

2. **配置Nginx**
   创建配置文件 `/etc/nginx/conf.d/cloud-clipboard.conf`：
   
   ```nginx
   server {
       listen 80;
       server_name your-domain.com;  # 替换为你的域名或IP
       
       root /path/to/frontend/dist;
       index index.html;
       
       location / {
           try_files $uri $uri/ /index.html;
       }
       
       # API代理
       location /api {
           proxy_pass http://localhost:3000;
           proxy_set_header Host $host;
           proxy_set_header X-Real-IP $remote_addr;
           proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
       }
       
       # 文件上传目录代理
       location /uploads {
           proxy_pass http://localhost:3000;
           proxy_set_header Host $host;
           proxy_set_header X-Real-IP $remote_addr;
           proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
       }
   }
   ```

3. **启动Nginx**
   ```bash
   sudo systemctl start nginx
   sudo systemctl enable nginx
   ```

4. **访问应用**
   浏览器访问 `http://服务器IP` 或你的域名

## 三、后端部署

### 1. 编译后端

```bash
cd backend
go build -o cloud-clipboard
```

### 2. 配置环境

确保服务器上有以下目录：

```bash
mkdir -p backend/uploads
mkdir -p backend/data
```

### 3. 启动后端服务

```bash
# 直接启动
./cloud-clipboard

# 或使用systemd管理（推荐）
```

#### 使用systemd管理服务

1. 创建服务文件 `/etc/systemd/system/cloud-clipboard.service`：
   
   ```ini
   [Unit]
   Description=Cloud Clipboard Service
   After=network.target
   
   [Service]
   Type=simple
   User=www-data
   WorkingDirectory=/path/to/backend
   ExecStart=/path/to/backend/cloud-clipboard
   Restart=on-failure
   
   [Install]
   WantedBy=multi-user.target
   ```

2. 启动服务
   ```bash
   sudo systemctl daemon-reload
   sudo systemctl start cloud-clipboard
   sudo systemctl enable cloud-clipboard
   ```

3. 查看服务状态
   ```bash
   sudo systemctl status cloud-clipboard
   ```

## 四、注意事项

1. **端口配置**
   - 默认后端端口为3000
   - 如果需要修改端口，请编辑`backend/app/config/config.go`文件中的`ServerConfig`部分

2. **CORS配置**
   - 后端已配置CORS，允许所有来源访问
   - 如果需要限制来源，请修改`backend/main.go`文件中的`cors.Config`部分

3. **文件存储**
   - 上传的文件默认存储在`backend/uploads`目录
   - 文件元数据存储在`backend/data/files.json`文件
   - 请确保这些目录有足够的存储空间和权限

4. **安全建议**
   - 建议使用HTTPS协议
   - 可以配置Nginx添加SSL证书
   - 定期备份文件和数据

## 五、常见问题

1. **前端访问API失败**
   - 检查API代理配置是否正确
   - 确保后端服务正在运行

2. **文件上传失败**
   - 检查上传目录权限
   - 确保文件大小不超过配置的限制（默认16MB）

3. **服务启动失败**
   - 检查端口是否被占用
   - 查看日志文件了解详细错误信息

## 六、升级说明

1. **前端升级**
   ```bash
   cd frontend
   git pull
   npm install
   npm run build
   cp -r dist/* ../backend/web/
   ```

2. **后端升级**
   ```bash
   cd backend
   git pull
   go build -o cloud-clipboard
   sudo systemctl restart cloud-clipboard
   ```
