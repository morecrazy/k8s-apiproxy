package rest

import "zhanghan/slb/api/base"
import "strconv"

type StopLoadBalancerListenerRequest struct {
	*base.RestApi
	ListenerPort   string
	LoadBalancerId string
}

func NewStopLoadBalancerListenerRequest(domain string, listenerPort int, loadBalancerId string) *StopLoadBalancerListenerRequest {
	api := base.NewRestApi(domain)
	var d = new(StopLoadBalancerListenerRequest)
	d.ListenerPort = strconv.Itoa(listenerPort)
	d.LoadBalancerId = loadBalancerId
	d.RestApi = api

	d.ApiName = "slb.aliyuncs.com.StopLoadBalancerListener.2014-05-15"
	d.EncodeParams["ListenerPort"] = d.ListenerPort
	d.EncodeParams["LoadBalancerId"] = d.LoadBalancerId

	return d
}
