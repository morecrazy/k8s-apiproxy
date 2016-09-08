package httpsrv

import (
	"codoon_ops/kubernetes-apiproxy/util"
	"encoding/json"
	"backend/common"
	"strings"
	"strconv"
	"codoon_ops/kubernetes-apiproxy/util/set"
	"net/http"
	"time"
	"fmt"
)

func getImagesList(fetcher util.RegistryResponseFetch) (int, []interface{}, error) {
	statusCode, response, err := fetcher.GetRegistryImageList()
	bytes := []byte(response)
	var res map[string]interface{}

	err = json.Unmarshal(bytes, &res)
	if err != nil {
		return statusCode, nil, err
	}

	list := res["repositories"].([]interface{})

	ret_list := []interface{}{}
	for _, value := range list {
		name := value.(string)
		mirror := registryPath + "/" + name
		ret_list = append(ret_list, map[string]interface{}{
			"name":     name,
			"mirror":  mirror,
			"latest": "latest",
		})
	}
	return statusCode, ret_list, err
}

func getImageTags(name string, fetcher util.RegistryResponseFetch) (int, []interface{}, error) {
	statusCode, response, err := fetcher.GetImageTags(name)

	bytes := []byte(response)
	var res map[string]interface{}
	err = json.Unmarshal(bytes, &res)
	if err != nil {
		return statusCode, nil, err
	}

	ret_list := []interface{}{}
	latestVersion := "latest"
	latestVersionNumber := 0
	versionNames := []string{}
	onlineVersionNames := []string{}
	tagList := []interface{}{}
	if res["tags"] != nil {
		tagList = res["tags"].([]interface{})
	}

	//兼容新老版本的tag命名方式,新命名:online_vXXX,老命名:vXXX
	for _, item := range tagList {
		version := item.(string)
		if strings.Contains(version, "online") {
			onlineVersionNames = append(onlineVersionNames, version)
		} else {
			versionNames = append(versionNames, version)
		}
	}

	if len(onlineVersionNames) != 0 {
		for _, item := range onlineVersionNames {
			version := util.Substr(item, 8, len(item))
			versionNumber, err := strconv.Atoi(version)
			if err != nil {
				common.Logger.Error("invalid tag")
				continue;
			}
			if versionNumber > latestVersionNumber {
				latestVersionNumber = versionNumber
			}
		}
		if latestVersionNumber != 0 {
			latestVersion = "online_v" + strconv.Itoa(latestVersionNumber)
		}

	} else {
		for _, item := range versionNames {
			common.Logger.Debug("The tag is: %v", item)
			version := util.Substr(item, 1, len(item))
			versionNumber, err := strconv.Atoi(version)
			if err != nil {
				common.Logger.Error("invalid tag")
				continue;
			}
			if versionNumber > latestVersionNumber {
				latestVersionNumber = versionNumber
			}
		}
		if latestVersionNumber != 0 {
			latestVersion = "v" + strconv.Itoa(latestVersionNumber)
		}
	}

	for _, item := range tagList {
		value := item.(string)
		var latest = false
		if value == latestVersion {
			latest = true
		}
		ret_list = append(ret_list, map[string]interface{}{
			"commit":  "commit",
			"tag":  value,
			"is_latest": latest,
		})
	}
	return statusCode, ret_list, err
}

func deleteImage(name string, fetcher util.RegistryResponseFetch) (int, string, error) {
	statusCode, response, err := fetcher.DeleteRegistryImage(name)
	return statusCode, response, err
}

func deleteTag(name, tag string, fetcher util.RegistryResponseFetch) (int, string, error) {
	statusCode, response, err := fetcher.DeleteImageTag(name, tag)

	bytes := []byte(response)
	var res map[string]interface{}

	err = json.Unmarshal(bytes, &res)
	if err != nil {
		return statusCode, response, err
	}

	blobs := set.New()
	fs := res["fsLayers"].([]interface{})
	for _, item := range fs {
		value := item.(map[string]interface{})
		blobs.Add(value["blobSum"].(string))
	}

	blobList := blobs.List()
	//遍历层数,依次删除
	for _, item := range blobList {

		statusCode, response, err := fetcher.DeleteImageTag(name, item.(string))

		if statusCode != http.StatusAccepted {
			return statusCode, response, err
		}
	}
	return statusCode, response, err
}

