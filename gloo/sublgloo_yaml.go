package gloo

var subGlooYaml_PreInstall = `
---
apiVersion: v1
kind: Namespace
metadata:
  name: gloo-system
  labels:
    app: gloo
  annotations:
    "helm.sh/hook": pre-install
`
var subGlooYaml_CRDInstall = `
---
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  name: settings.gloo.solo.io
  annotations:
    "helm.sh/hook": crd-install
  labels:
    gloo: settings
spec:
  group: gloo.solo.io
  names:
    kind: Settings
    listKind: SettingsList
    plural: settings
    shortNames:
      - st
  scope: Namespaced
  version: v1
---
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  name: gateways.gateway.solo.io
  annotations:
    "helm.sh/hook": crd-install
spec:
  group: gateway.solo.io
  names:
    kind: Gateway
    listKind: GatewayList
    plural: gateways
    shortNames:
      - gw
    singular: gateway
  scope: Namespaced
  version: v1
---
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  name: virtualservices.gateway.solo.io
  annotations:
    "helm.sh/hook": crd-install
spec:
  group: gateway.solo.io
  names:
    kind: VirtualService
    listKind: VirtualServiceList
    plural: virtualservices
    shortNames:
      - vs
    singular: virtualservice
  scope: Namespaced
  version: v1
---
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  name: proxies.gloo.solo.io
  annotations:
    "helm.sh/hook": crd-install
spec:
  group: gloo.solo.io
  names:
    kind: Proxy
    listKind: ProxyList
    plural: proxies
    shortNames:
      - px
    singular: proxy
  scope: Namespaced
  version: v1
---
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  name: upstreams.gloo.solo.io
  annotations:
    "helm.sh/hook": crd-install
spec:
  group: gloo.solo.io
  names:
    kind: Upstream
    listKind: UpstreamList
    plural: upstreams
    shortNames:
      - us
    singular: upstream
  scope: Namespaced
  version: v1
---
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  name: upstreamgroups.gloo.solo.io
  annotations:
    "helm.sh/hook": crd-install
spec:
  group: gloo.solo.io
  names:
    kind: UpstreamGroup
    listKind: UpstreamGroupList
    plural: upstreamgroups
    shortNames:
      - ug
    singular: upstreamgroup
  scope: Namespaced
  version: v1
`
var subGlooYaml_RBACInstall = `
---
kind: ClusterRole
apiVersion: rbac.authorization.k8s.io/v1
metadata:
    name: gloo-role-gateway
    labels:
        app: gloo
        gloo: rbac
rules:
- apiGroups: [""]
  resources: ["pods", "services", "secrets", "endpoints", "configmaps"]
  verbs: ["*"]
- apiGroups: [""]
  resources: ["namespaces"]
  verbs: ["get", "list", "watch"]
- apiGroups: ["apiextensions.k8s.io"]
  resources: ["customresourcedefinitions"]
  verbs: ["get", "create"]
- apiGroups: ["gloo.solo.io"]
  resources: ["settings", "upstreams","upstreamgroups", "proxies","virtualservices"]
  verbs: ["*"]
- apiGroups: ["gateway.solo.io"]
  resources: ["virtualservices", "gateways"]
  verbs: ["*"]
---
kind: ClusterRoleBinding
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: gloo-role-binding-gateway-gloo-system
  labels:
    app: gloo
    gloo: rbac
subjects:
- kind: ServiceAccount
  name: default
  namespace: gloo-system
roleRef:
  kind: ClusterRole
  name: gloo-role-gateway
  apiGroup: rbac.authorization.k8s.io
`
var subGlooYaml_MainInstall = `
---
apiVersion: v1
kind: Service
metadata:
  labels:
    app: gloo
    gloo: gloo
  name: gloo
  namespace: gloo-system
spec:
  ports:
  - name: grpc
    port: 9977
    protocol: TCP
  selector:
    gloo: gloo
  type: LoadBalancer
---
apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  labels:
    app: gloo
    gloo: gloo
  name: gloo
  namespace: gloo-system
spec:
  replicas: 1
  selector:
    matchLabels:
      gloo: gloo
  template:
    metadata:
      labels:
        gloo: gloo
      annotations:
        prometheus.io/path: /metrics
        prometheus.io/port: "9091"
        prometheus.io/scrape: "true"
    spec:
      containers:
      - image: "quay.io/solo-io/gloo:0.13.29"
        imagePullPolicy: Always
        name: gloo
        resources:
          requests:
            cpu: 1
            memory: 256Mi
        securityContext:
          readOnlyRootFilesystem: true
          allowPrivilegeEscalation: false
          runAsNonRoot: true
          runAsUser: 10101
          capabilities:
            drop:
            - ALL
        ports:
        - containerPort: 9977
          name: grpc
          protocol: TCP
        env:
          - name: POD_NAMESPACE
            valueFrom:
              fieldRef:
                fieldPath: metadata.namespace
          - name: START_STATS_SERVER
            value: "true"
---
apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  labels:
    app: gloo
    gloo: gateway
  name: gateway
  namespace: gloo-system
spec:
  replicas: 1
  selector:
    matchLabels:
      gloo: gateway
  template:
    metadata:
      labels:
        gloo: gateway
      annotations:
        prometheus.io/path: /metrics
        prometheus.io/port: "9091"
        prometheus.io/scrape: "true"
    spec:
      containers:
      - image: "quay.io/solo-io/gateway:0.13.29"
        imagePullPolicy: Always
        name: gateway
        securityContext:
          readOnlyRootFilesystem: true
          allowPrivilegeEscalation: false
          runAsNonRoot: true
          runAsUser: 10101
          capabilities:
            drop:
            - ALL
        env:
          - name: POD_NAMESPACE
            valueFrom:
              fieldRef:
                fieldPath: metadata.namespace
          - name: START_STATS_SERVER
            value: "true"
`
