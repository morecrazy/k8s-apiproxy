package rest

import "zhanghan/slb/api/base"
import "strconv"

type DescribeLoadBalancerHTTPListenerAttributeRequest struct {
	*base.RestApi
	ListenerPort   string
	LoadBalancerId string
}

func NewDescribeLoadBalancerHTTPListenerAttributeRequest(domain string, listenerPort int, loadBalancerId string) *DescribeLoadBalancerHTTPListenerAttributeRequest {
	api := base.NewRestApi(domain)
	var d = new(DescribeLoadBalancerHTTPListenerAttributeRequest)
	d.ListenerPort = strconv.Itoa(listenerPort)
	d.LoadBalancerId = loadBalancerId
	d.RestApi = api

	d.ApiName = "slb.aliyuncs.com.DescribeLoadBalancerHTTPListenerAttribute.2014-05-15"
	d.EncodeParams["ListenerPort"] = d.ListenerPort
	d.EncodeParams["LoadBalancerId"] = d.LoadBalancerId

	return d
}
