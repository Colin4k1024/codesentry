# PRD - Go Release Readiness

**Slug**: gorelease-readiness
**日期**: 2026-04-22
**状态**: draft → intake
**主责**: tech-lead

---

## 1. 背景与目标

### 需求背景

CodeSentry 已完成开源发布准备（Slice 1-3 验证完成），但尚未准备好通过 `go install` 安装。需要完成测试覆盖、完善发布配置，以便用户可以通过官方 Go 工具链安装使用。

### 目标

1. **完成所有程序测试**: 提升测试覆盖率至合理水平
2. **Go 官方仓库发布**: 配置 `goreleaser` 或手动发布，使 `go install github.com/Colin4k1024/codesentry/cmd/goreview@latest` 可用
3. **用户可本地安装**: 用户执行 `go install` 后可直接使用 `codesentry` 命令

### 成功标准

- [ ] `go test ./... -cover` 核心包覆盖率达标 (abcoder >70%, engine >60%)
- [ ] `go install github.com/Colin4k1024/codesentry/cmd/goreview@latest` 可成功安装
- [ ] 安装后 `codesentry version` 可正常输出版本
- [ ] `codesentry scan` 核心功能可用
- [ ] README.md 包含正确的 `go install` 使用说明
- [ ] 发布 v1.0.0 版本 tag

---

## 2. 用户故事

### 作为开发者

**我想要**通过 `go install` 安装 CodeSentry
**以便**快速开始代码安全扫描，无需手动编译

**验收标准:**
- `go install github.com/Colin4k1024/codesentry/cmd/goreview@latest` 成功
- `codesentry --help` 显示帮助信息
- `codesentry version` 显示版本号
- `codesentry scan ./... --security` 可正常扫描

---

## 3. 范围

### In Scope

| 分组 | 描述 | 优先级 |
|------|------|--------|
| A | 测试覆盖率提升 (abcoder, engine, cmd/goreview) | P0 |
| B | goreleaser 配置或手动发布配置 | P0 |
| C | 版本 v1.0.0 tag 创建 | P0 |
| D | README go install 说明更新 | P1 |
| E | CI/CD 发布流程 (GitHub Actions) | P1 |

### Out of Scope

- 不新增语言支持
- 不修改核心扫描逻辑
- 不改变输出格式
- 不添加 UI

---

## 4. 风险与约束

| 风险 | 影响 | 缓解 |
|------|------|------|
| 覆盖率提升工作量大 | 时间延长 | 聚焦核心包，边缘包接受现状 |
| goreleaser 配置不熟悉 | 发布失败 | 参考现有开源项目配置 |
| GitHub Actions 权限不足 | 无法自动发布 | 手动发布作为备选 |

### 约束

- Go 模块路径: `github.com/Colin4k1024/codesentry`
- Go 版本: 1.23+
- 发布平台: GitHub Releases
- 安装方式: `go install`

---

## 5. 待确认项

| # | 问题 | 状态 |
|---|------|------|
| 1 | 覆盖率目标是否合理 (abcoder 70%, engine 60%)? | 待确认 |
| 2 | 使用 goreleaser 还是手动发布? | 建议 goreleaser |
| 3 | 版本号从 v0.3.0 直接跳到 v1.0.0? | 待确认 |
| 4 | GitHub token 是否可用于发布? | 待确认 |
| 5 | README 中的仓库名是 `codesentry_refactor` 是否正确? | 待确认 |

---

## 6. 参与角色

| 角色 | 输入缺口 |
|------|----------|
| tech-lead | 整体协调、发布决策 |
| backend-engineer | 覆盖率提升、goreleaser 配置 |
| qa-engineer | 测试验证、覆盖率评估 |

---

## 7. 下一步

进入 `/team-plan` 制定详细的执行计划
