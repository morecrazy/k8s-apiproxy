package rest

import "zhanghan/slb/api/base"

type DescribeLoadBalancerAttributeRequest struct {
	*base.RestApi
	LoadBalancerId string
}

func NewDescribeLoadBalancerAttributeRequest(domain string, loadBalancerId string) *DescribeLoadBalancerAttributeRequest {
	api := base.NewRestApi(domain)
	var d = new(DescribeLoadBalancerAttributeRequest)
	d.LoadBalancerId = loadBalancerId
	d.RestApi = api

	d.ApiName = "slb.aliyuncs.com.DescribeLoadBalancerAttribute.2014-05-15"
	d.EncodeParams["LoadBalancerId"] = d.LoadBalancerId

	return d
}
