package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"fmt"
	"net/http"
	"net/url"
)

type RollingCli struct {
	Addr 	string
}

func NewRollingClient(addr string) *RollingCli {
	return &RollingCli{Addr: addr}
}

type RollingBuildInfo struct {
        BuildCmd        string          `json:"command"`
        Target          string          `json:"from"`
        Lang            string          `json:"lang"`
}

func (this *RollingCli) GetBuildInfo(project string) *RollingBuildInfo {
	fmt.Println("get build info from rolling")

	url := this.Addr + "/projects/" + project + "/build_info"
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Rolling Build")
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("get build info err:", err)
		return nil
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err !=nil {
		fmt.Println("read build info err:", err)
		return nil
	}

	r := &RollingBuildInfo{}
	err = json.Unmarshal(body, r)
	if err !=nil {
		fmt.Println("json UnMarshal err:", err)
		return nil
	}
	return r
}

type RollingBasicInfo struct {
	BuildCmd 	string 		`json:"build_cmd"`
	Target		string 		`json:"target"`
	UnzipDir 	string 		`json:"unzip_dir"`
	Lang		string 		`json:"lang"`
	BuildDependency 	string 		`json:"build_dependency"`
	StartCmd 		string 		`json:"start_cmd"`
	StopCmd 		string 		`json:"stop_cmd"`
	PreCmd 		string 		`json:"pre_cmd"`
	HttpPort	string 		`json:"http_port"`
	RpcPort	string 			`json:"rpc_port"`
}

func (this *RollingCli) GetBasicInfo(project string) *RollingBasicInfo {
	fmt.Println("get build info from rolling")

	url := this.Addr + "/projects/" + project + "/basic_info"
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Rolling Build")
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("get build info err:", err)
		return nil
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err !=nil {
		fmt.Println("read build info err:", err)
		return nil
	}

	r := &RollingBasicInfo{}
	err = json.Unmarshal(body, r)
	if err !=nil {
		fmt.Println("json UnMarshal err:", err)
		return nil
	}
	return r
}

func (this *RollingCli) CreateImage(project string, tag string) error {
	jsonStr := []byte(`{"Project":"` + project + `","Tag":"` + tag + `"}`)
	req, err := http.NewRequest("POST", this.Addr + "/image/create_image", bytes.NewBuffer(jsonStr))
	req.Header.Set("User-Agent", "Rolling Build")
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	} else {
		resp.Body.Close()
	}
	return nil
}

func (this *RollingCli) UpdateImage(project string, tag string, env string, deployStatus string, failLog string) error {
	jsonStr := []byte(`{"Project":"` + project + `","Tag":"` + tag + `","Env":"` + env + `","DeployStatus":"` + deployStatus + `"}`)
	if len(failLog) > 0 {
		failLog = url.QueryEscape(failLog)
		jsonStr = []byte(`{"Project":"` + project + `","Tag":"` + tag + `","Env":"` + env + `","DeployStatus":"` + deployStatus + `","MaintainPlan":"` + failLog + `"}`)
	}
	req, err := http.NewRequest("POST", this.Addr + "/image/update_image", bytes.NewBuffer(jsonStr))
	req.Header.Set("User-Agent", "Rolling Build")
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	} else {
		resp.Body.Close()
	}
	return nil
}