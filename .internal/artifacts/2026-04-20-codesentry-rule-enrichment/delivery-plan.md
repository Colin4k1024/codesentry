# Delivery Plan - CodeSentry 规则丰富与本地测试体系建设

**Slug:** `codesentry-rule-enrichment`
**创建日期:** 2026-04-20
**主责角色:** tech-lead
**状态:** plan (draft)

---

## 需求挑战会结论（Challenge Session Log）

### 假设 1：「标准 Go testing 足够」
- **质疑**：未解决 parser 重复代码反模式，会导致 10 处重复测试
- **替代路径**：标准库 + go-cmp，同时重构 parser 抽取 BaseParser
- **阻断条件**：不统一 parser 抽象就选测试工具，测试会复制错误设计

### 假设 2：「golden file 可验证所有规则」
- **质疑**：AST 规则声明与实现严重脱节（YAML 有 `type: ast`，但 parser 完全忽略）
- **替代路径**：短期：golden file 验证 hardcoded AST；中期：重建 YAML-driven AST engine
- **阻断条件**：`rules/go/goroutine_leak.yaml` 的 ast pattern 声明是谎言，测试设计不能基于未实现的功能

### 假设 3：「80% 覆盖率是合适指标」
- **质疑**：覆盖率与规则质量解耦，可被作弊，且关键路径（category filter）未覆盖
- **替代路径**：分模块覆盖目标（engine ≥ 85%, parser ≥ 50%）+ 规则有效性 golden file 误报率指标
- **阻断条件**：`engine.Scan()` 的 category filter 分支（143-144 行）未覆盖，关键安全功能无验证

---

## 技术架构决策

### ADR-001: Parser 抽象重构
**决策：** 在测试体系建设之前，先抽取 BaseParser 消除 10 个重复 parser

**影响：**
- 短期：测试跟着正确设计走，不会复制错误模式
- 中期：AST engine 实现时，只需要修改 BaseParser

### ADR-002: AST 规则实现策略
**决策：** 短期保持 hardcoded AST（GoParser），golden file 验证 hardcoded 行为；中期重构为 YAML-driven

**理由：** 当前 YAML `type: ast` 声明与实现脱节，golden file 测试应基于实际实现而非声明

### ADR-003: 覆盖率指标
**决策：** 分模块设定覆盖率阈值

| 模块 | 覆盖率目标 |
|------|-----------|
| engine | ≥ 85% |
| rules/loader | ≥ 80% |
| parser/registry | ≥ 70% |
| 各语言 parser | ≥ 50%（regex 逻辑简单）|

---

## Story Slices

### Slice 1: Parser 抽象重构
**目标：** 消除 10 个重复 parser，建立 BaseParser

**验收标准：**
- [ ] 10 个纯 regex parser 复用 BaseParser，代码行数减少 ≥ 60%
- [ ] 新增语言只需实现 `ParseRegex` 方法
- [ ] GoParser 保留 AST 逻辑，不受影响

**依赖：** 无
**Owner:** backend-engineer
**Handoff:** 交付代码 + parser 重构前后对比报告

---

### Slice 2: 测试框架选型与基础设施
**目标：** 建立测试基础设施，选择轻量测试工具

**验收标准：**
- [ ] 引入 `github.com/google/go-cmp/cmp` 用于结构比较
- [ ] 建立 `testdata/` 目录结构
- [ ] 测试可通过 `go test ./...` 运行

**依赖：** Slice 1 完成
**Owner:** backend-engineer
**Handoff:** 交付测试基础设施代码

---

### Slice 3: Engine 核心路径测试
**目标：** 覆盖 engine.Scan() 关键路径，特别是 category filter

**验收标准：**
- [ ] 测试 Security filter 分支（cfg.Security=true 时只走 security 规则）
- [ ] 测试 Performance filter 分支
- [ ] 测试无 filter 时全量规则行为
- [ ] 测试 Finding dedup 去重逻辑
- [ ] 测试 node_modules/vendor 跳过逻辑

**依赖：** Slice 2 完成
**Owner:** backend-engineer
**Handoff:** engine_test.go 覆盖 engine.go 关键分支

---

### Slice 4: Rules Loader 测试
**目标：** 覆盖 LoadRules 错误处理路径

**验收标准：**
- [ ] 测试无效 YAML 文件处理
- [ ] 测试缺失 ID 的规则跳过
- [ ] 测试空目录返回空 slice

