package main

import (
	"context"
	"github.com/containers/image/v5/types"
	"github.com/opencontainers/go-digest"
	"reflect"
	"testing"
)

func TestConfig_checkAuth(t *testing.T) {
	type fields struct {
		ctx    context.Context
		sysctx *types.SystemContext
	}
	type args struct {
		username     string
		password     string
		registryName string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				ctx:    tt.fields.ctx,
				sysctx: tt.fields.sysctx,
			}
			if err := c.checkAuth(tt.args.username, tt.args.password, tt.args.registryName); (err != nil) != tt.wantErr {
				t.Errorf("checkAuth() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConfig_getDigest(t *testing.T) {
	digestPrometheus, _ := digest.Parse("sha256:0eac377a90d361be9da35b469def699bcd5bb26eab8a6e9068516a9910717d58")

	type fields struct {
		ctx    context.Context
		sysctx *types.SystemContext
	}
	type args struct {
		imageName string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *digest.Digest
	}{
		{
			name: "Prometheus v2.23.0 right Digest",
			fields: fields{
				ctx: context.Background(),
				sysctx: &types.SystemContext{},
			},
			args: args{
				imageName: "quay.io/prometheus/prometheus:v2.23.0",
			},
			want: &digestPrometheus,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				ctx:    tt.fields.ctx,
				sysctx: tt.fields.sysctx,
			}
			if got := c.getDigest(tt.args.imageName); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getDigest() = %v, want %v", got, tt.want)
			}
		})
	}
}
