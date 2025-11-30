# LumaGate 演示指南（Hackathon 路演手册）

## 概览
- LumaGate：轻量级可插拔 API 网关（Go）+ 可视化 Dashboard（Vue）。
- 支持：基本路由、负载均衡（轮询/加权随机）、灰度流量、限流、CORS、API Key、IP 白名单、访问日志。
- 本指南用于现场路演：列出后端小服务接口、Dashboard 使用说明、以及每项功能的测试用例。

## 服务与端口
- 网关：`http://localhost:8080`
- 前端（开发模式）：`http://localhost:5173`（代理 `/admin/*` 到 `8080`）
- 下游服务（用于转发与演示）：
  - v1：`9001`（A）、`9002`（B）、`9003`（C）
  - v2：`9011`（A）、`9012`（B）、`9013`（C）

## 下游小服务接口清单
- 通用接口（每个服务都一致）：
  - `GET /user`：返回服务标识文本
    - v1：A→`"v1 A"`，B→`"v1 B"`，C→`"v1 C"`
    - v2：A→`"v2 A"`，B→`"v2 B"`，C→`"v2 C"`
  - `GET /health`：健康检查，返回 `200`
- 示例：
  - `curl -s http://localhost:9001/user` → `v1 A`
  - `curl -s -o /dev/null -w '%{http_code}\n' http://localhost:9013/health` → `200`

## 网关转发规则（默认）
- 路由前缀：`/api/users`
- 访问示例：`curl -s http://localhost:8080/api/users/user`
- 转发逻辑：先执行插件链，后依据灰度选择版本组（v1/v2），最终在组内按轮询或权重选择 A/B/C。

## 网关 Admin API（用于命令行备选）
- `GET  /admin/routes`：查询全部路由
- `POST /admin/routes/update`：更新路由（前缀、v1/v2、权重、灰度、插件配置）
- `GET  /admin/logs?tail=N`：返回访问日志尾部 N 行
- 示例：
```
curl -s http://localhost:8080/admin/routes
curl -s -X POST -H 'Content-Type: application/json' -d '{...}' http://localhost:8080/admin/routes/update
curl -s 'http://localhost:8080/admin/logs?tail=50'
```

## Dashboard 使用说明（http://localhost:5173）
- 顶部按钮：
  - `刷新`：从网关拉取最新路由配置；连接正常显示“网关已连接”。
  - `保存`：提交当前表单到网关；保存成功会提示，并写入本地缓存（便于断网时仍保留配置）。
- 配置区域：
  - 路由前缀：一般为 `/api/users`
  - 灰度 v1 百分比：按比例在 v1/v2 之间分流（如 20 表示 v1:20% / v2:80%）
  - 负载均衡权重启用：
    - 关闭：轮询（A→B→C）
    - 开启：加权随机（填写 v1/v2 的 A/B/C 权重）
  - API Key：启用后，请求需带 `X-API-Key: <key>`
  - IP 白名单：启用后，仅填写的 IP 可访问（逗号分隔，如 `127.0.0.1, ::1`）
  - 限流：设置“窗口秒数”和“该窗口最多请求数”（例如 `1 秒内最多 2 次`）
  - CORS 设置：启用后选择 `允许所有来源（*）` 或 `自定义来源`（逗号分隔 Origin）
- 日志页面：轮询尾部日志快速查看近期请求（`logs/access.log`）。

## 演示脚本与测试用例

### 1) 基本路由转发
- 目标：网关转发到 v1 组并轮询 A/B/C。
- 步骤：
```
# 连续访问 3 次，观察返回 v1 A/B/C
curl -s http://localhost:8080/api/users/user; echo
curl -s http://localhost:8080/api/users/user; echo
curl -s http://localhost:8080/api/users/user; echo
# 查看日志
tail -n 10 logs/access.log
```

