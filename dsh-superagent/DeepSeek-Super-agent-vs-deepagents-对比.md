# DeepSeek Harness Super-agent vs deepagents 对比

> 一份关于两个 Agent Harness 框架/架构的横向对比文档，从 deepagents 的功能维度出发，逐项映射到 Super-agent 的能力，并给出选型与组合建议。

---

## 0. 命名澄清（务必先读）

这两个名字容易被误认为同类，**实际上它们不是一个层级的东西**：

| | DeepSeek Harness Super-agent | deepagents |
|---|---|---|
| **本质** | 一套**架构设计探索 / 研究原型**（research prototype） | 一个**可安装使用的生产级 Python 库** |
| **出品方** | DeepSeek Harness 团队 | LangChain 团队（社区项目） |
| **官方定性** | "仅供参考，非官方产品，非生产就绪" | "production-ready AI agents 框架" |
| **代码仓库** | `deepseek-ai/dsh-superagent-prototype` | `deepagents`（PyPI / uv 安装） |

**一句话**：deepagents 是"拿来就能用"的 Agent 库；Super-agent 是"值得学习借鉴"的架构蓝图。

---

## 1. 基础信息对比

| 维度 | DeepSeek Harness Super-agent | deepagents |
|---|---|---|
| 定位 | "AI 调度 AI" 的持续在线自主智能体架构 | 开箱即用的通用 Agent 开发框架 |
| 技术栈/依赖 | DeepSeek Harness 插件化 bundle（GitHub + 飞书监听） | LangChain + LangGraph 之上的封装层 |
| 核心架构 | Router + Worker 双层，事件驱动、存算分离、Watchdog、投机执行 | 单 Agent 起步，可动态 spawn 子 Agent（`create_deep_agent`） |
| 语言 | TypeScript | Python |
| 许可证 | MIT | MIT |
| 成熟度 | ⚠️ 低：官方标注"可能有 bug、风格不一致、与内部版不同" | ✅ 高：约 30k+ stars、活跃维护 |
| 安装/上手 | 步骤多，需按 SETUP.md 由 agent 辅助配置 | `uv add deepagents` 秒装 |
| 模型支持 | 面向 DeepSeek Harness 生态 | 模型无关（OpenAI / Claude / Gemini / Llama / Ollama 等） |
| MCP 支持 | 未涉及 | ✅ 完整支持（`langchain-mcp-adapters`） |

---

## 2. 核心功能维度逐项对比

> 方法论：把 deepagents 的每一项库级能力，去问"Super-agent 有没有对应物、怎么实现、强在哪弱在哪"。

