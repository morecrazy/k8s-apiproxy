package httpsrv

import (
	"codoon_ops/kubernetes-apiproxy/util"
	. "backend/common"
	"time"
	"strings"
)

type KubeCmd interface {
	GetNodesIP(appName, appNamespace string) ([]byte, error)
	RollingUpdate(appName, appNamespace string, image string) (string, error)
	Restart(appName, appNamespace, oldVersion, oldImage string) (string, error)
}

type KubeCmdImpl struct {
}

func (kubeCmd *KubeCmdImpl) GetNodesIP(appName, appNamespace string) ([]byte, error) {
	cmd := "kubectl get pods -o wide --namespace=" + appNamespace + " | grep '^" + appName + "' | awk '{print $6}'"
	Logger.Debug("The cmd is: %v", cmd)
	var bytes []byte
	var err error

	timer1 := time.NewTicker(5 * time.Second)
	try:
	for {
		select {
		case <-timer1.C:
			bytes, err = util.ExecCommand(cmd)
			b := string(bytes)
			if b == "" || err != nil {
				Logger.Debug("unable to get ips, continue loop")
			} else {
				Logger.Debug("success get ips")
				break try
			}

		}
	}
	return bytes, err
}

func (kubeCmd *KubeCmdImpl) RollingUpdate(appName, appNamespace string, image string) (string, error) {
	Logger.Info("Starting rolling update app: %v", appName)
	cmd := "kubectl rolling-update " + appName + " --update-period=20s --namespace=" + appNamespace + " --image=" + image
	Logger.Debug("The cmd is: %v", cmd)
	go util.ExecCommand(cmd)
	return cmd, nil
}

func (kubeCmd *KubeCmdImpl) Restart(appName, appNamespace, oldVersion, oldImage string) (string, error) {
	Logger.Info("Starting restart app: %v", appName)
	t := time.Now()
	newVersion := t.Format("20060102150405")
	newVersion = "v" + newVersion

	newAppName := ""
	strs := strings.Split(appName, "-")
	for i :=0; i < len(strs) -1; i++ {
		newAppName += strs[i] + "-"
	}
	newAppName = newAppName + newVersion
	Logger.Debug("New App Name is: %v", newAppName)

	//滚动更新
	cmd := "kubectl get rc " + appName + " --namespace=" + appNamespace + " -o yaml " +
			" | sed 's/resourceVersion:.*/resourceVersion: ''/g' " +
			" | sed 's/name: " + appName + "/name: " + newAppName + "/g' " +   //替换rc名字
			" | sed 's/version: " + oldVersion + "/version: " + newVersion + "/g' " +		  //替换rc版本
			" | kubectl rolling-update " + appName + " --update-period=20s --namespace=" + appNamespace + " -f - "								//滚动更新

	Logger.Debug("The cmd is: %v", cmd)
	go util.ExecCommand(cmd)
	return cmd, nil
}