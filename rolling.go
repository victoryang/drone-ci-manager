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
	GitURL 		string 		`json:"gitUrl"`
	BuildCmd 	string 		`json:"buildCmd"`
	Target		string 		`json:"target"`
	UnzipDir 	string 		`json:"unzipDir"`
	Lang		string 		`json:"lang"`
	BuildDependency 	string 		`json:"buildDependency"`
	StartCmd 		string 		`json:"startCmd"`
	StopCmd 		string 		`json:"stopCmd"`
	PreCmd 		string 		`json:"preCmd"`
	HttpPort	string 		`json:"httpPort"`
	RpcPort	string 			`json:"rpcPort"`
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
	jsonStr := []byte(`{"project":"` + project + `","tag":"` + tag + `"}`)
	req, err := http.NewRequest("POST", this.Addr + "/image/create", bytes.NewBuffer(jsonStr))
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
	jsonStr := []byte(`{"project":"` + project + `","tag":"` + tag + `","env":"` + env + `","deployStatus":"` + deployStatus + `"}`)
	if len(failLog) > 0 {
		failLog = url.QueryEscape(failLog)
		jsonStr = []byte(`{"project":"` + project + `","tag":"` + tag + `","env":"` + env + `","deployStatus":"` + deployStatus + `","maintainPlan":"` + failLog + `"}`)
	}
	req, err := http.NewRequest("POST", this.Addr + "/image/update", bytes.NewBuffer(jsonStr))
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