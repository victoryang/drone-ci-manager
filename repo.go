package main

import (
	"errors"
	"strings"
	"sync"
)

func init() {
 	ORM.AutoMigrate(&Repository{})
}

type Repository struct {
	Mux 		*sync.Mutex `gorm:"-"`
	Name 		string 		`gorm:"column:name;primary_key"`
	GitUrl 		string 		`gorm:"column:git_url"`
	Projects 	string 		`gorm:"column:projects"`
}

func FirstOrCreateRepository(gitUrl string) (*Repository, error) {
	name, err := ParseGitUrl(gitUrl)
	if err!=nil {
		return nil, err
	}

	mux := new(sync.Mutex)
	mux.Lock()
	defer mux.Unlock()

	repo := &Repository{
		Mux: mux,
		GitUrl: gitUrl,
	}
	result := ORM.FirstOrCreate(repo, Repository{Name: name})

	return repo, result.Error
}

func (r *Repository) AddProject(project string) error {
	r.Mux.Lock()
	defer r.Mux.Unlock()

	var new []string
	var isExist bool = false
	if len(r.Projects) == 0 {
		new = []string{project}
	} else {
		old := strings.Split(r.Projects, ",")
		for _,v :=range old {
			if v == project {
				isExist = true
				break
			}
		}

		if isExist == false {
			new = append(old, project)
		} else {
			new = old
		}
	}
	r.Projects = strings.Join(new, ",")

	result := ORM.Save(r)
	return result.Error
}

func (r *Repository) RemoveProject(project string) error {
	r.Mux.Lock()
	defer r.Mux.Unlock()

	if len(r.Projects) == 0 {
		return errors.New("Not belongs to this repository")
	}

	var new []string
	old := strings.Split(r.Projects, ",")
	for _, v :=range old {
		if v == project {
			continue
		}

		new = append(new, v)
	}
	r.Projects = strings.Join(new, ",")

	result := ORM.Save(r)
	return result.Error
}

func GetRepository(gitUrl string) (*Repository, error) {
	name, err := ParseGitUrl(gitUrl)
	if err!=nil {
		return nil, err
	}

	repo := &Repository{}
	result := ORM.Where("name = ?", name).First(repo)

	return repo, result.Error
}

func GetProjectsByRepo(gitUrl string) ([]string, error) {
	name, err := ParseGitUrl(gitUrl)
	if err!=nil {
		return nil, err
	}

	repo := &Repository{
		Name: name,
	}

	result := ORM.First(repo)
	if result.Error!=nil {
		return nil, result.Error
	}

	return strings.Split(repo.Projects, ","), nil
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

		projectSlice := []string {
			pString,
			r.Projects,
		}

		pString = strings.Join(projectSlice, ",")
	}

	return strings.Split(pString, ","), nil
}