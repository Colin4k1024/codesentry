# Architecture Design - ABCoder Integration

**任务:** abcoder-integration
**日期:** 2026-04-22
**主责角色:** architect

---

## 1. 系统边界

### 1.1 外部依赖

| 依赖 | 用途 | 集成方式 |
|------|------|----------|
| `github.com/cloudwego/abcoder` | UniAST Parser | Go Library 直接引入 |

### 1.2 集成点

```
CodeSentry CLI
     │
     ▼
┌─────────────────────────────────────────────────────────┐
│  internal/engine/scan.go                                │
│  (检测漏洞，发现需要修复建议的节点)                       │
└─────────────────────────────────────────────────────────┘
     │
     ▼ 是否需要修复建议？
     │
     ├─ Yes ─────────────────────────────────────────────┐
     │                                                    │
     ▼                                                    ▼
┌──────────────────────┐    ┌──────────────────────────┐
│ internal/abcoder/    │───▶│ code-reviewer Skill      │
│ bridge.go            │    │ (修复代码生成)            │
│ (上下文获取)          │    └──────────────────────────┘
└──────────────────────┘             │
     │                                ▼
     │                     ┌──────────────────────────┐
     │                     │ security-reviewer Skill  │
     │                     │ (安全性验证)              │
     │                     └──────────────────────────┘
     │                                │
     └────────────────────────────────┘
              │
              ▼
     ┌─────────────────────┐
     │ 修复建议输出         │
     │ (before/after code) │
     └─────────────────────┘
     │
     ├─ No ──▶ YAML suggestion (回退)
```

### 1.3 边界内外划分

**在边界内:**
- abcoder Bridge 的 Go 代码解析能力
- 上下文传递协议
- Skill Agent 调用逻辑
- 回退触发机制

**在边界外:**
- abcoder 内部实现
- Skill Agent 内部生成逻辑
- 用户终端显示

---

## 2. 组件拆分

### 2.1 新增组件

#### `internal/abcoder/bridge.go`

**职责:** 封装 abcoder UniAST Parser，提供统一上下文接口

```go
package abcoder

import (
    "github.com/cloudwego/abcoder/lang"
    "github.com/cloudwego/abcoder/lang/uniast"
)

type Bridge struct {
    repo    *uniast.Repository
    repoPath string
}

type CodeContext struct {
    File        string
    Line       int
    Function    *uniast.Function
    Calls       []uniast.Dependency
    Variables   []uniast.Dependency
    Definitions map[string]*uniast.Var
}

func NewBridge(repoPath string) (*Bridge, error)
func (b *Bridge) Parse() error
func (b *Bridge) GetContext(file string, line int) (*CodeContext, error)
func (b *Bridge) GetFunctionContext(nodeID uniast.Identity) (*uniast.Function, error)
func (b *Bridge) GetVariableDefinition(varID uniast.Identity) (*uniast.Var, error)
```

### 2.2 修改组件

#### `internal/engine/scan.go`

**修改点:**
- 检测到漏洞时，调用 abcoder Bridge 获取上下文
- 将上下文传递给 Skill Agent

#### `cmd/codesentry/main.go`

**修改点:**
- 添加 `--with-context` flag，控制是否启用 abcoder 上下文

---

## 3. 关键数据流

### 3.1 检测 + 修复流程

```
1. 用户执行扫描
   $ codesentry scan ./src --security

2. Engine 检测漏洞
   Finding {
       Rule: "SQL_INJECTION",
       File: "db/query.go",
       Line: 42,
       Severity: "SEVERE"
   }

3. 检查 abcoder 可用性
   ├─ Go 文件 + abcoder 可用 ──▶ 获取上下文
   └─ 非 Go 文件 / abcoder 不可用 ──▶ 回退到 YAML suggestion

4. 获取上下文 (仅 Go 文件)
   CodeContext {
       Function: "execQuery",
       Calls: [db.Exec, strings.Join],
       Variables: [{Name: "query", Type: "string"}]
   }

5. 调用 Skill Agent
   Input: {
       vulnerability: "SQL_INJECTION",
       file: "db/query.go",
       line: 42,
       context: CodeContext
   }

6. Skill Agent 生成修复
   Output: {
       before: "query := \"SELECT * FROM users WHERE id=\" + id",
       after: "query := \"SELECT * FROM users WHERE id=?\", id",
       explanation: "使用参数化查询避免 SQL 注入"
   }

7. 输出结果
   [SEVERE] SQL Injection
     File: db/query.go:42
     Before: query := "SELECT * FROM users WHERE id=" + id
     After:  query := "SELECT * FROM users WHERE id=?", id
     Suggestion: 使用参数化查询
```

