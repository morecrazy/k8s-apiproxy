package rest

import "zhanghan/slb/api/base"
import "strconv"

type CreateLoadBalancerTCPListenerRequest struct {
	*base.RestApi
	ListenerPort           string
	BackendServerPort      string
	HealthCheckConnectPort string
	LoadBalancerId         string
	Bandwidth              string
}

func NewCreateLoadBalancerTCPListenerRequest(domain string, listenerPort int, backendServerPort int, loadBalancerId string) *CreateLoadBalancerTCPListenerRequest {
	api := base.NewRestApi(domain)
	var c = new(CreateLoadBalancerTCPListenerRequest)
	c.ListenerPort = strconv.Itoa(listenerPort)
	c.BackendServerPort = strconv.Itoa(backendServerPort)
	c.HealthCheckConnectPort = strconv.Itoa(backendServerPort)
	c.LoadBalancerId = loadBalancerId
	c.RestApi = api

	c.ApiName = "slb.aliyuncs.com.CreateLoadBalancerTCPListener.2014-05-15"
	c.EncodeParams["ListenerPort"] = c.ListenerPort
	c.EncodeParams["BackendServerPort"] = c.BackendServerPort
	c.EncodeParams["LoadBalancerId"] = c.LoadBalancerId
	c.EncodeParams["HealthCheckConnectPort"] = c.HealthCheckConnectPort
	c.EncodeParams["Bandwidth"] = strconv.Itoa(-1)

	return c
}
