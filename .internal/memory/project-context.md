# Project Context

**最后更新:** 2026-04-22

---

## 项目概览

| 字段 | 内容 |
|------|------|
| **项目名** | CodeSentry |
| **项目路径** | /Users/ailabuser1/Desktop/gitcode/codesentry |
| **技术栈** | Go 1.21+, Cobra CLI, YAML rules |
| **当前阶段** | plan |
| **关联任务** | 2026-04-20-codesentry-rule-enrichment, 2026-04-22-abcoder-integration |

---

## 当前任务

### 2026-04-20-codesentry-rule-enrichment

**目标：** 建立本地测试体系 + 丰富安全规则

**Phase:** handoff-ready

---

### 2026-04-22-abcoder-integration

**目标：** 集成 abcoder 框架，增强代码审查与修复建议能力

**Phase:** review → release ✅

**核心需求：**
- 集成 abcoder UniAST 统一代码表示
- 利用 Code-RAG 增强代码上下文理解
- 基于结构化代码上下文生成修复建议

**关键约束：**
- abcoder Go parser 可独立使用，其他语言需要 LSP
- 不改变现有 YAML 规则格式
- 修复建议为补充能力，不替代现有正则检测
- 性能约束：单次扫描延迟增加 < 2x

**已完成 Slices:**
- Slice 1: abcoder Bridge 实现 ✅
- Slice 2: 上下文传递层 ✅
- Slice 3: Skill Agent 集成 ✅
- Slice 4: 回退机制 ✅
- Slice 5: 端到端测试 ✅

**验收结论:** ✅ 建议放行

**方案决策：**
- 集成模式：Go Library 直接嵌入
- 上下文获取：abcoder UniAST Parser
- 修复生成：Subagent (code-reviewer skill)
- 回退策略：复用现有 YAML suggestion

---

## 技术架构

### 核心模块

| 模块 | 文件 | 复杂度 |
|------|------|--------|
| engine | internal/engine/engine.go (175行) | 高 - category filter 关键路径 |
| abcoder | internal/abcoder/ (新增) | 中 - abcoder Bridge |
| rules/loader | internal/rules/loader.go (75行) | 中 - 错误处理路径 |
| parser/registry | internal/parser/registry.go (59行) | 低 - 接口清晰 |
| 11 language parsers | langs/*/parser.go (各50-60行) | 低 - 10个重复，1个有AST |

### 架构决策

- **ADR-001:** Parser 抽象重构（BaseParser）
- **ADR-002:** AST 规则短期保持 hardcoded
- **ADR-003:** 覆盖率分级设定

---

## 风险项

| 风险 | 影响 | 状态 |
|------|------|------|
| Parser 重构破坏现有功能 | 高 | 监控中 |
| Golden file 维护成本 | 中 | 已识别 |
| AST engine 重构需求 | 低 | tech debt |
| abcoder API 不稳定 | 中 | 锁定版本 v0.3.1 |
| 非 Go 语言需要 LSP | 中 | 回退机制 |

---

## 下一步行动

0. **abcoder 集成：** ✅ 完成，QA 验收通过
1. **rule-enrichment：** handoff-ready，等待执行
2. **parser 重构：** 方案评审中
3. **测试体系：** 等待 Slice 1-6 执行

---

## 依赖

| 依赖 | 状态 |
|------|------|
| Go 1.23+ | ✅ 已安装（abcoder 要求） |
| go-cmp | ⏳ 待引入 |
| testify | ❌ 不引入（轻量方案） |
| abcoder v0.3.1 | ✅ 已引入 go.mod |
