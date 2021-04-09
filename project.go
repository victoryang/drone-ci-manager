package main

import (
	"os"
	"path"
)

const (
	ProjectDir = "projects"
)

func EnableRepository(name string) error {
	p := &Repository{}

	result := ORM.Model(p).Where("name = ?", name).Update("is_actived", true)
	return result.Error
}

func DisableRepository(name string) error {
	p := &Repository{}

	result := ORM.Model(p).Where("name = ?", name).Update("is_actived", false)
	return result.Error
}

func DeleteRepository(name string, giturl string) error {
	p := &Repository{
		Name: name,
	}

	result := ORM.Unscoped().Delete(p)
	return result.Error
}

func (p *Repository) AddRollingProject(projects []string) {
	for _, proj :=range projects {
		p.RollingProject = append(p.RollingProject, proj)
	}
}

func CreateProject(project string) error {
	dir := path.Join(ProjectDir, project)

	return os.Mkdir(dir, os.ModeDir)
}

func DeleteProject(project string) error {
	dir := path.Join(ProjectDir, project)

	return os.RemoveAll(dir)
}

func getProjects() ([]string, error){
	projDir, err := os.Open(ProjectDir)
	if err!=nil {
		return nil, err
	}

	return projDir.Readdirnames(-1)
}