package main

import (
	"flag"
	"fmt"
	_ "third/go-sql-driver/mysql"
	"backend/common"
	"runtime"
	"codoon_ops/kubernetes-apiproxy/httpsrv"
)

const (
	DEFAULT_CONF_FILE = "./apiproxy.conf"
)

var g_conf_file string
var g_config common.Configure

func init() {
	const usage = "kubernetes-apiproxy [-c config_file]"
	flag.StringVar(&g_conf_file, "c", "", usage)
}

func main() {
	//set runtime variable
	runtime.GOMAXPROCS(runtime.NumCPU())
	//get flag
	flag.Parse()

	if g_conf_file != "" {
		if err := common.InitConfigFile(g_conf_file, &g_config); err != nil {
			fmt.Println("init config err : ", err)
		}
	} else {
		addrs := []string{"http://etcd.in.codoon.com:2379"}
		if err := common.LoadCfgFromEtcd(addrs, "kubernetes-apiproxy", &g_config); err != nil {
			fmt.Println("init config from etcd err : ", err)
		}
	}

	g_logger, err := common.InitLogger(g_config.LogFile, "%{color}%{time:2006-01-02 15:04:05.000} %{level:.4s} %{id:03x} ▶ %{shortfunc}%{color:reset} %{message}")

	if err != nil {
		fmt.Println("init log error")
		return
	}

	err = httpsrv.InitDBPool(g_config.MysqlSetting["KubeMysqlSetting"])
	if err != nil {
		fmt.Println("init db error")
		return
	}

	g_logger.Debug("Start server...")
	httpsrv.InitExternalConfig(g_config)
	httpsrv.StartServer()
}
