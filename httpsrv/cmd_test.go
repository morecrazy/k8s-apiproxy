package httpsrv

import (
	"testing"
	"fmt"
)

func TestRollingUpdate(t *testing.T)  {
	appName := "test-service"
	appNamespace := "in"
	image := "dockerhub.codoon.com/backend/test:v1"
	kubeCmd := new(KubeCmdImpl)
	res, err := kubeCmd.RollingUpdate(appName, appNamespace, image)

	if err != nil {
		t.Error("error")
	}
	fmt.Printf(res)
}
