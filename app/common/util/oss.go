package util

import (
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
)

const (
	//GojsModel 设计图数据
	GojsModel = "gojs_model.json"
	//WorkflowMetadata 设计图入参数据
	WorkflowMetadata = "matedata.json"
	//WorkflowMetadataYaml yaml输入参数
	WorkflowMetadataYaml = "workflow_matedata.yaml"
	//WorkflowYaml 工作流yaml
	WorkflowYaml = "workflow.yaml"
)

//GetOssContent 并行获取Oss内容
func GetOssContent(urls map[string]string) (map[string]string, error) {
	result := map[string]string{}
	waitGroup := sync.WaitGroup{}
	var httpErr error
	for key, url := range urls {
		if !strings.HasPrefix(url, "http") {
			//如果不是url就直接赋值给返回数据
			result[key] = url
			continue
		}
		waitGroup.Add(1)
		go func(key, url string) {
			defer waitGroup.Done()
			response, err := http.Get(url)
			if err != nil {
				result[key] = ""
				httpErr = err
				return
			}
			bytes, err := ioutil.ReadAll(response.Body)
			result[key] = string(bytes)
		}(key, url)
	}
	waitGroup.Wait()
	if httpErr != nil {
		return nil, httpErr
	}
	return result, nil
}
