package main

import (
    "fmt"
    "os"
    "os/exec"
    "path"
    "strings"
    "text/template"
)

var (
    Environments = []string{"sep", "staging", "rc", "release"}
    EnvPrefix = "release-"

    TemplateFiles = []string{
        "deploy_1_stop.sh",
        "deploy_2_replace.sh",
        "deploy_3_start.sh",
        //"deploy_4_rollback.sh",
        "Dockerfile",
    }

    FromImages = map[string][]string {
        "java": []string {
            "docker.snowballfinance.com:5000/java8:v19",
            "docker.snowballfinance.com:5000/base-v2020java8",
        },
    }

    InputBase = "template"
    ProjectDir = "projects"
)

type Project struct {
    Project         string      `json:"-"`
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
    temp, err := template.ParseFiles(path.Join(InputBase, target))
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

    workingDir := getProjectDir(p.Project)
    _,err := os.Stat(workingDir)
    if err!=nil {
        return err
    }

    fmt.Println("Ready to generate scripts for ", p.Project)

    for _,e :=range envs {
        env := EnvPrefix + strings.ToLower(e)
        envDir := path.Join(workingDir, env)
        fmt.Println("\nGenerating ", envDir)

        // create sub directory
        err = os.Mkdir(envDir, os.ModeDir)
        if err!=nil {
            fmt.Println("mk sub dir err:", err)
            return err
        }
        os.Chmod(envDir, 0755)

        // generate files
        for _,file :=range TemplateFiles {

            if err = p.generateFile(file, envDir); err!=nil {
                fmt.Println(file, "is generated failed:", err)
                break
            }
            fmt.Printf("%v is generated successfully\n", file)
        }

        if err!=nil {
            break
        }
    }

    if err!=nil {
        fmt.Println("Generating files fails")
        defer func() {
            if err = os.RemoveAll(workingDir); err==nil {
                fmt.Printf("Clean up %v successuflly\n", workingDir)
            } else {
                fmt.Printf("Clean up %v fails, please handle it\n", workingDir)
            }
        }()
    } else {
        fmt.Println("\nFiles are generated successfully")
    }

    return err
}

func getProjectDir(project string) string {
    return path.Join(ProjectDir, project)
}

func CreateWorkingDir(project string) error {
    workingDir := getProjectDir(project)

    err := os.MkdirAll(workingDir, os.ModeDir)
    if err!=nil {
        return err
    }

    return os.Chmod(workingDir, 0755)
}

func GetProjectsFromDir() ([]string, error){
    projDir, err := os.Open(ProjectDir)
    if err!=nil {
        return nil, err
    }

    return projDir.Readdirnames(-1)
}

func DeleteProjectFromDir(project string) error {
    workingDir := path.Join(ProjectDir, project)

    return os.RemoveAll(workingDir)
}

func GetDockerfileFromBytes(project string, env string) string {
    workingDir := path.Join(ProjectDir, project)
    dir := path.Join(workingDir, EnvPrefix+strings.ToLower(env))

    fromBytes, _ := exec.Command("bash", "-c", `cd `+dir+` && head -n 1 Dockerfile  | awk -F'/' '{print $2}' | sed '{s/:/-/g}' | awk -F'_' '{print $NF}'`).Output()
    from := strings.TrimSpace(string(fromBytes))
    if len(from) == 0 {
        return ""
    }

    return from
}

func GetScripts(project string, env string) {
    workingDir := getProjectDir(project)
    outputDir := path.Join(workingDir, EnvPrefix + strings.ToLower(env))

    scripts := make(map[string]string)
}