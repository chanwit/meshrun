package gloo

import (
	"github.com/chanwit/meshrun/kubectl"
)

func preUninstall() error {
	return kubectl.Delete(subGlooYaml_PreInstall)
}

func crdUninstall() error {
	return kubectl.Delete(subGlooYaml_CRDInstall)
}

func roleUninstall() error {
	return kubectl.Delete(subGlooYaml_RBACInstall)
}

func mainUninstall() error {
	return kubectl.Delete(subGlooYaml_MainInstall)
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