### 2) 负载均衡（加权随机）
- 目标：在 v1 组内设置权重偏向 A。
- Dashboard：开启“负载均衡权重启用”，设置 v1 权重 A:5、B:1、C:1。
- 命令行备选：
```
# 启用权重并偏向 v1 A
curl -s -X POST -H 'Content-Type: application/json' -d '{
  "id":"users","prefix":"/api/users",
  "v1":[{"name":"A","url":"http://localhost:9001","weight":5},
         {"name":"B","url":"http://localhost:9002","weight":1},
         {"name":"C","url":"http://localhost:9003","weight":1}],
  "v2":[{"name":"A","url":"http://localhost:9011","weight":1},
         {"name":"B","url":"http://localhost:9012","weight":1},
         {"name":"C","url":"http://localhost:9013","weight":1}],
  "lbEnabled":true,
  "plugins":{"auth":{"enabled":false,"key":""},
             "ipWhitelist":{"enabled":false,"ips":[]},
             "rateLimit":{"enabled":false,"windowSec":1,"max":5},
             "cors":{"enabled":false,"allowAll":true,"origins":[]},
             "trafficSplit":{"enabled":false,"v1Percent":100}}
}' http://localhost:8080/admin/routes/update
# 多次访问，命中更偏向 v1 A（9001）
for i in 1 2 3 4 5; do curl -s http://localhost:8080/api/users/user; done; echo
```

### 3) 灰度流量（v1/v2 分流）
- 目标：设置 v1 20%，v2 80%。
- Dashboard：灰度 v1 百分比改为 `20`，保存。
- 命令行备选：
```
curl -s -X POST -H 'Content-Type: application/json' -d '{
  "id":"users","prefix":"/api/users",
  "v1":[{"name":"A","url":"http://localhost:9001","weight":1},
         {"name":"B","url":"http://localhost:9002","weight":1},
         {"name":"C","url":"http://localhost:9003","weight":1}],
  "v2":[{"name":"A","url":"http://localhost:9011","weight":1},
         {"name":"B","url":"http://localhost:9012","weight":1},
         {"name":"C","url":"http://localhost:9013","weight":1}],
  "lbEnabled":true,
  "plugins":{"auth":{"enabled":false,"key":""},
             "ipWhitelist":{"enabled":false,"ips":[]},
             "rateLimit":{"enabled":false,"windowSec":1,"max":5},
             "cors":{"enabled":true,"allowAll":true,"origins":[]},
             "trafficSplit":{"enabled":true,"v1Percent":20}}
}' http://localhost:8080/admin/routes/update
# 多次访问，结果大多数为 v2 *，少数为 v1 *
for i in 1 2 3 4 5; do curl -s http://localhost:8080/api/users/user; done; echo
```

### 4) 限流（固定窗口）
- 目标：设置 1 秒最多 2 次。
- Dashboard：启用限流，窗口秒数 `1`，最多请求 `2`。
- 命令行备选：
```
curl -s -X POST -H 'Content-Type: application/json' -d '{
  "id":"users","prefix":"/api/users",
  "v1":[{"name":"A","url":"http://localhost:9001","weight":1},
         {"name":"B","url":"http://localhost:9002","weight":1},
         {"name":"C","url":"http://localhost:9003","weight":1}],
  "v2":[{"name":"A","url":"http://localhost:9011","weight":1},
         {"name":"B","url":"http://localhost:9012","weight":1},
         {"name":"C","url":"http://localhost:9013","weight":1}],
  "lbEnabled":true,
  "plugins":{"auth":{"enabled":false,"key":""},
             "ipWhitelist":{"enabled":false,"ips":[]},
             "rateLimit":{"enabled":true,"windowSec":1,"max":2},
             "cors":{"enabled":false,"allowAll":true,"origins":[]},
             "trafficSplit":{"enabled":false,"v1Percent":100}}
}' http://localhost:8080/admin/routes/update
# 三次快速请求：200, 200, 429
curl -s -o /dev/null -w '%{http_code}\n' http://localhost:8080/api/users/user
curl -s -o /dev/null -w '%{http_code}\n' http://localhost:8080/api/users/user
curl -s -o /dev/null -w '%{http_code}\n' http://localhost:8080/api/users/user
```

