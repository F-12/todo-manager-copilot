# 个人Todo管理App

## 项目结构
- `backend/` Go后端服务，连接PostgreSQL，提供REST API
- `frontend/` React+TypeScript前端，ExpressJS开发服务器

## 启动方法

### 后端（Go）
1. 安装Go语言环境（推荐Go 1.21+）
2. 配置PostgreSQL数据库，默认连接参数：
   `host=localhost port=5432 user=postgres password=postgres dbname=todoapp sslmode=disable`
3. 进入`backend`目录，运行：
   ```sh
   go run main.go
   ```

### 前端（React+ExpressJS）
1. 进入`frontend`目录，安装依赖：
   ```sh
   npm install
   ```
2. 构建前端：
   ```sh
   npm run build
   ```
3. 启动Express服务器：
   ```sh
   npm start
   ```
4. 访问：http://localhost:3000

## 功能说明
- Todo列表展示、添加、删除、完成状态切换
- 前端通过`/api/todos`接口与后端交互
- ExpressJS代理API请求到Go后端

## 备注
- 如需自定义数据库连接，请设置环境变量`POSTGRES_CONN`
- 后端API需补充完整CRUD逻辑
- 前端UI可根据需求扩展美化
