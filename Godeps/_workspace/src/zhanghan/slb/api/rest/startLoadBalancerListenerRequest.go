package rest

import "zhanghan/slb/api/base"
import "strconv"

type StartLoadBalancerListenerRequest struct {
	*base.RestApi
	ListenerPort   string
	LoadBalancerId string
}

func NewStartLoadBalancerListenerRequest(domain string, listenerPort int, loadBalancerId string) *StartLoadBalancerListenerRequest {
	api := base.NewRestApi(domain)
	var d = new(StartLoadBalancerListenerRequest)
	d.ListenerPort = strconv.Itoa(listenerPort)
	d.LoadBalancerId = loadBalancerId
	d.RestApi = api

	d.ApiName = "slb.aliyuncs.com.StartLoadBalancerListener.2014-05-15"
	d.EncodeParams["ListenerPort"] = d.ListenerPort
	d.EncodeParams["LoadBalancerId"] = d.LoadBalancerId

	return d
}
