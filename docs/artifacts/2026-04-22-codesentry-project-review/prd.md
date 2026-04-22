# CodeSentry 项目审查 PRD

## 1. 背景与目标

**项目**: CodeSentry - 快速、可扩展的静态分析与AI辅助代码审查工具

**核心目标**:
- 提供多语言安全漏洞扫描能力
- 通过 abcoder UniAST 为 Go 代码提供上下文理解的修复建议
- 支持 YAML 规则扩展

**成功标准**:
- 11 种编程语言支持
- 所有语言解析器有单元测试覆盖
- 核心引擎测试覆盖率 > 80%
- 支持 Text/JSON/SARIF 输出格式

---

## 2. 当前项目状态

### 2.1 已完成功能

| 功能 | 状态 | 说明 |
|------|------|------|
| 多语言支持 | ✅ | Go, Python, TypeScript, JavaScript, Java, Ruby, Rust, C++, PHP, Swift, Kotlin |
| abcoder 集成 | ✅ | Go 代码上下文理解与修复建议生成 |
| YAML 规则引擎 | ✅ | 跨语言安全规则与语言特定规则 |
| 多输出格式 | ✅ | Text, JSON, SARIF |
| 单元测试 | ✅ | 9 个语言解析器 100% 覆盖 |

### 2.2 测试覆盖率

| 模块 | 覆盖率 | 备注 |
|------|--------|------|
| `langs/*` (9个解析器) | 100% | 最近添加的单元测试 |
| `internal/rules` | 81.2% | 最高覆盖率 |
| `internal/engine` | 47.3% | 需要提升 |
| `internal/parser` | 48.3% | 需要提升 |
| `internal/abcoder` | 45.7% | 需要提升 |
| `cmd/goreview` | 13.1% | 需要提升 |
| `internal/types` | N/A | 无测试文件 |
| `internal/output` | N/A | 无测试文件 |

### 2.3 Git 提交历史

```
08f7bc8 test: add parser unit tests for 9 languages
324942e docs: reorganize and polish documentation for开源 release
fd923ea feat: integrate abcoder for code context and fix suggestions
319ce42 Rename project to CodeSentry
655a66d Add documentation: README, CONTRIBUTING, RULES guide
30f4f4a Add Rust, C++, PHP, Swift, Kotlin parsers and security rules
9250f5f Add Python, TypeScript, Java, Ruby parsers and security rules
```

---

## 3. 待完善功能

### 3.1 测试覆盖

| 优先级 | 模块 | 当前状态 | 建议 |
|--------|------|----------|------|
| HIGH | `cmd/goreview` | 13.1% | 添加 CLI 命令测试 |
| HIGH | `internal/abcoder` | 45.7% | 添加更多端到端测试 |
| MEDIUM | `internal/engine` | 47.3% | 添加更多扫描场景测试 |
| MEDIUM | `internal/parser` | 48.3% | 添加解析器边界测试 |
| LOW | `internal/types` | 无测试 | 添加类型定义测试 |
| LOW | `internal/output` | 无测试 | 添加输出格式化测试 |

### 3.2 语言解析器

| 语言 | 状态 | 备注 |
|------|------|------|
| JavaScript | ⚠️ 无专用规则 | 仅使用跨语言规则 |
| Go | ⚠️ 无单元测试 | AST 方式不同，需专门测试策略 |
| Python | ✅ | 有测试 |
| TypeScript | ✅ | 有测试 |
| Java | ✅ | 有测试 |
| Ruby | ✅ | 有测试 |
| Rust | ✅ | 有测试 |
| C++ | ✅ | 有测试 |
| PHP | ✅ | 有测试 |
| Swift | ✅ | 有测试 |
| Kotlin | ✅ | 有测试 |

### 3.3 测试数据

| 测试类型 | 状态 | 备注 |
|----------|------|------|
| 单元测试 | ✅ | 9 个语言解析器 |
| Golden File 测试 | ⚠️ 稀疏 | 仅 HARDCODED_SECRET 有测试数据 |
| E2E 测试 | ✅ | abcoder 有完整 E2E 测试 |

---

## 4. 风险与约束

### 4.1 技术风险

| 风险 | 影响 | 缓解 |
|------|------|------|
| abcoder 仅支持 Go | 非 Go 语言无上下文理解 | 使用模板化修复建议作为回退 |
| JavaScript 无专用规则 | 扫描能力受限 | 可按需添加 |
| Go 解析器无单元测试 | 回归风险 | 可考虑添加集成测试 |

### 4.2 项目约束

- Go 1.23+ (abcoder 依赖)
- 单二进制部署，无外部依赖
- MIT 许可证

---

## 5. 需求挑战会候选分组

建议以下分组参与需求挑战会：

### 分组 A: 测试覆盖率提升

- **参与角色**: qa-engineer, backend-engineer
- **范围**: 提升 `cmd/goreview`、`internal/abcoder`、`internal/engine` 测试覆盖率
- **目标**: 核心包覆盖率 > 80%

### 分组 B: JavaScript 规则补充

- **参与角色**: backend-engineer, security-reviewer
- **范围**: 为 JavaScript 添加专用规则 (如 `JS_EVAL`, `JS_XSS`)
- **目标**: 完善 JS 扫描能力

### 分组 C: Go 解析器测试

- **参与角色**: backend-engineer
- **范围**: 为 Go AST 解析器设计测试策略
- **目标**: 可测试的 AST 规则覆盖

---

## 6. 参与角色清单

| 角色 | 输入缺口 | 优先级 |
|------|----------|--------|
| tech-lead | 项目方向确认 | HIGH |
| qa-engineer | 测试策略建议 | HIGH |
| backend-engineer | 实现能力评估 | HIGH |
| security-reviewer | 规则完整性审查 | MEDIUM |

---

## 7. 待确认项

1. **测试覆盖率目标**: 是否要求所有包 > 80% 覆盖率？
2. **JavaScript 规则**: 是否需要为 JS 添加语言专用规则？
3. **Go 解析器测试**: AST 解析器是否需要单独的测试策略？
4. **Golden File 测试**: 是否需要扩展测试数据集？
5. **文档持久化**: `docs/memory/` 是否需要重建？

---

## 8. 上线验收标准

- [ ] 所有新增代码通过测试
- [ ] `go build` 成功
- [ ] `go test ./...` 全部通过
- [ ] 无新增 lint 错误
- [ ] 代码符合 Go 编码规范

---

**创建日期**: 2026-04-22
**主责角色**: tech-lead
**状态**: intake
