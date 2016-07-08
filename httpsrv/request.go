package httpsrv

import (
	"third/gin"
	"encoding/json"
	"net/http"
	"io/ioutil"
	"strconv"
	"third/gorm"
	"codoon_ops/kubernetes-apiproxy/util"
)

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
	s := new(StatusResp)			//返回状态结构体
	r := new(Response)		//返回结果结构体

	//通过registry接口获取数据信息
	fetcher := new(util.RegistryResponseFetcher)
	statusCode, dataList, err := getImagesList(fetcher)

	if statusCode != http.StatusOK {
		s.Status.State = 1
		s.Status.Msg = err.Error()
		c.JSON(statusCode, s)
		return
	}
	if err != nil {
		s.Status.State = 1
		s.Status.Msg = err.Error()
		c.JSON(http.StatusInternalServerError, s)
	}

	r.Status.State = 0
	r.Status.Msg = "ok"
	r.Data.List = dataList
	c.JSON(http.StatusOK, r)
}

func DeleteImage(c * gin.Context) {
	s := new(StatusResp)			//返回状态结构体
	body, _ := ioutil.ReadAll(c.Request.Body)

	var requestData map[string]interface{}  //获取request body的数据
	err := json.Unmarshal(body, &requestData)
	name := requestData["name"].(string)

	fetcher := new(util.RegistryResponseFetcher)
	statusCode, response, err := deleteImage(name, fetcher)

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
	s := new(StatusResp)			//返回状态结构体
	r := new(Response)		//返回结果结构体

	//通过registry接口获取数据信息
	fetcher := new(util.RegistryResponseFetcher)
	statusCode, dataList, err := getImageTags(name, fetcher)

	if statusCode != http.StatusOK {
		s.Status.State = 1
		s.Status.Msg = err.Error()
		c.JSON(statusCode, s)
		return
	}
	if err != nil {
		s.Status.State = 1
		s.Status.Msg = err.Error()
		c.JSON(http.StatusInternalServerError, s)
	}

	r.Status.State = 0
	r.Status.Msg = "ok"
	r.Data.List = dataList
	c.JSON(http.StatusOK, r)
}

func DeleteTag(c * gin.Context) {
	s := new(StatusResp)			//返回状态结构体
	body, _ := ioutil.ReadAll(c.Request.Body)

	var requestData map[string]interface{}  //获取request body的数据
	err := json.Unmarshal(body, &requestData)

	name := requestData["name"].(string)
	tag := requestData["tag"].(string)

	//获取digest
	fetcher := new(util.RegistryResponseFetcher)
	statusCode, _, err := deleteTag(name, tag, fetcher)

	if statusCode != http.StatusOK {
		s.Status.State = 1
		s.Status.Msg = err.Error()
		c.JSON(statusCode, s)
		return
	}

	if err != nil {
		s.Status.State = 1
		s.Status.Msg = err.Error()
		c.JSON(http.StatusInternalServerError, s)
	}

	s.Status.State = 0
	s.Status.Msg = "OK"
	c.JSON(statusCode, s)
}

func UpdateImageTag (c * gin.Context) {
	s := new(StatusResp)			//返回状态结构体
	body, _ := ioutil.ReadAll(c.Request.Body)

	var requestData map[string]interface{}  //获取request body的数据
	if err := json.Unmarshal(body, &requestData); err != nil {
		s.Status.State = 1
		s.Status.Msg = err.Error()
		c.JSON(http.StatusInternalServerError, s)
	}

	/**
	name := requestData["name"].(string)
	commit := requestData["commit"].(string)
	**/
	s.Status.State = 0
	s.Status.Msg = "OK"
	c.JSON(http.StatusOK, s)

}

func GetClusterList(c * gin.Context) {
	s := new(StatusResp)			//返回状态结构体
	r := new(Response)		//返回结果结构体

	fetcher := new(util.KubeResponseFetcher)
	statusCode, dataList, err := getClusterList(fetcher)

	if statusCode != http.StatusOK {
		s.Status.State = 1
		s.Status.Msg = err.Error()
		c.JSON(statusCode, s)
		return
	}
	if err != nil {
		s.Status.State = 1
		s.Status.Msg = err.Error()
		c.JSON(http.StatusInternalServerError, s)
		return
	}

	r.Status.State = 0
	r.Status.Msg = "ok"
	r.Data.List = dataList
	c.JSON(http.StatusOK, r)
}

func GetServicesList(c * gin.Context) {
	s := new(StatusResp)
	r := new(Response)

	fetcher := new(util.KubeResponseFetcher)
	statusCode, dataList, err := getServiceList(fetcher)

	if statusCode != http.StatusOK {
		s.Status.State = 1
		s.Status.Msg = err.Error()
		c.JSON(statusCode, s)
		return
	}
	if err != nil {
		s.Status.State = 1
		s.Status.Msg = err.Error()
		c.JSON(http.StatusInternalServerError, s)
		return
	}

	r.Status.State = 0
	r.Status.Msg = "ok"
	r.Data.List = dataList
	c.JSON(http.StatusOK, r)
}

