package httpsrv

import (
	"codoon_ops/kubernetes-apiproxy/util"
	. "backend/common"
	"time"
)

type KubeCmd interface {
	GetNodesIP(appName, appNamespace string) ([]byte, error)
}

type KubeCmdImpl struct {
}

func (kubeCmd KubeCmdImpl) GetNodesIP(appName, appNamespace string) ([]byte, error) {
	cmd := "kubectl get pods -o wide --namespace=" + appNamespace + " | grep '^" + appName + "' | awk '{print $6}'"
	Logger.Debug("The cmd is: %v", cmd)
	var bytes []byte
	var err error

	timer1 := time.NewTicker(10 * time.Second)
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