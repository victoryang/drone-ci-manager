package main

import (
    "fmt"
    "os"
    "os/exec"
    "path"
    "strings"
    "text/template"

    scripts "drone-ci-manager/template"
)

var (
    ProjectBase = "projects"

    Environments = []string{"sep", "staging", "rc", "release"}
    EnvPrefix = "release-"

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
    GitUrl          string      `json:"gitUrl"`
    UnZipDir        string      `json:"unzipDir"`
    HTTPPort        string      `json:"httpPort"`
    RPCPort         string      `json:"rpcPort"`
    StartCmd        string      `json:"startCmd"`
    PreCmd          string      `json:"preCmd"`
    StopCmd         string      `json:"stopCmd"`
    FromImage       string      `json:"fromImage"`
    BuildDependency    string   `json:"buildDependency"`
}

func (p *Project) generateFile(target string, outputDir string) error {
    content,err := scripts.Asset(target)
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

    projectBaseDir,err := p.getProjectBaseDir()
    if err!=nil {
        return err
    }

    fmt.Println("Ready to generate scripts for ", p.Project)

    for _,e :=range envs {

        env := EnvPrefix + strings.ToLower(e)
        envDir := path.Join(projectBaseDir, env)

        fmt.Println("Generating ", envDir)
        // create sub directory
        err = Mkdir(envDir)
        if err!=nil {
            fmt.Println("Create subdir err:", err)
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

func (p *Project) getProjectBaseDir() (string,error) {
    repo, err := ParseGitUrl(p.GitUrl)
    if err!=nil {
        return "", err
    }

    repo = strings.Replace(repo, "/", "_", -1)
    baseDir := path.Join(ProjectBase, repo, p.Project)

    isExist, err := IsDirExist(baseDir)
    if err!=nil {
        return "", err
    }

    if isExist==false {
        err = MkdirAll(baseDir)
    }

    return baseDir, err
}

func getProjectDir(project string, gitUrl string) (string,error) {
    repo, err := ParseGitUrl(gitUrl)
    if err!=nil {
        return "", err
    }

    repo = strings.Replace(repo, "/", "_", -1)

    return path.Join(ProjectBase, repo, project), nil
}

func GetDockerfileFromBytes(project string, gitUrl string, env string) string {
    workingDir, err := getProjectDir(project, gitUrl)
    if err!=nil {
        return ""
    }

    dir := path.Join(workingDir, EnvPrefix+strings.ToLower(env))
    fromBytes, _ := exec.Command("bash", "-c", `cd `+dir+` && head -n 1 Dockerfile  | awk -F'/' '{print $2}' | sed '{s/:/-/g}' | awk -F'_' '{print $NF}'`).Output()
    from := strings.TrimSpace(string(fromBytes))
    if len(from) == 0 {
        return ""
    }

    return from
}