### 3.2 回退流程

```
检测到漏洞
     │
     ▼
abcoder 可用？ ── No ──▶ 使用 YAML suggestion
     │
    Yes
     │
     ▼
是 Go 文件？ ── No ──▶ 使用 YAML suggestion
     │
    Yes
     │
     ▼
解析上下文 ── 失败 ──▶ 使用 YAML suggestion
     │
   成功
     │
     ▼
调用 Skill Agent 生成修复
```

---

## 4. 接口约定

### 4.1 abcoder Bridge 接口

```go
// Bridge 初始化
func NewBridge(repoPath string) (*Bridge, error)

// 解析仓库（调用一次，结果可复用）
func (b *Bridge) Parse(ctx context.Context) error

// 获取指定位置的代码上下文
func (b *Bridge) GetContext(file string, line int) (*CodeContext, error)

// 获取函数完整信息
func (b *Bridge) GetFunction(nodeID string) (*uniast.Function, error)

// 获取变量定义
func (b *Bridge) GetVariable(nodeID string) (*uniast.Var, error)

// 获取调用链
func (b *Bridge) GetCallChain(nodeID string) ([]string, error)
```

### 4.2 上下文数据结构

```go
type CodeContext struct {
    // 位置信息
    File     string
    Line     int
    Column   int

    // 函数信息
    FunctionName string
    FunctionContent string
    FunctionSignature string

    // 调用关系
    FunctionCalls []CallInfo
    MethodCalls []CallInfo

    // 变量信息
    LocalVariables []VarInfo
    GlobalVariables []VarInfo

    // 依赖类型
    UsedTypes []TypeInfo
}

type CallInfo struct {
    Name    string
    PkgPath string
    File    string
    Line    int
}

type VarInfo struct {
    Name    string
    Type    string
    Content string
    Line    int
}
```

---

## 5. 技术选型

### 5.1 abcoder 版本锁定

```go
// go.mod
require github.com/cloudwego/abcoder v0.3.1
```

**原因:**
- v0.3.1 是当前稳定版本
- Go parser 可独立使用
- UniAST 结构稳定

### 5.2 Skill Agent 调用协议

**输入格式:**
```json
{
  "task": "generate_fix",
  "context": {
    "vulnerability": "SQL_INJECTION",
    "file": "db/query.go",
    "line": 42,
    "function": {
      "name": "execQuery",
      "content": "func execQuery(id string) { ... }",
      "calls": ["db.Exec", "strings.Join"]
    }
  }
}
```

**输出格式:**
```json
{
  "before": "original code",
  "after": "fixed code",
  "explanation": "why this is safe",
  "confidence": 0.95
}
```

---

## 6. 风险与约束

### 6.1 技术风险

| 风险 | 影响 | 缓解 |
|------|------|------|
| abcoder API 变更 | 需要更新 Bridge | 锁定版本，定期同步 |
| LSP 服务器缺失 | 非 Go 语言无法获取上下文 | 回退到 suggestion |
| 大型代码库解析慢 | 扫描时间增加 | 增量解析，缓存结果 |

### 6.2 约束

- **性能约束:** 单次扫描延迟增加 < 2x
- **内存约束:** abcoder 上下文不持久化，用完即释放
- **兼容性约束:** 不改变现有 YAML 规则格式

---

## 7. 目录结构

```
internal/
├── abcoder/
│   ├── bridge.go      # abcoder Bridge 实现
│   ├── context.go     # 上下文数据结构
│   └── bridge_test.go # 单元测试
├── engine/
│   ├── scan.go       # 修改：添加上下文获取
│   └── engine.go      # 保持不变
└── ...
```

---

*创建时间: 2026-04-22*
*主责角色: architect*
