package main

import (
	"errors"
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