package gcloud

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/chanwit/meshrun/kubectl"
	servingv1alpha1 "github.com/knative/serving/pkg/apis/serving/v1alpha1"
	"sigs.k8s.io/yaml"
)

func DeployV1alpha1Service(svc servingv1alpha1.Service) {
	spec := svc.Spec.DeprecatedRunLatest.Configuration.DeprecatedRevisionTemplate.Spec
	container := spec.DeprecatedContainer

	args := []string{"beta", "run", "deploy"}
	args = append(args, svc.ObjectMeta.Name)
	args = append(args, fmt.Sprintf("--namespace=%s", svc.ObjectMeta.Namespace))
	args = append(args, fmt.Sprintf("--region=%s", svc.ObjectMeta.Labels["cloud.googleapis.com/location"]))
	if v, exist := svc.ObjectMeta.Labels["meshrun.io/connectivity"]; exist {
		args = append(args, fmt.Sprintf("--connectivity=%s", v))
	}
	if svc.ObjectMeta.Labels["meshrun.io/allow-unauthenticated"] == "true" {
		args = append(args, "--allow-unauthenticated")
	}
	if svc.ObjectMeta.Labels["meshrun.io/async"] == "true" {
		args = append(args, "--async")
	}
	args = append(args, fmt.Sprintf("--image=%s", container.Image))
	if !container.Resources.Limits.Cpu().IsZero() {
		args = append(args, fmt.Sprintf("--cpu=%s", container.Resources.Limits.Cpu().String()))
	}
	if !container.Resources.Limits.Memory().IsZero() {
		args = append(args, fmt.Sprintf("--memory=%s", container.Resources.Limits.Memory().String()))
	}
	args = append(args, fmt.Sprintf("--concurrency=%d", spec.ContainerConcurrency))
	if spec.TimeoutSeconds != nil {
		args = append(args, fmt.Sprintf("--timeout=%d", *spec.TimeoutSeconds))
	}

	env := []string{}
	for _, e := range container.Env {
		env = append(env, e.Name+"="+e.Value)
	}
	args = append(args, "--set-env-vars="+strings.Join(env, ","))

	cmd := exec.Command("gcloud", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()

	RegisterUpstream(svc)
}

func RegisterUpstream(svc servingv1alpha1.Service) error {

	args := []string{"beta", "run", "services", "describe",
		"--namespace=" + svc.ObjectMeta.Namespace,
		"--region=" + svc.ObjectMeta.Labels["cloud.googleapis.com/location"],
		svc.ObjectMeta.Name,
	}

	cmd := exec.Command("gcloud", args...)
	// cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	output, err := cmd.Output()
	if err != nil {
		return err
	}

	var outputSvc servingv1alpha1.Service
	err = yaml.Unmarshal([]byte(output), &outputSvc)
	if err != nil {
		return err
	}

	url := outputSvc.Status.Address.Hostname
	parts := strings.SplitN(url, "//", 2)
	if len(parts) != 2 {
		return fmt.Errorf("URL format is invalid: %s", url)
	}

	hostname := parts[1]
	name := svc.ObjectMeta.Name
	namespace := svc.ObjectMeta.Namespace

	template := `
---
apiVersion: v1
kind: Namespace
metadata:
  name: %s
---
apiVersion: gloo.solo.io/v1
kind: Upstream
metadata:
  labels:
    service: %s
  name: %s
  namespace: %s
spec:
  discoveryMetadata: {}
  upstreamSpec:
    static:
      hosts:
      - addr: %s
        port: 443
`

	return kubectl.Apply(fmt.Sprintf(template, namespace, name, name, namespace, hostname))
}

func DeleteV1alpha1Service(svc servingv1alpha1.Service) {
	args := []string{"beta", "run", "services", "delete"}
	args = append(args, svc.ObjectMeta.Name)
	args = append(args, fmt.Sprintf("--namespace=%s", svc.ObjectMeta.Namespace))
	args = append(args, fmt.Sprintf("--region=%s", svc.ObjectMeta.Labels["cloud.googleapis.com/location"]))
	args = append(args, "--quiet")

	cmd := exec.Command("gcloud", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}
