package logic


import (
    "fmt"
    "golang.org/x/crypto/ssh"
    "github.com/mitchellh/go-homedir"
    "vfailover/myflag"
    "time"
    "log"
)

func MountVip(conf *myflag.OrchCfg, logger *log.Logger) {
    logger.Printf("Begin vip umount on: %v \n", conf.OldMaster)
    if conf.ClusterType == "kp" || conf.ClusterType == "keepalived" || conf.ClusterType == "orch" {
        logger.Printf("cluster type: %v \n", conf.ClusterType)
        sshconn := fmt.Sprintf("%v:22", o.OldMaster)
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
        _, err2 := session.CombinedOutput(conf.CmdVipDel)
        logger.Printf("exec command: %v \n", unionCmd)
        if err2 != nil {
            return
        }
    } else if conf.ClusterType == "mha" {
        logger.Printf("Begin vip umount on: %v \n", conf.OldMaster)
        if conf.ClusterType == "mha" {
            logger.Printf("cluster type: %v \n", conf.ClusterType)
            sshconn := fmt.Sprintf("%v:22", o.OldMaster)
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
            _, err2 := session.CombinedOutput(conf.UmountMHACmd)
            logger.Printf("exec command: %v \n", conf.UmountMHACmd)
            if err2 != nil {
                return
            }
        }
}