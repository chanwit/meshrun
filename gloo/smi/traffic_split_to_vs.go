package smi

import (
	"bytes"
	"text/template"

	splitv1alpha1 "github.com/deislabs/smi-sdk-go/pkg/apis/split/v1alpha1"
)

func ConvertTrafficSplitToVirtualService(object splitv1alpha1.TrafficSplit) string {
	vs := `
---
apiVersion: gateway.solo.io/v1
kind: VirtualService
metadata:
  name: {{ .ObjectMeta.Name }}
  namespace: {{ .ObjectMeta.Namespace }}
spec:
  virtualHost:
    domains:
    - {{ .Spec.Service }}{{ index .Labels "meshrun.io/domain-suffix" }}
    routes:
    - matcher:
        prefix: /
      routeAction:
        multi:
          destinations:
          {{- range $index, $element := .Spec.Backends }}
          - weight: {{ $element.Weight }}
            destination:
              upstream:
                name: {{ $element.Service }}
                namespace: {{ $.ObjectMeta.Namespace }}
          {{- end }}
`

	t := template.Must(template.New("vs").Parse(vs))
	var result bytes.Buffer
	t.Execute(&result, object)

	return result.String()
}
