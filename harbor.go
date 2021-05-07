package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"path"
	"sort"
	"time"
)

const (
	HTTPSchema = "https://"

	HarborBaseUrl = "hub.snowballfinance.com"
	HarborPublicProject = "cicd"

	HarborTimeout = 5
	HarborApiPath = "api/v2.0"

	HarborUser = "Harbor"
	HarborSecret = "HarborX123"
)

func BuildImageName(project string, tag string) string {

	return path.Join(HarborBaseUrl, HarborPublicProject, project) + ":" + tag
}

type HarborApi struct {
	ApiUrl 		string
	BaseUrl 	string
	HttpClient  *http.Client
}

func NewHarborApi(baseUrl string) *HarborApi {
	c := http.Client {
		Timeout: HarborTimeout * time.Second,
	}

	return &HarborApi {
		ApiUrl: path.Join(baseUrl, HarborApiPath),
		BaseUrl: baseUrl,
		HttpClient: &c,
	}
}

func getHarborRepositoryName(repository string) string {
	return path.Join("projects", HarborPublicProject, "repositories", repository)
}

func (h *HarborApi) doHttpRequest(request *http.Request) ([]byte, error) {
	resp, err := h.HttpClient.Do(request)
	if err!=nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

func (h *HarborApi) GetTagsByProj(project string) ([]string, error){

	repository := getHarborRepositoryName(project)
	url := HTTPSchema + path.Join(h.ApiUrl, repository, "artifacts")

	request, err := http.NewRequest("GET", url, nil)
	if err!=nil {
		return nil, err
	}

	body, err := h.doHttpRequest(request)
	if err!=nil {
		return nil, err
	}

	model := make([]map[string]interface{},0)
	err = json.Unmarshal(body, &model)
	if err!=nil{
		return nil, err
	}

	total_tags := make([]string, 0)
	for _,v :=range model {
		if tags,ok := v["tags"].([]interface{}); ok {
			for _,tag := range tags {
				if t,ok := tag.(map[string]interface{}); ok {
					total_tags = append(total_tags, t["name"].(string))
				}
			}
		}
	}

	sort.Sort(sort.Reverse(sort.StringSlice(total_tags)))

	return total_tags, nil
}