package main

import (
    "fmt"
    "os"
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
    OutputBase = "projects"
)

type Param struct {
    Project         string      `json:"-"`
    UnZipDir        string      `json:"unzip_dir"`
    HTTPPort        string      `json:"http_port"`
    RPCPort         string      `json:"rpc_port"`
    StartCmd        string      `json:"start_cmd"`
    PreCmd          string      `json:"pre_cmd"`
    StopCmd         string      `json:"stop_cmd"`
    FromImage       string      `json:"from_image"`
    BuildDependency    string   `json:"build_dependency"`
}

func (p *Param) generateFile(target string, outputDir string) error {
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

func (p *Param) generateFiles(envs []string) error {
    fmt.Println("Generating scripts...")

    var err error
    outputDir := path.Join(OutputBase, p.Project)
    err = os.Mkdir(outputDir, os.ModeDir)
    if err!=nil {
        fmt.Println("mkdir err:", err)
        return err
    }

    fmt.Println("Ready to generate scripts for ", p.Project)

    for _,e :=range envs {
        env := EnvPrefix + strings.ToLower(e)
        envDir := path.Join(outputDir, env)
        fmt.Println("\nGenerating ", envDir)

        // create sub directory
        err = os.Mkdir(envDir, os.ModeDir)
        if err!=nil {
            fmt.Println("mk sub dir err:", err)
            return err
        }

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
            if err = os.RemoveAll(outputDir); err==nil {
                fmt.Printf("Clean up %v successuflly\n", outputDir)
            } else {
                fmt.Printf("Clean up %v fails, please handle it\n", outputDir)
            }
        }()
    } else {
        fmt.Println("\nFiles are generated successfully")
    }

    return err
}

func getScripts(project string, env string) {
    outputDir := path.Join(OutputBase, project, EnvPrefix + strings.ToLower(e))

    scripts := make(map[string]string)
}

func main() {
    p := &Param{
        Project: "sce-rolling",
        UnZipDir: "sce-rolling",
        HTTPPort: "8000",
        StartScript: "start.sh",
        ExtraStartScript: "",
        StopScript: "stop.sh",
        FromImage: "docker.snowballfinance.com:5000/java8:v19",
        ExtraImageOperations: "",
    }

    p.generateFiles()
}