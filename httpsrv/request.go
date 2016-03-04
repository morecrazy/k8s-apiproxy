package httpsrv

import (
	"third/gin"
	"encoding/json"
	. "backend/common"
	"net/http"
	"codoon_ops/kubernetes-apiproxy/util"
	"zhanghan/slb/api/rest"
	"strings"
	"io/ioutil"
	"time"
	"strconv"
	"fmt"
	"errors"
	"codoon_ops/kubernetes-apiproxy/object"
	"third/gorm"
)

var hostPath = "dockerhub.codoon.com"
var hostPortA = ":5000"			//registry服务端口
var hostPortB = ":8080"			//kubernetes apiserver端口
var KubeDb *gorm.DB

type Response struct {
	Status Status `json:"status"`
	Data Data `json:"data"`
}

type StatusResp struct {
	Status Status `json:"status"`
}

type Status struct {
	State int `json:"state"`
	Msg string `json:"msg"`
}

type Data struct {
	List []interface{} `json:"list"`
}

func GetImagesList(c * gin.Context) {
	postFix := "/v1/search"
	url := "http://" + hostPath + hostPortA + postFix
	statusCode, response, err := SendRequest("GET", url, nil, nil, nil)
	Logger.Debug("statusCode: %d", statusCode)
	Logger.Debug("repsonse: %s", response)

	s := new(StatusResp)			//返回状态结构体
	r := new(Response)		//返回结果结构体

	if statusCode != http.StatusOK {
		s.Status.State = 1
		s.Status.Msg = err.Error()
		c.JSON(statusCode, s)
		return
	}

	bytes := []byte(response)
	var res map[string]interface{}

	err = json.Unmarshal(bytes, &res)
	if err != nil {
		s.Status.State = 1
		s.Status.Msg = err.Error()
		c.JSON(http.StatusInternalServerError, s)
	}

	list := res["results"].([]interface{})

	r.Status.State = 0
	r.Status.Msg = "ok"
	ret_list := []interface{}{}
	for _, value := range list {
		item := value.(map[string]interface{})
		name := item["name"].(string)
		mirror := hostPath + "/" + name
		ret_list = append(ret_list, map[string]interface{}{
			"name":     name,
			"mirror":  mirror,
			"latest": "latest",
		})
	}
	r.Data.List = ret_list

	c.JSON(http.StatusOK, r)
}

func DeleteImage(c * gin.Context) {
	body, _ := ioutil.ReadAll(c.Request.Body)

	var requestData map[string]interface{}  //获取request body的数据
	err := json.Unmarshal(body, &requestData)

	//name := c.PostForm("name")
	name := requestData["name"].(string)
	postFix := "/v1/repositories"
	url := "http://" + hostPath + hostPortA + postFix + "/" + name + "/"
	Logger.Debug("url: %v", url)
	statusCode, response, err := SendRequest("DELETE", url, nil, nil, nil)
	Logger.Debug("statusCode: %d", statusCode)
	Logger.Debug("repsonse: %s", response)

	s := new(StatusResp)			//返回状态结构体

	if statusCode != http.StatusOK {
		s.Status.State = 1
		s.Status.Msg = err.Error()
		c.JSON(statusCode, s)
		return
	}

	if response == "true" {
		s.Status.State = 0
		s.Status.Msg = "OK"
		c.JSON(statusCode, s)
	} else {
		s.Status.State = 1
		s.Status.Msg = "Wrong"
		c.JSON(statusCode, s)
	}

}

func GetImageTags(c * gin.Context) {
	name := c.Query("name")

	postFix := "/v1/repositories"
	url := "http://" + hostPath + hostPortA + postFix + "/" + name + "/tags"
	statusCode, response, err := SendRequest("GET", url, nil, nil, nil)
	Logger.Debug("statusCode: %d", statusCode)
	Logger.Debug("repsonse: %s", response)

	s := new(StatusResp)			//返回状态结构体
	r := new(Response)		//返回结果结构体

	if statusCode != http.StatusOK {
		s.Status.State = 1
		s.Status.Msg = err.Error()
		c.JSON(statusCode, s)
		return
	}

	bytes := []byte(response)
	var res map[string]interface{}

	err = json.Unmarshal(bytes, &res)
	if err != nil {
		s.Status.State = 1
		s.Status.Msg = err.Error()
		c.JSON(http.StatusInternalServerError, s)
	}

	r.Status.State = 0
	r.Status.Msg = "ok"
	ret_list := []interface{}{}

	latestVersionId := ""
	latestVersionNumber := 0
	for key,_ := range res {
		Logger.Debug("The tag is: %v", key)
		version := util.Substr(key, 1, len(key))
		versionNumber, err := strconv.Atoi(version)
		if err != nil {
			Logger.Error("invalid tag")
			continue;
		}
		if versionNumber > latestVersionNumber {
			latestVersionNumber = versionNumber
		}
	}

	if latestVersionNumber != 0 {
		latestVersion := "v" + strconv.Itoa(latestVersionNumber)
		latestVersionId = res[latestVersion].(string)
	}
	if res["latest"] != nil {
		latestVersionId = res["latest"].(string)
	}
	for key, value := range res {
		var latest = false
		if value == latestVersionId {
			latest = true
		}
		ret_list = append(ret_list, map[string]interface{}{
			"commit":    value,
			"tag":  key,
			"is_latest": latest,
		})
	}
	r.Data.List = ret_list
	c.JSON(http.StatusOK, r)
}

