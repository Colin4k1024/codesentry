# Release Plan - Go Release Readiness

**Slug**: gorelease-readiness
**日期**: 2026-04-22
**状态**: released
**主责**: tech-lead

---

## 发布信息

| 字段 | 值 |
|------|-----|
| 版本 | v1.0.0 |
| 发布类型 | 正式版 (Stable) |
| 发布日期 | 2026-04-22 |
| 发布渠道 | GitHub Releases + Go Package Index |
| 目标用户 | Go 开发者 |

---

## 变更与风险

### 主要变更

1. **go install 支持** - 用户可通过 `go install github.com/Colin4k1024/codesentry/cmd/goreview@latest` 安装
2. **版本号更新** - 0.3.0 → 1.0.0
3. **macOS 双架构支持** - darwin/amd64 + darwin/arm64

### 风险清单

| 风险 | 影响 | 缓解 |
|------|------|------|
| tree-sitter CGO 依赖 | Linux/Windows 构建复杂 | 仅发布 macOS |
| go install 路径问题 | 安装失败 | 已修复并验证 |
| GitHub Actions 权限 | 发布失败 | GITHUB_TOKEN 自动提供 |

### 非目标 (v1.0 不包含)

- Linux/Windows 二进制
- 容器镜像
- 其他包管理器

---

## 执行步骤

### 1. 发布前检查

```bash
# 确认 tag 存在
git tag -l "v1.0.0"

# 确认 tag 注解
git show v1.0.0 --quiet

# 确认构建成功
goreleaser build --snapshot --clean
```

### 2. 推送 Tag

```bash
# 推送到远程
git push origin v1.0.0
```

**Go-No-Go 判断点**: 推送前确认 tag 注解正确，构建成功

### 3. GitHub Actions 自动化

推送 tag 后自动触发:

1. `actions/checkout@v4` (fetch-depth: 0)
2. `actions/setup-go@v5` (go-version: 1.23)
3. `go test ./...`
4. `goreleaser/goreleaser-action@v6` (args: release --clean)
5. goreleaser 构建 → GitHub Release 创建 → 产物上传

### 4. 发布后验证

```bash
# 验证 GitHub Release
gh release view v1.0.0

# 验证 go install (等待 5-10 分钟让 Go 索引更新)
go install github.com/Colin4k1024/codesentry/cmd/goreview@v1.0.0
codesentry version
```

---

## 验证与监控

### 验证命令

| 验证项 | 命令 | 预期结果 |
|--------|------|----------|
| GitHub Release 存在 | `gh release view v1.0.0` | 显示 Release 信息 |
| 构建产物存在 | `gh release view v1.0.0 --json assets` | 包含 .tarball 文件 |
| go install 成功 | `go install ...@v1.0.0` | 无错误 |
| 版本正确 | `codesentry version` | 1.0.0 |

### 监控项

- GitHub Actions release job 状态
- GitHub Release 创建状态
- go install 下载量 (GitHub API)
- GitHub Issues 新建问题

---

## 回滚方案

### 触发条件

- Release 创建失败
- 构建产物损坏
- go install 无法工作

### 回滚步骤

```bash
# 1. 删除 GitHub Release (Web UI 或)
gh release delete v1.0.0 --yes

# 2. 删除远程 tag
git push origin :refs/tags/v1.0.0

# 3. 删除本地 tag
git tag -d v1.0.0

# 4. 修复问题后重新 tag 并推送
git tag -a v1.0.0 -m "CodeSentry v1.0.0 - First stable release"
git push origin v1.0.0
```

---

## 放行结论

| 条件 | 状态 | 说明 |
|------|------|------|
| 核心功能验证 | ✅ | CLI 可正常安装和使用 |
| 测试覆盖率达标 | ✅ | cmd/goreview 44.3% ≥ 40% |
| 构建验证通过 | ✅ | goreleaser snapshot 构建成功 |
| Launch Acceptance | ✅ | 所有验收标准满足 |
| 回滚方案明确 | ✅ | 可通过删除 Release 和 tag 回滚 |

**最终结论**: ✅ 同意发布 v1.0.0