func getClusterList(fetcher util.KubeResponseFetch) (int, []interface{}, error) {
	statusCode, response, err := fetcher.GetNodesList()

	bts := []byte(response)
	var res map[string]interface{}
	if err := json.Unmarshal(bts, &res); err != nil {
		return statusCode, nil, err
	}

	list := res["items"].([]interface{})
	ret_list := []interface{}{}
	group_map := map[string]string{}

	for _, value := range list {
		item := value.(map[string]interface{})
		metadata := item["metadata"].(map[string]interface{})
		labels := metadata["labels"].(map[string]interface{})
		if labels["groupname"] == nil {continue}
		key := labels["groupname"].(string)
		group_map[key] = "groupname"
	}

	for key,_ := range group_map {
		ret_list = append(ret_list, map[string]interface{}{
			"env_name":     key,
		})
	}
	return statusCode, ret_list, err
}

func getServiceList(fetcher util.KubeResponseFetch) (int, []interface{}, error) {
	statusCode, response, err := fetcher.GetRcList()

	bts := []byte(response)
	var res map[string]interface{}
	if err := json.Unmarshal(bts, &res); err != nil {
		return statusCode, nil, err
	}

	list := res["items"].([]interface{})
	ret_list := []interface{}{}

	for _, value := range list {
		item := value.(map[string]interface{})

		//获取服务名字,命名空间和创建时间
		metadata := item["metadata"].(map[string]interface{})
		name := metadata["name"].(string)
		namespace := metadata["namespace"].(string)
		createTime := metadata["creationTimestamp"].(string)
		t, _ := time.Parse(time.RFC3339, createTime)
		timestamp := fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d",
			t.Year(), t.Month(), t.Day(),
			t.Hour() + 8, t.Minute(), t.Second())

		//获取实例个数
		spec := item["spec"].(map[string]interface{})
		replicas := spec["replicas"].(float64)

		//获取服务运行状态
		state := ""
		if replicas == 0{
			state = "stopped"
		} else if replicas > 0 {
			state = "running"
		}

		//获取服务运行环境group
		template := spec["template"].(map[string]interface{})
		tSpec := template["spec"].(map[string]interface{})
		nodeSelector := tSpec["nodeSelector"].(map[string]interface{})
		groupname := nodeSelector["groupname"].(string)

		//获取实例中的镜像名
		images := ""
		containers := tSpec["containers"].([]interface{})
		for _, value := range containers {
			item := value.(map[string]interface{})
			images += item["image"].(string)
			images += ","
		}

		//获取所有实例(容器)运行状态
		status := ""
		cmd := "kubectl -s " + util.KubeUrl + " describe rc " + name + " --namespace=" + namespace
		bytes, err := util.ExecCommand(cmd)
		if err != nil {
			return statusCode, nil, err
		}
		lines := strings.Split(string(bytes), "\n")
		for _, value := range lines {
			cols := strings.Split(value, "\t")
			if cols[0] == "Replicas:" {
				status += cols[1] + " / "
			} else if cols[0] == "Pods Status:" {
				status += cols[1]
			}
		}

		//拼接返回数据
		ret_list = append(ret_list, map[string]interface{}{
			"app_name":     name,
			"app_namespace":  namespace,
			"time": timestamp,
			"env_name": groupname,
			"state": state,
			"status": status,
			"mirror": images,
		})
	}

	return statusCode, ret_list, err
}

