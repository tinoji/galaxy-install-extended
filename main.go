package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

const (
	GitHubEndpoint = "https://api.github.com"
)

type Role struct {
	Src     string `yaml:"src"`
	Version string `yaml:"version"`
}

type Release struct {
	TagNmae string `json:"tag_name"`
}

func main() {
	// args
	if os.Args[1] == "-h" || os.Args[1] == "--help" {
		fmt.Println("Usage: galaxy-install-extended -r FILE [options]\n\n" +
			"Options:\n" +
			"  -h, --help      show this help message and exit\n" +
			"  -r ROLE_FILE    A file containing a list of roles to be imported\n\n" +
			"  See 'ansible-galaxy install --help' for other options")

		os.Exit(0)
	}

	if os.Args[1] != "-r" {
		log.Fatal("First option must be -r ROLE_FILE")
	}
	reqFile := os.Args[2]

	otherOpts := strings.Join(os.Args[3:len(os.Args)], " ")

	// read yaml and conv "lastet" to latest release
	buf, err := ioutil.ReadFile(reqFile)
	if err != nil {
		log.Fatal(err)
	}

	var roles []Role
	err = yaml.Unmarshal(buf, &roles)
	if err != nil {
		log.Fatal(err)
	}

	for i, role := range roles {
		if role.Version != "latest" {
			continue
		}

		if u, err := url.ParseRequestURI(role.Src); err == nil {
			var endpoint string

			if u.Hostname() == "github.com" {
				endpoint = GitHubEndpoint
			} else {
				endpoint = "https://" + u.Hostname() + "/api/v3"
			}

			repo := getRepoName(role.Src)
			tag := getTagName(repo, endpoint)
			roles[i].Version = tag
		}
	}

	d, err := yaml.Marshal(&roles)

	// make tmp file
	tmpFile, _ := ioutil.TempFile("", "galaxy-install-extended_tmp")
	tmpFile.Write(([]byte)(d))

	tmpYaml := tmpFile.Name() + ".yml"
	err = os.Rename(tmpFile.Name(), tmpYaml)
	if err != nil {
		log.Fatal(err)
	}

	defer os.Remove(tmpYaml)

	// exec
	command := "ansible-galaxy install -r " + tmpYaml + " " + otherOpts
	out, err := exec.Command("sh", "-c", command).Output()
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(string(out))
}

func getRepoName(src string) string {
	s := strings.Split(src, "/")
	len := len(s)
	user := s[len-2]
	repo := s[len-1]
	repo = strings.Split(repo, ".")[0]

	return user + "/" + repo
}

func getTagName(repo, endpoint string) string {
	url := endpoint + "/repos/" + repo + "/releases/latest"

	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	var d Release
	err = json.Unmarshal(body, &d)
	if err != nil {
		log.Fatal(err)
	}

	return d.TagNmae
}
