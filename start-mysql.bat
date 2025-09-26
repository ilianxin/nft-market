@echo off
chcp 65001 >nul

echo 🚀 启动MySQL数据库服务...

:: 检查Docker是否已安装
docker --version >nul 2>&1
if errorlevel 1 (
    echo ❌ Docker 未安装，请先安装 Docker Desktop
    echo 或者手动安装MySQL 8.0并创建数据库 'nft_market'
    pause
    exit /b 1
)

:: 检查Docker Compose是否已安装
docker-compose --version >nul 2>&1
if errorlevel 1 (
    echo ❌ Docker Compose 未安装，请先安装 Docker Compose
    pause
    exit /b 1
)

:: 启动MySQL服务
echo 📦 启动MySQL容器...
docker-compose up -d mysql

:: 等待MySQL启动
echo ⏳ 等待MySQL服务启动...
timeout /t 10 >nul

:: 检查MySQL是否已启动
docker-compose ps mysql | findstr "Up" >nul
if errorlevel 1 (
    echo ❌ MySQL服务启动失败
    docker-compose logs mysql
    pause
    exit /b 1
) else (
    echo ✅ MySQL服务启动成功！
    echo.
    echo 📊 数据库连接信息：
    echo - 主机: localhost
    echo - 端口: 3306
    echo - 数据库: nft_market
    echo - 用户名: nft_user
    echo - 密码: nft_password
    echo - Root密码: rootpassword
    echo.
    echo 🌐 phpMyAdmin 管理界面: http://localhost:8081
    echo    用户名: root
    echo    密码: rootpassword
    echo.
    echo 💡 更新后端 .env 文件中的数据库连接字符串：
    echo DATABASE_URL=nft_user:nft_password@tcp^(localhost:3306^)/nft_market?charset=utf8mb4^&parseTime=True^&loc=Local
    echo.
    echo 🛑 停止服务: docker-compose down
    echo.
)

pause
