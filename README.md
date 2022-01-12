# rzap

Log rotate for uber-zap

## How to use

Create and set to global.

``` go
NewGlobalLogger([]zapcore.Core{
    NewCore(&lumberjack.Logger{
        Filename:   "/your/log/path/app.log",
        MaxSize:    10,   // 10 megabytes, defaults to 100 megabytes
        MaxAge:     10,   // 10 days, default is not to remove old log files
        MaxBackups: 10,   // 10 files, default is to retain all old log files
        Compress:   true, // compress to gzio, default is not to perform compression
    }, zap.InfoLevel),
})

zap.L().Info("some message", zap.Int("status", 0))
```

## Log to files at different levels

``` go
NewGlobalLogger([]zapcore.Core{
    NewCore(&lumberjack.Logger{
        Filename: "/path/to/info.log",
    }, zap.LevelEnablerFunc(func(level zapcore.Level) bool {
        return level <= zap.InfoLevel
    })),
    NewCore(&lumberjack.Logger{
        Filename: "/path/to/error.log",
    }, zap.LevelEnablerFunc(func(level zapcore.Level) bool {
        return level > zap.InfoLevel
    })),
})

zap.L().Info("some info message", zap.Int("status", 0))   // only output to /path/to/info.log
zap.L().Error("some error message", zap.Int("status", 1)) // only output to /path/to/error.log
```