| deepagents 功能维度 | deepagents 实现 | Super-agent 对应实现 | 差异与注意事项 |
|---|---|---|---|
| **规划工具 (Todos)** | `TodoListMiddleware` + `write_todos`，状态 pending/in_progress/completed，可选启用 | 无独立 todo 工具；靠 Router 的 `decide` 决策拆成 Ticket/Mandate/Assignment | 前者是"任务清单"粒度；后者是"工作项生命周期"粒度（Ticket 是持久实体） |
| **文件系统** | 可插拔后端（内存/本地/Store），`ls/read/write/edit/glob/grep/delete`，支持多模态 | 每个 Ticket 独立 Session + 独立目录；Worker 无"编辑外部文档"权限 | deepagents 文件工具更丰富；Super-agent 靠"独立目录"天然隔离 |
| **文件权限** | 声明式 `permissions` 规则（allow/deny，first-match-wins） | 无文件级 allowlist；靠角色级 `Pick` 接口 + Worker 工具白名单收窄能力 | 安全模型不同：deepagents 管"文件"，Super-agent 管"角色/工具" |
| **子代理 (Subagents)** | `task` 工具，同步/异步、isolated/fork 模式、临时子代理 | Worker（DshExecutor），每个 Ticket 一个独立 Session，可并发多 Worker | Super-agent 的 Worker 是常驻、可订阅消息、可恢复中断的长生命周期实体；deepagents 子代理是一次性临时的 |
| **持久化记忆** | `AGENTS.md` 文件，始终加载，可更新 | Policy（规则）+ 存算分离的 SQLite Store 持久化会话状态 | deepagents 记"约定/偏好"，Super-agent 记"策略/授权/任务状态" |
| **HITL** | `interrupt_on`，工具调用前审批/修改/拒绝 | 授权机制（Gate + Mandate）+ Panel 展示待处理 Ticket；人类用 `stopTicket` 撤回授权 | deepagents 是工具级拦截；Super-agent 是授权级 + 结果级验收 |
| **Skills 按需加载** | 渐进式披露，启动读 frontmatter，按需加载完整 skill | 无独立 Skills 机制；Router/Worker 靠 prompt 预设（`ROUTER_PRESET`） | deepagents 有成熟 Skills 生态；Super-agent 无此抽象 |
| **上下文管理/摘要** | 四层：Skills / Memory / 摘要卸载 / Prompt 缓存 | Inline & Spill：内联截断（1500/600 字符）+ `orchestrator_read` 按需读取 | 思路殊途同归；Super-agent 用分级读取指针，token 控制更精细 |
| **代码执行** | sandbox 后端 `execute` + `eval`（QuickJS 受限） | Worker 在独立目录的 Harness 里执行（复现 bug、跑测试、开 PR） | deepagents 是通用代码执行；Super-agent 绑定在"代码仓库"场景 |
| **模型支持** | 模型无关（Gemini/OpenAI/Claude/Ollama 等） | 面向 DeepSeek Harness 生态 | deepagents 选型更自由 |
| **MCP 支持** | ✅ 完整（`tools=` 传 MCP server 工具） | ❌ 未涉及 | 需要接 MCP 生态 → deepagents 完胜 |
| **流式输出** | 类型化事件流，`stream.subagents` 独立句柄 | Panel + Relay 机制 + DomainEvents（12 类事件） | deepagents 面向编程消费；Super-agent 面向面板可视化 |
| **可观测性/监控** | LangSmith 追踪 / 评估 | 自建 HealthBoard + Watchdog + Panel | deepagents 靠第三方；Super-agent 自建健康监控（Watchdog 检查无效循环） |

---

## 3. Super-agent 独有、deepagents 没有的能力

反过来问，能看清 Super-agent 的真正价值：

| Super-agent 独有能力 | 说明 | deepagents 是否有 |
|---|---|---|
| **事件驱动 + 两级过滤** | Gate（程序规则快筛）+ MasterAssessor（大模型慢判），无人 @ 也能主动识别 bug 并后台修复 | ❌ 无（被动响应式，等用户 invoke） |
| **多会话并发调度 (Router)** | 常驻全局协调节点，统一接收多渠道消息并分派 | ⚠️ 有子代理，但无"常驻 Router 持续监听" |
| **存算分离 + 持久在线** | SQLite 存状态，执行循环在云端，本地轻 daemon，离线也能继续 | ❌ 无（请求-响应式，进程结束即止） |
| **Watchdog 自愈** | 检查任务是否陷入无效循环并尝试修复 | ❌ 无 |
| **投机执行** | 结果不确定时后台先试，有价值再提交 | ❌ 无 |
| **授权治理 (Policy/Mandate)** | 谁允许启动工作、规则变更如何撤销授权、角色级接口隔离 | ⚠️ 有文件权限，但无"授权委托/撤回"治理层 |

---

## 3.5 Super-agent 功能 → 代码位置对照

> 补充：把上文提到的 Super-agent 各项能力，映射到原型仓库 `deepseek-ai/dsh-superagent-prototype` 中的具体源码文件，方便直接对照源码学习。源码为 TypeScript，`src/` 下分 `core`（纯领域逻辑，不 I/O）、`app`（应用/编排层）、`executor`（执行器）、`infra`（基础设施）、`sources`（多渠道接入）、`client`（Web 面板）。

### 3.5.1 目录职责总览

