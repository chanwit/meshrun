package gloo

import (
	"os"
	"os/exec"
	"strings"
)

func baseInstall(yaml string) error {
	cmd := exec.Command("kubectl","apply","-f","-")
	cmd.Stdin = strings.NewReader(yaml)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func preInstall() error {
	return baseInstall(subGlooYaml_PreInstall)
}

func crdInstall() error {
	return baseInstall(subGlooYaml_CRDInstall)
}

func roleInstall() error {
	return baseInstall(subGlooYaml_RBACInstall)
}

func mainInstall() error {
	return baseInstall(subGlooYaml_MainInstall)
}


func Install() error {
	err := preInstall()
	if err != nil {
		return err
	}

	err = crdInstall()
	if err != nil {
		return err
	}

	err = roleInstall()
	if err != nil {
		return err
	}

	err = mainInstall()
	if err != nil {
		return err
	}

	return nil
}
