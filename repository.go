package main

import (
	"os"
	"path"

	"github.com/jinzhu/gorm"
)

func init() {
	ORM.AutoMigrate(&Repository{})
}

type Repository struct {
	gorm.Model      `json:"-"`

	// e.g. ops/sce-rolling
	Name 			string 		`gorm:"column:name;size:64"`
	GitUrl			string 		`gorm:"column:git_url;size:64"`
	RollingProject 		[]string 	`gorm:"column:rolling_project;size:64"`
	IsActived			bool 		`gorm:"column:is_actived;size:64"`
}

func NewRepository(name string, gitSshUrl string) error {
	proj := &Repository{
		Name: name,
		GitUrl: gitSshUrl,
		IsActived: true,
	}

	if err := ORM.Create(proj).Error; err != nil {
		return err
	}

	return nil
}

func GetProjectsbyRepo(name string) []string {
	proj := &Repository{}

	err := ORM.First(proj, "name = ?", name).Error
	if err!=nil {
		return nil
	}

	return proj.RollingProject
}

func (p *Repository) AddRollingProjects(projects []string) {
	for _, proj :=range projects {
		p.RollingProject = append(p.RollingProject, proj)
	}
}

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