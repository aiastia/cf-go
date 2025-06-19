# Cloudflare DNS 管理器

一个用于管理 Cloudflare DNS 记录的命令行工具，支持查询、添加、修改、删除 DNS 记录，并提供排序和筛选功能。

## 功能特性

- 🌐 **批量查询**: 一次性查询所有域名的 DNS 记录
- 📝 **记录管理**: 添加、修改、删除 DNS 记录
- 🔍 **智能筛选**: 按域名、记录名、类型、内容筛选
- 📊 **灵活排序**: 支持多种排序方式（域名、记录名、类型等）
- 🎨 **交互式界面**: 美观的命令行界面，支持键盘操作
- ⚙️ **配置管理**: 支持配置文件和环境变量

## 安装

1. 确保已安装 Go 1.21 或更高版本
2. 克隆或下载项目
3. 安装依赖：

```bash
go mod tidy
```

4. 编译程序：

```bash
go build -o cf-dns-manager
```

## 配置

### 方法一：配置文件

创建 `config.yaml` 文件：

```yaml
cloudflare_token: "your_cloudflare_api_token"
account_id: "your_account_id"  # 可选
```

### 方法二：环境变量

```bash
export CF_API_TOKEN="your_cloudflare_api_token"
export CF_ACCOUNT_ID="your_account_id"  # 可选
```

### 获取 Cloudflare API Token

1. 登录 [Cloudflare Dashboard](https://dash.cloudflare.com/)
2. 进入 "My Profile" > "API Tokens"
3. 创建新的 Token，需要以下权限：
   - Zone:Zone:Read
   - Zone:DNS:Edit

## 使用方法

### 命令行模式

#### 列出所有 DNS 记录

```bash
# 列出所有记录
./cf-dns-manager list

# 按域名筛选
./cf-dns-manager list --filter-zone "example.com"

# 按记录类型筛选
./cf-dns-manager list --filter-type "A"

# 按记录名筛选
./cf-dns-manager list --filter-name "www"

# 按内容筛选
./cf-dns-manager list --filter-content "192.168.1.1"

# 排序（支持：name, type, zone, content, ttl, created, modified）
./cf-dns-manager list --sort-by "name" --ascending

# 组合使用
./cf-dns-manager list --filter-zone "example.com" --filter-type "A" --sort-by "name"
```

#### 添加 DNS 记录

```bash
# 添加 A 记录
./cf-dns-manager add example.com www A 192.168.1.1 --ttl 300 --proxied

# 添加 CNAME 记录
./cf-dns-manager add example.com api CNAME api.example.com --ttl 300

# 添加 MX 记录
./cf-dns-manager add example.com @ MX "mail.example.com" --ttl 300
```

#### 更新 DNS 记录

```bash
./cf-dns-manager update [记录ID] example.com www A 192.168.1.2 --ttl 600 --proxied
```

#### 删除 DNS 记录

```bash
./cf-dns-manager delete [记录ID] example.com
```

### 交互式模式

启动交互式界面：

```bash
./cf-dns-manager interactive
```

交互式界面操作：
- `↑↓` 或 `j/k`: 选择记录
- `a`: 添加记录
- `e`: 编辑记录
- `d`: 删除记录
- `s`: 切换排序方向
- `r`: 刷新记录
- `q`: 退出

## 支持的 DNS 记录类型

- A (IPv4 地址)
- AAAA (IPv6 地址)
- CNAME (规范名称)
- MX (邮件交换)
- TXT (文本记录)
- SRV (服务记录)
- NS (名称服务器)
- PTR (指针记录)
- CAA (证书颁发机构授权)
- 等等...

## 示例

### 查看所有 A 记录并按域名排序

```bash
./cf-dns-manager list --filter-type "A" --sort-by "zone" --ascending
```

### 查看特定域名的所有记录

```bash
./cf-dns-manager list --filter-zone "example.com"
```

### 添加网站记录

```bash
# 添加主域名 A 记录
./cf-dns-manager add example.com @ A 192.168.1.100 --proxied

# 添加 www 子域名
./cf-dns-manager add example.com www A 192.168.1.100 --proxied

# 添加邮件服务器记录
./cf-dns-manager add example.com @ MX "mail.example.com" --ttl 300
```

## 注意事项

1. **API 限制**: Cloudflare API 有速率限制，请避免频繁操作
2. **权限**: 确保 API Token 有足够的权限
3. **TTL**: TTL=1 表示自动（由 Cloudflare 管理）
4. **代理**: 启用代理后，流量会经过 Cloudflare CDN

## 故障排除

### 常见错误

1. **API Token 无效**
   - 检查 Token 是否正确
   - 确认 Token 有足够权限

2. **域名不存在**
   - 确认域名已在 Cloudflare 中正确配置
   - 检查域名拼写

3. **记录已存在**
   - 删除旧记录后再添加新记录
   - 或使用更新命令

### 调试模式

设置环境变量查看详细日志：

```bash
export DEBUG=1
./cf-dns-manager list
```

## 贡献

欢迎提交 Issue 和 Pull Request！

## 许可证

MIT License 