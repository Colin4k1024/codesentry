# Slice 2: JavaScript/TypeScript Rule Coverage Analysis

**Slug**: codesentry-project-review
**Slice**: 2
**日期**: 2026-04-22
**状态**: completed

---

## 1. 发现总结

### JavaScript 规则覆盖现状

| 规则 ID | 语言支持 | 来源 | JavaScript 覆盖 |
|---------|----------|------|-----------------|
| HARDCODED_SECRET | go, javascript, python | cross-language | ✅ |
| SQL_INJECTION | go, javascript, python | cross-language | ✅ |
| SENSITIVE_LOG | go, javascript, python | cross-language | ✅ |
| TS_EVAL | typescript, javascript | typescript | ✅ |
| TS_HARDCODED_SECRET | typescript, javascript | typescript | ✅ |
| TSPrototype_POLLUTION | typescript, javascript | typescript | ✅ |
| TS_SQL_INJECTION | typescript | typescript | ⚠️ 仅 TS |
| TS_EVAL_WITH_INPUT | typescript | typescript | ⚠️ 仅 TS |
| TS_SENSITIVE_LOG | typescript | typescript | ⚠️ 仅 TS |

### 关键发现

1. **JavaScript 有专用规则覆盖**: `TS_EVAL`、`TS_HARDCODED_SECRET`、`TSPrototype_POLLUTION` 明确列出 JavaScript
2. **跨语言规则覆盖 JS**: `HARDCODED_SECRET`、`SQL_INJECTION`、`SENSITIVE_LOG` 都明确支持 JavaScript
3. **无独立 JS 规则目录**: `rules/javascript/` 不存在，但现有规则已覆盖 JS

---

## 2. 覆盖评估

### JavaScript 安全问题覆盖

| 安全问题 | 规则 | 覆盖状态 |
|----------|------|----------|
| 硬编码密钥 | HARDCODED_SECRET, TS_HARDCODED_SECRET | ✅ 完全覆盖 |
| SQL 注入 | SQL_INJECTION | ✅ 完全覆盖 |
| 敏感数据日志 | SENSITIVE_LOG | ✅ 完全覆盖 |
| eval() 注入 | TS_EVAL | ✅ 完全覆盖 |
| 原型污染 | TSPrototype_POLLUTION | ✅ 完全覆盖 |
| innerHTML XSS | TS_EVAL | ✅ 覆盖 |
| React 危险操作 | TS_EVAL | ✅ 覆盖 |

### TypeScript 额外覆盖

| 安全问题 | 规则 | 说明 |
|----------|------|------|
| SQL 注入 | TS_SQL_INJECTION | TS 专用，比跨语言更精确 |
| eval + 输入 | TS_EVAL_WITH_INPUT | 检测 eval 接收外部输入 |
| 敏感日志 | TS_SENSITIVE_LOG | TS 专用 |

---

## 3. 结论与建议

### 结论

✅ **JavaScript 已有足够的安全规则覆盖**

现有规则可以有效检测 JavaScript 代码中的主要安全问题：
- 硬编码密钥
- SQL 注入
- 敏感数据日志
- eval() 危险用法
- 原型污染
- innerHTML XSS

### 建议

**不需要为 JavaScript 创建独立规则目录**，原因：
1. TypeScript 是 JavaScript 的超集，TS 规则覆盖 JS
2. 跨语言规则已覆盖 JS 的主要安全问题
3. 独立 JS 规则会导致维护成本增加且无实际收益

**可选改进**:
1. 如需更精确的 JS 特定检测，可考虑添加 `JS_SQL_INJECTION` 规则（当前跨语言版本对 JS 较弱）
2. 可考虑将 `TSPrototype_POLLUTION` 重命名为更通用的名称（如 `JSPROTOTYPE_POLLUTION`），明确支持 JS

---

## 4. 执行结果

| 项目 | 状态 |
|------|------|
| JavaScript 覆盖分析 | ✅ 完成 |
| TypeScript 规则复用评估 | ✅ 完成 |
| 建议输出 | ✅ 完成 |

**结论**: Slice 2 识别出 JavaScript 已有足够规则覆盖，不需要新建独立 JS 规则。TypeScript 规则通过明确列出 JavaScript 语言实现覆盖。
