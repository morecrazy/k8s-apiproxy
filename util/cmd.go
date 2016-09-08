package util

import (
	"os"
	"os/exec"
	"io/ioutil"
	"fmt"
	. "backend/common"
	"strings"
	"backend/common/protocol"
	"codoon_ops/kubernetes-apiproxy/util/set"
)

func GetRuntimeEnv() string {
	return os.Getenv("GOENV")
}

func Run(command , appName, appNamespace string) error {
	_, err := ExecCommand(command)
	kubeCmd := new(KubeCmdImpl)
	fetcher := new(KubeResponseFetcher)
	go RegisterDNS(appName, appNamespace, kubeCmd, fetcher)
	if err != nil {
		fmt.Errorf("exec command error: ", err.Error())
		return err
	}
	return nil
}

func ExecCommand(command string) ([]byte, error) {
	cmd := exec.Command(command)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		Logger.Error("StdoutPipe: " + err.Error())
		err := fmt.Errorf("StdoutPipe: " + err.Error())
		return nil, err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		err := fmt.Errorf("StderrPipe: ", err.Error())
		return nil, err
	}

	if err := cmd.Start(); err != nil {
		err := fmt.Errorf("Start: ", err.Error())
		return nil, err
	}

	bytesErr, err := ioutil.ReadAll(stderr)
	if err != nil {
		err := fmt.Errorf("ReadAll stderr: ", err.Error())
		return bytesErr, err
	}

	if len(bytesErr) != 0 {
		err := fmt.Errorf("stderr is not nil: %s", bytesErr)
		return bytesErr, err
	}

	bytes, err := ioutil.ReadAll(stdout)
	if err != nil {
		err := fmt.Errorf("ReadAll stdout: ", err.Error())
		return bytes, err
	}

	if err := cmd.Wait(); err != nil {
		err := fmt.Errorf("Wait: ", err.Error())
		return bytes, err
	}
	return bytes, nil
}

func BytesToDNSRequestJSON(bytes []byte, domain string) (req interface{}) {
	arrays := strings.Split(string(bytes), "\n")
	set := set.New()

	for _, item := range arrays {
		if item != "" {
			set.Add(item)
			Logger.Debug("the ip is: %v", item)
		}
	}
	urlList := set.List()

	urls := make([]protocol.RR, 0)
	for _, url := range urlList {
		item := protocol.RR{
			Host: url.(string),
		}
		urls = append(urls, item)
	}

	reqJson := protocol.SetDnsReq{
		URL: domain,
		RRs: urls,
	}

	//调用http接口,注册服务名到DNS server
	Logger.Debug("registry %v to DNS server, the ips is: %v", domain, urls)
	return reqJson
}