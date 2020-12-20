module github.com/uesyn/sample-metrics-server

go 1.15

require (
	github.com/kubernetes-sigs/custom-metrics-apiserver v0.0.0-20201216091021-1b9fa998bbaa
	github.com/prometheus/client_golang v1.9.0
	github.com/prometheus/common v0.15.0
	gopkg.in/yaml.v2 v2.4.0
	k8s.io/apimachinery v0.20.1
	k8s.io/component-base v0.20.1
	k8s.io/klog/v2 v2.4.0
	k8s.io/metrics v0.20.1
)
