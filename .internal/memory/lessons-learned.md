# Lessons Learned

**最后更新:** 2026-04-20

---

## 2026-04-20: CodeSentry 规则丰富与测试体系

### 架构决策

#### Parser 抽象重构的时机
**场景:** 10 个语言的 parser 存在重复代码

**学到的:**
- 重构应该先于测试体系建设，否则测试会复制错误设计
- 使用结构体嵌入（composition）而非继承，保持灵活性

#### 测试包位置选择
**场景:** `internal/parser` 无法导入 `langs` 包（循环依赖）

**学到的:**
- Go 的 internal 包保护机制导致测试代码组织受限
- 解决方案：将 registry 测试放在 `cmd/goreview` 而非 `internal/parser`
- 教训：设计 package 结构时要考虑测试可测试性

### 代码质量

#### 变量遮蔽问题
**场景:** `langs/golang/parser.go` 中 `p := fset.Position(pos)` 遮蔽了 receiver `p *GoParser`

**学到的:**
- Go 不允许同名变量遮蔽外部作用域变量，但 `:=` 在内层作用域会创建新变量
- 这种 bug 很难发现，建议使用 IDE 的 shadowing 检测或 linter
- 教训：内层循环变量使用不同命名（如 `pos2`）更安全

#### 死代码清理
**场景:** 保留 `ParseRegexOld` 作为"向后兼容"但从未使用

**学到的:**
- 标记为 deprecated 的代码如果从未使用，应该直接删除
- 死代码会误导维护者，增加理解成本
- 教训：没有调用方的"兼容"代码就是技术债务

### 安全规则

#### Rule ID 命名一致性
**场景:** `TSPrototype_POLLUTION` 使用混合大小写，与其他规则不一致

**学到的:**
- 所有规则 ID 必须遵循统一的命名规范（SCREAMING_SNAKE_CASE）
- lint 工具或自动化脚本可能依赖规则 ID 格式
- 教训：建立规则 ID 命名规范并在上游检测

#### AST Pattern 声明与实现脱节
**场景:** YAML 声明 `type: ast`，但 parser 完全忽略该类型

**学到的:**
- 文档化的接口必须与实现一致，否则会产生误导
- golden file 测试只验证"规则是否触发"，不验证"规则逻辑是否正确"
- 教训：YAML 中的 ast pattern 声明应删除或实现

### 测试设计

#### 测试需要导入副作用
**场景:** Engine 测试需要导入 langs 包才能触发 parser 注册

**学到的:**
- Go 的 init() 机制意味着测试可能依赖于导入链
- 跨包测试需要显式导入所有依赖包
- 教训：测试代码中的 `_ import` 是必要的，但需要注释说明原因

#### Golden File 测试框架
**场景:** `getRulesForLanguage` 是死代码 stub

**学到的:**
- 基础设施代码必须有实际功能或完全不存在，不允许"占位"实现
- stub 代码会误导其他开发者，以为功能已实现
- 教训：使用 `t.Skipf` 让测试优雅跳过，而非返回 nil

---

## 下次改进项

1. **Parser 单元测试** — 各语言 parser 应有独立测试
2. **output 包测试** — SARIF/JSON/Text 格式化应被测试
3. **AST YAML 引擎** — Phase 2 应实现完整的 YAML-driven AST pattern
4. **go.mod 依赖整理** — go-cmp 应移入 test require 块
