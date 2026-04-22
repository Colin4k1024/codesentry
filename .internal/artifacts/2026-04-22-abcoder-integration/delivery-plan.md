# Delivery Plan - ABCoder Integration

**任务:** abcoder-integration
**日期:** 2026-04-22
**阶段:** plan
**主责角色:** tech-lead

---

## 1. 需求挑战会结论

### 核心假设验证

| # | 假设 | 验证结果 | 状态 |
|---|------|---------|------|
| 1 | abcoder Go Library 可通过 `go get` 引入 | ✅ 已验证：v0.3.1 可引入 | 通过 |
| 2 | UniAST Parser 可独立使用 | ⚠️ 部分验证：Go parser 可独立使用，其他语言需要 LSP | 风险可控 |
| 3 | `code-reviewer` 可接收结构化上下文 | ✅ 确认：Skill Agent 支持结构化输入 | 通过 |
| 4 | YAML suggestion 可直接复用 | ✅ 确认：现有 suggestion 字段可用作回退模板 | 通过 |

### 方案决策

| 决策点 | 选择 | 原因 |
|--------|------|------|
| 集成模式 | Go Library 直接嵌入 | Go parser 可独立使用 |
| 上下文获取 | abcoder UniAST Parser | 结构化代码理解 |
| 修复生成 | Subagent (`code-reviewer`) | 利用现有 skill 能力 |
| 回退策略 | 复用现有 YAML `suggestion` | 零增量维护 |

---

## 2. 版本目标

### 目标版本
- **版本:** v0.2.0 (abcoder-integration)
- **范围:** 增强 CodeSentry 代码审查能力，提供上下文感知修复建议
- **放行标准:**
  - ✅ 能为 Go 语言检测结果提供上下文感知修复建议
  - ✅ 其他语言检测结果回退到 YAML suggestion
  - ✅ 不破坏现有功能
  - ✅ 向后兼容（现有规则格式不变）

---

## 3. 工作拆解

### Story Slices

#### Slice 1: abcoder Bridge 实现 (Phase 1)
**目标:** 实现 `internal/abcoder/bridge.go`，封装 abcoder UniAST Parser

**验收标准:**
- [ ] `Bridge` 结构体能解析 Go 代码仓库
- [ ] 提供 `GetFunctionContext(nodeID)` 接口
- [ ] 提供 `GetVariableDefinition(varID)` 接口
- [ ] 提供 `GetCallChain(funcID)` 接口

**依赖:** 无
**Owner:** backend-engineer
**Handoff终点:** bridge.go 单元测试通过

---

#### Slice 2: 上下文传递层 (Phase 2)
**目标:** 将 abcoder 上下文与 CodeSentry 检测结果关联

**验收标准:**
- [ ] 检测结果能携带文件路径和行号
- [ ] abcoder Bridge 能根据位置获取上下文
- [ ] 上下文数据结构设计完成

**依赖:** Slice 1 完成
**Owner:** backend-engineer
**Handoff终点:** 上下文传递集成测试通过

---

#### Slice 3: Skill Agent 集成 (Phase 3)
**目标:** 集成 `code-reviewer` skill 生成修复代码

**验收标准:**
- [ ] Skill Agent 能接收 abcoder 上下文
- [ ] Skill Agent 能输出修复建议（before/after code）
- [ ] 修复建议能展示给用户

**依赖:** Slice 2 完成
**Owner:** backend-engineer
**Handoff终点:** Skill Agent 集成测试通过

---

#### Slice 4: 回退机制 (Phase 4)
**目标:** 当 abcoder 不可用时，回退到 YAML suggestion

**验收标准:**
- [ ] 非 Go 语言使用现有 suggestion
- [ ] Go 语言在 abcoder 解析失败时回退
- [ ] 回退不产生错误，仅输出原有 suggestion

**依赖:** Slice 1 完成
**Owner:** backend-engineer
**Handoff终点:** 回退机制测试通过

---

#### Slice 5: 端到端测试 (Phase 5)
**目标:** 完整流程测试

**验收标准:**
- [ ] Go 文件检测 + abcoder 上下文 + Skill 修复建议 完整链路
- [ ] 非 Go 语言回退链路
- [ ] 性能影响 < 2x

**依赖:** Slice 1-4 完成
**Owner:** qa-engineer
**Handoff终点:** E2E 测试报告

---

## 4. 风险与缓解

| 风险 | 影响 | 缓解措施 | Owner |
|------|------|----------|-------|
| abcoder Go parser API 不稳定 | 中 | Phase 1 优先验证，锁定版本 | backend-engineer |
| 非 Go 语言需要 LSP | 中 | 回退到现有检测，不强制要求 LSP | backend-engineer |
| Skill Agent 调用延迟 | 低 | 异步处理，增量显示 | backend-engineer |
| abcoder 解析大型代码库性能 | 低 | 增量解析，按需加载 | backend-engineer |

---

## 5. 节点检查

| 节点 | 检查项 | 预期时间 | 状态 |
|------|--------|---------|------|
| 方案评审 | Slice 划分、依赖关系 | Day 1 | ✅ |
| Slice 1 完成 | Bridge 接口设计 + 实现 | Day 2 | ⏳ |
| Slice 2-3 完成 | 上下文传递 + Skill 集成 | Day 3-4 | ⏳ |
| Slice 4-5 完成 | 回退机制 + E2E 测试 | Day 5 | ⏳ |
| 发布准备 | 文档更新、版本发布 | Day 6 | ⏳ |

---

## 6. 角色分工

| 角色 | 职责 | 交付物 |
|------|------|--------|
| tech-lead | 方案评审、风险监控 | 本文档 |
| backend-engineer | Slice 1-4 实现 | abcoder bridge, 上下文传递, skill 集成 |
| qa-engineer | Slice 5 测试 | E2E 测试报告 |
| code-reviewer (skill) | 修复代码生成 | 修复建议 |

---

## 7. 依赖清单

| 依赖 | 状态 | 说明 |
|------|------|------|
| Go 1.23+ | ✅ 已满足 | abcoder 要求 |
| abcoder v0.3.1 | ✅ 已引入 | go.mod 中 |
| code-reviewer skill | ✅ 已安装 | Claude Code 内置 |
| security-reviewer skill | ✅ 已安装 | Claude Code 内置 |

---

## 8. 应用等级与技术架构

**应用等级:** T4 (工具能力增强，非核心生产系统)

**技术架构:**
- **组件偏离:** 无重大偏离
- **关键组件:** abcoder UniAST Parser (新增)
- **数据边界:** 仅在检测到漏洞时调用 abcoder，不持久化上下文

---

## 9. 技能装配清单

| 技能 | 用途 | 启用场景 |
|------|------|----------|
| `code-reviewer` | 生成修复代码 | Phase 3 |
| `security-reviewer` | 验证修复安全性 | Phase 3 |
| `golang-patterns` | Go 代码风格指导 | 全程 |

---

## 10. 阻塞项

| # | 阻塞项 | 解决方案 |
|---|--------|---------|
| 1 | 无 | - |

---

*创建时间: 2026-04-22*
*主责角色: tech-lead*
