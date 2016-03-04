package main

import (
	"flag"
	"fmt"
	_ "third/go-sql-driver/mysql"
	"backend/common"
	"runtime"
	"codoon_ops/kubernetes-apiproxy/httpsrv"
	//"zhanghan/slb/api"
)

const (
	DEFAULT_CONF_FILE = "./apiproxy.conf"
)

var g_conf_file string
var g_config common.Configure


func main() {
	//set runtime variable
	runtime.GOMAXPROCS(runtime.NumCPU())
	//get flag
	flag.Parse()

	if g_conf_file == "" {
		g_conf_file = DEFAULT_CONF_FILE
	}

	err := common.InitConfigFile(g_conf_file, &g_config)
	if err != nil {
		fmt.Println(err)
		return
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

	//api.SetAppInfo("KM0BWn5yGIjiYW3S", "VIECii4MYVv7QEVl5QDJbAxGH6nqH0")
	g_logger.Debug("Start server...")
	httpsrv.StartServer()

}
