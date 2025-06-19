# 快速开始指南

## 第一步：安装 Go

如果你还没有安装 Go，请先安装：

1. 访问 [https://golang.org/dl/](https://golang.org/dl/)
2. 下载适合你系统的 Go 安装包
3. 运行安装程序
4. 重启命令行窗口

## 第二步：编译程序

在项目目录中运行：

```bash
# Windows
install.bat

# 或者手动执行
go mod tidy
go build -o cf-dns-manager.exe
```

## 第三步：配置 API Token

1. 复制配置文件：
   ```bash
   copy config.yaml.example config.yaml
   ```

2. 编辑 `config.yaml`，填入你的 Cloudflare API Token：
   ```yaml
   cloudflare_token: "你的API_Token"
   ```

### 获取 Cloudflare API Token

1. 登录 [Cloudflare Dashboard](https://dash.cloudflare.com/)
2. 点击右上角头像 → "My Profile"
3. 左侧菜单选择 "API Tokens"
4. 点击 "Create Token"
5. 选择 "Custom token"
6. 设置权限：
   - Zone:Zone:Read
   - Zone:DNS:Edit
7. 设置 Zone Resources 为 "All zones"
8. 点击 "Continue to summary" → "Create Token"
9. 复制生成的 Token

## 第四步：开始使用

### 查看所有 DNS 记录
```bash
cf-dns-manager.exe list
```

### 启动交互式界面
```bash
cf-dns-manager.exe interactive
```

### 筛选特定域名的记录
```bash
cf-dns-manager.exe list --filter-zone "example.com"
```

### 添加 DNS 记录
```bash
cf-dns-manager.exe add example.com www A 192.168.1.1 --proxied
```

## 常用命令示例

```bash
# 查看所有 A 记录
cf-dns-manager.exe list --filter-type "A"

# 查看特定域名的所有记录并按名称排序
cf-dns-manager.exe list --filter-zone "example.com" --sort-by "name"

# 添加网站记录
cf-dns-manager.exe add example.com @ A 192.168.1.100 --proxied
cf-dns-manager.exe add example.com www A 192.168.1.100 --proxied

# 添加邮件记录
cf-dns-manager.exe add example.com @ MX "mail.example.com" --ttl 300
```

## 交互式界面操作

启动交互式界面后：
- `↑↓` 或 `j/k`: 选择记录
- `a`: 添加记录
- `e`: 编辑记录  
- `d`: 删除记录
- `s`: 切换排序方向
- `r`: 刷新记录
- `q`: 退出

## 故障排除

### 常见问题

1. **"Cloudflare API Token 未设置"**
   - 检查 `config.yaml` 文件是否存在
   - 确认 Token 已正确填入

2. **"获取域名列表失败"**
   - 检查 API Token 是否正确
   - 确认 Token 有足够权限

3. **"域名不存在"**
   - 确认域名已在 Cloudflare 中正确配置
   - 检查域名拼写

### 获取帮助

```bash
# 查看所有命令
cf-dns-manager.exe --help

# 查看特定命令帮助
cf-dns-manager.exe list --help
cf-dns-manager.exe add --help
``` 