package pom

import (
	"encoding/xml"
	"strconv"
	"strings"
)

type Dependency struct {
	GroupId    string `xml:"groupId"`
	ArtifactId string `xml:"artifactId"`
	Version    string `xml:"version"`
	Scope      string `xml:"scope"`
	Type       string `xml:"type"`
}

type DependencyManagement struct {
	Dependencies []Dependency `xml:"dependencies>dependency"`
}

type Parent struct {
	GroupId    string `xml:"groupId"`
	ArtifactId string `xml:"artifactId"`
	Version    string `xml:"version"`
}

type Properties struct {
	JenkinsVersion string `xml:"jenkins.version"`
}

type Project struct {
	XMLName              xml.Name             `xml:"project"`
	Parent               Parent               `xml:"parent"`
	DependencyManagement DependencyManagement `xml:"dependencyManagement"`
	Properties           Properties           `xml:"properties"`
}

type NamedVersion struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type ScoutResult struct {
	ModuleResults []Result `json:"modules"`
}

type Result struct {
	Path                string         `json:"path"`
	PluginParentVersion string         `json:"pluginParentVersion"`
	PluginBoms          []NamedVersion `json:"pluginBoms"`
	JenkinsVersion      string         `json:"jenkinsVersion"`
	RecommendedJava     string         `json:"recommendedJava"`
}

func IsJenkinsPluginPom(project Project) bool {
	return project.Parent.GroupId == "org.jenkins-ci.plugins" && project.Parent.ArtifactId == "plugin"
}

func PluginParentVersion(project Project) string {
	if IsJenkinsPluginPom(project) {
		return project.Parent.Version
	}
	return ""
}

func AppliedPluginBoms(project Project) []NamedVersion {
	found := make([]NamedVersion, 0)
	if !IsJenkinsPluginPom(project) {
		return found
	}
	for _, dep := range project.DependencyManagement.Dependencies {
		if dep.GroupId != "io.jenkins.tools.bom" {
			continue
		}
		if !strings.HasPrefix(dep.ArtifactId, "bom-") {
			continue
		}
		if dep.Scope != "import" {
			continue
		}
		if dep.Type != "pom" {
			continue
		}
		if len(dep.Version) == 0 {
			continue
		}
		found = append(found, NamedVersion{Name: dep.ArtifactId, Version: dep.Version})
	}
	return found
}

func JenkinsVersion(project Project) string {
	return project.Properties.JenkinsVersion
}

func RecommendJava(pluginParentVersion string) string {
	components := strings.Split(pluginParentVersion, ".")
	if len(components) < 2 {
		return "8"
	}
	major, err := strconv.Atoi(components[0])
	if err != nil {
		panic(err)
	}
	if major < 4 {
		return "8"
	}
	minor, err := strconv.Atoi(components[1])
	if err != nil {
		panic(err)
	}
	if minor >= 52 {
		return "11"
	}
	return "8"
}