| 目录 | 职责 | 关键点 |
|---|---|---|
| `src/core/` | 纯领域逻辑，**不 import 自身以外任何东西、不做 I/O**，Engine 方法同步、每命令一个事务 | 领域模型与规则 |
| `src/app/` | 应用/编排层：Store、Effects、Passes、Assessment、Runtime、会话等 | 持久化 + 调度 + 评估 |
| `src/executor/` | 执行器（Worker 的具体实现） | 真正干活的 Agent 会话 |
| `src/infra/` | 基础设施：健康监控、命令、环境、路径、文本 | Watchdog 相关 |
| `src/sources/` | 多渠道接入：GitHub、飞书 | 事件入口 |
| `src/client/` | Web 管理面板（React/tsx） | 人类验收界面 |
| 顶层 `src/*.ts` | `policy.ts`（策略）、`admission.ts`（准入）、`limits.ts`（限额）、`kv.ts`、`navigation.ts`、`setup.ts`、`deployment.ts` 等 | 横切关注点 |

### 3.5.2 功能 → 代码位置映射

| 功能 | 对应代码文件 | 代码层面体现 |
|---|---|---|
| **Router（主会话/全局协调）** | `src/core/engine.ts`（`EngineForMainSession` 角色）、`src/app/sessions.ts` | 路由决策 `orchestrator_route` 一次接受 1–20 个决策；命令用祈使语气（`stopTicket`），适配器事实用过去时（`assignmentStarted`） |
| **Worker（执行器）** | `src/executor/dsh.ts`（`DshExecutor`）、`src/executor/background.ts` | 每个 Ticket 独立 Session + 目录；`workerSessionId(ticket.id)` 派生；6 个 orchestrator 工具（read/ack/report/thread_unrelated/thread_follow/trigger_handback） |
| **事件驱动（DomainEvents/Effects）** | `src/app/effects.ts`（持久化副作用队列）、`src/app/passes.ts`（调度循环）、`src/core/ports.ts`（`DomainEvents` 定义） | 12 类事件（trigger.received / eval.started / decision.recorded / ticket.created / ticket.resumed / ticket.ended / assignment.started / assignment.ended / thread.subscribed / thread.unsubscribed / trigger.requeued / policy.changed）；Effects 在事务内 enqueue、提交后执行；同一时间仅一个 pass，每 2 秒运行一次（有命令请求时也会立即运行），处理到期 effects → `Engine.evaluateDue` |
| **两级过滤（Gate + Assessor）** | `src/policy.ts`（`gate`）、`src/admission.ts`（准入表）、`src/app/assessment.ts`（`MasterAssessor`） | 第一级 Gate 返回准入强度 + 依据（addressed/followup/admission/bot 四种 `via`）；第二级 MasterAssessor 大模型慢判 |
| **订阅直通 / Fan-out** | `src/core/intake.ts`、`src/core/tickets.ts`（`follows` 字段） | Trigger 一次存储、扇出到订阅的 Tickets；`followsDone` 规则决定 Ticket 何时结束 |
| **Inline & Spill 上下文** | `src/app/resources.ts`（`orchestrator_read` 资源）、`src/limits.ts` | 内联截断（assignment 1500 / trigger card 600 字符）+ 按需读取指针；文本页 4000–8000 字符 |
| **存算分离** | `src/app/store.ts`（唯一 SQLite 写入者）、`src/app/store-read.ts`、`src/app/migration.ts` | 单一 `orchestrator.sqlite` 文件，仅 Store 打开；单一写入者原则（每种状态只有一个写入者）；`core` 不做 I/O |
| **Watchdog / 自愈** | `src/infra/health.ts`（`HealthBoard`）、`src/app/retry-health.ts`（`RetryHealth`） | 每个部件认领（claim）健康条目，重复认领被拒绝；effect 连续失败三次报错（`delivery` 例外：其失败显示在对应 Ticket 上而非 health 条目）；`RetryHealth` 跟踪等待模型预算重试的 Session |
| **授权治理（Policy/Mandate）** | `src/policy.ts`、`src/core/evals.ts`、`src/app/role-models.ts` | `revalidate(mandate)` 检查授权是否仍有效；`policyChanged` 仅通知受影响 Ticket；角色级 `Pick` 接口（EngineForSources/ForMainSession/ForPeople/ForExecutor 等） |
| **Web 面板（Panel）** | `src/client/index.tsx`、`src/client/app.tsx`、`src/client/api.ts`、各阶段页面 `overview.tsx`/`ticket.tsx`/`settings.tsx` 等 | 展示需人类介入的 Ticket、当日流量、Executor 负载；面板专用操作（驳回 unavailable、重跑 setup 等） |
| **多渠道接入（Source）** | `src/sources/github*.ts`（github-channel/pull/read/trigger）、`src/sources/feishu*.ts`、`src/sources/registry.ts`、`src/sources/poll.ts` | GitHub 与飞书两个 Source；`registry.ts` 注册、`poll.ts` 轮询、`host.ts` 宿主 |
| **领域模型（Ticket/Assignment/Mandate）** | `src/core/tickets.ts`、`src/core/assignments.ts`、`src/core/evals.ts`、`src/core/intake.ts`、`src/core/replies.ts` | 核心实体与状态机 |
| **Prompt 预设** | `prompts/` 目录（顶层）、`src/app/prompts.ts` | Router 预设 `ROUTER_PRESET` 等 |

