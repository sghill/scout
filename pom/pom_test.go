package pom

import (
	"encoding/xml"
	"testing"

	"github.com/stretchr/testify/require"
)

var plugins = []string{
	`
	<project>
		<parent>
			<groupId>org.jenkins-ci.plugins</groupId>
			<artifactId>plugin</artifactId>
			<version>4.51</version>
		</parent>
	</project>`,
	`
	<project>
		<parent>
			<groupId>org.jenkins-ci.plugins</groupId>
			<artifactId>plugin</artifactId>
			<version>4.66</version>
			<relativePath />
		</parent>
	</project>`,
}

var notPlugins = []string{
	``,
	`<project></project>`,
	`
	<project>
	 	<parent>
			<groupId>net.sghill.parent</groupId>
			<artifactId>plugin</artifactId>
			<version>4.51</version>
		</parent>
	</project>`,
	`
	<project>
	 	<parent>
			<groupId>org.jenkins-ci.plugins</groupId>
			<artifactId>something-else</artifactId>
			<version>4.51</version>
		</parent>
	</project>`,
	`}not=xml{`,
}

func TestIsJenkinsPluginPom(t *testing.T) {
	for _, pom := range plugins {
		project := Project{}
		xml.Unmarshal([]byte(pom), &project)

		actual := IsJenkinsPluginPom(project)

		require.True(t, actual)
	}
}

func TestIsNotJenkinsPluginPom(t *testing.T) {
	for _, pom := range notPlugins {
		project := Project{}
		xml.Unmarshal([]byte(pom), &project)

		actual := IsJenkinsPluginPom(project)

		require.False(t, actual)
	}
}

func TestJenkinsPluginParentVersion(t *testing.T) {
	expected := []string{"4.51", "4.66"}
	for i, pom := range plugins {
		project := Project{}
		xml.Unmarshal([]byte(pom), &project)

		actual := PluginParentVersion(project)

		require.Equal(t, expected[i], actual)
	}
}

func TestJenkinsPluginParentVersionNonPlugin(t *testing.T) {
	for _, pom := range notPlugins {
		project := Project{}
		xml.Unmarshal([]byte(pom), &project)

		actual := PluginParentVersion(project)

		require.Empty(t, actual)
	}
}

var boms = []string{
	`
	<project>
		<parent>
			<groupId>org.jenkins-ci.plugins</groupId>
			<artifactId>plugin</artifactId>
			<version>4.51</version>
		</parent>
		<dependencyManagement>
			<dependencies>
				<dependency>
					<groupId>io.jenkins.tools.bom</groupId>
					<artifactId>bom-2.346.x</artifactId>
					<version>1607.va_c1576527071</version>
					<scope>import</scope>
					<type>pom</type>
				</dependency>
			</dependencies>
		</dependencyManagement>
	</project>`, `
	<project>
		<parent>
			<groupId>org.jenkins-ci.plugins</groupId>
			<artifactId>plugin</artifactId>
			<version>4.51</version>
		</parent>
		<dependencyManagement>
			<dependencies>
				<dependency>
					<groupId>io.jenkins.tools.bom</groupId>
					<artifactId>bom-2.346.x</artifactId>
					<version>1607.va_c1576527071</version>
					<scope>import</scope>
					<type>pom</type>
				</dependency>
				<dependency>
					<groupId>io.jenkins.tools.bom</groupId>
					<artifactId>bom-2.375.x</artifactId>
					<version>1948.veb_1fd345d3a_e</version>
					<scope>import</scope>
					<type>pom</type>
				</dependency>
			</dependencies>
		</dependencyManagement>
	</project>`,
	`
	<project>
		<parent>
			<groupId>org.jenkins-ci.plugins</groupId>
			<artifactId>plugin</artifactId>
			<version>4.51</version>
		</parent>
		<dependencyManagement>
			<dependencies>
				<dependency>
					<groupId>io.jenkins.tools.bom</groupId>
					<artifactId>bom-weekly</artifactId>
					<version>1948.veb_1fd345d3a_e</version>
					<scope>import</scope>
					<type>pom</type>
				</dependency>
			</dependencies>
		</dependencyManagement>
	</project>`,
}

