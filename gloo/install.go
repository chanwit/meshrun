package gloo

import (
	"github.com/chanwit/meshrun/kubectl"
)

func preInstall() error {
	return kubectl.Apply(subGlooYaml_PreInstall)
}

func crdInstall() error {
	return kubectl.Apply(subGlooYaml_CRDInstall)
}

func roleInstall() error {
	return kubectl.Apply(subGlooYaml_RBACInstall)
}

func mainInstall() error {
	return kubectl.Apply(subGlooYaml_MainInstall)
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
