package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"

	"sghill.net/scout/discovery"
	"sghill.net/scout/pom"

	flag "github.com/spf13/pflag"
)

func main() {
	var dir *string = flag.StringP("directory", "d", ".", "where to start the search from")
	var outFile *string = flag.StringP("out-file", "o", "results.json", "where to store the json results")
	flag.Parse()

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

	dat, err := json.Marshal(&pom.ScoutResult{ModuleResults: results})
	if err != nil {
		panic(err)
	}
	err = os.WriteFile(*outFile, dat, 0644)
	if err != nil {
		panic(err)
	}
	fmt.Println("results written to ", *outFile)
}
