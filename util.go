package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

func ParseGitUrl(gitSshUrl string) (string,error) {
	resList := strings.Split(gitSshUrl, "git.snowballfinance.com:")
	if len(resList) < 2 {
		return ``, errors.New("Git address error: " + gitSshUrl)
	}

	res := strings.TrimSuffix(resList[1], ".git")

	return res, nil
}

func GetOrCreateDir(dir string) error {
    isExist, _ := IsDirExist(dir)
    if isExist {
        return nil
    }

    return MkdirAll(dir)
}

func IsDirExist(dir string) (bool,error) {
	_,err := os.Stat(dir)
	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}

	return false, err
}

func Mkdir(dir string) error {
	err := os.Mkdir(dir, os.ModeDir)
	if err!=nil {
		fmt.Printf("MkdirAll for %v fails, please check\n", dir)
		return err
	}

	return os.Chmod(dir, 0755)
}

func MkdirAll(dir string) error {
	err := os.MkdirAll(dir, os.ModeDir)
	if err!=nil {
		fmt.Printf("MkdirAll for %v fails, please check\n", dir)
		return err
	}

	return os.Chmod(dir, 0755)
}

func RemoveAll(dir string) error {
    err := os.RemoveAll(dir)
    if err!=nil {
        fmt.Printf("Cleaning up %v fails, please check\n", dir)
        return err
    }

    fmt.Printf("Cleaning up %v successuflly\n", dir)
    return nil
}