// Command cvt-plugin-registry is a placeholder CVT RegistryProvider plugin
// that returns a hardcoded inline OpenAPI spec from FetchSchema and logs
// consumer-usage records to stderr. It exists to unblock end-to-end testing
// of CVT core's CLI wiring (sahina/cvt issue #83) before the central API
// registry's HTTP contract is defined. Replace with a real adapter once the
// contract lands.
package main

import (
	"context"

	"github.com/hashicorp/go-hclog"
	"github.com/sahina/cvt/pkg/cvtplugin"
	registrypb "github.com/sahina/cvt/pkg/cvtplugin/pb/registry/v1"
)

const pluginName = "cvt-plugin-registry"
const pluginVersion = "0.1.0"

const inlineSpec = `openapi: 3.0.0
info:
  title: hello
  version: 0.1.0
paths:
  /hello:
    get:
      summary: Hello world
      responses:
        '200':
          description: OK
          content:
            application/json:
              schema:
                type: object
                properties:
                  message:
                    type: string
`

type registry struct {
	log hclog.Logger
}

func (r *registry) FetchSchema(_ context.Context, req *registrypb.FetchSchemaRequest) (*registrypb.FetchSchemaResponse, error) {
	resolved := req.GetVersion()
	if resolved == "" {
		resolved = "0.1.0"
	}
	r.log.Info("fetch_schema",
		"schema_id", req.GetSchemaId(),
		"requested_version", req.GetVersion(),
		"resolved_version", resolved,
	)
	return &registrypb.FetchSchemaResponse{
		Spec:            []byte(inlineSpec),
		ResolvedVersion: resolved,
	}, nil
}

func (r *registry) RegisterConsumerUsage(_ context.Context, req *registrypb.RegisterConsumerUsageRequest) (*registrypb.RegisterConsumerUsageResponse, error) {
	endpoints := make([]map[string]string, 0, len(req.GetEndpoints()))
	for _, e := range req.GetEndpoints() {
		endpoints = append(endpoints, map[string]string{
			"method": e.GetMethod(),
			"path":   e.GetPath(),
		})
	}
	r.log.Info("register_consumer_usage",
		"consumer_id", req.GetConsumerId(),
		"schema_id", req.GetSchemaId(),
		"schema_version", req.GetSchemaVersion(),
		"environment", req.GetEnvironment(),
		"endpoints", endpoints,
	)
	return &registrypb.RegisterConsumerUsageResponse{Acknowledged: true}, nil
}

func main() {
	log := hclog.New(&hclog.LoggerOptions{
		Name:       pluginName,
		Level:      hclog.Info,
		JSONFormat: true,
	})
	cvtplugin.Serve(
		cvtplugin.PluginInfo{Name: pluginName, Version: pluginVersion},
		cvtplugin.WithRegistryProvider(&registry{log: log}),
		cvtplugin.WithLogger(log),
	)
}
