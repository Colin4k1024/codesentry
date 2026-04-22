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
| Slice 1: Go AST 验证 | ✅ 完成 | GOROUTINE_LEAK 已测试，RESOURCE_LEAK 规则缺失 |
| Slice 2: JS/TS 覆盖评估 | ✅ 完成 | JS 已有足够覆盖，无需独立 JS 规则 |
| Slice 3: 测试覆盖率提升 | ✅ 完成 | golang 68%, abcoder 65.2%, engine 48.8% |

## 覆盖率现状

| 包 | 覆盖率 | 备注 |
|-----|--------|------|
| langs/golang | 68% | 新增测试 |
| internal/abcoder | 65.2% | 新增测试 |
| internal/engine | 48.8% | 新增边界测试 |
| cmd/goreview | 13.1% | CLI 测试有限 |
| langs/* (其他) | 100% | 9 个解析器 |

## 已知缺陷

1. **RESOURCE_LEAK 规则缺失** - `rules/go/resource_leak.yaml` 不存在
2. **CONTEXT_LEAK 实现不匹配** - 代码检查 JWT_ERROR，但规则名是 CONTEXT_LEAK
3. **cmd/goreview 覆盖率低** - CLI 测试需要 cobra 命令环境

## 技术架构

```
cmd/goreview (CLI)
    └── internal/engine (扫描引擎)
            ├── internal/parser (解析器注册)
            ├── internal/rules (规则加载)
            ├── internal/abcoder (AI 上下文 - 仅 Go)
            └── langs/* (11 种语言解析器)
```

## 产出文档

- `docs/artifacts/2026-04-22-codesentry-project-review/prd.md`
- `docs/artifacts/2026-04-22-codesentry-project-review/delivery-plan.md`
- `docs/artifacts/2026-04-22-codesentry-project-review/execute-log.md`
- `docs/artifacts/2026-04-22-codesentry-project-review/slice2-js-coverage-analysis.md`
