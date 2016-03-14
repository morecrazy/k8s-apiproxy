package httpsrv

import . "backend/common"

type HttpResponseFetch interface {
	SendJsonRequest(http_method, urls string, req_body interface{}) (int, string, error)
}

type HttpResponseFetcher struct {}

func (fetcher HttpResponseFetcher) SendJsonRequest(http_method, urls string, req_body interface{})  (int, string, error) {
	Logger.Debug("Send request to DNS server url is: %v", urls)
	return SendJsonRequest(http_method, urls, req_body)
}