**依赖：** Slice 2 完成
**Owner:** backend-engineer
**Handoff:** loader_test.go 覆盖 loader.go 错误路径

---

### Slice 5: Golden File 测试框架
**目标：** 建立 golden file 测试验证规则正确性

**验收标准：**
- [ ] golden file 格式定义：`<rule_id>.input.<ext>` + `<rule_id>.golden.json`
- [ ] 至少 3 个 golden file 测试用例（hardcoded_secret, sql_injection, goroutine_leak）
- [ ] 测试可检测 regex 误报

**依赖：** Slice 2 完成
**Owner:** backend-engineer
**Handoff:** golden file 测试框架 + 至少 3 个测试样例

---

### Slice 6: Parser Registry 测试
**目标：** 覆盖 parser 注册与查询逻辑

**验收标准：**
- [ ] 测试 Register/Get 配对
- [ ] 测试 DetectFromPath 正确语言检测
- [ ] 测试未知扩展返回 nil

**依赖：** Slice 1 完成
**Owner:** backend-engineer
**Handoff:** registry_test.go 覆盖 registry.go

---

### Slice 7: 安全规则补全（Phase 1）
**目标：** 补全 Go/Python/JavaScript 安全规则

**验收标准：**
- [ ] Go: 现有规则保持，新增 3 条（exec, path_traversal, unsafe_deserialization）
- [ ] Python: 现有 4 条规则，新增 2 条（yaml_load, subprocess）
- [ ] JavaScript: 现有规则，新增 2 条（prototype_pollution, eval_with_input）

**依赖：** Slice 5 完成
**Owner:** backend-engineer + security-reviewer
**Handoff:** 新增规则 YAML 文件 + 对应 golden file 测试

---

## 角色分工

| Slice | Owner | Reviewer |
|-------|-------|----------|
| 1: Parser 重构 | backend-engineer | tech-lead |
| 2: 测试基础设施 | backend-engineer | tech-lead |
| 3: Engine 测试 | backend-engineer | qa-engineer |
| 4: Loader 测试 | backend-engineer | qa-engineer |
| 5: Golden File 框架 | backend-engineer | tech-lead |
| 6: Registry 测试 | backend-engineer | tech-lead |
| 7: 规则补全 | backend-engineer + security-reviewer | tech-lead |

---

## 风险与依赖

| 风险 | 影响 | 缓解措施 |
|------|------|---------|
| Parser 重构破坏现有功能 | 高 | Slice 1 后运行完整集成测试 |
| Golden file 测试维护成本 | 中 | 限制 golden file 数量，每个规则 ≤ 3 个样例 |
| AST engine 重构需求 | 低 | 记录为 tech debt，未来 phase 处理 |

---

## 检查节点

| 节点 | 触发条件 | 通过标准 |
|------|---------|---------|
| **CP1: Parser 重构完成** | Slice 1 完成后 | `go build ./...` 成功，`./codesentry scan ./...` 功能正常 |
| **CP2: 测试框架就绪** | Slice 2 完成后 | `go test ./... -cover` 显示覆盖率数据 |
| **CP3: 关键路径覆盖** | Slice 3-4 完成后 | engine.go 覆盖率 ≥ 85%，loader.go ≥ 80% |
| **CP4: 规则测试就绪** | Slice 5 完成后 | Golden file 测试全部通过 |
| **CP5: 规则补全完成** | Slice 7 完成后 | 每语言至少 5 条安全规则 |

---

## 技术栈

- **语言:** Go 1.21+
- **依赖:**
  - `github.com/spf13/cobra v1.8.0` (CLI)
  - `gopkg.in/yaml.v3 v3.0.1` (YAML parsing)
  - `github.com/google/go-cmp v0.6.0` (测试结构比较，新增)
- **测试:** 标准 Go testing + go-cmp
- **测试数据:** `testdata/` 目录

---

## 执行前提（Implementation Readiness）

- [x] 需求挑战会完成
- [x] Story slices 划分完成
- [ ] Parser 重构方案评审（建议 CP1 前）
- [ ] AST 规则实现策略确认（已确认：短期保持 hardcoded）
- [ ] 覆盖率阈值确认（已确认：分模块设定）

---

## 下一步

进入 `/team-execute` 阶段前，需确认：
1. Parser 重构方案已评审
2. 测试基础设施设计已确认
3. Backend-engineer 已接受 Slice 1-6 任务

**预计工期：** 2-3 周（并行工作可缩短）
