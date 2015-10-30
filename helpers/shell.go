package helpers

import (
	"os"
	"fmt"
	"strings"
	"os/exec"
	log "github.com/Sirupsen/logrus"

)


func Execute(command string) (string, error) {
	cmd := exec.Command("/bin/bash", "-c", fmt.Sprintf("''%s''", command))
	log.Println(strings.Join(cmd.Args, " "))
	output, err := cmd.CombinedOutput()
	if err != nil {
		return string(output) + err.Error(), err
	}
	return strings.TrimRight(string(output), "\n"), nil
}


func System(command string) bool {
	cmd := exec.Command("/bin/bash", "-c", fmt.Sprintf("''%s''", command))
	log.Println(strings.Join(cmd.Args, " "))
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	err := cmd.Run()
	if err != nil {
		return false
	}
	return true
}


func ExecRes(cmd string, a ...interface{}) map[string]string  {
	out, err := Execute(fmt.Sprintf(cmd, a...))
	log.Println("exec", out, err)
	if err != nil {
		return map[string]string{
			"status": "error", "message": out + err.Error()}
	}
	return map[string]string{"status": "ok", "message": out}
}


func FileExists(filename string) bool {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return false
	}
	return true
}


func WriteFile(filename string) bool {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return false
	}
	return true
}
