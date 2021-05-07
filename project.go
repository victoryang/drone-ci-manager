package main

import (
    "fmt"
    "os"
    "os/exec"
    "path"
    "strings"
    "text/template"

    scripts "sce-build-manager/template"
)

var (
    ProjectBase = "projects"

    ScriptTemplate = []string{
        "deploy_1_stop.sh",
        "deploy_2_replace.sh",
        "deploy_3_start.sh",
        //"deploy_4_rollback.sh",
        "Dockerfile",
    }
)

type Project struct {
    Project         string      `json:"project"`
    UnZipDir        string      `json:"unzipDir"`
    Lang            string      `json:"lang"`
    HTTPPort        string      `json:"httpPort"`
    RPCPort         string      `json:"rpcPort"`
    StartCmd        string      `json:"startCmd"`
    PreCmd          string      `json:"preCmd"`
    StopCmd         string      `json:"stopCmd"`
    FromImage       string      `json:"fromImage"`
    BuildDependency    string   `json:"buildDependency"`
}

func (p *Project) generateFile(target string, outputDir string) error {
    content,err := scripts.Asset(path.Join(p.Lang, target))
    if err!=nil {
        return err
    }

    temp, err := template.New(target).Parse(string(content))
    if err!=nil {
        return err
    } 
    t := template.Must(temp, err)

    //t.Execute(os.Stdout, p)

    output := path.Join(outputDir, target)
    f, err := os.Create(output)
    if err!=nil {
        return err
    }
    defer f.Close()

    err = t.Execute(f, p)

    if target == "Dockerfile" {
        return err
    }

    return os.Chmod(output, 0755)
}

func (p *Project) generateFiles(envs []string) error {
    fmt.Println("Generating scripts...")

    projectDir := path.Join(ProjectBase, p.Project)
    err := GetOrCreateDir(projectDir)
    if err!=nil {
        return err
    }

    for _,env :=range envs {

        envDir := path.Join(projectDir, "release-" + strings.ToLower(env))
        err = GetOrCreateDir(envDir)
        if err!=nil {
            return err
        }

        // generate files
        for _,target :=range ScriptTemplate {

            if err = p.generateFile(target, envDir); err!=nil {
                fmt.Println(target, "is generated failed:", err)
                break
            }
        }

        if err!=nil {
            break
        }
    }

    if err!=nil {
        fmt.Println("Generating files fails")
    } else {
        fmt.Println("\nFiles are generated successfully")
    }

    return err
}

func GetDockerfileFromBytes(project string, env string) string {
    dir := path.Join(ProjectBase, project, "release-" + strings.ToLower(env))

    fromBytes, _ := exec.Command("bash", "-c", `cd `+dir+` && head -n 1 Dockerfile  | awk -F'/' '{print $2}' | sed '{s/:/-/g}' | awk -F'_' '{print $NF}'`).Output()
    from := strings.TrimSpace(string(fromBytes))
    if len(from) == 0 {
        return ""
    }

    return "_" + from
}