### 3.5.3 读代码的建议路径

1. **先读文档**：`docs/000-philosophy.md` → `001-domain.md` → `002-authorization.md` → `007-architecture.md`，每篇只建立在前面基础上。
2. **再看领域核心**：`src/core/engine.ts`（调度决策的"大脑"）→ `src/core/tickets.ts` / `assignments.ts`（工作项生命周期）。
3. **然后看编排层**：`src/app/effects.ts` + `passes.ts`（事件驱动的"心跳"）→ `src/app/assessment.ts`（主动评估）→ `src/app/store.ts`（存算分离）。
4. **最后看外围**：`src/executor/dsh.ts`（Worker）、`src/sources/`（多渠道）、`src/infra/health.ts`（自愈）。

> 注意：官方明确说明该仓库是"发布时适配过"的版本，**与内部版本及技术文章描述的并非完全一致**，部分文章中的设计（如投机执行、异步 ReAct Loop）在原型里未实现或术语不同，读码时需对照 `docs/` 与文章判断。

---

## 4. 总结

### 4.1 核心判断

1. **定位不同导致功能维度天然错位**：deepagents 是"给一个 Agent 配齐干活工具"的**库**，Super-agent 是"让多个 Agent 持续在线协作"的**系统**。
   - deepagents 的强项（文件工具、Skills、MCP、模型无关）在 Super-agent 里要么没有、要么靠架构约束替代；
   - Super-agent 的强项（事件驱动主动性、并发调度、存算分离、自愈、授权治理）deepagents 完全没有。

2. **能直接对上的维度，实现哲学也不同**：
   - 子代理 → deepagents"临时、一次性" vs Super-agent 的 Worker"常驻、可订阅、可恢复"；
   - 上下文管理 → 前者"自动摘要" vs 后者"Inline & Spill 分级读取"；
   - HITL → 前者"工具级拦截" vs 后者"授权级 + 结果级验收"。

### 4.2 选型建议

| 你的需求 | 选择 |
|---|---|
| 快速构建单个能干活的 Agent（文件、记忆、子任务、MCP、可观测） | **deepagents**（生态成熟、MIT、模型无关、秒级上手） |
| 构建持续在线、多任务并发的团队级"管家型"系统（主动响应、并发调度、断点续跑、权限治理） | **借鉴 Super-agent 的架构思想**，但别直接用原型代码（官方明确非生产就绪） |

### 4.3 最务实的组合

用 **deepagents 作为每个 Worker 的执行内核**（拿到成熟的文件/记忆/子代理/HITL 能力），外面自建一层 **Router + 事件驱动 + Store** 来复刻 Super-agent 的调度与主动性 —— 各取所长。

### 4.4 最需警惕的一点

别因为名字里都有 "deep/agent" 就把 Super-agent 当 deepagents 的同类替换品 —— **一个是架构探索论文级产物，一个是工程库**，选型依据完全不同。

---

## 附：参考资料

- DeepSeek Super-agent 技术文章：《Super-agent 设计探索：迈向更自主的 Agent 系统》
- Super-agent 原型仓库：`https://github.com/deepseek-ai/dsh-superagent-prototype`
- deepagents 官网：`https://deepagents.org/`
- deepagents 文档：`https://docs.langchain.com/oss/python/deepagents/overview`