func getServiceMeta(appName, appNamespace string, fetcher util.KubeResponseFetch) (int, map[string]interface{}, error) {
	statusCode, response, err := fetcher.GetKubeAppContent(appName, appNamespace)

	bts := []byte(response)
	var res map[string]interface{}
	if err := json.Unmarshal(bts, &res); err != nil {
		return statusCode, nil, err
	}

	//获取更新时间
	metadata := res["metadata"].(map[string]interface{})

	labels := metadata["labels"].(map[string]interface{})
	svcName := ""
	if labels["app"] != nil {
		svcName = labels["app"].(string)
	}

	createTime := metadata["creationTimestamp"].(string)
	t, _ := time.Parse(time.RFC3339, createTime)
	timestamp := fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d",
		t.Year(), t.Month(), t.Day(),
		t.Hour() + 8, t.Minute(), t.Second())

	//获取容器个数
	spec := res["spec"].(map[string]interface{})
	replicas := spec["replicas"].(float64)

	//判断服务运行状态,根据rc个数和容器个数进行判断
	state := ""
	if replicas == 0 {
		state = "stopped"
	}else if replicas > 0 {
		state = "running"
	}

	//获取服务运行环境group
	spec = res["spec"].(map[string]interface{})
	template := spec["template"].(map[string]interface{})
	tSpec := template["spec"].(map[string]interface{})
	nodeSelector := tSpec["nodeSelector"].(map[string]interface{})
	groupname := nodeSelector["groupname"].(string)

	//获取服务容器实例中的镜像名
	images := ""
	containers := tSpec["containers"].([]interface{})
	for _, value := range containers {
		item := value.(map[string]interface{})
		images += item["image"].(string)
		images += ","
	}

	//获取实例运行状态
	status := ""
	cmd := "kubectl -s " + util.KubeUrl + " describe rc " + appName + " --namespace=" + appNamespace
	bts, err = util.ExecCommand(cmd)
	if err != nil {
		return statusCode, nil, err
	}
	lines := strings.Split(string(bts), "\n")
	for _, value := range lines {
		cols := strings.Split(value, "\t")
		if cols[0] == "Replicas:" {
			status += cols[1] + " / "
		} else if cols[0] == "Pods Status:" {
			status += cols[1]
		}
	}

	//获取svc name
	svcNameList := []string{}
	cmd = "kubectl -s " + util.KubeUrl + " get svc -l app=" + svcName + " --namespace=" + appNamespace
	common.Logger.Debug("The cmd is: %v", cmd)
	bts, err = util.ExecCommand(cmd)
	if err != nil {
		return statusCode, nil, err
	}
	lines = strings.Split(string(bts), "\n")
	for index, value := range lines {
		if index == 0 {continue}
		cols := strings.Split(value, " ")
		if cols[0] == "" {continue}
		svcNameList = append(svcNameList, cols[0])
	}

	//获取服务信息,地址,端口,协议类型
	urls := []interface{}{}

	for _, value := range svcNameList {
		svcName := value

		statusCode, response, err := fetcher.GetKubeSvcContent(svcName, appNamespace)

		if statusCode != http.StatusOK {
			return statusCode, nil, err
		}

		bts = []byte(response)

		err = json.Unmarshal(bts, &res)
		if err != nil {
			return statusCode, nil, err
		}

		metadata := res["metadata"].(map[string]interface{})
		portMap := map[string]interface{}{}
		if metadata["annotations"] != nil {
			portMap = metadata["annotations"].(map[string]interface{})
		}

		spec = res["spec"].(map[string]interface{})
		ports := spec["ports"].([]interface{})

		domain := ""
		switch appNamespace {
		case "in":
			domain = appName + ".in.codoon.com"
		case "default":
			domain = appName + ".codoon.com"
		}
		//查询DNS服务,判断域名是否注册上
		//var fetcher HttpResponseFetcher
		//statusCode, response, err = checkServiceDNS(domain, fetcher)
		//if statusCode != http.StatusOK {
		//	s.Status.State = 1
		//	s.Status.Msg = err.Error()
		//	c.JSON(statusCode, s)
		//	return
		//}

		//如果域名没有注册上,则使用cluterIp
		//返回值"OK"表示DNS server没有相应的服务
		if response == "OK" {
			domain = spec["clusterIp"].(string)
		}

		for _, value := range ports {
			item := value.(map[string]interface{})

			port := item["port"].(float64)
			name := strconv.FormatFloat(port, 'g', 5, 64)

			protocol := "HTTP"
			if portMap[name] != nil {
				protocol = portMap[name].(string)
			}
			listenerPort := strconv.FormatFloat(port, 'g', 5, 64)
			url := protocol + "://" + domain + ":" + listenerPort
			urls = append(urls, url)
		}
	}

	data := map[string]interface{}{
		"app_name":     appName,
		"app_namespace":  appNamespace,
		"time": timestamp,
		"env_name": groupname,
		"state": state,
		"status": status,
		"mirror": images,
		"urls": urls,
	}
	return statusCode, data, err
}

