package main

import (
	"fmt"
	"os"
	"os/exec"
)

func pkg(tos, arch string, done chan bool) {
	cmd := exec.Command("make", "package")
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env,
		fmt.Sprintf("GOOS=%s", tos),
		fmt.Sprintf("GOARCH=%s", arch),
		"CGO_ENABLED=0",
	)
	err := cmd.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	done <- true
}

func main() {
	targets := []struct {
		OS   string
		Arch string
	}{
		{"darwin", "amd64"},
		{"darwin", "arm64"},
		{"linux", "amd64"},
		{"linux", "386"},
		{"linux", "arm64"},
		{"linux", "arm"},
		{"windows", "amd64"},
		{"windows", "386"},
	}

	done := make(chan bool)
	for _, t := range targets {
		fmt.Printf("%s-%s\n", t.OS, t.Arch)
		go pkg(t.OS, t.Arch, done)
	}

	for range targets {
		<-done
	}
}