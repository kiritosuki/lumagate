目录结构

- lumagate/
  - gateway/                      # Go API 网关（后端）
    - main.go
    - router.go                   # 路由前缀匹配 + 反向代理
    - admin.go                    # Admin API：路由与插件配置、日志读取
    - upstream/
      - load_balancer.go          # 轮询 + 加权随机（三服务权重）
    - plugins/
      - auth.go                   # X-API-Key 认证
      - ip_whitelist.go           # IP 白名单
      - rate_limit.go             # 固定窗口限流
      - cors.go                   # CORS 允许/拒绝
      - traffic.go                # v1/v2 灰度比例选择
      - logging.go                # 访问日志写文件（logs/access.log）
    - config/
      - types.go                  # Route/Upstream/PluginConfig 结构体
  - services/                     # 下游小服务（Go）
    - v1/
      - serviceA/main.go          # 9001: GET /user → "v1 A"; GET /health
      - serviceB/main.go          # 9002: GET /user → "v1 B"; GET /health
      - serviceC/main.go          # 9003: GET /user → "v1 C"; GET /health
    - v2/
      - serviceA/main.go          # 9011: GET /user → "v2 A"; GET /health
      - serviceB/main.go          # 9012: GET /user → "v2 B"; GET /health
      - serviceC/main.go          # 9013: GET /user → "v2 C"; GET /health
  - dashboard/                    # 前端（Vue 3 + Vite + Element Plus）
    - src/
      - pages/
        - Routes.vue              # 路由列表与编辑（前缀、灰度比例、权重、插件参数）
        - Logs.vue                # 日志查看（轮询尾部）
      - api/
        - admin.ts                # Axios 调 /admin 接口
      - store/
        - routes.ts               # Pinia 管理路由配置状态
      - App.vue / main.ts
    - vite.config.ts
  - logs/
    - access.log                  # 网关访问日志文件

后端实现要点（Go）

- 路由匹配：前缀匹配（prefix）；命中后进入插件链：auth → ipWhitelist → rateLimit → cors → trafficSplit → proxy
- 负载均衡：
  - 未开启权重控制 → 轮询（round robin）
  - 开启权重控制 → 组内按权重加权随机（Weighted Random）
  - 灰度：先按 v1/v2 比例选择版本组，再在组内选择具体服务
- 日志：每次请求写入 logs/access.log（时间戳、方法、路径、状态码、上游、耗时、客户端 IP）
- Admin API：
  - POST `/admin/routes`           新增/替换路由（含 groups、lb、plugins）
  - GET  `/admin/routes`           查询所有路由
  - POST `/admin/routes/update`    更新路由（权重、灰度、插件开关与参数）
  - DELETE `/admin/routes/:id`     删除路由
  - GET  `/admin/logs?tail=N`      返回日志文件尾部 N 行
- 配置存储：内存（启动时给出默认示例），Admin API 调整后热更新；无需数据库

数据结构（简要）

- Route { id, prefix, groups: { v1: []Upstream, v2: []Upstream }, lb: { enabled bool, weights map[string]int }, plugins: { auth{enabled,key}, ipWhitelist{enabled,ips[]}, rateLimit{enabled,windowSec,max}, cors{enabled,allowAll,origins[]}, trafficSplit{enabled,v1Percent} } }
- Upstream { name, url, weight }

前端 Dashboard（Vue）

- 路由列表与编辑：
  - 前缀（prefix）
  - 灰度比例（v1Percent slider）
  - 三服务权重 sliders（当 lb.enabled=true）
  - 插件开关与参数：auth key、IP 列表、限流窗口与最大请求数、CORS 允许/拒绝
  - 保存后调用 `/admin/routes/update`
- 日志页面：轮询 `/admin/logs?tail=200` 显示最新访问记录

下游服务（Go）

- v1/v2 各三个服务，端口预设为 9001/9002/9003 与 9011/9012/9013
- 每个服务：
  - GET `/user` 返回对应文本（用于观察转发命中）
  - GET `/health` 返回 200（健康检查可选）

执行步骤（两小时内）

1. 初始化 gateway：main + router + admin + plugins 空体 + upstream LB
2. 实现插件最小逻辑：auth、ipWhitelist、rateLimit（固定窗口）、cors、trafficSplit、logging（写文件）
3. 完成 Admin API：路由增改查删、权重与灰度更新、日志读取
4. 实现 6 个下游服务（v1×3、v2×3）
5. 初始化 Dashboard：Routes.vue 与 Logs.vue + 简易表单与 sliders
6. 联调演示：
   - 新建路由 `/api/users` → groups(v1:9001/9002/9003, v2:9011/9012/9013)
   - 调整 v1/v2 灰度比例与三服务权重，观察命中
   - 依次开启限流（429）、API Key（401）、IP 白名单（403）、CORS

如果确认执行，我将按此目录与步骤开始编码，并在每阶段完成后进行基本验证（本地请求与日志校验）。