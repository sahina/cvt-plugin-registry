package main

import (
	"context"
	"strings"
	"testing"

	"github.com/hashicorp/go-hclog"
	registrypb "github.com/sahina/cvt/pkg/cvtplugin/pb/registry/v1"
)

func newTestRegistry() *registry {
	return &registry{log: hclog.NewNullLogger()}
}

func TestFetchSchema_ReturnsInlineSpec(t *testing.T) {
	r := newTestRegistry()
	resp, err := r.FetchSchema(context.Background(), &registrypb.FetchSchemaRequest{
		SchemaId: "anything",
		Version:  "1.2.3",
	})
	if err != nil {
		t.Fatalf("FetchSchema: %v", err)
	}
	if !strings.Contains(string(resp.GetSpec()), "openapi: 3.0.0") {
		t.Errorf("spec missing openapi marker: %q", resp.GetSpec())
	}
	if resp.GetResolvedVersion() != "1.2.3" {
		t.Errorf("ResolvedVersion = %q, want %q", resp.GetResolvedVersion(), "1.2.3")
	}
}

func TestFetchSchema_DefaultVersionWhenEmpty(t *testing.T) {
	r := newTestRegistry()
	resp, err := r.FetchSchema(context.Background(), &registrypb.FetchSchemaRequest{SchemaId: "x"})
	if err != nil {
		t.Fatalf("FetchSchema: %v", err)
	}
	if resp.GetResolvedVersion() != "0.1.0" {
		t.Errorf("ResolvedVersion = %q, want %q", resp.GetResolvedVersion(), "0.1.0")
	}
}

func TestRegisterConsumerUsage_Acknowledges(t *testing.T) {
	r := newTestRegistry()
	resp, err := r.RegisterConsumerUsage(context.Background(), &registrypb.RegisterConsumerUsageRequest{
		ConsumerId:    "order-service",
		SchemaId:      "pet-api",
		SchemaVersion: "2.0.0",
		Environment:   "ci",
		Endpoints: []*registrypb.EndpointUsage{
			{Method: "GET", Path: "/pets"},
		},
	})
	if err != nil {
		t.Fatalf("RegisterConsumerUsage: %v", err)
	}
	if !resp.GetAcknowledged() {
		t.Errorf("Acknowledged = false, want true")
	}
}
