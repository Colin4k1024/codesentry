# Decisions Log

**最后更新:** 2026-04-20

---

## 2026-04-20: Parser 抽象重构技术决策

### D1: 使用组合而非继承

**决策：** 10 个 parser 嵌入 `BaseRegexParser` 而非继承

**背景：** Go 不支持多继承，但支持结构体嵌入（composition）

**替代方案考虑：** 继承会限制未来灵活性，组合更灵活

**影响：** Parser 解析逻辑统一到 BaseRegexParser，后续 AST engine 可扩展

---

### D2: GoParser 保留 AST 逻辑

**决策：** GoParser 保留原有 AST 检查（checkAST 方法），不迁移到 BaseRegexParser

**背景：** AST 检查是 Go 特有的，与 regex 逻辑正交

**影响：** 保持 Go 特殊性的同时复用 regex 逻辑

---

### D3: 测试框架选型

**决策：** 标准 Go testing + go-cmp（轻量方案）

**背景：** 避免重型依赖（testify），但需要解决结构比较痛点

**影响：** 新增 `github.com/google/go-cmp` 依赖

---

### D4: 测试包位置

**决策：** Registry 测试放在 `cmd/goreview` 而非 `internal/parser`

**背景：** `internal/parser` 无法导入 langs 包（循环依赖），`cmd/goreview` 已导入所有 langs

**影响：** 测试按包分散，但逻辑清晰

---

### D5: JavaScript/TypeScript 扩展重叠

**决策：** JavaScript parser 只保留 .js, .jsx, .mjs, .cjs，TypeScript parser 处理 .ts, .tsx

**背景：** 原 JavaScript parser 包含 .ts, .tsx，导致 .tsx 被错误识别为 JavaScript

**影响：** 修复了扩展检测 bug
