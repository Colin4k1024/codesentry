# Brainstorm - ABCoder Integration

**日期:** 2026-04-22
**参与者:** tech-lead
**状态:** 已收敛

---

## 问题框定

**核心问题：** CodeSentry 当前只能检测漏洞，无法提供基于代码上下文的修复建议。

**关键假设质疑：**

| # | 假设 | 质疑点 |
|---|------|--------|
| 1 | abcoder 可作为 Go Library 直接集成 | 需验证 Go SDK API |
| 2 | UniAST 能覆盖 11 种语言 | abcoder Writer 仅支持 Go |
| 3 | 修复生成不需要 LLM | Skill Agent 作为生成器 |
| 4 | 单 binary 架构可行 | Go Library 嵌入 vs MCP Server |

---

## 方案对比

### 方案 A：MCP Server Sidecar
- **概述:** abcoder 独立进程，stdio 通信
- **优势:** 完全解耦
- **风险:** 需安装两个组件
- **排除原因:** 违背单一 binary 原则

### 方案 B：Go Library 直接嵌入 ✅
- **概述:** abcoder 作为 Go 依赖直接引入
- **优势:** 单一 binary，用户体验最佳
- **风险:** abcoder 版本耦合
- **选择原因:** 用户体验最优

### 方案 C：渐进增强模式
- **概述:** 按需调用 abcoder
- **优势:** 风险可控
- **风险:** 决策逻辑复杂
- **排除原因:** 增加不必要复杂性

---

## 最终决策

| 决策点 | 选择 | 原因 |
|--------|------|------|
| 集成模式 | Go Library 直接嵌入 | 单一 binary |
| 上下文获取 | abcoder UniAST Parser | 结构化代码理解 |
| 修复生成 | Subagent (`code-reviewer`) | 利用现有 skill |
| 修复验证 | `security-reviewer` | 安全性验证 |
| 回退策略 | 复用现有 YAML `suggestion` | 零增量维护 |

---

## 架构总览

```
CodeSentry (单 binary)
├── Engine (检测漏洞)
├── abcoder Bridge (上下文)
├── Skill Agent (修复生成)
└── 回退 → YAML suggestion
```

---

## 关键假设（待验证）

| # | 假设 | 验证方式 |
|---|------|----------|
| 1 | abcoder Go Library 可通过 `go get` 引入 | 需实际验证 import 路径 |
| 2 | UniAST Parser 可独立使用 | 需研究 abcoder API |
| 3 | `code-reviewer` 可接收结构化上下文 | 需测试 agent 接口 |
| 4 | YAML suggestion 可直接复用 | 需确认内容质量 |

---

## 实施阶段

### Phase 1：基础集成
- [ ] 引入 abcoder 依赖
- [ ] 实现 abcoder Bridge
- [ ] 验证 UniAST 解析

### Phase 2：上下文理解
- [ ] abcoder → CodeSentry 上下文传递
- [ ] 实现上下文接口

### Phase 3：Skill 生成修复
- [ ] 集成 code-reviewer skill
- [ ] 集成 security-reviewer 验证

### Phase 4：多语言 + 回退
- [ ] 启用支持语言的 abcoder
- [ ] 实现模板化回退

---

## 风险与缓解

| 风险 | 影响 | 缓解 |
|------|------|------|
| abcoder API 不稳定 | 高 | Phase 1 验证，锁定版本 |
| Subagent 调用延迟 | 中 | 异步处理 |
| 多语言回退质量差 | 低 | 逐步丰富 suggestion |

---

*创建时间: 2026-04-22*
