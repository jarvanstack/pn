package utils

import "testing"

func TestRandomFileName(t *testing.T) {
	type args struct {
		fileName string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			args: args{
				fileName: "test.txt",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RandomFileName(tt.args.fileName)
			t.Logf("[got] %v \n", got)
		})
	}
}
