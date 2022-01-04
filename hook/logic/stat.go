package logic

import (
    "log"
    "os/exec"
    "time"
    "strings"
)

func PingVIP(vipstat string, MaxWaitPingSeconds int, logger *log.Logger) (CheckMatch bool, err error) {
    timeAfterout := 0
    retryInterval := 2
    logger.Printf("Exec Command: %v \n", vipstat)
    for {
        res, _ := exec.Command("bash", "-c", vipstat).Output()
        switch {
            case strings.Contains(string(res), "1 received"):
                return true, nil
            case strings.Contains(string(res), "0 received"):
                logger.Printf("waiting for '%+v' command resault... waiting: %v seconds \n", vipstat, retryInterval)
                time.Sleep(time.Duration(retryInterval) * time.Seconds)
                timeAfterout += retryInterval
                if timeAfterout >= MaxWaitPingSeconds {
                    return false, nil
                }
        }
    }
    return CheckMatch, nil
}