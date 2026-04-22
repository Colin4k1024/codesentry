# Project Context - CodeSentry

**最后更新**: 2026-04-22

## 项目概览

- **项目名**: CodeSentry
- **项目类型**: 开源静态代码分析与 AI 辅助审查工具
- **技术栈**: Go 1.23+, cloudwego/abcoder v0.3.1
- **许可证**: MIT

## 当前任务

**2026-04-22**: codesentry-project-review

### 分组状态

| 分组 | 状态 | 备注 |
|------|------|------|
| A: 测试覆盖率提升 | pending | 目标: 核心包 > 80% |
| B: JavaScript 规则补充 | blocked by A | 先评估 TS 规则覆盖 |
| C: Go AST 测试策略 | blocked | 需先验证 AST 实现正确性 |

### 关键依赖

- Go AST 检测逻辑 (`langs/golang/parser.go`) 从未测试
- `internal/engine` 是所有 parser 基础，需优先提升覆盖率

## 技术架构

```
cmd/goreview (CLI)
    └── internal/engine (扫描引擎)
            ├── internal/parser (解析器注册)
            ├── internal/rules (规则加载)
            ├── internal/abcoder (AI 上下文 - 仅 Go)
            └── langs/* (11 种语言解析器)
```

## 已知风险

1. **Go AST 实现未验证** - GOROUTINE_LEAK, CONTEXT_LEAK, RESOURCE_LEAK 逻辑无测试
2. **JavaScript 无专用规则** - 仅依赖跨语言通用规则
3. **测试覆盖率分布不均** - 核心包 45-47%, CLI 仅 13%

## 下一步行动

1. **Slice 1**: 验证 Go AST 实现正确性 (backend-engineer)
2. **Slice 2**: 评估 JS/TS 规则覆盖关系 (security-reviewer)
3. **Slice 3**: 提升核心包测试覆盖率 (qa-engineer)

## 产出文档

- `docs/artifacts/2026-04-22-codesentry-project-review/prd.md`
- `docs/artifacts/2026-04-22-codesentry-project-review/delivery-plan.md`
