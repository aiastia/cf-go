@echo off
echo ========================================
echo Cloudflare DNS 管理器安装脚本
echo ========================================
echo.

echo 检查 Go 是否已安装...
go version >nul 2>&1
if %errorlevel% neq 0 (
    echo Go 未安装，正在下载...
    echo 请访问 https://golang.org/dl/ 下载并安装 Go
    echo 安装完成后重新运行此脚本
    pause
    exit /b 1
)

echo Go 已安装，版本信息：
go version
echo.

echo 安装依赖...
go mod tidy
if %errorlevel% neq 0 (
    echo 安装依赖失败
    pause
    exit /b 1
)

echo 编译程序...
go build -o cf-dns-manager.exe
if %errorlevel% neq 0 (
    echo 编译失败
    pause
    exit /b 1
)

echo.
echo ========================================
echo 安装完成！
echo ========================================
echo.
echo 下一步：
echo 1. 复制 config.yaml.example 为 config.yaml
echo 2. 编辑 config.yaml，填入你的 Cloudflare API Token
echo 3. 运行 cf-dns-manager.exe 开始使用
echo.
echo 使用示例：
echo   cf-dns-manager.exe list
echo   cf-dns-manager.exe interactive
echo.
pause 