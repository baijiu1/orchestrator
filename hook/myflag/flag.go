package myflag

import (
    "database/sql"
    "encoding/base64"
    "flag"
    "fmt"
    "io/ioutil"
    "os"
    "strings"
    "time"

    "github.com/mitchellh/go-homedir"
    "golang.org/x/crypto/ssh"
    _ "github.com/go-sql-driver/mysql"
)

var (
    o OrchCfg
    m MetaDBInfo
    c ClusterInfo
)

type OrchCfg struct {
    ClusterName string
    DeadStatus string
    OldMaster string
    OldMasterPort string
    NewMaster string
    NewMasterPort string
    CmdVipAdd string
    CmdVipDel string
    CmdVipstat string
    GateWayCmd string
    GateWay string
    MHAKey string
    MHAKeyCmd string
    Mask string
    ArPingCmdb string
    InterFace string
    SSHUser string
    SSHPublicKeys string
    VipAddr string
    ClusterType string
    MountMHACmd string
    UmountMHACmd string
    MaxWaitPing int
}

type MetaDBInfo struct {
    MetaUser string
    MetaPasswd []byte
    MetaDBName string
    MetaHost string
    MetaPort int
    MetaNetWork string
    MetaEncodePasswd string
    MetaDsn string
}

type ClusterInfo struct {
    HaType string
    VipAddr string
}

func NewFlagArgs() (*OrchCfg, error) {
    flag.StringVar(&o.ClusterName, "ClusterName", "ClusterName", "cluster name")
    flag.StringVar(&o.DeadStatus, "DeadStatus", "DeadStatus", "dead status")
    flag.StringVar(&o.OldMaster, "OldMaster", "OldMaster", "old master")
    flag.StringVar(&o.OldMasterPort, "OldMasterPort", "OldMasterPort", "old master port")
    flag.StringVar(&o.NewMaster, "NewMaster", "NewMaster", "new master host")
    flag.StringVar(&o.NewMasterPort, "NewMasterPort", "NewMasterPort", "new master port")
    flag.StringVar(&o.VipAddr, "VipAddr", "VipAddr", "vip addr, env: {failureClusterDomain}")
    flag.StringVar(&o.InterFace, "InterFace", "InterFace", "net interface, like eth0 bond0 ...")
    flag.Parse()
    meatInfo()
    // FromMetaGetInfo()
    o.MaxWaitPing = 30
    o.SSHUser = "mysql"
    o.SSHPublicKeys = "/home/mysql/.ssh/id_rsa"
    c.HaType = "kp"
    // ho, _ := os.Hostname()
    // if strings.Contains(ho, "stg") || strings.Contains(o.OldMaster, "192.168") || strings.Contains(o.NewMaster, "192.168") {
    //     o.InterFace = "eth0"
    //     o.Mask = "24"
    // } else if (strings.Contains(o.NewMaster, "yue") && strings.Contains(o.OldMaster, "yue")) || (strings.Contains(o.NewMaster, "10.26") && strings.Contains(o.OldMaster, "10.26")) {
    //     o.InterFace = "bond1"
    //     o.Mask = "23"
    // } else if (strings.Contains(o.NewMaster, "yun") && strings.Contains(o.OldMaster, "yun")) || (strings.Contains(o.NewMaster, "10.24") && strings.Contains(o.OldMaster, "10.24")){
    //     o.InterFace = "bond0"
    //     o.Mask = "23"
    // }

    // keepalived
    o.CmdVipAdd = fmt.Sprintf("source /etc/profile;/usr/bin/sudo ip addr add %v/32 dev %v", c.VipAddr, o.InterFace)
    o.CmdVipDel = fmt.Sprintf("source /etc/profile;/usr/bin/sudo ip addr del %v/32 dev %v", c.VipAddr, o.InterFace)
    o.CmdVipstat = fmt.Sprintf("source /etc/profile;ping -c 1 -W 1 %v", c.VipAddr)

    // mha
    o.GateWayCmd = fmt.Sprintf("source /etc/profile;route -n | grep 'UG[ \\t]' | head -n 1 | awk '{print $2}'")
    o.MHAKeyCmd = fmt.Sprintf("source /etc/profile;source ~/.bash_profile;ip a | grep -v 127.0.0.1 | grep -i %s | grep -v inet6 | awk -F 'eth|bond' '{print $2}' | awk -F ':' '{print $2}'", c.VipAddr)
    getGateWay()
    o.ArPingCmdb = fmt.Sprintf("source /etc/profile;/usr/bin/sudo /sbin/arping -I %s -c 1 -s %s %s", o.InterFace, c.VipAddr, o.Mask)
    o.MountMHACmd = fmt.Sprintf("source /etc/profile;/usr/bin/sudo /sbin/ifconfig %v:%v %v/%v", o.InterFace, o.MHAKey, c.VipAddr, o.Mask)
    o.UmountMHACmd = fmt.Sprintf("source /etc/profile;/usr/bin/sudo /sbin/ifconfig %v:%v down", o.InterFace, o.Mask)
    return &OrchCfg{ClusterName: o.ClusterName, DeadStatus: o.DeadStatus, OldMaster: o.OldMaster, OldMasterPort: o.OldMasterPort, NewMaster: o.NewMaster, NewMasterPort: o.NewMasterPort, VipAddr: c.VipAddr, Mask: o.Mask, InterFace: o.InterFace, CmdVipAdd: o.CmdVipAdd, CmdVipDel: o.CmdVipDel, CmdVipstat: o.CmdVipstat, SSHUser: o.SSHUser, SSHPublicKeys: o.SSHPublicKeys, GateWayCmd: o.GateWayCmd, MHAKey: o.MHAKey, ArPingCmdb: o.ArPingCmdb, ClusterType: c.HaType, MountMHACmd: o.MountMHACmd, UmountMHACmd: o.UmountMHACmd, MaxWaitPing: o.MaxWaitPing}, nil
}

