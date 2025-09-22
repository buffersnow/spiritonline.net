package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"strings"
	"syscall"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Programs []struct {
		Name string   `yaml:"name"`
		Args []string `yaml:"args"`
	} `yaml:"programs"`
}

var longestName = 0

var colorsList = map[string]string{
	"gsp":     "\033[38;5;62m",  /// Hex: #5f5fd7
	"myspace": "\033[38;5;208m", /// Hex: #ff8700
	"proxy":   "\033[38;5;93m",  /// Hex: #8700ff
	"qr":      "\033[38;5;191m", /// Hex: #d7ff5f
	"router":  "\033[38;5;111m", /// Hex: #87afff
	"wfc":     "\033[38;5;212m", /// Hex: #ff87d7
}

func runService(name, command string, args ...string) {
	cmd := exec.Command(command, args...)
	stdout, _ := cmd.StdoutPipe()

	if err := cmd.Start(); err != nil {
		panic(fmt.Sprintf("failed to start %s: %v\n", name, err))
	}

	go func(gname string, reader io.ReadCloser) {
		scanner := bufio.NewScanner(reader)
		for scanner.Scan() {
			txt := scanner.Text()
			if strings.HasPrefix(txt, "\033[38;5;61m") {
				continue // skip version text
			}

			color, ok := colorsList[gname]
			if !ok {
				color = "\033[38;5;61m"
			}

			fmt.Printf("%s%+*s | %s\n", color, longestName, gname, txt)
		}
	}(name, stdout)

	if err := cmd.Wait(); err != nil {
		fmt.Printf("%s exited with err: %v\n", name, err)
	}
}

func main() {

	file := "taxi.yaml"
	if len(os.Args) > 1 {
		file = "taxi." + os.Args[1] + ".yaml"
	}

	yamlFile, err := os.ReadFile(file)
	if err != nil {
		panic(fmt.Sprintf("failed to load yaml file: %v", err))
	}

	var cfg Config
	if err := yaml.Unmarshal([]byte(yamlFile), &cfg); err != nil {
		panic(fmt.Sprintf("failed to unmarshal yaml: %v", err))
	}

	fileExt := ".exe"
	if runtime.GOOS != "windows" {
		fileExt = ".lxb"
	}

	for _, prg := range cfg.Programs {
		if longestName < len(prg.Name) {
			longestName = len(prg.Name)
		}
	}

	for _, prg := range cfg.Programs {
		go runService(prg.Name, fmt.Sprintf("bin/%s%s", prg.Name, fileExt), prg.Args...)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	if sig := <-c; sig != nil {
		println("quitting")
		os.Exit(0)
	}
}
