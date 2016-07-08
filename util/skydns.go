package util

import (
	. "backend/common"
	"strings"
	"backend/common/protocol"
	"net/http"
	"codoon_ops/kubernetes-apiproxy/util/set"
)

var SkyDnsUrl = ""

func RegisterDNS(appName, appNamespace string, kubeCmd KubeCmd, fetcher KubeResponseFetch) (int, error) {
	domain := ""
	switch appNamespace {
	case "default":
		domain = appName + ".codoon.com"
	case "in":
		domain = appName + ".in.codoon.com"
	}
	//调用kube cmd命令获取数据
	bytes, err := kubeCmd.GetNodesIP(appName, appNamespace)

	if err != nil {
		return http.StatusInternalServerError, err
	}

	reqJson := bytesToDNSRequestJSON(bytes, domain)
	statusCode, _, err := fetcher.SendJsonRequest("POST", SkyDnsUrl, reqJson)

	if err != nil {
		return statusCode, err
	}
	return http.StatusOK, nil
}

func bytesToDNSRequestJSON(bytes []byte, domain string) (req interface{}) {
	arrays := strings.Split(string(bytes), "\n")
	set := set.New()

	for _, item := range arrays {
		if item != "" {
			set.Add(item)
			Logger.Debug("the ip is: %v", item)
		}
	}
	urlList := set.List()

	urls := make([]protocol.RR, 0)
	for _, url := range urlList {
		item := protocol.RR{
			Host: url.(string),
		}
		urls = append(urls, item)
	}

	reqJson := protocol.SetDnsReq{
		URL: domain,
		RRs: urls,
	}

	//调用http接口,注册服务名到DNS server
	Logger.Debug("registry %v to DNS server, the ips is: %v", domain, urls)
	return reqJson
}