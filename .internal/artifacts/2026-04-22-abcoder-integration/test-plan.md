# Test Plan - ABCoder Integration

**任务:** abcoder-integration
**日期:** 2026-04-22
**阶段:** review
**主责角色:** qa-engineer

---

## 1. 测试范围

### 1.1 功能测试

| # | 测试项 | 类型 | 状态 |
|---|--------|------|------|
| 1 | abcoder Bridge 初始化 | 单元测试 | ✅ 通过 |
| 2 | abcoder Repository 解析 | 单元测试 | ✅ 通过 |
| 3 | CodeContext 获取 | 单元测试 | ✅ 通过 |
| 4 | Skill Agent 修复生成 | 单元测试 | ✅ 通过 |
| 5 | Fallback Handler 回退 | 单元测试 | ✅ 通过 |
| 6 | 多语言回退 (Go/Py/JS/TS/Java/Ruby/Rust) | E2E 测试 | ✅ 通过 |
| 7 | 修复输出格式化 (JSON/Text) | E2E 测试 | ✅ 通过 |
| 8 | 并发访问 Bridge | E2E 测试 | ✅ 通过 |

### 1.2 回归测试

| # | 测试项 | 状态 |
|---|--------|------|
| 1 | cmd/goreview | ✅ 通过 |
| 2 | internal/engine | ✅ 通过 |
| 3 | internal/parser | ✅ 通过 |
| 4 | internal/rules | ✅ 通过 |

---

## 2. 测试矩阵

| 场景 | 前置条件 | 输入 | 预期结果 |
|------|---------|------|---------|
| Go 文件 + abcoder 可用 | Go 文件 | 检测结果 + 上下文请求 | 返回完整上下文 |
| 非 Go 文件 | Python/JS/Java 等 | 检测结果 | 回退到 YAML suggestion |
| abcoder 解析失败 | 解析错误 | - | 回退到模板化建议 |
| Skill Agent 生成修复 | 上下文可用 | ruleID, suggestion | 生成修复前后代码 |
| Skill Agent 生成修复 | 上下文不可用 | ruleID, suggestion | 使用模板化修复 |

---

## 3. 风险评估

### 3.1 高风险路径

| # | 风险 | 影响 | 缓解 |
|---|------|------|------|
| 1 | abcoder 解析大型代码库性能 | 扫描延迟增加 | 增量解析，按需加载 |
| 2 | abcoder API 不稳定 | 升级可能导致破坏 | 锁定版本 v0.3.1 |

### 3.2 中风险路径

| # | 风险 | 影响 | 缓解 |
|---|------|------|------|
| 1 | 非 Go 语言上下文缺失 | 修复建议质量下降 | 回退到模板化建议 |
| 2 | Skill Agent 修复置信度 | 修复建议可能不准确 | 标注置信度，用户判断 |

### 3.3 低风险路径

| # | 风险 | 影响 | 缓解 |
|---|------|------|------|
| 1 | 多余日志输出 | 轻微性能影响 | 可接受 |

---

## 4. 已验证项

### 4.1 单元测试覆盖

```
internal/abcoder/
├── bridge.go         - TestNewBridge, TestBridgeParse, TestBridgeGetContext
├── skill.go          - TestSkillAgent_GenerateFix, TestSkillOutput_ToJSON, TestSkillOutput_FormatFix
├── fallback.go       - TestNewFallbackHandler, TestFallbackHandler_GetFix, TestIsFallbackNeeded
└── e2e_test.go      - TestE2E_CompleteFlow, TestE2E_MultipleLanguages, TestE2E_SkillOutputFormat
```

### 4.2 构建验证

```
✅ go build ./...        - 成功
✅ go test ./...        - 全部通过
✅ Binary 构建          - codesentry (13MB)
```

---

## 5. 阻塞项

| # | 阻塞项 | 状态 |
|---|--------|------|
| 1 | 无 | - |

---

## 6. 放行建议

| 检查项 | 状态 |
|--------|------|
| 单元测试通过 | ✅ |
| 回归测试通过 | ✅ |
| E2E 测试通过 | ✅ |
| Binary 构建成功 | ✅ |
| 阻塞项 | 无 |

**建议：✅ 建议放行**

---

*创建时间: 2026-04-22*
*主责角色: qa-engineer*
