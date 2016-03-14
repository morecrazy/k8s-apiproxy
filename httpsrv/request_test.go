package httpsrv
import (
	"net/http"
	"testing"
)

type stubKubeCmd struct {}
type stubHttpResonpseFetcher struct {}

func (kubeCmd *stubKubeCmd) GetNodesIP(appName, appNamespace string) ([]byte, error) {
	ipRaw := "192.168.1.10\n192.168.1.11\n192.168.1.21"
	ips := []byte(ipRaw)

	return ips, nil
}

func (fetcher *stubHttpResonpseFetcher) SendJsonRequest(http_method, urls string, req_body interface{}) (int, string, error) {
	return http.StatusOK, "", nil
}


func testRegisterDNS(t *testing.T)  {
	var kubecmd stubKubeCmd
	var httpResponseFetcher stubHttpResonpseFetcher

	registerDNS("test", "default", kubecmd, httpResponseFetcher)
}