package httpsrv
import (
	"zhanghan/slb/api/rest"
	"backend/common"
	"encoding/json"
	"errors"
)
func createListener(listenerPort, backendPort int, loadBalancerId, protocol string) (bool, error) {
	common.Logger.Info("Begin creating slb listener: listenerPort->%s, backendPort->%s, loadBalancerId->%s, protocl->%s",
		listenerPort, backendPort, loadBalancerId, protocol )
	var res string
	switch protocol {
	case "TCP":
		createApi := rest.NewCreateLoadBalancerTCPListenerRequest("http://slb.aliyuncs.com/", listenerPort, backendPort, loadBalancerId)
		res = createApi.GetResponse("", "60")
	case "HTTP":
		createApi := rest.NewCreateLoadBalancerHTTPListenerRequest("http://slb.aliyuncs.com/", listenerPort, backendPort, loadBalancerId)
		res = createApi.GetResponse("", "60")
	case "HTTPS":
		createApi := rest.NewCreateLoadBalancerHTTPSListenerRequest("http://slb.aliyuncs.com/", listenerPort, backendPort, loadBalancerId)
		res = createApi.GetResponse("", "60")
	case "UDP":
		createApi := rest.NewCreateLoadBalancerUDPListenerRequest("http://slb.aliyuncs.com/", listenerPort, backendPort, loadBalancerId)
		res = createApi.GetResponse("", "60")
	}

	common.Logger.Info("the create result is : %s", res)
	var dat map[string]interface{}
	json.Unmarshal([]byte(res), &dat)

	if dat["Code"] != nil {
		return false, errors.New(dat["Message"].(string))
	}
	return true, nil
}

func startListener(listenerPort int, loadBalancerId string) (bool, error) {
	common.Logger.Debug("Begin starting listener")
	startApi := rest.NewStartLoadBalancerListenerRequest("http://slb.aliyuncs.com/", listenerPort, loadBalancerId)
	res := startApi.GetResponse("", "60")
	common.Logger.Debug("the result is : %s", res)

	var dat map[string]interface{}
	json.Unmarshal([]byte(res), &dat)

	if dat["Code"] != nil {
		return false, errors.New(dat["Message"].(string))
	}

	return true, nil
}
