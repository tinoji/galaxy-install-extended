package main

import (
	"encoding/json"
	"fmt"
	// flags "github.com/jessevdk/go-flags"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

// const (
// 	GitHubEndpoint = "https://api.github.com"
// )

type Role struct {
	Src     string `yaml:"src"`
	Version string `yaml:"version"`
}

type Release struct {
	TagNmae string `json:"tag_name"`
}

var endpoint = "https://api.github.com"

func main() {
	// args
	if os.Args[1] == "-h" || os.Args[1] == "--help" {
		fmt.Println("Usage: galaxy-install-extended -r FILE [options]\n\n" +
			"Options:\n" +
			"  -h, --help               show this help message and exit\n" +
			"  -r ROLE_FILE             A file containing a list of roles to be imported\n" +
			"  -e GITHUB_API_ENDPOINT   API endpoint for GitHub Enterprise. The default is\n" +
			"                           https://api.github.com\n\n" +
			"  See 'ansible-galaxy install --help' for other options")

		os.Exit(0)
	}

	if os.Args[1] != "-r" {
		log.Fatal("First option must be -r ROLE_FILE")
	}
	reqFile := os.Args[2]

	var otherOpts string

	if os.Args[3] == "-e" {
		endpoint = os.Args[4]
		otherOpts = strings.Join(os.Args[5:len(os.Args)], " ")
	} else {
		otherOpts = strings.Join(os.Args[3:len(os.Args)], " ")
	}

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
		if role.Version == "latest" {
			repo := getRepoName(role.Src)
			tag := getTagName(repo)
			roles[i].Version = tag
		}
	}

	d, err := yaml.Marshal(&roles)

	// make tmp file
	tmpFile, _ := ioutil.TempFile("", "tmp")
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
	return s[len-2] + "/" + s[len-1]
}

func getTagName(repo string) string {
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
