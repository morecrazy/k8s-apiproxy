package util

import "backend/common"

var KubeUrl = ""

type KubeResponseFetch interface {
	SendJsonRequest(http_method, urls string, req_body interface{}) (int, string, error)
	GetKubeAppContent(appName, appNamespace string) (int, string, error)
	GetKubeSvcContent(svcName, appNamespace string) (int, string, error)
	GetNodesList() (int, string, error)
	GetRcList() (int, string, error)
	CreateSvc(appNamespace string, svcRequestJson interface{}) (int, string, error)
	CreateRc(appNamespace string, rcRequestJson interface{}) (int, string, error)
}

type KubeResponseFetcher struct {}

func (fetcher *KubeResponseFetcher) SendJsonRequest(http_method, urls string, req_body interface{})  (int, string, error) {
	common.Logger.Info("Send request to DNS server url is: %v", urls)
	return common.SendJsonRequest(http_method, urls, req_body)
}

func (fetcher *KubeResponseFetcher) GetKubeAppContent(appName, appNamespace string) (int, string, error) {
	common.Logger.Info("Starting get app content from kube-apiserver")
	postFix := "/api/v1/namespaces/" + appNamespace + "/replicationcontrollers" + "/" + appName
	url := KubeUrl + postFix

	statusCode, response, err := common.SendRawRequest("GET", url, nil)
	return statusCode, response, err
}

func (fetcher *KubeResponseFetcher) GetKubeSvcContent(svcName, appNamespace string) (int, string, error) {
	common.Logger.Debug("Svc Name is: %v", svcName)
	postFix := "/api/v1/namespaces/" + appNamespace + "/services" + "/" + svcName
	url := KubeUrl + postFix
	statusCode, response, err := common.SendRawRequest("GET", url, nil)
	return statusCode, response, err
}

func (fetcher *KubeResponseFetcher) GetNodesList() (int, string, error) {
	postFix := "/api/v1/nodes"
	url := KubeUrl + postFix
	statusCode, response, err := common.SendRawRequest("GET", url, nil)
	return statusCode, response, err
}

func (fetcher *KubeResponseFetcher) GetRcList() (int, string, error) {
	postFix := "/api/v1/replicationcontrollers"
	url := KubeUrl + postFix
	statusCode, response, err := common.SendRawRequest("GET", url, nil)
	return statusCode, response, err
}

func (fetcher *KubeResponseFetcher) CreateSvc(appNamespace string, svcRequestJson interface{}) (int, string, error) {
	postFix := "/api/v1/namespaces/" + appNamespace + "/services"
	url := KubeUrl + postFix
	statusCode, response, err := common.SendJsonRequest("POST", url, svcRequestJson)
	common.Logger.Debug("statusCode: %d", statusCode)
	common.Logger.Debug("repsonse: %s", response)

	return statusCode, response, err
}

func (fetcher *KubeResponseFetcher) CreateRc(appNamespace string, rcRequestJson interface{}) (int, string, error) {
	postFix := "/api/v1/namespaces/" + appNamespace + "/replicationcontrollers"
	url := KubeUrl + postFix
	statusCode, response, err := common.SendJsonRequest("POST", url, rcRequestJson)
	common.Logger.Debug("statusCode: %d", statusCode)
	common.Logger.Debug("repsonse: %s", response)

	return statusCode, response, err
}