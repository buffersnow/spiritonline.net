package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type projectConfig struct {
	ServiceLong  string
	ServiceShort string
}

func generateProject(cfg projectConfig) error {
	files := map[string]string{
		"src/cmd/%s/main.go":                    "src/tools/pd/templates/main.go.tmpl",
		"src/internal/%s/controllers/.gitkeep":  "src/tools/pd/templates/.gitkeep.tmpl",
		"src/internal/%s/gateway/.gitkeep":      "src/tools/pd/templates/.gitkeep.tmpl",
		"src/internal/%s/handlers/server.go":    "src/tools/pd/templates/handler.go.tmpl",
		"src/internal/%s/protocol/.gitkeep":     "src/tools/pd/templates/.gitkeep.tmpl",
		"src/internal/%s/repositories/.gitkeep": "src/tools/pd/templates/.gitkeep.tmpl",
	}

	for outputFile, tmplPath := range files {
		tmpl, err := template.ParseFiles(tmplPath)
		if err != nil {
			return err
		}

		outputFile = fmt.Sprintf(outputFile, cfg.ServiceLong)
		outputPath := filepath.Dir(outputFile)
		if err = os.MkdirAll(outputPath, os.ModePerm); err != nil {
			return err
		}

		file := filepath.Dir(outputFile)
		file = strings.ReplaceAll(file, "\\", "/")
		if file != outputPath {
			if err = os.MkdirAll(file, os.ModePerm); err != nil {
				return err
			}
		}

		f, err := os.Create(outputFile)
		if err != nil {
			return err
		}
		defer f.Close()

		err = tmpl.Execute(f, cfg)
		if err != nil {
			return err
		}
	}
	return nil
}

func main() {
	long := flag.String("long", "myservice", "full service name")
	short := flag.String("short", "mysvc", "short service name")
	flag.Parse()

	// print(os.Getwd())

	cfg := projectConfig{
		ServiceLong:  *long,
		ServiceShort: *short,
	}
	err := generateProject(cfg)
	if err != nil {
		log.Fatalf("pd: %v", err)
	}
	fmt.Printf("Generation complete: %s\n", cfg.ServiceLong)
}
