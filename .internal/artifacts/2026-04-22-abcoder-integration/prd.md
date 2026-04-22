# PRD - ABCoder Integration for Code Review & Fix Suggestions

## 1. 背景与动机

### 1.1 当前状态
CodeSentry 是一个静态分析工具，支持 11 种编程语言的代码安全扫描。当前仅能检测问题，**无法提供基于代码上下文的修复建议**。

### 1.2 目标
集成 abcoder (cloudwego/abcoder) 框架，增强 CodeSentry 的代码审查能力：
- 利用 UniAST（Universal AST）统一代码表示
- 提供 Code-RAG 能力，实现精准的代码库理解
- 基于结构化代码上下文生成修复建议

### 1.3 触发原因
用户需要 CodeSentry 不仅能发现问题，还能：
1. 理解代码的语义上下文（而不仅是正则匹配）
2. 提供智能化的修复建议
3. 支持多语言的统一代码处理

## 2. 目标与成功标准

### 2.1 业务目标
- CodeSentry 能在检测到安全漏洞时，基于完整代码上下文提供修复方案
- 修复建议应具有语言无关性（通过 UniAST）
- 支持 Go、Python、TypeScript/JavaScript、Java、Rust 等主要语言

### 2.2 技术目标
- 集成 abcoder 作为代码解析和上下文理解引擎
- 利用 abcoder 的 Code-RAG MCP 工具增强代码库理解能力
- 利用 abcoder 的 General Writer 将修复建议以代码形式输出

### 2.3 成功指标
| 指标 | 目标 |
|------|------|
| 支持语言数 | ≥ 6 种主流语言 |
| 修复建议准确率 | 能生成语法正确的修复代码 |
| 上下文理解深度 | 能追踪变量定义、使用链路 |
| 集成方式 | MCP Server / Go Library |

## 3. 用户故事

### 3.1 核心用户故事
**作为** 安全工程师
**我希望** CodeSentry 在发现漏洞时能提供具体修复代码
**以便** 我可以直接应用修复，而无需手动研究修复方案

### 3.2 扩展用户故事
- **作为** 开发者 **我希望** 能看到修复前后的代码对比 **以便** 理解为什么这样做是安全的
- **作为** 安全团队 **我希望** 修复建议能考虑项目代码风格 **以便** 修复代码能自然融入现有代码库

## 4. 范围

### 4.1 In Scope
- abcoder Go SDK 集成到 CodeSentry
- UniAST 解析能力用于主要支持语言
- 基于 abcoder 的代码上下文理解
- 修复建议生成（基于检测到的漏洞）
- MCP Server 模式支持（可选）

### 4.2 Out of Scope
- abcoder CLI 的完整功能复制
- abcoder 的 Agent 模式（独立 AI Agent）
- abcoder 不支持的语言（C++, PHP, Swift, Kotlin 的高级上下文理解）
- 独立的 abcoder Web UI

### 4.3 非目标
- 不改变现有的 YAML 规则格式
- 不替代现有的正则匹配检测逻辑（作为补充而非替换）
- 不实现完整的代码生成能力（仅限修复建议）

## 5. 风险与依赖

### 5.1 关键依赖
| 依赖 | 状态 | 说明 |
|------|------|------|
| abcoder Go SDK | 已发布 | github.com/cloudwego/abcoder |
| abcoder UniAST 规范 | 稳定 | 需理解 UniAST schema |
| Go 1.21+ | 已满足 | 项目当前使用 |

### 5.2 技术风险
| 风险 | 影响 | 缓解措施 |
|------|------|----------|
| abcoder 解析性能 | 可能影响扫描速度 | 使用流式解析，按需加载 |
| UniAST 内存占用 | 大型代码库可能高 | 实现增量解析和缓存 |
| 修复建议质量 | 依赖 LLM 能力 | 提供结构化上下文，控制生成参数 |
| abcoder 版本兼容性 | 外部依赖可能 API 变更 | 锁定版本，定期同步 |

### 5.3 已知约束
- abcoder Writer 仅支持 Go（其他语言 Coming Soon）
- Code-RAG 功能需要 MCP Server 模式运行
- abcoder 解析需要项目根目录配置

## 6. 企业治理待确认项

### 6.1 应用等级
- 本次为**工具能力增强**，不属于核心生产系统
- 建议评级：T4（可简化架构复杂度，遵守安全底线）

### 6.2 数据/合规风险
- 本次不涉及用户数据处理
- 不涉及跨境数据传输
- 不涉及监管合规要求

### 6.3 集团组件约束
- 无需使用集团 PAAS/SAAS 能力
- 纯本地 SDK 集成，无外部服务依赖

## 7. 领域技能包启用建议

| 技能 | 用途 | 启用原因 |
|------|------|----------|
| `security-review` | 安全模式检测 | 增强现有安全规则与 abcoder 上下文结合 |
| `code-reviewer` | 代码审查 | 审查 abcoder 集成代码质量 |
| `golang-patterns` | Go 专项 | abcoder SDK 为 Go 实现，需遵循 Go 惯用写法 |

## 8. UI 范围与质量门禁

### 8.1 是否涉及 UI
**否** - 本次为纯 CLI 工具能力增强

### 8.2 终端假设
- 主要终端：开发者 CLI
- 输出形式：命令行修复建议 + JSON 结构化输出
- 无桌面/移动端需求

### 8.3 质量门禁
- 修复建议必须是语法正确的代码
- 不破坏现有功能
- 向后兼容（现有规则格式不变）

## 9. 需求挑战会候选分组

### 9.1 建议参与角色
| 角色 | 关注点 |
|------|--------|
| `backend-engineer` | abcoder SDK 集成实现 |
| `security-reviewer` | 修复建议安全性审查 |
| `qa-engineer` | 测试用例设计 |

### 9.2 待讨论议题
1. abcoder UniAST 与现有检测规则的结合方式
2. 修复建议生成的技术方案选型（MCP vs 直接 SDK）
3. 多语言支持优先级
4. 性能与内存占用的权衡

## 10. 待确认项

| # | 待确认项 | Owner | 优先级 |
|---|---------|-------|--------|
| 1 | abcoder 集成的最小可行范围（MVP） | tech-lead | P0 |
| 2 | 修复建议生成是否需要 LLM API（如需要，选型） | tech-lead | P0 |
| 3 | abcoder Writer 对其他语言的规划时间线 | cloudwego | P1 |
| 4 | 现有检测结果如何传递给 abcoder 进行上下文分析 | backend-engineer | P1 |
| 5 | Code-RAG 是否需要独立部署 MCP Server | backend-engineer | P2 |

## 11. 输入依据

- abcoder 官方仓库：https://github.com/cloudwego/abcoder
- abcoder 支持语言：Go ✅, Rust ✅, C ✅, Python ✅, JS/TS ✅, Java ✅
- abcoder License: Apache-2.0

## 12. 当前阶段

- **阶段**: intake
- **目标阶段**: requirement-challenge → design-swarm → design-review → handoff-ready

---

*创建时间: 2026-04-22*
*主责角色: tech-lead*