### 5) API Key 认证
- 目标：启用后仅携带正确 `X-API-Key` 的请求通过。
- Dashboard：启用 API Key，填写例如 `abc123`。
- 命令行备选：
```
curl -s -X POST -H 'Content-Type: application/json' -d '{
  "id":"users","prefix":"/api/users",
  "v1":[{"name":"A","url":"http://localhost:9001","weight":1},
         {"name":"B","url":"http://localhost:9002","weight":1},
         {"name":"C","url":"http://localhost:9003","weight":1}],
  "v2":[{"name":"A","url":"http://localhost:9011","weight":1},
         {"name":"B","url":"http://localhost:9012","weight":1},
         {"name":"C","url":"http://localhost:9013","weight":1}],
  "lbEnabled":true,
  "plugins":{"auth":{"enabled":true,"key":"abc123"},
             "ipWhitelist":{"enabled":false,"ips":[]},
             "rateLimit":{"enabled":false,"windowSec":1,"max":5},
             "cors":{"enabled":false,"allowAll":true,"origins":[]},
             "trafficSplit":{"enabled":false,"v1Percent":100}}
}' http://localhost:8080/admin/routes/update
# 不带头：401；带头：200
curl -s -o /dev/null -w '%{http_code}\n' http://localhost:8080/api/users/user
curl -s -H 'X-API-Key: abc123' -o /dev/null -w '%{http_code}\n' http://localhost:8080/api/users/user
```

### 6) IP 白名单
- 目标：只允许本机访问（示例：`127.0.0.1, ::1`）。
- Dashboard：启用 IP 白名单并填写 IP。
- 命令行备选：
```
curl -s -X POST -H 'Content-Type: application/json' -d '{
  "id":"users","prefix":"/api/users",
  "v1":[{"name":"A","url":"http://localhost:9001","weight":1},
         {"name":"B","url":"http://localhost:9002","weight":1},
         {"name":"C","url":"http://localhost:9003","weight":1}],
  "v2":[{"name":"A","url":"http://localhost:9011","weight":1},
         {"name":"B","url":"http://localhost:9012","weight":1},
         {"name":"C","url":"http://localhost:9013","weight":1}],
  "lbEnabled":true,
  "plugins":{"auth":{"enabled":false,"key":""},
             "ipWhitelist":{"enabled":true,"ips":["127.0.0.1","::1"]},
             "rateLimit":{"enabled":false,"windowSec":1,"max":5},
             "cors":{"enabled":false,"allowAll":true,"origins":[]},
             "trafficSplit":{"enabled":false,"v1Percent":100}}
}' http://localhost:8080/admin/routes/update
```

### 7) CORS（跨域）
- 目标：允许所有来源或自定义来源（用于前端跨域演示）。
- Dashboard：启用 CORS，选择“允许所有来源（*）”，或填写自定义 Origin（逗号分隔）。

### 8) 日志查看
- 目标：展示访问信息（时间戳、IP、方法、状态码、上游、耗时、路径）。
- 命令：
```
tail -n 20 logs/access.log
curl -s 'http://localhost:8080/admin/logs?tail=50'
```

## 常见问题与诊断
- 前端提示“未连接”，但接口可用：点击“刷新”后查看是否变为“网关已连接”；或直接保存并用 `GET /admin/routes` 验证。
- 端口占用导致服务退出：重启前先 `pkill -f <进程名>`，再 `nohup` 启动。
- 代理异常：Dashboard 已内置“代理失败→直连网关”回退逻辑，保存与刷新仍可用。

## 启停与重启建议（现场操作顺序）
1. 构建：`go build -o ./bin/lumagate-gateway ./gateway && go build -o ./bin/serviceA_v1 ./services/v1/serviceA && ...`
2. 启动 6 个后端：各用 `nohup ./bin/serviceX_vY >/tmp/serviceX_vY.log 2>&1 &`
3. 启动网关：`nohup ./bin/lumagate-gateway >/tmp/gateway.log 2>&1 &`
4. 验证健康：逐个 `curl -s -o /dev/null -w '%{http_code}\n' http://localhost:<port>/health`
5. 启动前端：`cd dashboard && npm install && npm run dev`，打开 `http://localhost:5173`

> 以上命令与步骤仅用于现场演示；Dashboard 可视化与命令行互为备选，任一方式都能更新配置并验证。

