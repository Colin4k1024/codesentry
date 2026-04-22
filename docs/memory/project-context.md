# Project Context - CodeSentry

**最后更新**: 2026-04-22

## 项目概览

- **项目名**: CodeSentry
- **项目类型**: 开源静态代码分析与 AI 辅助审查工具
- **技术栈**: Go 1.23+, cloudwego/abcoder v0.3.1
- **许可证**: MIT

## 当前任务

**2026-04-22**: gorelease-readiness (执行中)

### 任务状态

| 阶段 | 状态 | 说明 |
|------|------|------|
| Slice 1: go install 路径修复 | 待开始 | P0 阻断项 |
| Slice 2: 测试覆盖率提升 | 待开始 | CLI ≥40%, engine ≥55% |
| Slice 3: goreleaser 配置 | 待开始 | macOS only |
| Slice 4: v1.0.0 发布 | 待开始 | 依赖 S1, S3 |
| Slice 5: README 更新 | 待开始 | 依赖 S4 |

### 已知风险

1. **go install 路径阻断** - 模块根目录无 main 包，安装会失败
2. **CGO 依赖** - tree-sitter 跨平台编译复杂

---

## 已完成任务

**2026-04-22**: codesentry-project-review (已关闭 ✅)

### 分组状态

| 分组 | 状态 | 备注 |
|------|------|------|
| Slice 1: Go AST 验证 | ✅ 完成 | GOROUTINE_LEAK, RESOURCE_LEAK, JWT_ERROR 已测试 |
| Slice 2: JS/TS 覆盖评估 | ✅ 完成 | JS 已有足够覆盖，无需独立 JS 规则 |
| Slice 3: 测试覆盖率提升 | ✅ 完成 | golang 80%, rules 81.2%, abcoder 65.2% |

## 覆盖率现状

| 包 | 覆盖率 | 备注 |
|-----|--------|------|
| langs/golang | 80% | ✅ 达标 |
| internal/rules | 81.2% | ✅ 达标 |
| langs/* (其他 9 个) | 100% | ✅ 全部达标 |
| internal/abcoder | 65.2% | 接受（后续迭代提升） |
| internal/engine | 48.8% | 接受（parser 层已覆盖） |
| cmd/goreview | 13.1% | 接受（CLI 测试优先级低） |

## 技术架构

```
cmd/goreview (CLI)
    └── internal/engine (扫描引擎)
            ├── internal/parser (解析器注册)
            ├── internal/rules (规则加载)
            ├── internal/abcoder (AI 上下文 - 仅 Go)
            └── langs/* (11 种语言解析器)
```

## 已知缺陷 (已关闭)

1. ~~**RESOURCE_LEAK 规则缺失**~~ - ✅ 已修复：创建 `rules/go/resource_leak.yaml`
2. ~~**JWT_ERROR 规则 ID 不一致**~~ - ✅ 已修复：修正 `langs/golang/parser.go` 中的规则引用

## 遗留项 (低优先级)

- `cmd/goreview` CLI 测试覆盖率提升
- `internal/engine` 覆盖率提升至 60%+
- `internal/abcoder` 覆盖率提升至 80%+

## 产出文档

- `docs/artifacts/2026-04-22-codesentry-project-review/prd.md`
- `docs/artifacts/2026-04-22-codesentry-project-review/delivery-plan.md`
- `docs/artifacts/2026-04-22-codesentry-project-review/execute-log.md`
- `docs/artifacts/2026-04-22-codesentry-project-review/slice2-js-coverage-analysis.md`
- `docs/artifacts/2026-04-22-codesentry-project-review/closeout-summary.md`