var almostBoms = []string{
	// type not specified
	`
	<project>
		<parent>
			<groupId>org.jenkins-ci.plugins</groupId>
			<artifactId>plugin</artifactId>
			<version>4.51</version>
		</parent>
		<dependencyManagement>
			<dependencies>
				<dependency>
					<groupId>io.jenkins.tools.bom</groupId>
					<artifactId>bom-2.346.x</artifactId>
					<version>1607.va_c1576527071</version>
					<scope>import</scope>
				</dependency>
			</dependencies>
		</dependencyManagement>
	</project>`,
	// wrong type
	`
	<project>
		<parent>
			<groupId>org.jenkins-ci.plugins</groupId>
			<artifactId>plugin</artifactId>
			<version>4.51</version>
		</parent>
		<dependencyManagement>
			<dependencies>
				<dependency>
					<groupId>io.jenkins.tools.bom</groupId>
					<artifactId>bom-2.346.x</artifactId>
					<version>1607.va_c1576527071</version>
					<scope>import</scope>
					<type>jar</type>
				</dependency>
			</dependencies>
		</dependencyManagement>
	</project>`,
	// no scope
	`
	<project>
		<parent>
			<groupId>org.jenkins-ci.plugins</groupId>
			<artifactId>plugin</artifactId>
			<version>4.51</version>
		</parent>
		<dependencyManagement>
			<dependencies>
				<dependency>
					<groupId>io.jenkins.tools.bom</groupId>
					<artifactId>bom-2.346.x</artifactId>
					<version>1607.va_c1576527071</version>
					<type>jar</type>
				</dependency>
			</dependencies>
		</dependencyManagement>
	</project>`,
	// wrong scope
	`
	<project>
		<parent>
			<groupId>org.jenkins-ci.plugins</groupId>
			<artifactId>plugin</artifactId>
			<version>4.51</version>
		</parent>
		<dependencyManagement>
			<dependencies>
				<dependency>
					<groupId>io.jenkins.tools.bom</groupId>
					<artifactId>bom-2.346.x</artifactId>
					<version>1607.va_c1576527071</version>
					<scope>runtime</scope>
					<type>jar</type>
				</dependency>
			</dependencies>
		</dependencyManagement>
	</project>`,
	// no version
	`
	<project>
		<parent>
			<groupId>org.jenkins-ci.plugins</groupId>
			<artifactId>plugin</artifactId>
			<version>4.51</version>
		</parent>
		<dependencyManagement>
			<dependencies>
				<dependency>
					<groupId>io.jenkins.tools.bom</groupId>
					<artifactId>bom-2.346.x</artifactId>
					<scope>import</scope>
					<type>pom</type>
				</dependency>
			</dependencies>
		</dependencyManagement>
	</project>`,
	// wrong artifactId
	`
	<project>
		<parent>
			<groupId>org.jenkins-ci.plugins</groupId>
			<artifactId>plugin</artifactId>
			<version>4.51</version>
		</parent>
		<dependencyManagement>
			<dependencies>
				<dependency>
					<groupId>io.jenkins.tools.bom</groupId>
					<artifactId>bad-2.346.x</artifactId>
					<version>1607.va_c1576527071</version>
					<scope>import</scope>
					<type>pom</type>
				</dependency>
			</dependencies>
		</dependencyManagement>
	</project>`,
	// wrong group
	`
	<project>
		<parent>
			<groupId>org.jenkins-ci.plugins</groupId>
			<artifactId>plugin</artifactId>
			<version>4.51</version>
		</parent>
		<dependencyManagement>
			<dependencies>
				<dependency>
					<groupId>io.jenkins.tools.bad</groupId>
					<artifactId>bom-2.346.x</artifactId>
					<version>1607.va_c1576527071</version>
					<scope>import</scope>
					<type>pom</type>
				</dependency>
			</dependencies>
		</dependencyManagement>
	</project>`,
	// none listed
	`
	<project>
		<parent>
			<groupId>org.jenkins-ci.plugins</groupId>
			<artifactId>plugin</artifactId>
			<version>4.51</version>
		</parent>
		<dependencyManagement>
			<dependencies>
			</dependencies>
		</dependencyManagement>
	</project>`,
	// list missing
	`
	<project>
		<parent>
			<groupId>org.jenkins-ci.plugins</groupId>
			<artifactId>plugin</artifactId>
			<version>4.51</version>
		</parent>
		<dependencyManagement>
		</dependencyManagement>
	</project>`,
	// section missing
	`
	<project>
		<parent>
			<groupId>org.jenkins-ci.plugins</groupId>
			<artifactId>plugin</artifactId>
			<version>4.51</version>
		</parent>
	</project>`,
}

func TestAppliedPluginBoms(t *testing.T) {
	expected := [][]NamedVersion{
		{
			{Name: "bom-2.346.x", Version: "1607.va_c1576527071"},
		},
		{
			{Name: "bom-2.346.x", Version: "1607.va_c1576527071"},
			{Name: "bom-2.375.x", Version: "1948.veb_1fd345d3a_e"},
		},
		{
			{Name: "bom-weekly", Version: "1948.veb_1fd345d3a_e"},
		},
	}
	for i, pom := range boms {
		project := Project{}
		xml.Unmarshal([]byte(pom), &project)

		actual := AppliedPluginBoms(project)

		require.Equal(t, expected[i], actual)
	}
}

func TestAppliedPluginBomsSkips(t *testing.T) {
	bad := []string{}
	for _, almost := range almostBoms {
		bad = append(bad, almost)
	}
	for _, nonPlugin := range notPlugins {
		bad = append(bad, nonPlugin)
	}
	for _, pom := range bad {
		project := Project{}
		xml.Unmarshal([]byte(pom), &project)

		actual := AppliedPluginBoms(project)

		require.Empty(t, actual)
	}
}

func TestJenkinsVersion(t *testing.T) {
	pom := `
		<project>
			<properties>
				<jenkins.version>2.346.1</jenkins.version>
			</properties>
		</project>
	`
	project := Project{}
	xml.Unmarshal([]byte(pom), &project)

	actual := JenkinsVersion(project)

	require.Equal(t, "2.346.1", actual)
}

func TestRecommendJava(t *testing.T) {
	eleven := []string{"4.52", "4.66"}
	for _, v := range eleven {
		require.Equal(t, "11", RecommendJava(v))
	}

	eight := []string{"4.51", "3.456", ""}
	for _, v := range eight {
		require.Equal(t, "8", RecommendJava(v))
	}
}
