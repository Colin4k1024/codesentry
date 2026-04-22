# Delivery Plan - CodeSentry 项目完善

**Slug**: codesentry-project-review
**日期**: 2026-04-22
**状态**: draft → plan
**主责**: tech-lead

---

## 1. 需求挑战会结论

### 核心假设验证

| 分组 | 原假设 | 质疑结论 | 决策 |
|------|--------|----------|------|
| A | 测试覆盖率 > 80% | 风险分级覆盖率更合理 | 核心包(engine, abcoder, parser core) > 80%; 边缘包可 40-60% |
| B | JavaScript 需要专用规则 | JS 规则目录不存在，TS 规则可能有重叠 | 先评估 TS 规则对 JS 的覆盖，再决定是否新建 |
| C | Go 需要单元测试 | AST 代码从未验证，golden test 可能不够 | 先验证 Go AST 实现正确性，再设计测试策略 |

### 关键依赖关系

```
Group C (Go AST) ──[blocked]──→ Go AST 代码未验证
                                     ↓
Group B (JS 规则) ←── 依赖 ──→ langs/javascript/parser.go 架构确认
                                     ↓
Group A (覆盖率) ←── engine 是所有 parser 基础 ──→ internal/engine
```

---

## 2. Story Slices

### Slice 1: Go AST 实现验证
- **Owner**: backend-engineer
- **目标**: 验证 `langs/golang/parser.go` 的 AST 检测逻辑正确性
- **验收**: 输出 Go AST 测试用例，覆盖 GOROUTINE_LEAK, CONTEXT_LEAK, RESOURCE_LEAK
- **依赖**: 无
- **风险**: AST 逻辑可能有隐藏 bug

### Slice 2: JavaScript/TypeScript 规则覆盖评估
- **Owner**: security-reviewer
- **目标**: 评估现有 TS 规则对 JavaScript 的覆盖程度
- **验收**: 输出覆盖分析报告，建议 JS 规则补充方案或确认复用 TS 规则
- **依赖**: Slice 1 完成
- **风险**: 若需新建 JS parser 架构，工作量增加

### Slice 3: 分层测试覆盖率提升
- **Owner**: qa-engineer
- **目标**: 提升核心包覆盖率
  - `internal/engine`: 47% → 80%+
  - `internal/abcoder`: 45% → 80%+
  - `cmd/goreview`: 13% → 40%+ (CLI wrapper，降低优先级)
- **验收**: 覆盖率达标，新增测试通过 CI
- **依赖**: Slice 1 完成（Go parser 是 engine 的一部分）
- **风险**: 测试桩可能与真实检测能力脱节

### Slice 4 (可选): JavaScript 规则补充
- **Owner**: backend-engineer
- **目标**: 根据 Slice 2 结论，补充 JS 专用规则
- **验收**: JS 扫描能发现对应安全漏洞
- **依赖**: Slice 2 结论确认需要 JS 规则
- **风险**: 可能需要 parser 架构升级

---

## 3. 执行计划

| 阶段 | 分组 | 主责角色 | 依赖 | 输出 |
|------|------|----------|------|------|
| S1 | Go AST 验证 | backend-engineer | 无 | Go parser 测试用例 |
| S2 | JS/TS 覆盖评估 | security-reviewer | S1 | 覆盖分析报告 |
| S3 | 分层覆盖率提升 | qa-engineer | S1 | 覆盖率报告 |
| S4 | JS 规则补充 | backend-engineer | S2 | 新增 JS 规则 |

**并行性**: S1 独立执行，S2 可与 S1 并行，S3 依赖 S1，S4 依赖 S2

---

## 4. 风险与依赖清单

| 风险 | 影响 | 缓解 |
|------|------|------|
| Go AST bug 未被发现 | 检测结果不可靠 | S1 先验证正确性 |
| JS 需 parser 架构升级 | 工作量大幅增加 | S2 先评估再决策 |
| 测试桩与真实检测脱节 | 覆盖率无实际意义 | 按 architect 建议分层 |
| 三组相互依赖 | 串行化导致总工期增加 | S1 独立，S2 可并行 |

---

## 5. 关键决策

1. **覆盖率目标**: 核心包 > 80%，边缘包 40-60%
2. **JS 规则策略**: 先评估再决定，不盲目新建
3. **Go 测试策略**: golden test + 状态机测试结合
4. **Parser 能力矩阵**: 作为技术债务记录

---

## 6. 非目标 (Out of Scope)

- 不新增语言支持
- 不修改 abcoder 集成架构
- 不改变输出格式 (Text/JSON/SARIF)
- 不添加 UI

---

## 7. 验收标准

- [ ] Go AST 检测逻辑有测试验证
- [ ] JS/TS 规则覆盖关系明确
- [ ] 核心包覆盖率达标
- [ ] 所有新增测试通过
- [ ] `go build` 和 `go test ./...` 通过

---

**下一步**: 进入 /team-execute 执行 Slice 1-3
