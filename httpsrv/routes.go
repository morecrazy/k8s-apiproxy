package httpsrv

import (
	"third/gin"
)

func SetupRoutes(router *gin.Engine) {
	image_group := router.Group("/api/v1/image")
	{
		image_group.GET("/list", GetImagesList) 		//获取镜像列表(registry v2版本接口可以实现)
		image_group.POST("/delete", DeleteImage)		//删除镜像

	}
	tag_group := router.Group("/api/v1/tag")
	{
		tag_group.GET("/list", GetImageTags)		//获取镜像版本列表(registry v2版本接口可以实现)
		tag_group.POST("/delete", DeleteTag)		//删除镜像版本(registry v2版本接口可以实现)
		tag_group.POST("/latest", UpdateImageTag)		//设置镜像版本为latest
	}
	cluster_group := router.Group("/api/v1/cluster")
	{
		cluster_group.GET("/list", GetClusterList)      //获取运行环境(集群)列表
	}
	service_group := router.Group("/api/v1/service")
	{
		service_group.GET("/list", GetServicesList) 		//获取服务列表
		service_group.GET("/metadata", GetServiceMetadata) 		//获取服务详情
		service_group.GET("/configdata", GetServiceConfig)			//获取服务配置和环境信息
		service_group.POST("/updatestate", UpdateServiceStatus) 	//更新服务状态,启动,停止
		service_group.POST("/create", CreateService)				//创建服务
		service_group.POST("/scale", ScaleReplicas)				//伸缩服务(加减容器个数)
		service_group.POST("/delivery", DeliveryRelease) 		//发布版本(更新)
		service_group.POST("/quota", UpdateServiceQuota)		//更新服务资源配额(cpu,内存)
		service_group.POST("/command", UpdateServiceCmd)		//更新服务执行命名
	}
}
