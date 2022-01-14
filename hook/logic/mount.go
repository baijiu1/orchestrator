package logic

import (
    "fmt"
    "golang.org/x/crypto/ssh"
    "github.com/mitchellh/go-homedir"
    "vfailover/myflag"
    "time"
    "io/ioutil"
    "log"
)

func MountVip(conf *myflag.OrchCfg, logger *log.Logger) {
    logger.Printf("Begin vip mount on: %v \n", conf.NewMaster)
    if conf.ClusterType == "kp" || conf.ClusterType == "keepalived" || conf.ClusterType == "orch" {
        logger.Printf("cluster type: %v \n", conf.ClusterType)
        sshconn := fmt.Sprintf("%v:22", o.NewMaster)
        client, err := ssh,Dial("tcp", sshconn, &ssh.ClientConfig{
            Timeout: time.Seconds,
            User: o.SSHUser,
            Auth: []ssh.AuthMethod{publicKeyAuthFunc(o.SSHPublicKeys)},
            HostKeyCallback: ssh.InsecureIgnoreHostKey(),
        })
        if err != nil {
            return
        }
        session, err1 := client.NewSession()
        if err1 != nil {
            return
        }
        defer session.Close()
        unionCmd := fmt.Sprintf("%v;%v", conf.CmdVipAdd)
        mhaKey, err2 := session.CombinedOutput(o.MHAKeyCmd)
        if err2 != nil {
            return
        }
        MHAKEY := strings.Replace(string(if), "\n", "", -1)
        o.MHAKey = string(MHAKEY)
    }
}