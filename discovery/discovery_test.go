package discovery

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPomDiscovery(t *testing.T) {
	root, err := os.MkdirTemp("", "example")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(root)
	fmt.Println("created " + root)
	rootPom := filepath.Join(root, "pom.xml")
	submod := filepath.Join(root, "submod")
	nestedPom := filepath.Join(submod, "pom.xml")
	subLevel := filepath.Join(root, "sub", "level")
	deeplyNestedPom := filepath.Join(subLevel, "pom.xml")
	noPom := filepath.Join(root, "no-pom")
	ignore := filepath.Join(noPom, "README.txt")
	for _, dir := range []string{submod, subLevel, noPom} {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			panic(err)
		}
	}
	for _, f := range []string{rootPom, nestedPom, deeplyNestedPom, ignore} {
		fmt.Println("creating " + f)
		_, err := os.Create(f)
		if err != nil {
			panic(err)
		}
	}
	expected := []string{".", "./sub/level", "./submod"}

	actual := FindMavenModules(root)

	require.Equal(t, expected, actual)
}