func getServiceConfig(appName, appNamespace string, fetcher util.KubeResponseFetch) (int, map[string]interface{}, error) {
	statusCode, response, err := fetcher.GetKubeAppContent(appName, appNamespace)

	if statusCode != http.StatusOK {
		return statusCode, nil, err
	}

	bytes := []byte(response)
	var res map[string]interface{}

	if err = json.Unmarshal(bytes, &res); err != nil {
		return statusCode, nil, err
	}

	//获取容器个数
	spec := res["spec"].(map[string]interface{})
	replicas := spec["replicas"].(float64)

	//获取运行环境
	template := spec["template"].(map[string]interface{})
	tSpec := template["spec"].(map[string]interface{})
	nodeSelector := tSpec["nodeSelector"].(map[string]interface{})
	groupname := nodeSelector["groupname"].(string)

	//获取集群ip列表:ips
	ips, _, err := getClusterIPs(groupname)
	if err != nil {
		return statusCode, nil, err
	}

	containers := tSpec["containers"].([]interface{})
	//TODO 只返回一个容器的信息
	container := containers[0].(map[string]interface{})
	image := container["image"].(string)

	//获取镜像信息
	ims := strings.Split(image, ":")
	imagePath := ims[0]
	imageUrls := strings.Split(imagePath, "/")
	imageName := ""
	if len(imageUrls) <=2 {
		imageName += "library" + "/"
	}
	for index,value := range imageUrls {
		if index == 0 {
			continue
		}else if index < len(imageUrls) - 1 {
			imageName += value + "/"
		} else {
			imageName += value
		}
	}

	tag := "latest" //默认为latest
	if len(ims) == 2 {
		tag = ims[1]
	}

	//获取执行命令信息
	cmd := ""
	if container["command"] != nil {
		cmdList := container["command"].([]interface{})
		for _, value := range cmdList {
			cmd += value.(string) + " "
		}
	}

	//获取资源信息
	cpu := ""
	mem := ""
	if container["resources"] != nil {
		resource := container["resources"].(map[string]interface{})
		limits := resource["limits"].(map[string]interface{})
		if limits["cpu"] != nil {
			cpu = limits["cpu"].(string)
			l := len(cpu)
			if util.Substr(cpu, l - 1, 1) == "m" {
				cpuNum := util.Substr(cpu, 0, l - 1)
				c,_ := strconv.ParseFloat(cpuNum, 10)
				r := c / 1000
				cpu = strconv.FormatFloat(r, 'g', 2, 64)
			}
		}
		if limits["memory"] != nil {
			mem = limits["memory"].(string)
		}
	}

	var ports = []interface{}{}
	if container["ports"] != nil {
		portsList := container["ports"].([]interface{})
		for _, value := range portsList {
			port := value.(map[string]interface{})
			common.Logger.Debug("container port is: %v", port["containerPort"].(float64))
			containerPort := 0.0
			protocol := ""
			typ := ""
			if port["containerPort"] != nil {
				containerPort = port["containerPort"].(float64)
			}
			if port["protocol"] != nil {
				protocol = port["protocol"].(string)
				typ = protocol
			}
			ports = append(ports, map[string]interface{}{
				"port": containerPort,
				"protocol": protocol,
				"type": typ,
			})
		}
	}
	var envValues = []interface{}{}
	if container["env"] != nil {
		envValues = container["env"].([]interface{})
	}

	data := map[string]interface{}{
		"env": map[string]interface{}{
			"app_name":     appName,
			"app_namespace":  appNamespace,
			"image_name": imageName,
			"tag": tag,
			"env_name": groupname,
			"core": cpu,
			"memory": mem,
			"count": replicas,
			"code": cmd,
			"ips": ips,
		},
		"config": map[string]interface{}{
			"ports": ports,
			"extra_env": envValues,
		},
	}
	return statusCode, data, err
}

