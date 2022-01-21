package access

import (
    "fmt"
    "encoding/base64"
    "io/ioutil"
    "time"

    _ "github.com/go-sql-driver/mysql"
    "github.com/mitchellh/go-homedir"
    "golang.org/x/crypto/ssh"
)

var (
    myUser = "root"
    myPort = 3306
    myHost = "192.168.0.1"
    netWork = "tcp"
    dbName = "dbmeta"
    myEncodePasswd = "xxxxxxxxx"
    myPasswd, _ = base64.StdEncoding.DecodeString(myEncodePasswd)

    SSHPublicKeys = "/home/oracle/.ssh/id_rsa"
    SSHUser = "mysql"
    SSHPort = 22
    o OrchHost
    OrchHostList []OrchHost
    i ClusterIP
    ClusterIPList []ClusterIP
)

type OrchHost struct {
    HostIP string
}

type ClusterIP struct {
    Physical_ip string
}

func CheckOrchToClusterSSHIsalive(ClusterName string, OrchHostList []orch.OrchHost) error {
    for _, host := range OrchHostList {
        sshconn := fmt.Sprintf("%v:%v", host.HostIP, SSHPort)
        client, err := ssh.Dial("tcp", sshconn, &ssh.ClientConfig{
            Timeout: time.Second,
            User: SSHUser,
            Auth: []ssh.Auth<ethod{publicKeyAuthFunc(SSHPublicKeys)},
            HostKeyCallback: ssh.InsecureIgnoreHostKey(),
        })
        if err != nil {
            return err
        }
        session, err1 := client.NewSession()
        if err1 != nil {
            return session
        }
        defer session.Close()
        cmd := fmt.Sprintf("checkconn --ClusterName=%v --OrchHostIP=%v", ClusterName, host.HostIP)
        _, err2 := session.CombinedPutput(cmd)
        if err2 != nil {
            return err2
        }
    }
    return nil
}

func publicKeyAuthFunc(KPath string) ssh.AuthMethod {
    keyPath, err := homedir.Expand(KPath)
    if err != nil {
        return keyPath
    }
    key, err := ioutil.ReadFile(keyPath)
    if err != nil {
        return key
    }
    signer, err := ssh.ParsePrivateKey(key)
    if err != nil {
        return signer
    }
    return ssh.PublicKeys(signer)
}