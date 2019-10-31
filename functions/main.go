package functions

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

func FmtPath(path string) (err error) {
	cmd := exec.Command("go", "fmt", path+"/", "./...")
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		return
	}
	return
}

func ImportsPath(path string) (err error) {
	cmd := exec.Command("goimports", "-w", path+"/")
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		return
	}
	return
}

func CurrentPath() (path string, err error) {
	path, err = os.Getwd()
	if err != nil {
		return
	}
	return
}