func DeleteTag(c * gin.Context) {
	body, _ := ioutil.ReadAll(c.Request.Body)

	var requestData map[string]interface{}  //获取request body的数据
	err := json.Unmarshal(body, &requestData)

	name := requestData["name"].(string)
	tag := requestData["tag"].(string)

	postFix := "/v1/repositories"
	url := "http://" + hostPath + hostPortA + postFix + "/" + name + "/tags" + "/" + tag
	Logger.Debug("url: %v", url)
	statusCode, response, err := SendRequest("DELETE", url, nil, nil, nil)
	Logger.Debug("statusCode: %d", statusCode)
	Logger.Debug("repsonse: %s", response)

	s := new(StatusResp)			//返回状态结构体

	if statusCode != http.StatusOK {
		s.Status.State = 1
		s.Status.Msg = err.Error()
		c.JSON(statusCode, s)
		return
	}

	if response == "true" {
		s.Status.State = 0
		s.Status.Msg = "OK"
		c.JSON(statusCode, s)
	} else {
		s.Status.State = 1
		s.Status.Msg = "Wrong"
		c.JSON(statusCode, s)
	}
}

func UpdateImageTag (c * gin.Context) {
	body, _ := ioutil.ReadAll(c.Request.Body)

	var requestData map[string]interface{}  //获取request body的数据
	err := json.Unmarshal(body, &requestData)

	name := requestData["name"].(string)
	commit := requestData["commit"].(string)

	Logger.Debug("name is: %v", name)
	Logger.Debug("commit is: %v", commit)

	commitId := "\"" + commit + "\""
	postFix := "/v1/repositories"
	url := "http://" + hostPath + hostPortA + postFix + "/" + name + "/tags/latest"
	statusCode, response, err := SendRequest("PUT", url, nil, nil, commitId)
	Logger.Debug("statusCode: %d", statusCode)
	Logger.Debug("repsonse: %s", response)

	s := new(StatusResp)			//返回状态结构体

	if statusCode != http.StatusOK {
		s.Status.State = 1
		s.Status.Msg = err.Error()
		c.JSON(statusCode, s)
		return
	}

	if response == "true" {
		s.Status.State = 0
		s.Status.Msg = "OK"
		c.JSON(statusCode, s)
	} else {
		s.Status.State = 1
		s.Status.Msg = "Wrong"
		c.JSON(statusCode, s)
	}
}

func GetClusterList(c * gin.Context) {
	postFix := "/api/v1/nodes"
	url := "http://" + hostPath + hostPortB + postFix
	statusCode, response, err := SendRequest("GET", url, nil, nil, nil)
	Logger.Debug("statusCode: %d", statusCode)
	Logger.Debug("repsonse: %s", response)

	s := new(StatusResp)			//返回状态结构体
	r := new(Response)		//返回结果结构体

	if statusCode != http.StatusOK {
		s.Status.State = 1
		s.Status.Msg = err.Error()
		c.JSON(statusCode, s)
		return
	}

	byte := []byte(response)
	var res map[string]interface{}

	err = json.Unmarshal(byte, &res)
	if err != nil {
		s.Status.State = 1
		s.Status.Msg = err.Error()
		c.JSON(http.StatusInternalServerError, s)
	}

	list := res["items"].([]interface{})

	r.Status.State = 0
	r.Status.Msg = "ok"
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

	r.Data.List = ret_list

	c.JSON(http.StatusOK, r)
}

