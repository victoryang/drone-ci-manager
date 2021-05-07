package main

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/jinzhu/gorm"
)

func init() {
 	ORM.AutoMigrate(&Repository{})
}

var (
	Mux = new(sync.Mutex)
)

type Repository struct {
	gorm.Model
	Name 		string 		`gorm:"column:name;primary_key"`
	GitUrl 		string 		`gorm:"column:git_url"`
	Projects 	string 		`gorm:"column:projects"`
}

func FirstOrCreateRepository(gitUrl string) (*Repository, error) {
	name, err := ParseGitUrl(gitUrl)
	if err!=nil {
		return nil, err
	}

	Mux.Lock()
	defer Mux.Unlock()

	repo := &Repository{
		GitUrl: gitUrl,
	}
	result := ORM.FirstOrCreate(repo, Repository{Name: name})

	return repo, result.Error
}

func (r *Repository) AddProject(project string) error {
	Mux.Lock()
	defer Mux.Unlock()

	var latest []string
	var isExist bool = false
	if len(r.Projects) == 0 {
		latest = []string{project}
	} else {
		old := strings.Split(r.Projects, ",")
		for _,v :=range old {
			if v == project {
				isExist = true
				break
			}
		}

		if isExist == false {
			latest = append(old, project)
		} else {
			latest = old
		}
	}
	r.Projects = strings.Join(latest, ",")

	result := ORM.Save(r)
	return result.Error
}

func (r *Repository) RemoveProject(project string) error {
	Mux.Lock()
	defer Mux.Unlock()

	if len(r.Projects) == 0 {
		return errors.New("Not belongs to this repository")
	}

	var latest []string
	old := strings.Split(r.Projects, ",")
	for _, v :=range old {
		if v == project {
			continue
		}

		latest = append(latest, v)
	}

	if len(latest) == 0 {
		result := ORM.Delete(r)
		return result.Error
	}

	r.Projects = strings.Join(latest, ",")

	result := ORM.Save(r)
	return result.Error
}

func GetRepository(gitUrl string) (*Repository, error) {
	name, err := ParseGitUrl(gitUrl)
	if err!=nil {
		return nil, err
	}

	repo := &Repository{
		GitUrl: gitUrl,
	}
	result := ORM.Where("name = ?", name).First(repo)

	return repo, result.Error
}

func GetProjectsByRepo(gitUrl string) ([]string, string,error) {
	name, err := ParseGitUrl(gitUrl)
	if err!=nil {
		fmt.Println("could not resolve git url:", err)
		return nil, name, err
	}

	repo := &Repository{
		Name: name,
	}

	result := ORM.First(repo)
	if result.Error!=nil {
		return nil, name, result.Error
	}

	projects := strings.Split(repo.Projects, ",")
	if len(projects) == 0 {
		return nil, name, errors.New("No projects found")
	}

	return projects, name, nil
}

func GetAllProjects() ([]string, error) {
	repos := make([]Repository, 0)

	result := ORM.Find(&repos)
	if result.Error!=nil {
		return nil, result.Error
	}

	var pString string = ""
	for _, r :=range repos {
		if pString == "" {
			pString = r.Projects
			continue
		}

		if r.Projects == "" {
			continue
		}

		projectSlice := []string {
			pString,
			r.Projects,
		}

		pString = strings.Join(projectSlice, ",")
	}

	if pString == "" {
		return nil, nil
	}

	return strings.Split(pString, ","), nil
}
