package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
)

func updateGithubWorkflow(teamId int) error {
	filePath := ".github/workflows/build.yml"

	bytes, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	contents := string(bytes)

	re := regexp.MustCompile("TEAM_NUM: team[0-9]+")
	newString := re.ReplaceAllString(contents, fmt.Sprintf("TEAM_NUM: team%d", teamId))

	err = os.WriteFile(filePath, []byte(newString), 0644)
	if err != nil {
		return err
	}

	fmt.Printf("Update File: %s\n", filePath)

	return nil
}

func main() {
	var (
		g = flag.String("g", "", "Github Name")
		t = flag.Int("t", -1, "Team ID")
	)
	flag.Parse()

	if (*g) == "" || (*t) < 0 {
		fmt.Println("Please input both Github name and Team id")
		return
	}

	err := filepath.Walk("./backend", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if !(filepath.Ext(path) == ".go" || filepath.Ext(path) == ".mod") {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		contents, err := io.ReadAll(file)
		if err != nil {
			return err
		}

		re := regexp.MustCompile("github.com/.+/mecari-build-hackathon-2023")
		newString := re.ReplaceAllString(string(contents), fmt.Sprintf("github.com/%s/mecari-build-hackathon-2023", *g))

		if string(contents) == newString {
			return nil
		}

		err = os.WriteFile(path, []byte(newString), 0644)
		if err != nil {
			return err
		}

		fmt.Printf("Update File: %s\n", info.Name())

		return nil
	})

	if err != nil {
		fmt.Println(err)
	}

	err = updateGithubWorkflow(*t)
	if err != nil {
		fmt.Println(err)
	}
}
