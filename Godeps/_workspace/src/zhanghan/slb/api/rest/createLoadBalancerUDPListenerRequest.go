package rest

import "zhanghan/slb/api/base"
import "strconv"

type CreateLoadBalancerUDPListenerRequest struct {
	*base.RestApi
	ListenerPort      string
	BackendServerPort string
	LoadBalancerId    string
	Bandwidth         string
}

func NewCreateLoadBalancerUDPListenerRequest(domain string, listenerPort int, backendServerPort int, loadBalancerId string) *CreateLoadBalancerUDPListenerRequest {
	api := base.NewRestApi(domain)
	var c = new(CreateLoadBalancerUDPListenerRequest)
	c.ListenerPort = strconv.Itoa(listenerPort)
	c.BackendServerPort = strconv.Itoa(backendServerPort)
	c.LoadBalancerId = loadBalancerId
	c.RestApi = api

	c.ApiName = "slb.aliyuncs.com.CreateLoadBalancerUDPListener.2014-05-15"
	c.EncodeParams["ListenerPort"] = c.ListenerPort
	c.EncodeParams["BackendServerPort"] = c.BackendServerPort
	c.EncodeParams["LoadBalancerId"] = c.LoadBalancerId
	c.EncodeParams["Bandwidth"] = strconv.Itoa(-1)

	return c
}
