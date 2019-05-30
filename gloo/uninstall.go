package gloo

import (
	"os"
	"os/exec"
	"strings"
)

func baseUninstall(yaml string) error {
	cmd := exec.Command("kubectl","delete","-f","-")
	cmd.Stdin = strings.NewReader(yaml)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func preUninstall() error {
	return baseUninstall(subGlooYaml_PreInstall)
}

func crdUninstall() error {
	return baseUninstall(subGlooYaml_CRDInstall)
}

func roleUninstall() error {
	return baseUninstall(subGlooYaml_RBACInstall)
}

func mainUninstall() error {
	return baseUninstall(subGlooYaml_MainInstall)
}

func Uninstall() error {
	err := mainUninstall()
	if err != nil {
		return err
	}

	err = roleUninstall()
	if err != nil {
		return err
	}

	err = crdUninstall()
	if err != nil {
		return err
	}

	err = preUninstall()
	if err != nil {
		return err
	}

	return nil
}