func GetServiceMetadata(c * gin.Context) {
	appName := c.Query("app_name")
	appNamespace := c.Query("app_namespace")

	s := new(StatusResp)			//返回状态结构体

	fetcher := new(util.KubeResponseFetcher)
	statusCode, data, err := getServiceMeta(appName, appNamespace, fetcher)
	if statusCode != http.StatusOK {
		s.Status.State = 1
		s.Status.Msg = err.Error()
		c.JSON(statusCode, s)
		return
	}
	if err != nil {
		s.Status.State = 1
		s.Status.Msg = err.Error()
		c.JSON(http.StatusInternalServerError, s)
		return
	}

	//拼接返回数据
	st := map[string]interface{}{
		"state": 0,
		"msg": "ok",
	}
	r := map[string]interface{}{
		"status": st,
		"data": data,
	}
	c.JSON(http.StatusOK, r)
}

func UpdateServiceStatus(c * gin.Context) {
	s := new(StatusResp)
	body, _ := ioutil.ReadAll(c.Request.Body)

	var requestData map[string]interface{}  //获取request body的数据
	if err := json.Unmarshal(body, &requestData); err != nil {
		s.Status.State = 1
		s.Status.Msg = err.Error()
		c.JSON(http.StatusInternalServerError, s)
	}

	appName := requestData["app_name"].(string)
	appNamespace := requestData["app_namespace"].(string)
	operation := requestData["operation"].(string)

	kubeCmd := new(util.KubeCmdImpl)
	if err := updateServiceStatus(appName, appNamespace, operation, kubeCmd); err != nil {
		s.Status.State = 1
		s.Status.Msg = err.Error()
		c.JSON(http.StatusInternalServerError, s)
	}

	s.Status.State = 0
	s.Status.Msg = "updated"
	c.JSON(http.StatusOK, s)
}

