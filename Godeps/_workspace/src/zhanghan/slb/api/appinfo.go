package api

type AppInfo struct {
	AccessKeyId     string
	AccessKeySecret string
}

var appInfo *AppInfo

func GetAppInfo() *AppInfo {
	return appInfo
}

func SetAppInfo(keyId string, keySecret string) {
	appInfo = new(AppInfo)
	appInfo.AccessKeyId = keyId
	appInfo.AccessKeySecret = keySecret
}
