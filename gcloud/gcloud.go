package gcloud

import (
	"os"
	"os/exec"
	"fmt"

	servingv1alpha1 "github.com/knative/serving/pkg/apis/serving/v1alpha1"
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

	cmd := exec.Command("gcloud", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
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