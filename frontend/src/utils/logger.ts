// 前端日志工具
interface LogLevel {
  DEBUG: 0;
  INFO: 1;
  WARN: 2;
  ERROR: 3;
}

const LOG_LEVELS: LogLevel = {
  DEBUG: 0,
  INFO: 1,
  WARN: 2,
  ERROR: 3,
};

class Logger {
  private currentLevel: number;
  private isDevelopment: boolean;

  constructor() {
    this.isDevelopment = process.env.NODE_ENV === 'development';
    this.currentLevel = this.isDevelopment ? LOG_LEVELS.DEBUG : LOG_LEVELS.INFO;
  }

  private formatMessage(level: string, message: string, data?: any): string {
    const timestamp = new Date().toISOString();
    const logData = data ? ` | Data: ${JSON.stringify(data)}` : '';
    return `[${timestamp}] [${level}] ${message}${logData}`;
  }

  private shouldLog(level: number): boolean {
    return level >= this.currentLevel;
  }

  debug(message: string, data?: any) {
    if (this.shouldLog(LOG_LEVELS.DEBUG)) {
      console.debug(this.formatMessage('DEBUG', message, data));
    }
  }

  info(message: string, data?: any) {
    if (this.shouldLog(LOG_LEVELS.INFO)) {
      console.info(this.formatMessage('INFO', message, data));
    }
  }

  warn(message: string, data?: any) {
    if (this.shouldLog(LOG_LEVELS.WARN)) {
      console.warn(this.formatMessage('WARN', message, data));
    }
  }

  error(message: string, error?: Error | any, data?: any) {
    if (this.shouldLog(LOG_LEVELS.ERROR)) {
      const errorInfo = error instanceof Error 
        ? { message: error.message, stack: error.stack }
        : error;
      
      console.error(this.formatMessage('ERROR', message, { error: errorInfo, ...data }));
      
      // 在生产环境中，可以将错误发送到日志服务
      if (!this.isDevelopment && error instanceof Error) {
        this.sendToLogService(message, error, data);
      }
    }
  }

  // API调用日志
  logApiCall(method: string, url: string, data?: any) {
    this.debug(`API调用: ${method} ${url}`, data);
  }

  // API响应日志
  logApiResponse(method: string, url: string, status: number, data?: any) {
    if (status >= 400) {
      this.error(`API错误: ${method} ${url} - Status: ${status}`, null, data);
    } else {
      this.debug(`API成功: ${method} ${url} - Status: ${status}`, data);
    }
  }

  // 用户操作日志
  logUserAction(action: string, data?: any) {
    this.info(`用户操作: ${action}`, data);
  }

  // 发送到日志服务（生产环境）
  private sendToLogService(message: string, error: Error, data?: any) {
    // 这里可以集成第三方日志服务，如 Sentry, LogRocket 等
    // 暂时只是控制台输出
    console.error('发送到日志服务:', { message, error: error.message, data });
  }

  // 设置日志级别
  setLevel(level: keyof LogLevel) {
    this.currentLevel = LOG_LEVELS[level];
  }
}

// 创建全局日志实例
const logger = new Logger();

// 全局错误处理
window.addEventListener('error', (event) => {
  logger.error('全局JavaScript错误', event.error, {
    filename: event.filename,
    lineno: event.lineno,
    colno: event.colno,
  });
});

window.addEventListener('unhandledrejection', (event) => {
  logger.error('未处理的Promise拒绝', event.reason);
});

export default logger;
export { LOG_LEVELS };