func updateServiceStatus(appName, appNamespace, operation string, kubeCmd util.KubeCmd) error {
	var err error = nil
	file := "/tmp/" + appName + ".yaml"
	common.Logger.Info("tmp file is: %v", file)

	if operation == "1" {
		//停止服务操作:设置容器个数为0
		err = kubeCmd.Stop(appName, appNamespace, file)
	} else if operation == "0" {
		err = kubeCmd.Create(appName, appNamespace, file)
	}
	return err
}

func updateServiceQuota(appName, appNamespace, cpu, mem string, fetcher util.KubeResponseFetch, kubeCmd util.KubeCmd) (int, error) {
	statusCode, response, err := fetcher.GetKubeAppContent(appName, appNamespace)

	bytes := []byte(response)
	var res map[string]interface{}

	if err := json.Unmarshal(bytes, &res); err != nil {
		return statusCode, err
	}

	//获取老版本信息
	metadata := res["metadata"].(map[string]interface{})
	labels := metadata["labels"].(map[string]interface{})
	oldVersion := ""
	if labels["version"] != nil {
		oldVersion = labels["version"].(string)
	}

	//获取老版本配额
	spec := res["spec"].(map[string]interface{})
	template := spec["template"].(map[string]interface{})
	tSpec := template["spec"].(map[string]interface{})
	containers := tSpec["containers"].([]interface{})

	//TODO 暂时支持单个容器
	container := containers[0].(map[string]interface{})
	oldCpu := ""
	oldMem := ""
	if container["resources"] != nil {
		resource := container["resources"].(map[string]interface{})
		limits := resource["limits"].(map[string]interface{})
		if limits["cpu"] != nil {
			oldCpu = limits["cpu"].(string)
		}
		if limits["memory"] != nil {
			oldMem = limits["memory"].(string)
		}
	}
	if _, err := kubeCmd.UpdateQuota(appName, appNamespace, oldVersion, oldCpu, oldMem, cpu, mem); err != nil {
		return statusCode, err
	}
	return statusCode, err
}

func updateServiceCmd(appName, appNamespace, newCmd string, fetcher util.KubeResponseFetch, kubeCmd util.KubeCmd) (int, error) {
	statusCode, response, err := fetcher.GetKubeAppContent(appName, appNamespace)

	bytes := []byte(response)
	var res map[string]interface{}

	if err := json.Unmarshal(bytes, &res); err != nil {
		return statusCode, err
	}

	//获取老版本信息
	metadata := res["metadata"].(map[string]interface{})
	labels := metadata["labels"].(map[string]interface{})
	oldVersion := ""
	if labels["version"] != nil {
		oldVersion = labels["version"].(string)
	}

	//获取老版本执行命令
	spec := res["spec"].(map[string]interface{})
	template := spec["template"].(map[string]interface{})
	tSpec := template["spec"].(map[string]interface{})
	containers := tSpec["containers"].([]interface{})

	//TODO 暂时支持单个容器
	container := containers[0].(map[string]interface{})
	oldCmd := ""
	if container["command"] != nil {
		cmdList := container["command"].([]interface{})
		for _, value := range cmdList {
			oldCmd += value.(string) + " "
		}
	}

	if _, err := kubeCmd.UpdateCmd(appName, appNamespace, oldVersion, oldCmd, newCmd); err != nil {
		return statusCode, err
	}
	return statusCode, err
}

