package rest

import "zhanghan/slb/api/base"
import "strconv"

type DescribeLoadBalancerTCPListenerAttributeRequest struct {
	*base.RestApi
	ListenerPort   string
	LoadBalancerId string
}

func NewDescribeLoadBalancerTCPListenerAttributeRequest(domain string, listenerPort int, loadBalancerId string) *DescribeLoadBalancerTCPListenerAttributeRequest {
	api := base.NewRestApi(domain)
	var d = new(DescribeLoadBalancerTCPListenerAttributeRequest)
	d.ListenerPort = strconv.Itoa(listenerPort)
	d.LoadBalancerId = loadBalancerId
	d.RestApi = api

	d.ApiName = "slb.aliyuncs.com.DescribeLoadBalancerTCPListenerAttribute.2014-05-15"
	d.EncodeParams["ListenerPort"] = d.ListenerPort
	d.EncodeParams["LoadBalancerId"] = d.LoadBalancerId

	return d
}
