package httpsrv
import (
	"net/http"
	"testing"
)


type stubHttpResonpseFetcher struct {}
type stubKubeCmd struct {}

func (kubeCmd *stubKubeCmd) GetNodesIP(appName, appNamespace string) ([]byte, error) {
	ipRaw := "192.168.1.10\n192.168.1.11\n192.168.1.21"
	ips := []byte(ipRaw)

	return ips, nil
}

func (kubeCmd *stubKubeCmd) Restart(appName, appNamespace, oldVersion, oldImage string) (string, error) {
	return "", nil
}


func (kubeCmd *stubKubeCmd) RollingUpdate(appName, appNamespace string, image string) (string, error) {
	return "", nil
}

func (fetcher *stubHttpResonpseFetcher) SendJsonRequest(http_method, urls string, req_body interface{}) (int, string, error) {
	return http.StatusOK, "", nil
}

func (fetcher *stubHttpResonpseFetcher) GetKubeAppContent(appName, appNamespace string) (int, string, error) {
	return http.StatusOK, "", nil
}

func TestRegisterDNS(t *testing.T)  {
	kubecmd := new(stubKubeCmd)
	httpResponseFetcher := new(stubHttpResonpseFetcher)

	registerDNS("test", "default", kubecmd, httpResponseFetcher)
}