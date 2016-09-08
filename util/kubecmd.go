package util

import (
	"backend/common"
	"time"
	"strings"
)

type KubeCmd interface {
	GetNodesIP(appName, appNamespace string) ([]byte, error)
	Delivery(appName, appNamespace, image string) (string, error)
	Restart(appName, appNamespace, oldVersion, oldImage string) (string, error)
	UpdateQuota(appName, appNamespace, oldVersion, oldcpu, oldmem, cpu, mem string) (string, error)
	UpdateCmd(appName, appNamespace, oldVersion, oldCmd, newCmd string) (string, error)
	Stop(appName, appNamespace, file string) error
	Create(appName, appNamespace, file string) error
	Scale(appName, appNamespace, replicas string) error
}

type KubeCmdImpl struct {
}

func (kubeCmd *KubeCmdImpl) GetNodesIP(appName, appNamespace string) ([]byte, error) {
	cmd := "kubectl -s " + KubeUrl + " get pods -o wide --namespace=" + appNamespace + " | grep '^" + appName + "' | awk '{print $6}'"
	common.Logger.Info("The cmd is: %v", cmd)
	var bytes []byte
	var err error

	timer1 := time.NewTicker(5 * time.Second)
	try:
	for {
		select {
		case <-timer1.C:
			bytes, err = ExecCommand(cmd)
			b := string(bytes)
			if b == "" || err != nil {
				common.Logger.Debug("unable to get ips, continue loop")
			} else {
				common.Logger.Debug("success get ips")
				break try
			}
		}
	}
	return bytes, err
}

func (kubeCmd *KubeCmdImpl) Delivery(appName, appNamespace, image string) (string, error) {
	cmd := "kubectl -s " + KubeUrl + " rolling-update " + appName + " --update-period=5s --namespace=" + appNamespace + " --image=" + image
	common.Logger.Info("The cmd is: %v", cmd)
	go Run(cmd, appName, appNamespace)
	return cmd, nil
}

func (kubeCmd *KubeCmdImpl) Restart(appName, appNamespace, oldVersion, oldImage string) (string, error) {
	t := time.Now()
	newVersion := t.Format("20060102150405")
	newVersion = "v" + newVersion

	newAppName := ""
	strs := strings.Split(appName, "-v")
	newAppName = strs[0]
	/**
	for i :=0; i < len(strs) -1; i++ {
		newAppName += strs[i] + "-"
	}**/
	newAppName = newAppName + "-" + newVersion
	common.Logger.Debug("New App Name is: %v", newAppName)

	//滚动更新
	cmd := "kubectl -s " + KubeUrl + " get rc " + appName + " --namespace=" + appNamespace + " -o yaml " +
	" | sed 's/resourceVersion:.*/resourceVersion: ''/g' " +
	" | sed 's/name: " + appName + "/name: " + newAppName + "/g' " +   //替换rc名字
	" | sed 's/version: " + oldVersion + "/version: " + newVersion + "/g' " +		  //替换rc版本
	" | kubectl rolling-update " + appName + " --update-period=5s --namespace=" + appNamespace + " -f - "	//滚动更新

	common.Logger.Info("The cmd is: %v", cmd)
	go Run(cmd, appName, appNamespace)
	return cmd, nil
}

func (kubeCmd *KubeCmdImpl) UpdateQuota(appName, appNamespace, oldVersion, oldCpu, oldMem, cpu, mem string) (string, error) {
	//生成新版本号,以日期作为标致
	t := time.Now()
	newVersion := t.Format("20060102150405")
	newVersion = "v" + newVersion

	newAppName := ""
	strs := strings.Split(appName, "-v")
	newAppName = strs[0]

	newAppName = newAppName + "-" + newVersion
	common.Logger.Info("New App Name is: %v", newAppName)

	//滚动更新
	cmd := "kubectl -s " + KubeUrl + " get rc " + appName + " --namespace=" + appNamespace + " -o yaml " +
	" | sed 's/resourceVersion:.*/resourceVersion: ''/g' " +
	" | sed 's/name: " + appName + "/name: " + newAppName + "/g' " +   //替换rc名字
	" | sed 's/version: " + oldVersion + "/version: " + newVersion + "/g' " +		  //替换rc版本
	" | sed 's/cpu: \"" + oldCpu + "\"/cpu: " + cpu + "/g' " + 							 //替换cpu
	" | sed 's/memory: " + oldMem + "/memory: " + mem + "/g' " + 					//替换mem
	" | kubectl rolling-update " + appName + " --update-period=5s --namespace=" + appNamespace + " -f - "								//滚动更新

	common.Logger.Info("The cmd is: %v", cmd)
	go Run(cmd, appName, appNamespace)
	return cmd, nil
}

func (kubeCmd *KubeCmdImpl) UpdateCmd(appName, appNamespace, oldVersion, oldCmd, newCmd string) (string, error) {
	//生成新版本号,以日期作为标致
	t := time.Now()
	newVersion := t.Format("20060102150405")
	newVersion = "v" + newVersion

	newAppName := ""
	strs := strings.Split(appName, "-v")
	newAppName = strs[0]

	newAppName = newAppName + "-" + newVersion
	common.Logger.Debug("New App Name is: %v", newAppName)

	//滚动更新
	cmd := "kubectl -s " + KubeUrl + " get rc " + appName + " --namespace=" + appNamespace + " -o yaml " +
	" | sed 's/resourceVersion:.*/resourceVersion: ''/g' " +
	" | sed 's/name: " + appName + "/name: " + newAppName + "/g " +   //替换rc名字
	" | sed 's/version: " + oldVersion + "/version: " + newVersion + "/g " +		  //替换rc版本
	" | sed 's/command: " + oldCmd + "/cmd: " + newCmd + "/g " + 					//替换mem
	" | kubectl rolling-update " + appName + " --update-period=5s --namespace=" + appNamespace + " -f - "

	common.Logger.Info("The cmd is: %v", cmd)
	go Run(cmd, appName, appNamespace)
	return cmd, nil
}

func (kubeCmd *KubeCmdImpl) Stop(appName, appNamespace, file string) error {
	common.Logger.Info("Stopping service: %v", appName)
	//先保存当前配置到特地文件中
	cmd := "kubectl -s " + KubeUrl + " get rc " + appName + " --namespace=" + appNamespace + " -o yaml >" + file
	if _, err := ExecCommand(cmd); err != nil {
		return err
	}

	//缩减容器个数为0既而停止服务
	cmd = "kubectl -s " + KubeUrl + " scale rc " + appName + " --namespace=" + appNamespace + " --replicas=0"
	if _, err := ExecCommand(cmd); err != nil {
		return err
	}
	return nil
}

func (kubeCmd *KubeCmdImpl) Create(appName, appNamespace, file string) error {
	common.Logger.Info("Starting service: %v", appName)
	//删除久的rc
	cmd := "kubectl -s " + KubeUrl + " delete rc " + appName + " --namespace=" + appNamespace
	if _, err := ExecCommand(cmd); err != nil {
		return err
	}

	cmd = "kubectl -s " + KubeUrl + " create -f " + file
	if _, err := ExecCommand(cmd); err != nil {
		return err
	}
	return nil
}

func (kubeCmd *KubeCmdImpl) Scale(appName, appNamespace, replicas string) error {
	cmd := "kubectl -s " + KubeUrl + " scale rc " + appName + " --namespace=" + appNamespace + " --replicas=" + replicas
	common.Logger.Info("The cmd is: %v", cmd)

	if _, err := ExecCommand(cmd); err != nil {
		return err
	}
	return nil
}