package main

import (
	"bufio"
	"fmt"
	"os"

	splitv1alpha1 "github.com/deislabs/smi-sdk-go/pkg/apis/split/v1alpha1"
	servingv1alpha1 "github.com/knative/serving/pkg/apis/serving/v1alpha1"
	servingv1beta1 "github.com/knative/serving/pkg/apis/serving/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilYaml "k8s.io/apimachinery/pkg/util/yaml"
	"sigs.k8s.io/yaml"

	"github.com/urfave/cli"

	"github.com/chanwit/meshrun/gcloud"
	"github.com/chanwit/meshrun/gloo"
)

func Apply_serving_v1alpha1_Service(object servingv1alpha1.Service) {
	// fmt.Println(">> Service v1alpha1")
	// spew.Dump(object)
	gcloud.DeployV1alpha1Service(object)
}

func Delete_serving_v1alpha1_Service(object servingv1alpha1.Service) {
	// fmt.Println(">> Service v1alpha1")
	// spew.Dump(object)
	gcloud.DeleteV1alpha1Service(object)
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
		{
			Name: "gloo",
			Usage: "gloo install and uninstall",
			Subcommands: []cli.Command {
				{
					Name: "install",
					Usage: "install Gloo",
					Action: func(c *cli.Context) error {
						return gloo.Install()
					},
				},
				{
					Name: "uninstall",
					Usage: "uninstall Gloo",
					Action: func(c *cli.Context) error {
						return gloo.Uninstall()
					},
				},
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
			return err
		}

		var typeMeta metav1.TypeMeta
		err = yaml.Unmarshal(data, &typeMeta)
		if err != nil {
			return err
		}

		key := typeMeta.APIVersion + "/" + typeMeta.Kind
		switch key {
		case "serving.knative.dev/v1alpha1/Service":
			var svc servingv1alpha1.Service
			err := yaml.Unmarshal(data, &svc)
			if err != nil {
				return err
			}
			Delete_serving_v1alpha1_Service(svc)

		case "serving.knative.dev/v1beta1/Service":
			var svc servingv1beta1.Service
			err := yaml.Unmarshal(data, &svc)
			if err != nil {
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
			return err
		}

		var typeMeta metav1.TypeMeta
		err = yaml.Unmarshal(data, &typeMeta)
		if err != nil {
			return err
		}

		key := typeMeta.APIVersion + "/" + typeMeta.Kind
		switch key {
		case "split.smi-spec.io/v1alpha1/TrafficSplit":
			var ts splitv1alpha1.TrafficSplit
			err := yaml.Unmarshal(data, &ts)
			if err != nil {
				return err
			}
			Apply_smi_v1alpha1_TrafficSplit(ts)

		case "serving.knative.dev/v1alpha1/Service":
			var svc servingv1alpha1.Service
			err := yaml.Unmarshal(data, &svc)
			if err != nil {
				return err
			}
			Apply_serving_v1alpha1_Service(svc)

		case "serving.knative.dev/v1beta1/Service":
			var svc servingv1beta1.Service
			err := yaml.Unmarshal(data, &svc)
			if err != nil {
				return err
			}
			Apply_serving_v1beta1_Service(svc)

		case "serving.knative.dev/v1alpha1/Revision":
			var rev servingv1alpha1.Revision
			err := yaml.Unmarshal(data, &rev)
			if err != nil {
				return err
			}
			Apply_serving_v1alpha1_Revision(rev)

		case "serving.knative.dev/v1beta1/Revision":
			var rev servingv1beta1.Revision
			err := yaml.Unmarshal(data, &rev)
			if err != nil {
				return err
			}
			Apply_serving_v1beta1_Revision(rev)
		}
	}

	return nil
}
