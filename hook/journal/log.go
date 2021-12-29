package journal

import (
    "fmt"
    "log"
    "os"
    "time"
)

var (
    logger *log.Logger
    file *os.File
    err error
)

fun InitLog(ClusterName string) *log.Logger {
    timeStr := time.Now().Format("2006-01-02 15:04:05")
    filename := fmt.Sprintf("/home/mysql/orch/logs/%v_%v_switchover.log", ClusterName, timeStr)
    file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return
    }
    logger = log.New(file, "", log.LstdFlags)
    logger.SetPrefix(ClusterName + "- ")
    logger.SetFlags(log.LstdFlags | log.Lshortfile)
    return logger
}