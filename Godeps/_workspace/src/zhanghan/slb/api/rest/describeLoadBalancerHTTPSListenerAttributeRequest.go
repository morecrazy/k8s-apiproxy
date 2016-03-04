package rest

import "zhanghan/slb/api/base"
import "strconv"

type DescribeLoadBalancerHTTPSListenerAttributeRequest struct {
	*base.RestApi
	ListenerPort   string
	LoadBalancerId string
}

func NewDescribeLoadBalancerHTTPSListenerAttributeRequest(domain string, listenerPort int, loadBalancerId string) *DescribeLoadBalancerHTTPSListenerAttributeRequest {
	api := base.NewRestApi(domain)
	var d = new(DescribeLoadBalancerHTTPSListenerAttributeRequest)
	d.ListenerPort = strconv.Itoa(listenerPort)
	d.LoadBalancerId = loadBalancerId
	d.RestApi = api

	d.ApiName = "slb.aliyuncs.com.DescribeLoadBalancerHTTPSListenerAttribute.2014-05-15"
	d.EncodeParams["ListenerPort"] = d.ListenerPort
	d.EncodeParams["LoadBalancerId"] = d.LoadBalancerId

	return d
}
