package httpsrv

import . "backend/common"

type HttpResonseFetcher interface {
	SendJsonRequest(http_method, urls string, req_body interface{}) (int, string, error)
}

type Fetcher struct {

}

func (fetcher Fetcher) SendJsonRequest(http_method, urls string, req_body interface{})  (int, string, error) {
	 return SendJsonRequest(http_method, urls, req_body)
}