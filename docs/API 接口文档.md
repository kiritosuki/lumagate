# LumaGate API 接口文档

本文档面向开发人员，描述网关的 Admin API、路由与插件的配置结构、转发行为与常见返回码。

## 服务地址
- 网关：`http://<host>:8080`
- Dashboard 开发：`http://<host>:5173`（仅本地开发）
- Dashboard 部署（Nginx）：`http://<host>/`（`/admin/*` 由 Nginx 代理到网关）

## Admin API

### 获取全部路由
- `GET /admin/routes`
- 响应：`Route[]`

### 新建或替换路由
- `POST /admin/routes`
- 请求体：`Route`
- 响应：`{"status":"ok"}`

### 更新路由（推荐）
- `POST /admin/routes/update`
- 请求体：`Route`
- 响应：`{"status":"ok"}`

### 删除路由
- `DELETE /admin/routes/:id`
- 响应：`{"status":"ok"}`

### 获取访问日志尾部
- `GET /admin/logs?tail=N`
- 响应：`string[]`（最后 N 行）

## 路由与配置结构（Route）
```json
{
  "id": "users",
  "prefix": "/api/users",
  "v1": [
    { "name": "A", "url": "http://serviceA_v1:9001", "weight": 1 },
    { "name": "B", "url": "http://serviceB_v1:9002", "weight": 1 },
    { "name": "C", "url": "http://serviceC_v1:9003", "weight": 1 }
  ],
  "v2": [
    { "name": "A", "url": "http://serviceA_v2:9011", "weight": 1 },
    { "name": "B", "url": "http://serviceB_v2:9012", "weight": 1 },
    { "name": "C", "url": "http://serviceC_v2:9013", "weight": 1 }
  ],
  "lbEnabled": false,
  "plugins": {
    "auth": { "enabled": false, "key": "" },
    "ipWhitelist": { "enabled": false, "ips": [] },
    "rateLimit": { "enabled": false, "windowSec": 1, "max": 5 },
    "cors": { "enabled": false, "allowAll": true, "origins": [] },
    "trafficSplit": { "enabled": true, "v1Percent": 33 }
  }
}
```

### 字段说明
- `id`：路由唯一标识，建议与前缀关联（例如 `users`）
- `prefix`：前缀路由（最长前缀优先）。命中后进入插件链与负载均衡。
- `v1` / `v2`：上游服务列表（URL 与权重）。
- `lbEnabled`：是否启用加权随机；关闭时为轮询。
- `plugins`：插件配置集合。

## 插件配置说明

### API Key（auth）
- 作用：校验请求头 `X-API-Key`。
- 字段：
  - `enabled`：是否启用
  - `key`：密钥字符串
- 返回码：不正确或缺失 → `401`。

### IP 白名单（ipWhitelist）
- 作用：仅允许列表中的来源 IP 访问。
- 字段：
  - `enabled`：是否启用
  - `ips`：字符串数组，逗号分隔在 Dashboard 中填写（支持 IPv4 与 IPv6）
- 返回码：不在白名单 → `403`。
- 说明：当前实现基于 `RemoteAddr` 主机部分做精确匹配；部署在反向代理后，需将代理的出站 IP 配入白名单。

### 限流（rateLimit）
- 作用：固定时间窗口内的最大请求数。
- 字段：
  - `enabled`：是否启用
  - `windowSec`：窗口秒数（整数）
  - `max`：该窗口允许的最大请求数（整数）
- 返回码：超限 → `429`。

### CORS（cors）
- 作用：控制浏览器前端的跨域脚本访问。
- 字段：
  - `enabled`：是否启用
  - `allowAll`：允许所有来源（`Access-Control-Allow-Origin: *`）
  - `origins`：自定义允许的来源列表（启用时按需扩展）
- 说明：仅影响浏览器中的跨域脚本；命令行或后端服务不受限制。
- OPTIONS 预检：返回 `204`，开启时带允许头。

### 灰度（trafficSplit）
- 作用：按比例在 v1/v2 版本组进行流量分配。
- 字段：
  - `enabled`：是否启用
  - `v1Percent`：v1 的百分比（0–100）。
- 行为：先版本组分流，再在组内按轮询或权重选具体上游。

## 负载均衡
- 轮询：`lbEnabled=false` 时，按组内顺序 A→B→C 轮询。
- 加权随机：`lbEnabled=true` 时，按上游 `weight` 做加权随机选择。
- 健康：演示服务提供 `GET /health`，网关当前不做主动健康探测；可在后续扩展。

## 转发与路径改写
- 命中 `prefix` 后，网关将请求路径剥离前缀再转发到上游：
  - 例如：`/api/users/user` → 上游接收 `/user`
