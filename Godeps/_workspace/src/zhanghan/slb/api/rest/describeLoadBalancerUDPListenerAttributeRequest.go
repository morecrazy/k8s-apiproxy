package rest

import "zhanghan/slb/api/base"
import "strconv"

type DescribeLoadBalancerUDPListenerAttributeRequest struct {
	*base.RestApi
	ListenerPort   string
	LoadBalancerId string
}

func NewDescribeLoadBalancerUDPListenerAttributeRequest(domain string, listenerPort int, loadBalancerId string) *DescribeLoadBalancerUDPListenerAttributeRequest {
	api := base.NewRestApi(domain)
	var d = new(DescribeLoadBalancerUDPListenerAttributeRequest)
	d.ListenerPort = strconv.Itoa(listenerPort)
	d.LoadBalancerId = loadBalancerId
	d.RestApi = api

	d.ApiName = "slb.aliyuncs.com.DescribeLoadBalancerUDPListenerAttribute.2014-05-15"
	d.EncodeParams["ListenerPort"] = d.ListenerPort
	d.EncodeParams["LoadBalancerId"] = d.LoadBalancerId

	return d
}
