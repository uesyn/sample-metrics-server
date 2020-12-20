package main

import (
	"flag"
	"os"

	basecmd "github.com/kubernetes-sigs/custom-metrics-apiserver/pkg/cmd"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/component-base/logs"
	"k8s.io/klog/v2"
)

var (
	configPath string
)

func main() {
	logs.InitLogs()
	defer logs.FlushLogs()

	cmd := basecmd.AdapterBase{}
	cmd.Flags().StringVar(&configPath, "config", "", "config path")
	cmd.Flags().AddGoFlagSet(flag.CommandLine)
	cmd.Flags().Parse(os.Args)

	config, err := NewConfig(configPath)
	if err != nil {
		klog.Fatal(err)
	}

	provider, err := NewPrometheusProvider(config)
	if err != nil {
		klog.Fatal(err)
	}
	cmd.WithExternalMetrics(provider)

	klog.Infof("starting sample metrics server")
	if err := cmd.Run(wait.NeverStop); err != nil {
		klog.Fatalf("unable to run custom metrics adapter: %v", err)
	}
}
