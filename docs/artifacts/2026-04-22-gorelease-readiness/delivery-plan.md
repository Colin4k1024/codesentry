# Delivery Plan - Go Release Readiness

**Slug**: gorelease-readiness
**日期**: 2026-04-22
**状态**: draft → plan
**主责**: tech-lead

---

## 1. 需求挑战会结论

### 核心假设验证

| # | 原假设 | 质疑结论 | 决策 |
|---|--------|----------|------|
| 1 | `go install github.com/Colin4k1024/codesentry/cmd/goreview@latest` 可正常工作 | **阻断条件** - 模块根目录无 main 包，安装会失败 | 必须在发布前修复 |
| 2 | 覆盖率目标 (abcoder 70%, engine 60%) | 合理但需聚焦 CLI 集成测试 | 聚焦 CLI 和核心路径 |
| 3 | goreleaser 配置可本周期完成 | 无 goreleaser 配置，CGO 依赖增加复杂度 | 先完成基本配置，多平台后续 |
| 4 | v0.3.0 → v1.0.0 版本跳跃 | 语义版本暗示 API 稳定性 | 明确声明 v1.0.0 为稳定版 |

### 阻断条件

**ABORT if**: `go install github.com/Colin4k1024/codesentry/cmd/goreview@latest` 无法在发布 tag 前修复。

### 替代路径

1. **安装路径修复**: 重构使 `cmd/goreview` 可独立安装
2. **平台范围缩减**: v1.0.0 仅发布 macOS (Darwin)，Linux/Windows 延后
3. **覆盖率接受**: 接受当前覆盖率，明确 v1.0.0 为 beta

---

## 2. Story Slices

### Slice 1: 修复 go install 安装路径 (P0 - 阻断)

**Owner**: backend-engineer
**目标**: 修复 `go install github.com/Colin4k1024/codesentry/cmd/goreview@latest` 使其可正常工作
**验收**:
- `go install github.com/Colin4k1024/codesentry/cmd/goreview@latest` 成功
- 安装后 `codesentry version` 正常输出
**依赖**: 无
**风险**: 需要在 `cmd/goreview/` 添加 `main.go` 入口

### Slice 2: 测试覆盖率提升 (P0)

**Owner**: qa-engineer
**目标**: 提升核心包测试覆盖率
**验收**:
- `cmd/goreview` 覆盖率 ≥ 40% (CLI 集成测试)
- `internal/engine` 覆盖率 ≥ 55%
- `internal/abcoder` 覆盖率 ≥ 70%
**依赖**: Slice 1
**风险**: 测试桩可能与真实检测能力脱节

### Slice 3: goreleaser 配置 (P1)

**Owner**: backend-engineer
**目标**: 配置 goreleaser 支持多平台构建
**验收**:
- `.goreleaser.yaml` 存在并配置正确
- `goreleaser check` 通过
- 支持 macOS (darwin/amd64, darwin/arm64)
**依赖**: Slice 1 完成
**风险**: CGO 依赖 (tree-sitter) 增加跨平台编译复杂度

### Slice 4: v1.0.0 版本发布 (P0)

**Owner**: tech-lead
**目标**: 创建 v1.0.0 tag 并发布
**验收**:
- `git tag v1.0.0` 已创建
- GitHub Release 已发布
- `go install` 验证通过
**依赖**: Slice 1, Slice 3
**风险**: GitHub token 权限

### Slice 5: README 和文档更新 (P1)

**Owner**: tech-lead
**目标**: 更新 README 包含正确的 go install 说明
**验收**:
- README 包含正确的 `go install` 命令
- 版本号更新为 v1.0.0
**依赖**: Slice 4

---

## 3. 执行计划

| 阶段 | Slice | 主责角色 | 依赖 | 输出 |
|------|-------|----------|------|------|
| S1 | go install 路径修复 | backend-engineer | 无 | cmd/goreview/main.go |
| S2 | 测试覆盖率提升 | qa-engineer | S1 | 测试覆盖率报告 |
| S3 | goreleaser 配置 | backend-engineer | S1 | .goreleaser.yaml |
| S4 | v1.0.0 发布 | tech-lead | S1, S3 | v1.0.0 tag |
| S5 | README 更新 | tech-lead | S4 | README.md 更新 |

**并行性**: S2 和 S3 可并行 (都依赖 S1)

---

## 4. 风险与依赖清单

| 风险 | 影响 | 缓解 |
|------|------|------|
| go install 路径失败 | 发布无法满足核心成功标准 | Slice 1 必须先完成 |
| CGO tree-sitter 跨平台编译 | goreleaser 配置复杂 | 先支持 macOS，多平台后续 |
| 测试覆盖率提升工作量大 | 时间延长 | 聚焦 CLI 核心路径 |
| GitHub token 权限不足 | 无法自动发布 | 手动发布作为备选 |

---

## 5. 关键技术决策

1. **go install 路径**: 在 `cmd/goreview/` 添加标准 `main.go` 入口
2. **goreleaser 范围**: v1.0.0 仅 macOS (darwin/amd64 + darwin/arm64)
3. **覆盖率目标**: cmd/goreview ≥ 40%, engine ≥ 55%, abcoder ≥ 70%
4. **版本策略**: v1.0.0 标记为第一个稳定版本

---

## 6. 非目标 (Out of Scope)

- 不发布 Linux/Windows 版本 (延后到 v1.1)
- 不改变核心扫描逻辑
- 不新增语言支持
- 不修改输出格式

---

## 7. 验收标准

- [ ] `go install github.com/Colin4k1024/codesentry/cmd/goreview@latest` 成功
- [ ] `codesentry version` 正常输出
- [ ] `codesentry scan ./... --security` 核心功能可用
- [ ] CLI 测试覆盖率 ≥ 40%
- [ ] `.goreleaser.yaml` 配置完成
- [ ] v1.0.0 tag 已创建
- [ ] GitHub Release 已发布

---

## 8. 下一步

进入 /team-execute 执行 Slice 1-5
