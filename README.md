# LumaGate — 轻量可插拔 API 网关

LumaGate 是一个用于快速演示与落地的轻量级 API 网关系统，包含：
- 路由转发（前缀匹配）
- 负载均衡（轮询 / 加权随机）
- 灰度流量（v1/v2 百分比）
- 插件：限流、CORS、API Key、IP 白名单
- 访问日志（文件）
- 可视化 Dashboard（Vue3 + Element Plus）

## 技术选型
- 网关后端：Go 1.21+，`net/http`
- 插件体系：统一接口 + 链式执行
- 前端：Vue3 + Vite + TypeScript + Element Plus
- 日志：文件 `logs/access.log`
- 部署：Docker / Docker Compose；前端使用 Nginx 静态托管并反向代理到网关

## 目录结构
```
lumagate/
├─ gateway/                 # Go API 网关
│  ├─ main.go
│  ├─ router.go             # 路由匹配 + 反向代理
│  ├─ admin.go              # Admin API（配置与日志）
│  ├─ state.go              # 内存路由与插件配置
│  ├─ load_balancer.go      # 轮询 / 加权随机
│  ├─ plugins_*.go          # 插件：auth、ip_whitelist、rate_limit、cors、traffic、logging
│  └─ Dockerfile
├─ services/                # 下游演示服务（Go）
│  ├─ v1/serviceA,B,C       # 9001/9002/9003
│  ├─ v2/serviceA,B,C       # 9011/9012/9013
│  └─ 每个目录含 Dockerfile
├─ dashboard/               # 前端 Dashboard
│  ├─ src/
│  ├─ nginx/default.conf    # Nginx 配置（/admin 代理到网关）
│  └─ Dockerfile
├─ docs/                    # 文档
│  ├─ LumaGate 演示指南.md
│  └─ API 接口文档.md
├─ docker-compose.yml       # 编排：网关、6个服务、Dashboard
├─ LICENSE
└─ go.mod
```

## 如何构建（本地）
- 构建网关与 6 个服务：
```
# 在仓库根目录
go build -o ./bin/lumagate-gateway ./gateway
# v1
go build -o ./bin/serviceA_v1 ./services/v1/serviceA
go build -o ./bin/serviceB_v1 ./services/v1/serviceB
go build -o ./bin/serviceC_v1 ./services/v1/serviceC
# v2
go build -o ./bin/serviceA_v2 ./services/v2/serviceA
go build -o ./bin/serviceB_v2 ./services/v2/serviceB
go build -o ./bin/serviceC_v2 ./services/v2/serviceC
```
- 启动 6 个服务与网关：
```
nohup ./bin/serviceA_v1 >/tmp/serviceA_v1.log 2>&1 &
nohup ./bin/serviceB_v1 >/tmp/serviceB_v1.log 2>&1 &
nohup ./bin/serviceC_v1 >/tmp/serviceC_v1.log 2>&1 &
nohup ./bin/serviceA_v2 >/tmp/serviceA_v2.log 2>&1 &
nohup ./bin/serviceB_v2 >/tmp/serviceB_v2.log 2>&1 &
nohup ./bin/serviceC_v2 >/tmp/serviceC_v2.log 2>&1 &
nohup ./bin/lumagate-gateway >/tmp/gateway.log 2>&1 &
```
- 启动前端（开发模式）：
```
cd dashboard
npm install
npm run dev
```
- 访问：
  - Dashboard：`http://localhost:5173`
  - 网关 Admin：`http://localhost:8080/admin/routes`
  - 转发示例：`http://localhost:8080/api/users/user`

## 如何部署（Docker / Compose）
- 单机部署：
```
docker compose build
docker compose up -d
```
- 访问：
  - Dashboard：`http://<服务器IP>/`
  - 网关 Admin：`http://<服务器IP>:8080/admin/routes`
- 说明：
  - `dashboard/nginx/default.conf` 将 `/admin` 请求反向代理到 `lumagate-gateway:8080/admin`
  - 网关默认上游地址已指向容器服务名：
    - v1：`serviceA_v1:9001`、`serviceB_v1:9002`、`serviceC_v1:9003`
    - v2：`serviceA_v2:9011`、`serviceB_v2:9012`、`serviceC_v2:9013`

## 项目内容与特性
- 路由：前缀匹配（最长前缀优先），默认 `/api/users`；命中后进入插件链，再负载到上游。
- 负载均衡：
  - 轮询（默认关闭权重时）
  - 加权随机（开启权重时，依据 `Upstream.Weight`）
- 灰度流量：`v1Percent` 控制 v1/v2 的分流比例（先版本组，再组内 LB）。
- 插件：
  - API Key：`X-API-Key` 认证（错误/缺失返回 401）
  - IP 白名单：精确匹配来源 IP（非白名单返回 403）
  - 限流：固定窗口（超限返回 429）
  - CORS：跨域允许头（影响浏览器前端脚本；非浏览器不受限）
- 日志：文件 `logs/access.log`（时间戳、IP、方法、状态码、上游、耗时、路径）

## Dashboard 使用
- 刷新：从网关读取路由配置；连接后顶部显示“网关已连接”。
- 保存：提交表单到网关；写入成功提示并缓存到浏览器本地。
- 配置项：灰度比例、权重、API Key、IP 白名单（逗号分隔）、限流（秒内最多次数）、CORS（允许所有或自定义 Origin）。

## Admin API 与数据结构
详见 `docs/API 接口文档.md`。

## 许可
本项目遵循仓库内 `LICENSE` 所示的协议。
