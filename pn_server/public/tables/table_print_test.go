package tables

import (
	"fmt"
	"testing"
)

type User struct {
	Name string
	Age  int
}

func TestTable(t *testing.T) {
	type args struct {
		i interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "1结构体",
			args: args{
				i: User{
					Name: "test",
					Age:  18,
				},
			},
		},
		{
			name: "2结构体数组",
			args: args{
				i: []User{
					{
						"test",
						18,
					},
					{
						"test2",
						18,
					},
				},
			},
		},
		{
			name: "3Map",
			args: args{
				i: map[string]interface{}{
					"Name": "test",
					"Age":  18,
				},
			},
		},
		{
			name: "4Map数组",
			args: args{
				i: []map[string]interface{}{
					{
						"Name": "test",
						"Age":  18,
					},
					{
						"Name": "test",
						"Age":  18,
					},
				},
			},
		},
		{
			name: "5string",
			args: args{
				i: "test",
			},
		},
		{
			name: "5string数组",
			args: args{
				i: []string{
					"test",
					"test2",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := TableString(tt.args.i)
			fmt.Println()
			fmt.Println(got)
		})
	}
}
