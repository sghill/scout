package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"runtime/debug"

	"sghill.net/scout/discovery"
	"sghill.net/scout/pom"

	"github.com/go-git/go-git/v5"

	flag "github.com/spf13/pflag"
)

func main() {
	var dir *string = flag.StringP("directory", "d", ".", "where to start the search from")
	var outFile *string = flag.StringP("out-file", "o", "results.json", "where to store the json results")
	var repo *string = flag.StringP("repository", "r", os.Getenv("SUBJECT_REPO"), "the scouted repo")
	var contextUri *string = flag.StringP("context-uri", "u", os.Getenv("CIRCLE_BUILD_URL"), "where to get more context on this run")
	flag.Parse()

	indexerVersion := "0.0.1"
	buildInfo, ok := debug.ReadBuildInfo()
	if !ok {
		fmt.Println("could not read build info of module")
	} else {
		indexerVersion = buildInfo.Main.Version
	}

	r, err := git.PlainOpen(*dir)
	if err != nil {
		panic(err)
	}
	ref, err := r.Head()
	if err != nil {
		panic(err)
	}
	branch := ref.Name().Short()
	commitId := ref.Hash().String()

	modules := discovery.FindMavenModules(*dir)

	results := make([]pom.Result, len(modules))
	for i, module := range modules {
		fmt.Println("for module ", module)
		path := filepath.Join(*dir, module)
		dat, err := os.ReadFile(filepath.Join(path, "pom.xml"))
		if err != nil {
			panic(err)
		}
		parsed := pom.Project{}
		err = xml.Unmarshal(dat, &parsed)
		if err != nil {
			panic(err)
		}

		result := &pom.Result{
			Path:                module,
			PluginParentVersion: pom.PluginParentVersion(parsed),
			PluginBoms:          pom.AppliedPluginBoms(parsed),
			JenkinsVersion:      pom.JenkinsVersion(parsed),
			RecommendedJava:     pom.RecommendJava(parsed.Parent.Version),
		}

		results[i] = *result
	}

	dat, err := json.Marshal(&pom.ScoutResult{
		Repo:           *repo,
		Branch:         branch,
		IndexerVersion: indexerVersion,
		ExecUri:        *contextUri,
		CommitId:       commitId,
		ModuleResults:  results,
	})
	if err != nil {
		panic(err)
	}
	err = os.WriteFile(*outFile, dat, 0644)
	if err != nil {
		panic(err)
	}
	fmt.Println("results written to ", *outFile)
}
