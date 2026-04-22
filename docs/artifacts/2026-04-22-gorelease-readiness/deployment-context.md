# Deployment Context - Go Release Readiness

**Slug**: gorelease-readiness
**日期**: 2026-04-22
**状态**: released
**主责**: devops-engineer

---

## 环境清单

### 生产环境

| 环境 | 用途 | 访问入口 |
|------|------|----------|
| GitHub Releases | 官方发布渠道 | https://github.com/Colin4k1024/codesentry/releases |
| Go Package Index | go install 索引 | https://pkg.go.dev/github.com/Colin4k1024/codesentry/cmd/goreview |

### 构建环境

| 环境 | 用途 | 配置 |
|------|------|------|
| GitHub Actions | 自动构建发布 | `.github/workflows/release.yml` |
| goreleaser | 多平台构建 | `.goreleaser.yaml` |

---

## 部署入口

### 主入口 (自动)

```
git push origin v1.0.0
```

推送 tag 触发 `.github/workflows/release.yml` → goreleaser release --clean → GitHub Release

### 手工入口

```bash
# 本地构建
goreleaser release --clean

# 本地快照构建 (不发布)
goreleaser build --snapshot --clean
```

### 回退入口

1. 删除 GitHub Release (GitHub UI)
2. 删除本地 tag: `git tag -d v1.0.0`
3. 推送删除: `git push origin :refs/tags/v1.0.0`
4. 创建新 tag 重新发布

---

## 配置与密钥

### 环境变量

| 变量 | 来源 | 说明 |
|------|------|------|
| GITHUB_TOKEN | GitHub Actions 自动提供 | 用于发布到 GitHub Releases |
| CGO_ENABLED=1 | goreleaser.yaml | tree-sitter 需要 CGO |

### 密钥来源

- 无需额外密钥配置
- GITHUB_TOKEN 由 GitHub Actions 自动注入

---

## 运行保障

### Feature Flag / 灰度

无 (直接全量发布)

### 监控与告警

| 监控项 | 来源 | 告警阈值 |
|--------|------|----------|
| GitHub Release 发布状态 | GitHub Actions | 失败即告警 |
| go install 下载量 | GitHub API | 无自动告警 |
| Issues 反馈 | GitHub Issues | 人工关注 |

### 值守安排

- 发布后 24 小时内由 tech-lead 观察 GitHub Actions 状态
- 如有异常，通过 GitHub Notifications 通知

### 观察窗口

- 发布后 24 小时：主要观察 GitHub Actions 是否成功、Release 是否正确创建
- 发布后 1 周：关注 go install 使用反馈和 GitHub Issues

---

## 恢复能力

### 回滚触发条件

- GitHub Release 创建失败
- 构建产物无法下载
- go install 无法正常工作

### 回滚路径

1. 删除已创建的 GitHub Release
2. 删除本地和远程 tag
3. 修复问题后重新 tag 并推送

### 验证方法

```bash
# 验证 go install
go install github.com/Colin4k1024/codesentry/cmd/goreview@v1.0.0
codesentry version

# 验证构建产物
gh release view v1.0.0
```

---

## 企业内控补充

| 字段 | 值 |
|------|-----|
| 应用等级 | T4 (开源工具，低风险) |
| 技术架构等级 | 简单 CLI 工具 |
| 资源隔离 | 无 (单进程 CLI) |
| 关键组件偏离 | 无 |
| 资产文档入口 | https://github.com/Colin4k1024/codesentry |