func GetServicesList(c * gin.Context) {
	postFix := "/api/v1/replicationcontrollers"
	url := "http://" + hostPath + hostPortB + postFix
	statusCode, response, err := SendRequest("GET", url, nil, nil, nil)
	Logger.Debug("statusCode: %d", statusCode)
	Logger.Debug("repsonse: %s", response)

	s := new(StatusResp)
	r := new(Response)

	if statusCode != http.StatusOK {
		s.Status.State = 1
		s.Status.Msg = err.Error()
		c.JSON(statusCode, s)
		return
	}

	bytes := []byte(response)
	var res map[string]interface{}

	err = json.Unmarshal(bytes, &res)
	if err != nil {
		s.Status.State = 1
		s.Status.Msg = err.Error()
		c.JSON(http.StatusInternalServerError, s)
		return
	}

	list := res["items"].([]interface{})

	r.Status.State = 0
	r.Status.Msg = "ok"
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
		cmd := "kubectl describe rc " + name + " --namespace=" + namespace
		bytes, err := util.ExecCommand(cmd)
		if err != nil {
			s.Status.State = 1
			s.Status.Msg = err.Error()
			c.JSON(http.StatusInternalServerError, s)
			return
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

	r.Data.List = ret_list
	c.JSON(http.StatusOK, r)

}

func GetServiceMetadata(c * gin.Context) {
	appName := c.Query("app_name")
	appNamespace := c.Query("app_namespace")

	/**
	rcNamespace := appNamespace
	//获取rc的名字和命令空间:rcName,rcNamespace
	rcName, rcNum, statusCode, err := getReplicationControllerName(appName, appNamespace)
	**/
	s := new(StatusResp)			//返回状态结构体

	var res = map[string]interface{}{}

	postFix := "/api/v1/namespaces"
	url := "http://" + hostPath + hostPortB + postFix + "/" + appNamespace + "/replicationcontrollers" + "/" + appName

	statusCode, response, err := SendRequest("GET", url, nil, nil, nil)
	Logger.Debug("statusCode: %d", statusCode)
	Logger.Debug("repsonse: %s", response)


	if statusCode != http.StatusOK {
		s.Status.State = 1
		s.Status.Msg = err.Error()
		c.JSON(statusCode, s)
		return
	}

	bytes := []byte(response)

	err = json.Unmarshal(bytes, &res)
	if err != nil {
		s.Status.State = 1
		s.Status.Msg = err.Error()
		c.JSON(http.StatusInternalServerError, s)
		return
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
	cmd := "kubectl describe rc " + appName + " --namespace=" + appNamespace
	bytes, err = util.ExecCommand(cmd)
	if err != nil {
		s.Status.State = 1
		s.Status.Msg = err.Error()
		c.JSON(http.StatusInternalServerError, s)
		return
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

	//获取svc name
	svcNameList := []string{}
	cmd = "kubectl get svc -l app=" + svcName + " --namespace=" + appNamespace
	Logger.Debug("The cmd is: %v", cmd)
	bytes, err = util.ExecCommand(cmd)
	if err != nil {
		s.Status.State = 1
		s.Status.Msg = err.Error()
		c.JSON(http.StatusInternalServerError, s)
		return
	}
	lines = strings.Split(string(bytes), "\n")
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
		Logger.Debug("Svc Name is: %v", svcName)
		postFix = "/api/v1/namespaces/" + appNamespace + "/services" + "/" + svcName
		url = "http://" + hostPath + hostPortB + postFix
		statusCode, response, err = SendRequest("GET", url, nil, nil, nil)
		Logger.Debug("statusCode: %d", statusCode)
		Logger.Debug("repsonse: %s", response)

		if statusCode != http.StatusOK {
			s.Status.State = 1
			s.Status.Msg = err.Error()
			c.JSON(statusCode, s)
			return
		}

		bytes = []byte(response)

		err = json.Unmarshal(bytes, &res)
		if err != nil {
			s.Status.State = 1
			s.Status.Msg = err.Error()
			c.JSON(http.StatusInternalServerError, s)
			return
		}

		metadata := res["metadata"].(map[string]interface{})
		portMap := map[string]interface{}{}
		if metadata["annotations"] != nil {
			portMap = metadata["annotations"].(map[string]interface{})
		}

		spec = res["spec"].(map[string]interface{})
		ports := spec["ports"].([]interface{})
		svcType := spec["type"].(string)
		ip := ""
		//内部服务
		if svcType == "ClusterIp" {
			ip = spec["clusterIp"].(string)
		} else if svcType == "NodePort" {
			//外部服务
			//查询数据库,根据集群名字获取slb的ip
			slb := object.SLB{
				GroupName: groupname,
			}
			err = slb.Fetch(KubeDb)
			ip = slb.Ip
			Logger.Debug("SLB ip is: %v", ip)
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
			url := protocol + "://" + ip + ":" + listenerPort
			urls = append(urls, url)
		}
	}

	//拼接返回数据
	st := map[string]interface{}{
		"state": 0,
		"msg": "ok",
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

	r := map[string]interface{}{
		"status": st,
		"data": data,
	}
	c.JSON(http.StatusOK, r)
}

func UpdateServiceStatus(c * gin.Context) {

	body, _ := ioutil.ReadAll(c.Request.Body)

	var requestData map[string]interface{}  //获取request body的数据
	err := json.Unmarshal(body, &requestData)

	appName := requestData["app_name"].(string)
	appNamespace := requestData["app_namespace"].(string)
	operation := requestData["operation"].(string)

	s := new(StatusResp)
	if err != nil {
		s.Status.State = 1
		s.Status.Msg = err.Error()
		c.JSON(http.StatusInternalServerError, s)
	}

	file := "/tmp/" + appName + ".yaml"
	Logger.Debug("tmp file is: %v", file)
	//停止服务操作:设置容器个数为0
	if operation == "1" {
		Logger.Debug("Stopping service: %v", appName)
		//先保存当前配置到特地文件中
		cmd := "kubectl get rc " + appName + " --namespace=" + appNamespace + " -o yaml >" + file
		_, err := util.ExecCommand(cmd)
		if err != nil {
			s.Status.State = 1
			s.Status.Msg = err.Error()
			c.JSON(http.StatusInternalServerError, s)
			return
		}
		//缩减容器个数为0既而停止服务
		cmd = "kubectl scale rc " + appName + " --namespace=" + appNamespace + " --replicas=0"
		_, err = util.ExecCommand(cmd)
		if err != nil {
			s.Status.State = 1
			s.Status.Msg = err.Error()
			c.JSON(http.StatusInternalServerError, s)
			return
		}
	} else if operation == "0" {
		Logger.Debug("Starting service: %v", appName)
		//删除久的rc
		cmd := "kubectl delete rc " + appName + " --namespace=" + appNamespace
		_, err := util.ExecCommand(cmd)
		if err != nil {
			s.Status.State = 1
			s.Status.Msg = err.Error()
			c.JSON(http.StatusInternalServerError, s)
			return
		}
		cmd = "kubectl create -f " + file
		_, err = util.ExecCommand(cmd)
		if err != nil {
			s.Status.State = 1
			s.Status.Msg = err.Error()
			c.JSON(http.StatusInternalServerError, s)
			return
		}
	}

	s.Status.State = 0
	s.Status.Msg = "updated"
	c.JSON(http.StatusOK, s)
}

func GetServiceConfig(c * gin.Context) {
	appName := c.Query("app_name")
	appNamespace := c.Query("app_namespace")

	/**
	rcNamespace := appNamespace
	//获取rc的名字和命名空间
	rcName, _, statusCode, err := getReplicationControllerName(appName, appNamespace)
	**/

	s := new(StatusResp)			//返回状态结构体

	postFix := "/api/v1/namespaces"
	url := "http://" + hostPath + hostPortB + postFix + "/" + appNamespace + "/replicationcontrollers" + "/" + appName

	statusCode, response, err := SendRequest("GET", url, nil, nil, nil)
	Logger.Debug("statusCode: %d", statusCode)
	Logger.Debug("repsonse: %s", response)

	if statusCode != http.StatusOK {
		s.Status.State = 1
		s.Status.Msg = err.Error()
		c.JSON(statusCode, s)
		return
	}

	bytes := []byte(response)
	var res map[string]interface{}

	err = json.Unmarshal(bytes, &res)
	if err != nil {
		s.Status.State = 1
		s.Status.Msg = err.Error()
		c.JSON(http.StatusInternalServerError, s)
		return
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
	ips, code, err := getClusterIPs(groupname)
	if err != nil {
		s.Status.State = 1
		s.Status.Msg = err.Error()
		c.JSON(code, s)
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
			Logger.Debug("container port is: %v", port["containerPort"].(float64))
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

	//拼接返回数据
	state := map[string]interface{}{
		"state": 0,
		"msg": "ok",
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

	r := map[string]interface{}{
		"status": state,
		"data": data,
	}
	c.JSON(http.StatusOK, r)

}

/**
创建服务: 包括创建svc和rc两个部分.首先创建svc,然后创建rc
 */
func CreateService(c * gin.Context) {
	body, _ := ioutil.ReadAll(c.Request.Body)

	var res map[string]interface{}
	err := json.Unmarshal(body, &res)

	s := new(StatusResp)			//返回状态结构体

	if err != nil {
		s.Status.State = 1
		s.Status.Msg = err.Error()
		c.JSON(http.StatusInternalServerError, s)
		return
	}

	//获取基本环境数据
	env := res["env"].(map[string]interface{})

	imageName := env["name"].(string)
	imageTag := env["tag"].(string)
	//拼接镜像的完整路径
	imagePath:= hostPath + "/" + imageName + ":" + imageTag
	appName := env["app_name"].(string)
	appNamespace := env["app_namespace"].(string)
	group := env["env_name"].(string) //集群名(组名)
	cpu := env["core"].(string)		//cpu限制
	mem := env["memory"].(string)	//内存限制
	replicas := env["count"].(float64)	//容器个数

	command :=[]string{}
	if env["code"].(string) != "" {
		command = append(command, env["code"].(string))
	}//执行命命令

	//获取配置数据
	config := res["config"].(map[string]interface{})

	ports := config["ports"].([]interface{})
	envValues := config["extra_env"].([]interface{})

	innerSvcPortsList := []interface{}{}	//内部svc端口列表
	outterSvcPortsList := []interface{}{}	//外部svc端口列表
	//svcPortsList := []interface{}{}  //svc端口列表,包含端口的协议和访问方式
	rcPortsList := []interface{}{}	//rc端口列表
	innerPortMap := map[string]string{} //内部port->访问类型映射
	outterPortMap := map[string]string{} //外部port->访问类型映射

	innerPortNum := 0
	outterPortNum := 0

	svcType := "ClusterIP"

	for _,value := range ports {
		item := value.(map[string]interface{})
		port := item["port"].(float64)
		name := strconv.FormatFloat(port, 'g', 5, 64)
		//如果是外部端口,则添加到外部svc端口列表
		if item["access"].(string) == "outter" || item["access"].(string) == "all" {
			outterSvcPortsList = append(outterSvcPortsList, map[string]interface{}{
				"name": name,
				"port": port,
				"targetPort": port,
				"nodePort": port,
				"protocol": 	item["protocol"].(string),
			})
			outterPortMap[name] = item["type"].(string)
			outterPortNum ++
		} else if item["access"].(string) == "inner" {
			innerSvcPortsList = append(innerSvcPortsList, map[string]interface{}{
				"name": name,
				"port": port,
				"targetPort": port,
				"protocol": 	item["protocol"].(string),
			})
			innerPortMap[name] = item["type"].(string)
			innerPortNum ++
		}
		rcPortsList = append(rcPortsList, map[string]interface{}{
			"containerPort": port,
			"protocol": item["protocol"].(string),
		})
	}
	/**
	//如果全是内部端口,还是需要创建svc,但是nodetype为普通类型
	if outterPortsNum == 0 {
		for _, value := range ports {
			item := value.(map[string]interface{})
			port := item["port"].(float64)
			name := strconv.FormatFloat(port, 'g', 5, 64)
			svcPortsList = append(svcPortsList, map[string]interface{}{
				"name":     name,
				"port":     port,
				"targetPort": 		port,
				"protocol": 		item["protocol"].(string),
			})
			portMap[port] = item["type"].(string)
		}
	} else {
		//既有内部端口,也有外部端口.则创建svc时只需要包含外部端口且设置nodetype类型为NodePort
		for _, value := range ports {
			svcType = "NodePort"
			item := value.(map[string]interface{})
			port := item["port"].(float64)
			name := strconv.FormatFloat(port, 'g', 5, 64)
			if item["access"].(string) == "outter" || item["access"].(string) == "all" {
				svcPortsList = append(svcPortsList, map[string]interface{}{
					"name":    name,
					"port":     port,
					"targetPort": 		item["port"].(float64),
					"protocol": 		item["protocol"].(string),
				})
				portMap[port] = item["type"].(string)
			}
		}
	}
	**/

	envValuesList := []interface{}{}	//环境变量列表
	for _,value := range envValues {
		item := value.(map[string]interface{})
		envValuesList = append(envValuesList, map[string]interface{}{
			"name":     item["name"].(string),
			"value":	item["value"].(string),
		})
	}

	/**
	首先创建svc
	 */

	svcName := appName + "-inner"
	//请求body(json格式)
	if innerPortNum > 0 {
		svcRequestJson := map[string]interface{}{
			"apiVersion": "v1",
			"kind": "Service",
			"metadata": map[string]interface{}{
				"name": svcName,
				"namespace": appNamespace,
				"labels": map[string]interface{}{
					"app": appName,
				},
				"annotations": innerPortMap,
			},
			"spec": map[string]interface{}{
				"type": svcType,
				"ports": innerSvcPortsList,
				"selector": map[string]interface{}{
					"app": appName,
				},
			},
		}

		b, _ := json.Marshal(svcRequestJson)
		Logger.Debug("svc request json is: %v", string(b))

		postFix := "/api/v1/namespaces/" + appNamespace + "/services"
		url := "http://" + hostPath + hostPortB + postFix
		statusCode, response, err := SendRequest("POST", url, svcRequestJson, nil, nil)
		Logger.Debug("statusCode: %d", statusCode)
		Logger.Debug("repsonse: %s", response)

		if statusCode != http.StatusCreated {
			s.Status.State = 1
			s.Status.Msg = err.Error()
			c.JSON(statusCode, s)
			return
		}
	}
	svcName = appName + "-outter"
	if outterPortNum > 0 {
		svcType = "NodePort"
		svcRequestJson := map[string]interface{}{
			"apiVersion": "v1",
			"kind": "Service",
			"metadata": map[string]interface{}{
				"name": svcName,
				"namespace": appNamespace,
				"labels": map[string]interface{}{
					"app": appName,
				},
				"annotations": outterPortMap,
			},
			"spec": map[string]interface{}{
				"type": svcType,
				"ports": outterSvcPortsList,
				"selector": map[string]interface{}{
					"app": appName,
				},
			},
		}

		b, _ := json.Marshal(svcRequestJson)
		Logger.Debug("svc request json is: %v", string(b))

		postFix := "/api/v1/namespaces/" + appNamespace + "/services"
		url := "http://" + hostPath + hostPortB + postFix
		statusCode, response, err := SendRequest("POST", url, svcRequestJson, nil, nil)
		Logger.Debug("statusCode: %d", statusCode)
		Logger.Debug("repsonse: %s", response)

		if statusCode != http.StatusCreated {
			s.Status.State = 1
			s.Status.Msg = err.Error()
			c.JSON(statusCode, s)
			return
		}
	}

	/*
	然后创建rc
	*/
	rcRequestJson := map[string]interface{}{
		"kind": "ReplicationController",
		"apiVersion": "v1",
		"metadata": map[string]interface{}{
			"name": appName,
			"namespace": appNamespace,
		},
		"spec": map[string]interface{}{
			"replicas": replicas,
			"selector": map[string]interface{}{
				"app": appName,
				"version": "v1",
			},
			"template": map[string]interface{}{
				"metadata": map[string]interface{}{
					"labels": map[string]interface{}{
						"app": appName,
						"version": "v1",
					},
				},
				"spec": map[string]interface{}{
					"containers": []interface{}{
						//TODO :需要支持多容器
						map[string]interface{}{
							"name": appName,
							"image": imagePath,
							"command": command,
							"ports": rcPortsList,
							"env": envValuesList,
							"resources": map[string]interface{}{
								"limits": map[string]interface{}{
									"cpu": cpu,
									"memory": mem,
								},
							},
							"volumeMounts": []interface{}{
								map[string]interface{}{
									"mountPath": "/var/log/go_log",
									"name": "log",
								},
								map[string]interface{}{
									"mountPath": "/etc/localtime",
									"name": "time-zone",
								},
							},
							"imagePullPolicy": "Always",
						},
					},
					"volumes": []interface{}{
						map[string]interface{}{
							"name": "log",
							"hostPath": map[string]interface{}{
								"path": "/var/log/go_log",
							},
						},
						map[string]interface{}{
							"name": "time-zone",
							"hostPath": map[string]interface{}{
								"path": "/etc/localtime",
							},
						},
					},
					"imagePullSecrets": []interface{}{
						map[string]interface{}{
							"name": "dockerhub.codoon.com.key",
						},
					},
					"nodeSelector": map[string]interface{}{
						"groupname": group,
					},
				},
			},
		},
	}

	b, _ := json.Marshal(rcRequestJson)
	Logger.Debug("rc request json is: %v", string(b))

	postFix := "/api/v1/namespaces/" + appNamespace + "/replicationcontrollers"
	url := "http://" + hostPath + hostPortB + postFix
	statusCode, response, err := SendRequest("POST", url, rcRequestJson, nil, nil)
	Logger.Debug("statusCode: %d", statusCode)
	Logger.Debug("repsonse: %s", response)

	if statusCode != http.StatusCreated {
		s.Status.State = 1
		s.Status.Msg = err.Error()
		c.JSON(statusCode, s)
		return
	}

	/*
	创建slb listener
	1. 调用kubernetes的api,获取svc的port list,包括listenerport, backendport, protocol
	1 .如果svcType == "ClusterIp",说明是内部服务,则不需要创建slb listener,直接返回svc的IP和port即可
	2. 如果svcType 等于 "NodePort",说明是外部服务,则需要
	   * 创建slb listener,根据groupname 查询数据库获取slb Id
	   * 遍历svc port list,依次创建listener,
	   	  * 返回slb的外部ip和listenerPort即可
	*/

	/**
	*外部服务需要创建slb listener
	if outterPortNum > 0 {
		slb := object.SLB{
			GroupName: group,
		}
		err = slb.Fetch(KubeDb)
		slbId := slb.LoadBalancerId

		postFix := "/api/v1/namespaces/" + appNamespace + "/services" + "/" + svcName
		url := "http://" + hostPath + hostPortB + postFix
		statusCode, response, err := SendRequest("GET", url, nil, nil, nil)
		Logger.Debug("statusCode: %d", statusCode)
		Logger.Debug("repsonse: %s", response)

		if statusCode != http.StatusOK {
			s.Status.State = 1
			s.Status.Msg = err.Error()
			c.JSON(statusCode, s)
			return
		}
		bytes := []byte(response)

		err = json.Unmarshal(bytes, &res)
		if err != nil {
			s.Status.State = 1
			s.Status.Msg = err.Error()
			c.JSON(http.StatusInternalServerError, s)
			return
		}

		spec := res["spec"].(map[string]interface{})
		ports := spec["ports"].([]interface{})

		//遍历outter svc port list
		for _, value := range ports {
			item := value.(map[string]interface{})

			name := item["name"].(string)
			port := item["port"].(float64)
			nodePort := item["nodePort"].(float64)
			protocol := outterPortMap[name]		//根据port获取访问协议

			listenerPort := int(port)
			backendPort := int(nodePort)

			//创建slb listener
			b, err := createListener(listenerPort, backendPort, slbId, protocol)
			if b != true {
				r := StatusResp{}

				r.Status.State = 1
				r.Status.Msg = err.Error()
				c.JSON(http.StatusInternalServerError, r)
			}

			b, err = startListener(listenerPort, slbId)
			if b != true {
				r := StatusResp{}

				r.Status.State = 1
				r.Status.Msg = err.Error()
				c.JSON(http.StatusInternalServerError, r)
			}
		}
	}
	*/
	r := StatusResp{}

	r.Status.State = 0
	r.Status.Msg = "successful created!"
	c.JSON(statusCode, r)
}

func UpdateServiceQuota(c *gin.Context) {
	body, _ := ioutil.ReadAll(c.Request.Body)

	var requestData map[string]interface{}  //获取request body的数据
	err := json.Unmarshal(body, &requestData)

	appName := requestData["app_name"].(string)
	appNamespace := requestData["app_namespace"].(string)
	cpu := requestData["core"].(string)
	mem := requestData["memory"].(string)

	s := new(StatusResp)
	if err != nil {
		s.Status.State = 1
		s.Status.Msg = err.Error()
		c.JSON(http.StatusInternalServerError, s)
	}

	//获取当前服务的版本和配额
	postFix := "/api/v1/namespaces"
	url := "http://" + hostPath + hostPortB + postFix + "/" + appNamespace + "/replicationcontrollers" + "/" + appName

	statusCode, response, err := SendRequest("GET", url, nil, nil, nil)
	Logger.Debug("statusCode: %d", statusCode)
	Logger.Debug("repsonse: %s", response)

	if statusCode != http.StatusOK {
		s.Status.State = 1
		s.Status.Msg = err.Error()
		c.JSON(statusCode, s)
		return
	}

	bytes := []byte(response)
	var res map[string]interface{}

	err = json.Unmarshal(bytes, &res)
	if err != nil {
		s.Status.State = 1
		s.Status.Msg = err.Error()
		c.JSON(http.StatusInternalServerError, s)
		return
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

	//生成新版本号,以日期作为标致
	t := time.Now()
	newVersion := t.Format("20060102150405")
	newVersion = "v" + newVersion

	newAppName := ""
	strs := strings.Split(appName, "-")
	for i :=0; i < len(strs) -1; i++ {
		newAppName += strs[i] + "-"
	}
	newAppName = newAppName + newVersion
	Logger.Debug("New App Name is: %v", newAppName)

	//滚动更新
	cmd := "kubectl get rc " + appName + " --namespace=" + appNamespace + " -o yaml " +
	        " | sed 's/resourceVersion:.*/resourceVersion: ''/g' " +
			" | sed 's/name: " + appName + "/name: " + newAppName + "/g' " +   //替换rc名字
			" | sed 's/version: " + oldVersion + "/version: " + newVersion + "/g' " +		  //替换rc版本
			" | sed 's/cpu: " + oldCpu + "/cpu: " + cpu + "/g' " + 							 //替换cpu
			" | sed 's/memory: " + oldMem + "/memory: " + mem + "/g' " + 					//替换mem
			" | kubectl rolling-update " + appName + " --update-period=20s --namespace=" + appNamespace + " -f - "								//滚动更新

	Logger.Debug("The cmd is: %v", cmd)
	go util.ExecCommand(cmd)
	/**
	if err != nil {
		s.State = 1
		s.Msg = err.Error()
		c.JSON(http.StatusInternalServerError, s)
		return
	}**/
	s.Status.State = 0
	s.Status.Msg = "updated"
	c.JSON(http.StatusOK, s)
}

func UpdateServiceCmd(c *gin.Context) {
	body, _ := ioutil.ReadAll(c.Request.Body)

	var requestData map[string]interface{}  //获取request body的数据
	err := json.Unmarshal(body, &requestData)

	appName := requestData["app_name"].(string)
	appNamespace := requestData["app_namespace"].(string)
	newCmd := requestData["command"].(string)

	s := new(StatusResp)
	if err != nil {
		s.Status.State = 1
		s.Status.Msg = err.Error()
		c.JSON(http.StatusInternalServerError, s)
	}

	//获取当前服务的版本和配额
	postFix := "/api/v1/namespaces"
	url := "http://" + hostPath + hostPortB + postFix + "/" + appNamespace + "/replicationcontrollers" + "/" + appName

	statusCode, response, err := SendRequest("GET", url, nil, nil, nil)
	Logger.Debug("statusCode: %d", statusCode)
	Logger.Debug("repsonse: %s", response)

	if statusCode != http.StatusOK {
		s.Status.State = 1
		s.Status.Msg = err.Error()
		c.JSON(statusCode, s)
		return
	}

	bytes := []byte(response)
	var res map[string]interface{}

	err = json.Unmarshal(bytes, &res)
	if err != nil {
		s.Status.State = 1
		s.Status.Msg = err.Error()
		c.JSON(http.StatusInternalServerError, s)
		return
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

	//生成新版本号,以日期作为标致
	t := time.Now()
	newVersion := t.Format("20060102150405")
	newVersion = "v" + newVersion

	newAppName := ""
	strs := strings.Split(appName, "-")
	for i :=0; i < len(strs) -1; i++ {
		newAppName += strs[i] + "-"
	}
	newAppName = newAppName + newVersion
	Logger.Debug("New App Name is: %v", newAppName)

	//滚动更新
	cmd := "kubectl get rc " + appName + " --namespace=" + appNamespace + " -o yaml " +
			" | sed 's/resourceVersion:.*/resourceVersion: ''/g' " +
			" | sed 's/name: " + appName + "/name: " + newAppName + "/g " +   //替换rc名字
			" | sed 's/version: " + oldVersion + "/version: " + newVersion + "/g " +		  //替换rc版本
			" | sed 's/command: " + oldCmd + "/cmd: " + newCmd + "/g " + 					//替换mem
	        " | kubectl rolling-update " + appName + " --update-period=20s --namespace=" + appNamespace + " -f - "

	go util.ExecCommand(cmd)
	/**
	if err != nil {
		s.State = 1
		s.Msg = err.Error()
		c.JSON(http.StatusInternalServerError, s)
		return
	}
	**/
	s.Status.State = 0
	s.Status.Msg = "updated"
	c.JSON(http.StatusOK, s)
}

func ScaleReplicas(c *gin.Context) {
	body, _ := ioutil.ReadAll(c.Request.Body)

	var requestData map[string]interface{}  //获取request body的数据
	err := json.Unmarshal(body, &requestData)

	appName := requestData["app_name"].(string)
	appNamespace := requestData["app_namespace"].(string)
	count := requestData["count"].(float64)
	counts := uint64(count)
	replicas := strconv.FormatUint(counts, 10)

	s := new(StatusResp)
	if err != nil {
		s.Status.State = 1
		s.Status.Msg = err.Error()
		c.JSON(http.StatusInternalServerError, s)
	}

	cmd := "kubectl scale rc " + appName + " --namespace=" + appNamespace + " --replicas=" + replicas
	Logger.Debug("The cmd is: %v", cmd)

	_, err = util.ExecCommand(cmd)
	if err != nil {
		s.Status.State = 1
		s.Status.Msg = err.Error()
		c.JSON(http.StatusInternalServerError, s)
		return
	}
	s.Status.State = 0
	s.Status.Msg = "updated"
	c.JSON(http.StatusOK, s)

}

func DeliveryRelease(c *gin.Context) {
	body, _ := ioutil.ReadAll(c.Request.Body)

	var requestData map[string]interface{}  //获取request body的数据
	err := json.Unmarshal(body, &requestData)

	appName := requestData["app_name"].(string)
	appNamespace := requestData["app_namespace"].(string)
	imageName := requestData["image_name"].(string)
	tag := requestData["tag"].(string)

	s := new(StatusResp)
	if err != nil {
		s.Status.State = 1
		s.Status.Msg = err.Error()
		c.JSON(http.StatusInternalServerError, s)
	}

	image := hostPath + "/" + imageName + ":" + tag
	cmd := "kubectl rolling-update " + appName + " --update-period=20s --namespace=" + appNamespace + " --image=" + image
	Logger.Debug("The cmd is: %v", cmd)

	go util.ExecCommand(cmd)
	/**
	_, err = util.ExecCommand(cmd)
	if err != nil {
		s.State = 1
		s.Msg = err.Error()
		c.JSON(http.StatusInternalServerError, s)
		return
	}
	**/
	s.Status.State = 0
	s.Status.Msg = "updated"
	c.JSON(http.StatusOK, s)
}

/**
返回值:
1: replicationController名字
2: replicationCOntroller个数
 */
func getReplicationControllerName(name, namespace string) (string, int, int, error) {
	postFix := "/api/v1/namespaces"
	url := "http://" + hostPath + hostPortB + postFix + "/" + namespace + "/services" + "/" + name
	statusCode, response, err := SendRequest("GET", url, nil, nil, nil)
	Logger.Debug("statusCode: %d", statusCode)
	Logger.Debug("repsonse: %s", response)

	if statusCode != http.StatusOK {
		return "", 0, statusCode, err
	}

	bytes := []byte(response)
	var res map[string]interface{}

	err = json.Unmarshal(bytes, &res)
	if err != nil {
		return "", 0, http.StatusInternalServerError, err
	}

	//获取selector: app=XXX.以此作为筛选rc的标记
	selector := ""
	spec := res["spec"].(map[string]interface{})
	if spec["selector"] != nil {
		selectorMap := spec["selector"].(map[string]interface{})
		for key, _ := range selectorMap {
			selector += " -l " + key + "=" + selectorMap[key].(string)
		}
	}
	/**
	获取rc name
	 */
	cmd := "kubectl get rc " + selector + " --namespace=" + namespace
	bytes, err = util.ExecCommand(cmd)
	if err != nil {
		return "", 0, http.StatusInternalServerError, err
	}
	lines := strings.Split(string(bytes), "\n")
	rcName := ""
	rcNum := 0
	for _, line := range lines {
		cols := strings.Split(line, " ")
		Logger.Debug("rcName is: %v", cols[0])
		if cols[0] != "" {
			rcName = cols[0]
			rcNum++
		}
	}

	Logger.Debug("the latest rc name is: %v", rcName)

	return rcName,rcNum, http.StatusOK, nil
}

func createListener(listenerPort, backendPort int, loadBalancerId, protocol string) (bool, error) {
	Logger.Info("Begin creating slb listener: listenerPort->%s, backendPort->%s, loadBalancerId->%s, protocl->%s",
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

	Logger.Info("the create result is : %s", res)
	var dat map[string]interface{}
	json.Unmarshal([]byte(res), &dat)

	if dat["Code"] != nil {
		return false, errors.New(dat["Message"].(string))
	}
	return true, nil
}

func startListener(listenerPort int, loadBalancerId string) (bool, error) {
	Logger.Debug("Begin starting listener")
	startApi := rest.NewStartLoadBalancerListenerRequest("http://slb.aliyuncs.com/", listenerPort, loadBalancerId)
	res := startApi.GetResponse("", "60")
	Logger.Debug("the result is : %s", res)

	var dat map[string]interface{}
	json.Unmarshal([]byte(res), &dat)

	if dat["Code"] != nil {
		return false, errors.New(dat["Message"].(string))
	}

	return true, nil
}


func InitDBPool(setting MysqlConfig) error {
	var err error
	KubeDb, err = InitGormDbPool(&setting, true)
	if nil != err {
		Logger.Error("InitGormDbPool err :%v,%v", setting, err)
		return err
	}

	return err
}

func getClusterIPs(groupname string) ([]string, int, error) {
	ips := []string{}
	postFix := "/api/v1/nodes"
	url := "http://" + hostPath + hostPortB + postFix
	statusCode, response, err := SendRequest("GET", url, nil, nil, nil)
	Logger.Debug("statusCode: %d", statusCode)
	Logger.Debug("repsonse: %s", response)


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