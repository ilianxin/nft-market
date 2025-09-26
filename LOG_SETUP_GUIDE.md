# 📋 日志系统配置完成指南

## ✅ 已完成的日志配置

### 🔧 后端日志系统

1. **日志库**: 使用 `logrus` + `lumberjack` 实现结构化日志和文件轮转
2. **日志文件位置**: `backend/logs/app.log`
3. **日志配置**:
   - 格式: JSON格式
   - 输出: 同时输出到控制台和文件
   - 轮转: 单文件最大100MB，保留10个备份，保存30天
   - 压缩: 自动压缩旧日志文件

### 📁 日志文件位置

```
backend/
├── logs/
│   ├── app.log              # 主日志文件
│   ├── app.log.1.gz         # 轮转的压缩日志文件
│   └── app.log.2.gz         # 更旧的日志文件
└── ...
```

### 🔍 如何查看日志

#### 1. 实时查看日志
```bash
# Windows PowerShell
cd backend
Get-Content logs/app.log -Wait -Tail 50

# 或者使用 type 命令
type logs\app.log
```

#### 2. 查看特定错误
```bash
# 查找创建订单相关的错误
Select-String -Path "logs/app.log" -Pattern "创建订单失败"

# 查找所有ERROR级别的日志
Select-String -Path "logs/app.log" -Pattern "ERROR"
```

### 📊 日志内容示例

**成功创建订单的日志**:
```json
{
  "level": "info",
  "msg": "订单创建成功",
  "order_id": 123,
  "user_address": "0x1234...",
  "collection_address": "0x5678...",
  "token_id": "1",
  "order_type": 1,
  "price": 0.1,
  "time": "2023-12-26 15:04:05"
}
```

**创建订单失败的日志**:
```json
{
  "level": "error",
  "msg": "创建订单失败",
  "error": "database connection failed",
  "user_address": "0x1234...",
  "collection_address": "0x5678...",
  "token_id": "1",
  "order_type": 1,
  "price": 0.1,
  "time": "2023-12-26 15:04:05"
}
```

### 🌐 前端日志系统

1. **日志工具**: `frontend/src/utils/logger.ts`
2. **功能**:
   - API调用日志
   - 用户操作日志
   - 错误日志
   - 全局错误捕获

#### 前端日志查看位置
- **开发环境**: 浏览器开发者工具 Console
- **生产环境**: 可集成第三方日志服务（如 Sentry）

## 🚀 启动应用并查看日志

### 1. 启动后端服务
```bash
cd backend
go run main.go
```

启动后会看到类似日志：
```
{"level":"info","msg":"日志系统初始化完成","time":"2023-12-26 15:04:05"}
{"level":"info","msg":"配置加载完成","port":"8080","env":"development","time":"2023-12-26 15:04:05"}
{"level":"info","msg":"数据库连接成功","time":"2023-12-26 15:04:05"}
{"level":"info","msg":"服务器启动","port":"8080","mode":"debug","time":"2023-12-26 15:04:05"}
```

### 2. 测试创建订单并查看日志

当创建订单失败时，可以在以下位置查看详细错误信息：

1. **控制台输出**: 实时显示
2. **日志文件**: `backend/logs/app.log`
3. **浏览器控制台**: 前端错误信息

## 🔧 日志配置自定义

### 修改日志配置
在 `backend/main.go` 中可以自定义日志配置：

```go
logConfig := &logger.Config{
    Level:      "debug",        // 日志级别: debug, info, warn, error
    Format:     "json",         // 格式: json, text
    Output:     "both",         // 输出: file, console, both
    FilePath:   "logs/app.log", // 日志文件路径
    MaxSize:    100,            // 单文件最大尺寸(MB)
    MaxBackups: 10,             // 保留旧文件数量
    MaxAge:     30,             // 保留天数
    Compress:   true,           // 是否压缩
}
```

## 🔍 常见错误定位方法

### 1. 创建订单失败
查看日志中的错误信息：
```bash
Select-String -Path "backend/logs/app.log" -Pattern "创建订单失败" -Context 2
```

### 2. 数据库连接问题
```bash
Select-String -Path "backend/logs/app.log" -Pattern "数据库" -Context 2
```

### 3. API调用错误
前端控制台会显示详细的API调用和响应日志。

### 4. 区块链交互错误
```bash
Select-String -Path "backend/logs/app.log" -Pattern "区块链" -Context 2
```

## 📝 日志最佳实践

1. **开发阶段**: 使用 `debug` 级别查看详细信息
2. **生产环境**: 使用 `info` 或 `warn` 级别减少日志量
3. **定期清理**: 日志会自动轮转，无需手动清理
4. **错误追踪**: 每个错误都包含足够的上下文信息用于定位问题

现在您可以启动应用并查看详细的日志信息了！
