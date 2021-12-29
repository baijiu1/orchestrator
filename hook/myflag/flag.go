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
    flag.Parse()
    meatInfo()
    FromMetaGetInfo()
    o.MaxWaitPing = 30
    o.SSHUser = "mysql"
    o.SSHPublicKeys = "/home/mysql/.ssh/id_rsa"
    ho, _ := os.Hostname()
    if strings.Contains(ho, "stg") || strings.Contains(o.OldMaster, "192.168") || strings.Contains(o.NewMaster, "192.168") {
        o.InterFace = "eth0"
        o.Mask = "24"
    } else if (strings.Contains(o.NewMaster, "yue") && strings.Contains(o.OldMaster, "yue")) || (strings.Contains(o.NewMaster, "10.26") && strings.Contains(o.OldMaster, "10.26")) {
        o.InterFace = "bond1"
        o.Mask = "23"
    } else if (strings.Contains(o.NewMaster, "yun") && strings.Contains(o.OldMaster, "yun")) || (strings.Contains(o.NewMaster, "10.24") && strings.Contains(o.OldMaster, "10.24")){
        o.InterFace = "bond0"
        o.Mask = "23"
    }

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
    
}