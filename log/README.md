# logging

一个可以显示文件行数，和颜色的log系统，基于golang

# Usage

```gotemplate
// import
import (
    "github.com/dashjay/logging"
)
// Usage

    logging.Info("info test")
    logging.Debug("debug test")
    logging.Error("error test")
    // 这个Fatal会触发fatal，程序会停止
    // logging.Fatal call log.Fatal，the main process will exit
    // logging.Fatal("fatal test") 
    logging.Warn("warn test")
    logging.Info("sso server start")
```

# 代码部分来源于

- <https://github.com/eddycjy/go-gin-example>
