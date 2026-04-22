# Arch Design - CodeSentry Parser 抽象重构

**Slug:** `codesentry-rule-enrichment`
**创建日期:** 2026-04-20
**主责角色:** architect

---

## 系统边界

### 外部依赖
- `gopkg.in/yaml.v3` - YAML 规则解析
- `github.com/spf13/cobra` - CLI 框架
- `github.com/google/go-cmp` - 测试结构比较（新增）

### 内部模块边界

```
cmd/
├── codesentry/           # CLI 入口
│   ├── root.go          # 根命令
│   ├── scan.go          # 扫描命令
│   └── langs.go         # 语言注册
│
internal/
├── engine/              # 扫描引擎（调用 parser）
│   └── engine.go        # Scan() 方法
├── parser/             # Parser 注册中心
│   └── registry.go     # Parser 接口 + 注册
├── rules/              # 规则加载
│   ├── loader.go       # LoadRules()
│   └── types.go        # Rule 结构定义
│
langs/                  # 语言解析器
├── golang/             # 特殊：含 AST 检查
├── python/             # BaseRegexParser
├── typescript/         # BaseRegexParser
└── ...                 # 共 11 种语言
```

---

## 组件拆分

### 当前问题：10 个重复 Parser

| Parser | 行数 | 模式 |
|--------|------|------|
| cpp | ~60 | pure regex loop |
| java | ~60 | pure regex loop |
| javascript | ~60 | pure regex loop |
| kotlin | ~60 | pure regex loop |
| php | ~60 | pure regex loop |
| python | ~60 | pure regex loop |
| ruby | ~60 | pure regex loop |
| rust | ~60 | pure regex loop |
| swift | ~60 | pure regex loop |
| typescript | ~60 | pure regex loop |
| **golang** | ~200 | regex + hardcoded AST |

**重复代码模式（10 个 parser 相同）：**
```go
func (p *Parser) Parse(filePath string, content []byte, langRules []rules.Rule) ([]parserpkg.Finding, error) {
    var findings []parserpkg.Finding
    text := string(content)
    lines := strings.Split(text, "\n")

    for _, rule := range langRules {
        for _, pattern := range rule.Patterns {
            if pattern.Type != "regex" {
                continue
            }
            re, err := regexp.Compile(pattern.Pattern)
            if err != nil {
                continue
            }
            for lineNum, line := range lines {
                if re.MatchString(line) {
                    findings = append(findings, parserpkg.Finding{...})
                }
            }
        }
    }
    return findings, nil
}
```

---

## 目标架构：BaseRegexParser

### 核心接口

```go
// Parser 接口（保持不变）
type Parser interface {
    Language() string
    Extensions() []string
    Parse(filePath string, content []byte, langRules []rules.Rule) ([]Finding, error)
}

// BaseRegexParser：10 个语言的通用实现
type BaseRegexParser struct{}

func (p *BaseRegexParser) Parse(filePath string, content []byte, langRules []rules.Rule) ([]Finding, error) {
    // 通用 regex 逻辑
}

// 各语言只需定义：
func (p *PythonParser) Language() string { return "python" }
func (p *PythonParser) Extensions() []string { return []string{".py", ".pyw", ".pyi"} }
```

### GoParser 特殊处理

GoParser 保留 AST 逻辑（goroutine leak, context leak, resource leak），不继承 BaseRegexParser：

```go
type GoParser struct {
    BaseRegexParser  // 嵌入，复用 regex 逻辑
    // Go 特定 AST 检查
}

func (p *GoParser) Parse(filePath string, content []byte, langRules []rules.Rule) ([]Finding, error) {
    // 1. 调用 BaseRegexParser.Parse() 处理 regex 规则
    // 2. 调用 p.checkAST() 处理 AST 规则
    // 3. 合并结果
}
```

---

## 关键数据流

### Scan 流程

```
engine.Scan(paths, cfg)
  │
  ├─► filepath.Walk() 遍历文件
  │     跳过 node_modules/.git/vendor
  │
  ├─► parser.DetectFromPath(filePath)
  │     └─► parser.Registry[lang]
  │
  ├─► parser.Parse(filePath, content, langRules)
  │     ├─► regex 规则 → BaseRegexParser.Parse()
  │     └─► AST 规则 → GoParser.checkAST()
  │
  ├─► Finding dedup（fmt.Sprintf("%s:%d:%s", file, line, ruleID)）
  │
  └─► 按 category 过滤（cfg.Security / cfg.Performance）
```

### 测试数据流

```
go test ./...
  │
  ├─► *_test.go 文件
  │     ├─► engine_test.go       → 测试 engine.Scan()
  │     ├─► loader_test.go        → 测试 LoadRules()
  │     ├─► registry_test.go      → 测试 Parser 注册
  │     └─► golden file 对比
  │
  └─► testdata/
        ├─► hardcoded_secret.input.py
        ├─► hardcoded_secret.golden.json
        └─► ...
```

---

## 技术选型

### 测试框架
| 选项 | 决策 | 理由 |
|------|------|------|
| 标准 Go testing | ✅ 采用 | 够用，不引入重型依赖 |
| testify | ❌ 不采用 | 增加迁移成本 |
| go-cmp | ✅ 新增 | 轻量，解决结构比较痛点 |
| golden | ❌ 不采用 | 手动实现 golden file 对比更灵活 |

### Golden File 格式
```json
// testdata/<rule_id>.golden.json
{
  "rule_id": "HARDCODED_SECRET",
  "findings": [
    {
      "line": 10,
      "column": 1,
      "message": "Possible hardcoded secret"
    }
  ]
}
```

---

## 风险与约束

### 风险
| 风险 | 影响 | 缓解 |
|------|------|------|
| Parser 重构破坏现有功能 | 高 | Slice 1 后完整集成测试 |
| AST 规则与 BaseRegex 耦合 | 中 | GoParser 显式组合，不继承 |
| Golden file 格式变更 | 低 | 版本化格式，变更需评审 |

### 约束
- **向后兼容：** `./codesentry scan ./...` 行为不变
- **零停机：** 重构不影响现有 CLI 接口
- **测试隔离：** 测试不依赖真实扫描结果，只测单元逻辑

---

## 未来扩展

### Phase 2: YAML-driven AST Engine
**触发条件：** AST 规则数量 ≥ 5 条时

**设计：**
```go
type ASTPattern struct {
    NodeType string  // "CallExpr", "SelectorExpr", etc.
    Field    string
    Value    string
}
```

**注意：** 当前 `rules/go/goroutine_leak.yaml` 的 `type: ast` 声明是**无效的**（parser 忽略），Phase 2 需要实现真正的 AST query 解析器。

---

## 验证标准

- [ ] 10 个 regex parser 代码减少 ≥ 60%（从 ~600 行 → ~240 行）
- [ ] `go build ./...` 成功
- [ ] `./codesentry scan ./...` 功能正常
- [ ] 新增语言只需实现 3 个方法（Language, Extensions, 可选 ParseRegex override）