func getGateWay() {
    // 连接老实例 获取MH A网关、拿到MHAKEYMHAKEY
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
    gateWay, err2 := session.CombinedOutput(o.GateWayCmd)
    if err2 != nil {
        return
    }
    gateway := strings.Replace(string(gateWay), "\n", "", -1)
    o.GateWay = string(gateway)
    if c.HaType == "mha" {
        // 连接老实例 拿到MHAKEY 1
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
        mhaKey, err2 := session.CombinedOutput(o.MHAKeyCmd)
        if err2 != nil {
            return
        }
        MHAKEY := strings.Replace(string(if), "\n", "", -1)
        o.MHAKey = string(MHAKEY)
    }
}

func meatInfo() {
    m.MetaEncodePasswd = "xxxxxx"
    m.MetaPasswd, _ = base64.StdEncoding.DecodeString(m.MetaEncodePasswd)
    m.MetaUser = "xxx"
    m.MetaDBName = "metadb"
    hn, _ = os.Hostname()
    if strings.Contains(hn, "prod") {
        m.MetaHost = "xxxx"
    } else {
        m.MetaHost = "xxxx"
    }
    m.MetaPort = 3306
    m.MetaNetWork = "tcp"
    m.MetaDsn = fmt.Sprintf("%s:%s@%s(%s:%s\d)/%s", m.MetaUser, m.MetaPasswd, m.MetaNetWork, m.MetaHost, m.MetaPort, m.MetaDBName)
}

func FromMetaGetInfo() {
    // 通过cmdb获取到集群类型 以兼容MHA和KP的VIP漂移 自己修改自己环境所需
    MessSelectSQL := fmt.Sprintf("select ha_type,vip from meta_info where db_type = 'mysql' and cluster_name = '%v' and vip is not null and vip <> ''", o.ClusterName)
    DB, err := sql.Open("mysql", m.MetaDsn)
    if err != nil {
        return
    }
    defer DB.Close()
    err1 := DB.QueryRow(MessSelectSQL).Scan(&c.HaType, &c.VipAddr)
    if err1 != nil {
        return
    }
}

fun publicKeyAuthFunc(KPath string) ssh.AuthMethod {
    KeyPath, err := homedir.Expand(KPath)
    if err != nil {
        return
    }
    key, err := ioutil.ReadFile(KeyPath)
    if err != nil {
        return
    }
    signer, err := ssh.ParsePrivateKey(key)
    if err != nil {
        return
    }
    return ssh.PublicKeys(signer)
}