func scaleReplicas(appName, appNamespace, replicas string, fetcher util.KubeResponseFetch, kubeCmd util.KubeCmd) (int, error) {
	if err := kubeCmd.Scale(appName, appNamespace, replicas); err != nil {
		return http.StatusInternalServerError, err
	}
	//更新DNS
	strs := strings.Split(appName, "-v")
	appName = strs[0]
	go util.RegisterDNS(appName, appNamespace, kubeCmd, fetcher)
	return http.StatusOK, nil
}

//版本发布
//appName: 应用名字
//appNamespace: 应用命名空间
//image: 镜像
//fetcher: 访问kube-apiserver获取数据
//kubeCmd: 调用kubectl命令
func deliveryRelease(appName, appNamespace, image string, fetcher util.KubeResponseFetch, kubeCmd util.KubeCmd) (int, error) {
	//获取当前服务的版本和镜像名
	statusCode, response, err := fetcher.GetKubeAppContent(appName, appNamespace)

	bytes := []byte(response)
	var res map[string]interface{}

	err = json.Unmarshal(bytes, &res)

	//获取老版本信息
	metadata := res["metadata"].(map[string]interface{})
	labels := metadata["labels"].(map[string]interface{})
	oldVersion := ""
	if labels["version"] != nil {
		oldVersion = labels["version"].(string)
	}

	//获取老版本镜像
	spec := res["spec"].(map[string]interface{})
	template := spec["template"].(map[string]interface{})
	tSpec := template["spec"].(map[string]interface{})
	containers := tSpec["containers"].([]interface{})

	//TODO 暂时支持单个容器
	container := containers[0].(map[string]interface{})
	oldImage := ""
	if container["image"] != nil {
		oldImage = container["image"].(string)
	}

	if image != oldImage {
		//发布版本
		common.Logger.Info("Starting delivery app, and update the old image %s to the new image %s", oldImage, image)
		_, err = kubeCmd.Delivery(appName, appNamespace, image)
	} else if image == oldImage {
		//应用重启
		common.Logger.Info("Starting restart app: %v", appName)
		_, err = kubeCmd.Restart(appName, appNamespace, oldVersion, oldImage)
	}
	return statusCode, err

}

func getClusterIPs(groupname string) ([]string, int, error) {
	ips := []string{}
	postFix := "/api/v1/nodes"
	url := "http://" + kubeApiserverPath + kubeApiserverPort + postFix
	statusCode, response, err := common.SendRawRequest("GET", url, nil)
	common.Logger.Debug("statusCode: %d", statusCode)
	common.Logger.Debug("repsonse: %s", response)


	if statusCode != http.StatusOK {
		return nil, statusCode, err
	}

	byte := []byte(response)
	var res map[string]interface{}

	err = json.Unmarshal(byte, &res)
	if err != nil {
		return nil, statusCode, err
	}

	list := res["items"].([]interface{})

	for _, value := range list {
		item := value.(map[string]interface{})

		metadata := item["metadata"].(map[string]interface{})
		labels := metadata["labels"].(map[string]interface{})
		if labels["groupname"] == nil {continue}
		key := labels["groupname"].(string)

		if key == groupname {
			status := item["status"].(map[string]interface{})
			addresses := status["addresses"].([]interface{})

			for _,value := range addresses {
				item := value.(map[string]interface{})
				if item["type"].(string) == "LegacyHostIP" {
					ips = append(ips, item["address"].(string))
				}
			}
		}

	}
	return ips, http.StatusOK, nil
}