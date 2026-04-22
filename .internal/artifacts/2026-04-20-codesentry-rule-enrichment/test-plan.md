# Test Plan - CodeSentry 规则丰富与本地测试体系

**Slug:** `codesentry-rule-enrichment`
**评审日期:** 2026-04-20
**主责角色:** qa-engineer

---

## 测试范围

### 已测试组件

| 组件 | 测试用例数 | 覆盖率 |
|------|-----------|--------|
| internal/engine | 9 | 84.7% |
| internal/rules/loader | 5 | 81.2% |
| internal/parser | 6 | 48.3% |
| cmd/goreview (含 golden tests) | 6 | 13.1% |

### 测试类型

| 类型 | 描述 |
|------|------|
| 单元测试 | Engine, Parser, Rules Loader 核心逻辑 |
| Golden File 测试 | HARDCODED_SECRET, SQL_INJECTION, GOROUTINE_LEAK |
| 集成测试 | 完整 scan 流程（文件遍历 → 解析 → 规则匹配 → 去重） |

---

## 测试矩阵

| 场景 | 类型 | 前置条件 | 预期结果 |
|------|------|---------|---------|
| 基本扫描 | 集成 | 有效 Go 文件含 hardcoded secret | 检测到 1 个 SEVERE |
| Security Filter | 单元 | cfg.Security=true | 只返回 security 类别 |
| Performance Filter | 单元 | cfg.Performance=true | 只返回 performance 类别 |
| 无 Filter | 单元 | cfg 为空 | 返回全部规则 |
| 跳过目录 | 单元 | node_modules/.git/vendor 存在 | 文件被跳过 |
| 去重 | 单元 | 同一行多个匹配 | 只保留 1 个 finding |
| 未知扩展 | 单元 | .xyz 文件 | 返回 0 findings |
| 读取错误 | 单元 | 损坏的 symlink | 优雅处理，继续扫描 |

---

## 风险与测试覆盖

| 风险 | 严重度 | 测试覆盖 | 备注 |
|------|--------|---------|------|
| Category filter 分支逻辑 | HIGH | ✅ 已覆盖 | TestEngine_Scan_SecurityFilter |
| Dedup 去重逻辑 | HIGH | ✅ 已覆盖 | TestEngine_Scan_Dedup |
| Parser 注册与发现 | MEDIUM | ✅ 已覆盖 | TestParserRegistry, TestDetectFromPath |
| Golden file 格式正确性 | MEDIUM | ✅ 已覆盖 | TestGoldenFile_* |
| AST-based 检测 (Go) | MEDIUM | ⚠️ 部分覆盖 | GOROUTINE_LEAK 已测，其他未测 |
| AST pattern YAML 声明 | LOW | ❌ 未覆盖 | 规则中存在但未实现 |

---

## 已接受风险

| 风险 | 理由 | 缓解措施 |
|------|------|---------|
| AST pattern 在 YAML 中声明但未实现 | 短期保持 hardcoded，待 Phase 2 重构 | 文档化，标记为 tech debt |
| Parser 单元测试缺失 | regex 逻辑在 BaseRegexParser 已覆盖 | 后续 phase 补充 |
| output 包无测试 | 现有测试覆盖核心逻辑 | 后续 phase 补充 |

---

## 阻塞项

| 项 | 状态 | 说明 |
|----|------|------|
| 无 | - | - |

---

## 放行建议

**建议：放行 (Go)**

- 所有测试通过 ✅
- 构建成功 ✅
- 核心路径覆盖率 ≥ 80% ✅
- 无阻塞问题

---

## 前端质量门禁

本次变更不涉及前端组件。
