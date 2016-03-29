package httpsrv

import . "backend/common"

type HttpResponseFetch interface {
	SendJsonRequest(http_method, urls string, req_body interface{}) (int, string, error)
	GetKubeAppContent(appName, appNamespace string) (int, string, error)
}

type HttpResponseFetcher struct {}

func (fetcher *HttpResponseFetcher) SendJsonRequest(http_method, urls string, req_body interface{})  (int, string, error) {
	Logger.Info("Send request to DNS server url is: %v", urls)
	return SendJsonRequest(http_method, urls, req_body)
}

func (fetcher *HttpResponseFetcher) GetKubeAppContent(appName, appNamespace string) (int, string, error) {
	Logger.Info("Starting get app content from kube-apiserver")

	postFix := "/api/v1/namespaces"
	url := "http://" + kubeApiserverPath + kubeApiserverPort + postFix + "/" + appNamespace + "/replicationcontrollers" + "/" + appName

	statusCode, response, err := SendRawRequest("GET", url, nil)
	Logger.Debug("statusCode: %d", statusCode)
	Logger.Debug("repsonse: %s", response)

	return statusCode, response, err
}

