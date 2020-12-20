package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/kubernetes-sigs/custom-metrics-apiserver/pkg/provider"
	promapi "github.com/prometheus/client_golang/api"
	promv1 "github.com/prometheus/client_golang/api/prometheus/v1"
	prommodel "github.com/prometheus/common/model"
	"gopkg.in/yaml.v2"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/metrics/pkg/apis/external_metrics"
)

type Queries map[string]string

type Config struct {
	PrometheusURL string  `yaml:"url"`
	Queries       Queries `yaml:"queries"`
}

func NewConfig(path string) (*Config, error) {
	c := Config{}
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	if err := yaml.Unmarshal(data, &c); err != nil {
		return nil, err
	}
	return &c, nil
}

type PrometheusProvider struct {
	provider.ExternalMetricsProvider
	config      *Config
	promClient  promv1.API
	metricInfos []provider.ExternalMetricInfo
}

func NewPrometheusProvider(c *Config) (*PrometheusProvider, error) {
	client, err := promapi.NewClient(
		promapi.Config{
			Address: c.PrometheusURL,
		},
	)
	if err != nil {
		return nil, err
	}

	metricInfos := make([]provider.ExternalMetricInfo, len(c.Queries))
	i := 0
	for k := range c.Queries {
		metricInfos[i].Metric = k
		i++
	}

	return &PrometheusProvider{
		config:      c,
		promClient:  promv1.NewAPI(client),
		metricInfos: metricInfos,
	}, nil
}

func (p *PrometheusProvider) GetExternalMetric(_ string, selector labels.Selector, info provider.ExternalMetricInfo) (*external_metrics.ExternalMetricValueList, error) {
	query := p.config.Queries[info.Metric]
	val, _, err := p.promClient.Query(context.TODO(), query, time.Now())
	if err != nil {
		return nil, err
	}

	samples, ok := val.(prommodel.Vector)
	if !ok {
		return nil, fmt.Errorf("couldn't get metrics")
	}

	items := make([]external_metrics.ExternalMetricValue, 0, len(samples))
	for i := range samples {
		ls := labels.Set{}
		for k, v := range samples[i].Metric {
			ls[string(k)] = string(v)
		}

		if !selector.Matches(ls) {
			continue
		}

		item := external_metrics.ExternalMetricValue{
			MetricName:   info.Metric,
			Value:        resource.MustParse(samples[i].Value.String()),
			MetricLabels: map[string]string(ls),
			Timestamp: metav1.Time{
				Time: samples[i].Timestamp.Time(),
			},
		}
		items = append(items, item)
	}

	return &external_metrics.ExternalMetricValueList{
		Items: items,
	}, nil
}

func (p *PrometheusProvider) ListAllExternalMetrics() []provider.ExternalMetricInfo {
	return p.metricInfos
}
