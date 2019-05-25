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
)

func Visit_serving_v1beta1_Revision(object servingv1beta1.Revision) {
	fmt.Println(">> Revision v1beta1")
	fmt.Println(object)
}

func Visit_serving_v1alpha1_Revision(object servingv1alpha1.Revision) {
	fmt.Println(">> Revision v1alpha1")
	fmt.Println(object)
}

func Visit_smi_v1alpha1_TrafficSplit(object splitv1alpha1.TrafficSplit) {
	fmt.Println(">> TrafficSplit v1alpha1")
	fmt.Println(object)
}

func main() {
	f, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}

	yamlReader := utilYaml.NewYAMLReader(bufio.NewReader(f))

	for {
		data, err := yamlReader.Read()
		if data == nil {
			break
		}
		if err != nil {
			fmt.Printf("err: %v\n", err)
			return
		}

		var typeMeta metav1.TypeMeta
		err = yaml.Unmarshal(data, &typeMeta)
		if err != nil {
			fmt.Printf("err: %v\n", err)
			return
		}

		key := typeMeta.APIVersion + "/" + typeMeta.Kind
		switch key {
		case "split.smi-spec.io/v1alpha1/TrafficSplit":
			var ts splitv1alpha1.TrafficSplit
			err := yaml.Unmarshal(data, &ts)
			if err != nil {
				fmt.Printf("err: %v\n", err)
				return
			}
			Visit_smi_v1alpha1_TrafficSplit(ts)

		case "serving.knative.dev/v1alpha1/Revision":
			var rev servingv1alpha1.Revision
			err := yaml.Unmarshal(data, &rev)
			if err != nil {
				fmt.Printf("err: %v\n", err)
				return
			}
			Visit_serving_v1alpha1_Revision(rev)

		case "serving.knative.dev/v1beta1/Revision":
			var rev servingv1beta1.Revision
			err := yaml.Unmarshal(data, &rev)
			if err != nil {
				fmt.Printf("err: %v\n", err)
				return
			}
			Visit_serving_v1beta1_Revision(rev)
		}
	}

}
