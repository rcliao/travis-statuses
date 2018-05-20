package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

var (
	orgName       = flag.String("orgName", "", "The organization name to look under in Github")
	repoNamesFile = flag.String("repoNamesFile", "", "The repository files to look for")
	token         = flag.String("token", "", "Authorization token for TravisCI")
)

type BuildsDTO struct {
	Builds []BuildDTO `json:"builds"`
}
type BuildDTO struct {
	State string `json:"state"`
}

func main() {
	flag.Parse()

	client := &http.Client{}
	repoNames := readRepoNamesFile(*repoNamesFile)
	for _, repoName := range repoNames {
		fmt.Printf(
			"Latest state of %s: %s\n",
			repoName,
			getLatestBuildState(client, *token, *orgName, repoName),
		)
	}
}

func readRepoNamesFile(filepath string) []string {
	f, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	result := []string{}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		result = append(result, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return result
}

func getLatestBuildState(client *http.Client, token, orgName, repoName string) string {
	url := fmt.Sprintf(
		"https://api.travis-ci.com/repo/%s%%2F%s/builds?limit=1&sort_by=finished_at:desc",
		orgName,
		repoName,
	)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Travis-API-Version", "3")
	req.Header.Set("User-Agent", "GoLang CLI App")
	req.Header.Set("Authorization", "token "+token)
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	var builds BuildsDTO

	err = json.NewDecoder(res.Body).Decode(&builds)
	if err != nil {
		log.Fatal(err)
	}
	if len(builds.Builds) > 0 {
		return builds.Builds[0].State
	}
	return "no-build"
}
