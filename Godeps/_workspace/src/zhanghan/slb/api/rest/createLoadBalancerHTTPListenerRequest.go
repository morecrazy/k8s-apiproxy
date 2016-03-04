package rest

import "zhanghan/slb/api/base"
import "strconv"

type CreateLoadBalancerHTTPListenerRequest struct {
	*base.RestApi
	ListenerPort      string
	BackendServerPort string
	LoadBalancerId    string
	Bandwidth         string
	StickySession     string
	HealthCheck       string
}

func NewCreateLoadBalancerHTTPListenerRequest(domain string, listenerPort int, backendServerPort int, loadBalancerId string) *CreateLoadBalancerHTTPListenerRequest {
	api := base.NewRestApi(domain)
	var c = new(CreateLoadBalancerHTTPListenerRequest)
	c.ListenerPort = strconv.Itoa(listenerPort)
	c.BackendServerPort = strconv.Itoa(backendServerPort)
	c.LoadBalancerId = loadBalancerId
	c.RestApi = api

	c.ApiName = "slb.aliyuncs.com.CreateLoadBalancerHTTPListener.2014-05-15"
	c.EncodeParams["ListenerPort"] = c.ListenerPort
	c.EncodeParams["BackendServerPort"] = c.BackendServerPort
	c.EncodeParams["LoadBalancerId"] = c.LoadBalancerId
	c.EncodeParams["Bandwidth"] = strconv.Itoa(-1)
	c.EncodeParams["StickySession"] = "off"
	c.EncodeParams["HealthCheck"] = "off"

	return c
}
