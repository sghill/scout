package discovery

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
)

const pomFileName string = "pom.xml"

func FindMavenModules(rootDir string) []string {
	results := make([]string, 0)
	filepath.WalkDir(rootDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}
		if d.Name() == pomFileName {
			relative := strings.TrimPrefix(path, rootDir)
			parent := strings.TrimSuffix(relative, pomFileName)
			noTrailingSlash := strings.TrimSuffix(parent, string(filepath.Separator))
			results = append(results, noTrailingSlash)
		}
		return nil
	})
	return results
}
