package rest

import "zhanghan/slb/api/base"
import "strconv"

type DeleteLoadBalancerListenerRequest struct {
	*base.RestApi
	ListenerPort   string
	LoadBalancerId string
}

func NewDeleteLoadBalancerListenerRequest(domain string, listenerPort int, loadBalancerId string) *DeleteLoadBalancerListenerRequest {
	api := base.NewRestApi(domain)
	var d = new(DeleteLoadBalancerListenerRequest)
	d.ListenerPort = strconv.Itoa(listenerPort)
	d.LoadBalancerId = loadBalancerId
	d.RestApi = api

	d.ApiName = "slb.aliyuncs.com.DeleteLoadBalancerListener.2014-05-15"
	d.EncodeParams["ListenerPort"] = d.ListenerPort
	d.EncodeParams["LoadBalancerId"] = d.LoadBalancerId

	return d
}