func GetServiceConfig(c * gin.Context) {
	appName := c.Query("app_name")
	appNamespace := c.Query("app_namespace")

	s := new(StatusResp)			//返回状态结构体

	fetcher := new(util.KubeResponseFetcher)
	statusCode, data, err := getServiceConfig(appName, appNamespace, fetcher)

	if statusCode != http.StatusOK {
		s.Status.State = 1
		s.Status.Msg = err.Error()
		c.JSON(statusCode, s)
		return
	}
	if err != nil {
		s.Status.State = 1
		s.Status.Msg = err.Error()
		c.JSON(http.StatusInternalServerError, s)
		return
	}

	//拼接返回数据
	state := map[string]interface{}{
		"state": 0,
		"msg": "ok",
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
	imagePath:= registryPath + "/" + imageName + ":" + imageTag
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

		fetcher := new(util.KubeResponseFetcher)
		statusCode, _, err := fetcher.CreateSvc(appNamespace, svcRequestJson)
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

		fetcher := new(util.KubeResponseFetcher)
		statusCode, _, err := fetcher.CreateSvc(appNamespace, svcRequestJson)
		if statusCode != http.StatusCreated {
			s.Status.State = 1
			s.Status.Msg = err.Error()
			c.JSON(statusCode, s)
			return
		}
	}

	//创建rc
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
								"requests": map[string]interface{}{
									"cpu": "250m",
									"memory": "250Mi",
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
								map[string]interface{}{
									"mountPath": "/etc/hosts",
									"name": "hosts",
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
						map[string]interface{}{
							"name": "hosts",
							"hostPath": map[string]interface{}{
								"path": "/etc/hosts",
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

	fetcher := new(util.KubeResponseFetcher)
	statusCode, _, err := fetcher.CreateRc(appNamespace, rcRequestJson)

	if statusCode != http.StatusCreated {
		s.Status.State = 1
		s.Status.Msg = err.Error()
		c.JSON(statusCode, s)
		return
	}

	//注册服务
	//两种方式:注册到slb或者注册到DNS.当前版本采用注册到DNS的方法

	//注册到DNS
	//if outterPortNum > 0 {
	kubeCmd := new(util.KubeCmdImpl)
	go util.RegisterDNS(appName, appNamespace, kubeCmd, fetcher)
	/**
	statusCode, err = registerDNS(appName, appNamespace, kubeCmd, httpFetcher)
	if statusCode != http.StatusOK {
		s.Status.State = 1
		s.Status.Msg = err.Error()
		c.JSON(statusCode, s)
		return
	}
	**/
	//}

	//创建slb listener
	//1. 调用kubernetes的api,获取svc的port list,包括listenerport, backendport, protocol
	//2 .如果svcType == "ClusterIp",说明是内部服务,则不需要创建slb listener,直接返回svc的IP和port即可
	//3. 如果svcType 等于 "NodePort",说明是外部服务,则需要
	//   * 创建slb listener,根据groupname 查询数据库获取slb Id
	//   * 遍历svc port list,依次创建listener,
	//   	  * 返回slb的外部ip和listenerPort即可

	//注册到SLB
	/*外部服务需要创建slb listener
	if outterPortNum > 0 {
		slb := object.SLB{
			GroupName: group,
		}
		err = slb.Fetch(KubeDb)
		slbId := slb.LoadBalancerId

		postFix := "/api/v1/namespaces/" + appNamespace + "/services" + "/" + svcName
		url := "http://" + kubeApiserverPath + kubeApiserverPort + postFix
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

//更改服务资源配额
func UpdateServiceQuota(c *gin.Context) {
	body, _ := ioutil.ReadAll(c.Request.Body)

	s := new(StatusResp)
	var requestData map[string]interface{}  //获取request body的数据
	if err := json.Unmarshal(body, &requestData); err != nil {
		s.Status.State = 1
		s.Status.Msg = err.Error()
		c.JSON(http.StatusInternalServerError, s)
	}

	appName := requestData["app_name"].(string)
	appNamespace := requestData["app_namespace"].(string)
	cpu := requestData["core"].(string)
	mem := requestData["memory"].(string)

	fetcher := new(util.KubeResponseFetcher)
	kubeCmd := new(util.KubeCmdImpl)
	statusCode, err := updateServiceQuota(appName, appNamespace, cpu, mem, fetcher, kubeCmd)

	if statusCode != http.StatusOK {
		s.Status.State = 1
		s.Status.Msg = err.Error()
		c.JSON(statusCode, s)
		return
	}
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

//更新服务执行的命令
func UpdateServiceCmd(c *gin.Context) {
	body, _ := ioutil.ReadAll(c.Request.Body)
	s := new(StatusResp)
	var requestData map[string]interface{}  //获取request body的数据
	if err := json.Unmarshal(body, &requestData); err != nil {
		s.Status.State = 1
		s.Status.Msg = err.Error()
		c.JSON(http.StatusInternalServerError, s)
	}

	appName := requestData["app_name"].(string)
	appNamespace := requestData["app_namespace"].(string)
	newCmd := requestData["command"].(string)

	fetcher := new(util.KubeResponseFetcher)
	kubeCmd := new(util.KubeCmdImpl)
	statusCode, err := updateServiceCmd(appName, appNamespace, newCmd, fetcher, kubeCmd)

	if statusCode != http.StatusOK {
		s.Status.State = 1
		s.Status.Msg = err.Error()
		c.JSON(statusCode, s)
		return
	}
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

func ScaleReplicas(c *gin.Context) {
	body, _ := ioutil.ReadAll(c.Request.Body)
	s := new(StatusResp)
	var requestData map[string]interface{}  //获取request body的数据
	if err := json.Unmarshal(body, &requestData); err != nil {
		s.Status.State = 1
		s.Status.Msg = err.Error()
		c.JSON(http.StatusInternalServerError, s)
	}

	//获取参数
	//appName: 应用名字
	//appNamespace: 应用命名空间
	//replicas: 实例个数
	appName := requestData["app_name"].(string)
	appNamespace := requestData["app_namespace"].(string)
	count := requestData["count"].(float64)
	counts := uint64(count)
	replicas := strconv.FormatUint(counts, 10)

	fetcher := new(util.KubeResponseFetcher)
	kubeCmd := new(util.KubeCmdImpl)
	statusCode, err := scaleReplicas(appName, appNamespace, replicas, fetcher, kubeCmd)

	if statusCode != http.StatusOK {
		s.Status.State = 1
		s.Status.Msg = err.Error()
		c.JSON(statusCode, s)
		return
	}
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
	//解析request body数据
	body, _ := ioutil.ReadAll(c.Request.Body)
	var requestData map[string]interface{}
	err := json.Unmarshal(body, &requestData)

	//appName: 应用名字
	//appNamespace: 应用命名空间
	//image: 应用镜像
	appName := requestData["app_name"].(string)
	appNamespace := requestData["app_namespace"].(string)
	imageName := requestData["image_name"].(string)
	tag := requestData["tag"].(string)
	image := registryPath + "/" + imageName + ":" + tag

	s := new(StatusResp)
	if err != nil {
		s.Status.State = 1
		s.Status.Msg = err.Error()
		c.JSON(http.StatusInternalServerError, s)
		return
	}

	kubeCmd := new(util.KubeCmdImpl)
	fetcher := new(util.KubeResponseFetcher)
	statusCode, err := deliveryRelease(appName, appNamespace, image, fetcher, kubeCmd)

	if err != nil {
		s.Status.State = 1
		s.Status.Msg = err.Error()
		c.JSON(http.StatusInternalServerError, s)
		return
	}

	if statusCode != http.StatusOK {
		s.Status.State = 1
		s.Status.Msg = err.Error()
		c.JSON(statusCode, s)
		return
	}

	s.Status.State = 0
	s.Status.Msg = "updated"
	c.JSON(http.StatusOK, s)
}





