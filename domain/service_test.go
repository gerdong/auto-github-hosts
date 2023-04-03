package domain

import (
	"github.com/gerdong/auto-github-hosts/config"
	"os"
	"testing"
)

func TestUpdateHosts(t *testing.T) {
	os.Args = append(os.Args, "-c", "../config/auto-github-hosts.toml")
	config.Init()

	tests := []struct {
		name string
	}{
		{name: "test_UpdateHosts_01"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := UpdateHosts(); got == "" {
				t.Errorf("UpdateHosts() is empty")
			}
		})
	}
}

func Test_getIP(t *testing.T) {
	type args struct {
		host string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wanterr bool
	}{
		{
			name:    "test_getIP_01",
			args:    args{host: "github.com"},
			want:    "20.205.243.166", // 执行测试时，github.com的ip可能会变化
			wanterr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := getIP(tt.args.host); got != tt.want && tt.wanterr != (err != nil) {
				t.Errorf("getIP() = %v, want %v", got, tt.want)
			}
		})
	}
}
