package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"

	splitv1alpha1 "github.com/deislabs/smi-sdk-go/pkg/apis/split/v1alpha1"
	servingv1alpha1 "github.com/knative/serving/pkg/apis/serving/v1alpha1"
	servingv1beta1 "github.com/knative/serving/pkg/apis/serving/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilYaml "k8s.io/apimachinery/pkg/util/yaml"
	"sigs.k8s.io/yaml"

	"github.com/urfave/cli"
)

func Apply_serving_v1alpha1_Service(object servingv1alpha1.Service) {
	// fmt.Println(">> Service v1alpha1")
	// spew.Dump(object)
	apply_v1alpha1_Service_To_Gcloud_Deploy(object)
}

func Delete_serving_v1alpha1_Service(object servingv1alpha1.Service) {
	// fmt.Println(">> Service v1alpha1")
	// spew.Dump(object)
	delete_v1alpha1_Service_To_Gcloud_Deploy(object)
}

func apply_v1alpha1_Service_To_Gcloud_Deploy(svc servingv1alpha1.Service) {
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

func delete_v1alpha1_Service_To_Gcloud_Deploy(svc servingv1alpha1.Service) {
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

func Delete_serving_v1beta1_Service(object servingv1beta1.Service) {
	fmt.Println(">> Service v1beta1")
	// spew.Dump(object)
	// delete_v1alpha1_Service_To_Gcloud_Deploy(object)
}

func Apply_serving_v1beta1_Service(object servingv1beta1.Service) {
	fmt.Println(">> Service v1beta1")
	//spew.Dump(object)
}

func Apply_serving_v1beta1_Revision(object servingv1beta1.Revision) {
	fmt.Println(">> Revision v1beta1")
	//spew.Dump(object)
}

func Apply_serving_v1alpha1_Revision(object servingv1alpha1.Revision) {
	fmt.Println(">> Revision v1alpha1")
	//spew.Dump(object)
}

func Apply_smi_v1alpha1_TrafficSplit(object splitv1alpha1.TrafficSplit) {
	fmt.Println(">> TrafficSplit v1alpha1")
	//spew.Dump(object)
}

func main() {
	app := cli.NewApp()
	app.Commands = []cli.Command{
		{
			Name:  "apply",
			Usage: "apply",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "filename, f",
					Usage: "YAML filename",
				},
			},
			Action: func(c *cli.Context) error {
				filename := c.String("filename")
				return apply(filename)
			},
		},
		{
			Name:  "delete",
			Usage: "delete",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "filename, f",
					Usage: "YAML filename",
				},
			},
			Action: func(c *cli.Context) error {
				filename := c.String("filename")
				return doDelete(filename)
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		panic(err)
	}
}

func doDelete(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}

	yamlReader := utilYaml.NewYAMLReader(bufio.NewReader(f))

	for {
		data, err := yamlReader.Read()
		if data == nil {
			break
		}
		if err != nil {
			// fmt.Printf("err: %v\n", err)
			return err
		}

		var typeMeta metav1.TypeMeta
		err = yaml.Unmarshal(data, &typeMeta)
		if err != nil {
			// fmt.Printf("err: %v\n", err)
			return err
		}

		key := typeMeta.APIVersion + "/" + typeMeta.Kind
		switch key {
		case "serving.knative.dev/v1alpha1/Service":
			var svc servingv1alpha1.Service
			err := yaml.Unmarshal(data, &svc)
			if err != nil {
				// fmt.Printf("err: %v\n", err)
				return err
			}
			Delete_serving_v1alpha1_Service(svc)

		case "serving.knative.dev/v1beta1/Service":
			var svc servingv1beta1.Service
			err := yaml.Unmarshal(data, &svc)
			if err != nil {
				// fmt.Printf("err: %v\n", err)
				return err
			}
			Delete_serving_v1beta1_Service(svc)
		}
	}

	return nil
}

func apply(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}

	yamlReader := utilYaml.NewYAMLReader(bufio.NewReader(f))

	for {
		data, err := yamlReader.Read()
		if data == nil {
			break
		}
		if err != nil {
			// fmt.Printf("err: %v\n", err)
			return err
		}

		var typeMeta metav1.TypeMeta
		err = yaml.Unmarshal(data, &typeMeta)
		if err != nil {
			// fmt.Printf("err: %v\n", err)
			return err
		}

		key := typeMeta.APIVersion + "/" + typeMeta.Kind
		switch key {
		case "split.smi-spec.io/v1alpha1/TrafficSplit":
			var ts splitv1alpha1.TrafficSplit
			err := yaml.Unmarshal(data, &ts)
			if err != nil {
				// fmt.Printf("err: %v\n", err)
				return err
			}
			Apply_smi_v1alpha1_TrafficSplit(ts)

		case "serving.knative.dev/v1alpha1/Service":
			var svc servingv1alpha1.Service
			err := yaml.Unmarshal(data, &svc)
			if err != nil {
				// fmt.Printf("err: %v\n", err)
				return err
			}
			Apply_serving_v1alpha1_Service(svc)

		case "serving.knative.dev/v1beta1/Service":
			var svc servingv1beta1.Service
			err := yaml.Unmarshal(data, &svc)
			if err != nil {
				// fmt.Printf("err: %v\n", err)
				return err
			}
			Apply_serving_v1beta1_Service(svc)

		case "serving.knative.dev/v1alpha1/Revision":
			var rev servingv1alpha1.Revision
			err := yaml.Unmarshal(data, &rev)
			if err != nil {
				fmt.Printf("err: %v\n", err)
				return err
			}
			Apply_serving_v1alpha1_Revision(rev)

		case "serving.knative.dev/v1beta1/Revision":
			var rev servingv1beta1.Revision
			err := yaml.Unmarshal(data, &rev)
			if err != nil {
				// fmt.Printf("err: %v\n", err)
				return err
			}
			Apply_serving_v1beta1_Revision(rev)
		}
	}

	return nil
}