- 若未命中任何前缀，返回 `404`。

## 常见返回码
- `200`：成功（由上游决定）
- `401`：API Key 认证失败
- `403`：IP 白名单拒绝
- `404`：未命中路由前缀或根路径（默认不提供欢迎页）
- `429`：限流超限
- `502`：上游错误或不可用

## 示例：更新为加权与 100% v1
```bash
curl -s -X POST -H 'Content-Type: application/json' -d '{
  "id":"users","prefix":"/api/users",
  "v1":[{"name":"A","url":"http://serviceA_v1:9001","weight":5},
         {"name":"B","url":"http://serviceB_v1:9002","weight":1},
         {"name":"C","url":"http://serviceC_v1:9003","weight":1}],
  "v2":[{"name":"A","url":"http://serviceA_v2:9011","weight":1},
         {"name":"B","url":"http://serviceB_v2:9012","weight":1},
         {"name":"C","url":"http://serviceC_v2:9013","weight":1}],
  "lbEnabled":true,
  "plugins":{"auth":{"enabled":false,"key":""},
             "ipWhitelist":{"enabled":false,"ips":[]},
             "rateLimit":{"enabled":false,"windowSec":1,"max":5},
             "cors":{"enabled":false,"allowAll":true,"origins":[]},
             "trafficSplit":{"enabled":false,"v1Percent":100}}
}' http://<host>:8080/admin/routes/update
```

## 示例：限流（1 秒最多 2 次）与 33% 灰度
```bash
curl -s -X POST -H 'Content-Type: application/json' -d '{
  "id":"users","prefix":"/api/users",
  "v1":[{"name":"A","url":"http://serviceA_v1:9001","weight":1},
         {"name":"B","url":"http://serviceB_v1:9002","weight":1},
         {"name":"C","url":"http://serviceC_v1:9003","weight":1}],
  "v2":[{"name":"A","url":"http://serviceA_v2:9011","weight":1},
         {"name":"B","url":"http://serviceB_v2:9012","weight":1},
         {"name":"C","url":"http://serviceC_v2:9013","weight":1}],
  "lbEnabled":false,
  "plugins":{"auth":{"enabled":false,"key":""},
             "ipWhitelist":{"enabled":false,"ips":[]},
             "rateLimit":{"enabled":true,"windowSec":1,"max":2},
             "cors":{"enabled":true,"allowAll":true,"origins":[]},
             "trafficSplit":{"enabled":true,"v1Percent":33}}
}' http://<host>:8080/admin/routes/update
# 三次快速请求：200, 200, 429
curl -s -o /dev/null -w '%{http_code}\n' http://<host>:8080/api/users/user
curl -s -o /dev/null -w '%{http_code}\n' http://<host>:8080/api/users/user
curl -s -o /dev/null -w '%{http_code}\n' http://<host>:8080/api/users/user
```

## 示例：启用 API Key 与 IP 白名单
```bash
curl -s -X POST -H 'Content-Type: application/json' -d '{
  "id":"users","prefix":"/api/users",
  "v1":[{"name":"A","url":"http://serviceA_v1:9001","weight":1},
         {"name":"B","url":"http://serviceB_v1:9002","weight":1},
         {"name":"C","url":"http://serviceC_v1:9003","weight":1}],
  "v2":[{"name":"A","url":"http://serviceA_v2:9011","weight":1},
         {"name":"B","url":"http://serviceB_v2:9012","weight":1},
         {"name":"C","url":"http://serviceC_v2:9013","weight":1}],
  "lbEnabled":true,
  "plugins":{"auth":{"enabled":true,"key":"abc123"},
             "ipWhitelist":{"enabled":true,"ips":["127.0.0.1","::1"]},
             "rateLimit":{"enabled":false,"windowSec":1,"max":5},
             "cors":{"enabled":false,"allowAll":true,"origins":[]},
             "trafficSplit":{"enabled":false,"v1Percent":100}}
}' http://<host>:8080/admin/routes/update
# 认证：不带头 401；带头 200
curl -s -o /dev/null -w '%{http_code}\n' http://<host>:8080/api/users/user
curl -s -H 'X-API-Key: abc123' -o /dev/null -w '%{http_code}\n' http://<host>:8080/api/users/user
```

## 日志
- 文件：`logs/access.log`
- Admin：`GET /admin/logs?tail=N`

## 备注
- 插件执行顺序：`auth → ipWhitelist → rateLimit → cors → trafficSplit → proxy`
- 如需 CIDR 白名单 / 主动健康检查 / 指标采集等，可在现有模块上扩展实现。
