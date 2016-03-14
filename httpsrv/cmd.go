package httpsrv

import (
	"codoon_ops/kubernetes-apiproxy/util"
	. "backend/common"
)

type KubeCmd interface {
	GetNodesIP(appName, appNamespace string) ([]byte, error)
}

type KubeCmdImpl struct {
}

func (kubeCmd KubeCmdImpl) GetNodesIP(appName, appNamespace string) ([]byte, error) {
	cmd := "kubectl get pods -o wide --namespace=" + appNamespace + " | grep '^" + appName + "' | awk '{print $6}'"
	Logger.Debug("The cmd is: %v", cmd)

	bytes, err := util.ExecCommand(cmd)

	return bytes, err
}