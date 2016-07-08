package util

import (
	"backend/common"
)

var RegistryUrl = ""

type RegistryResponseFetch interface {
	GetRegistryImageList() (int, string, error)
	DeleteRegistryImage(name string) (int, string, error)
	GetImageTags(name string) (int, string, error)
	DeleteImageTag(name, tag string) (int, string, error)
}

type RegistryResponseFetcher struct {}

func (fetcher *RegistryResponseFetcher) GetRegistryImageList() (int, string, error) {
	postFix := "/v2/_catalog"
	url := RegistryUrl + postFix
	statusCode, response, err := common.SendRawRequest("GET", url, nil)
	common.Logger.Debug("statusCode: %d", statusCode)
	common.Logger.Debug("repsonse: %s", response)
	return statusCode, response, err
}

func (fetcher *RegistryResponseFetcher) DeleteRegistryImage(name string) (int, string, error) {
	postFix := "/v1/repositories"
	url := RegistryUrl + postFix + "/" + name + "/"
	common.Logger.Info("url: %v", url)
	statusCode, response, err := common.SendRawRequest("DELETE", url, nil)
	return statusCode, response, err
}

func (fetcher *RegistryResponseFetcher) GetImageTags(name string) (int, string, error) {
	postFix := "/v2"
	url := RegistryUrl + postFix + "/" + name + "/tags/list"
	statusCode, response, err := common.SendRawRequest("GET", url, nil)
	return statusCode, response, err
}

func (fetcher *RegistryResponseFetcher) DeleteImageTag(name, tag string) (int, string, error) {
	//获取digest
	postFix := "/v2"
	url := RegistryUrl + postFix + "/" + name + "/manifests" + "/" + tag
	common.Logger.Info("url: %v", url)
	statusCode, response, err := common.SendRawRequest("DELETE", url, nil)
	return statusCode, response, err
}