package utils

import (
	"fmt"
	"os/exec"
	"path"
)

const (
	localPath = "./scripts/"
)

func ExecuteCmd(name string, env map[string]string, log Logger) {
	pt := path.Join(localPath, name)
	cmd := exec.Command(pt)

	cenv := cmd.Env
	for k, v := range env {
		s := fmt.Sprintf("%s=%s", k, v)
		cenv = append(cenv, s)
	}
	cmd.Env = cenv

	log.Info("Executing command: %s, with environmant: %v", pt, cenv)

	output, err := cmd.Output()
	if err != nil {
		log.Error("Command %s executing error: %v", name, err)
	} else {
		log.Info("Command %s executed", pt)
		log.Info("Output: %s", string(output))
	}
}
