package main

import (
	_ "net/http/pprof"
	"testing"
)

func Test_main(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "Test_main",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			main()
		})
	}
}
