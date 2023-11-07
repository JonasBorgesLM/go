package main

import "testing"

func TestIsCEPValid(t *testing.T) {
	type args struct {
		cep string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Teste 1",
			args: args{
				cep: "",
			},
			want: false,
		},
		{
			name: "Teste 2",
			args: args{
				cep: "99999-999",
			},
			want: true,
		},
		{
			name: "Teste 3",
			args: args{
				cep: "99999999",
			},
			want: true,
		},
		{
			name: "Teste 4",
			args: args{
				cep: "999-99999",
			},
			want: false,
		},
		{
			name: "Teste 5",
			args: args{
				cep: "99999",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsCEPValid(tt.args.cep); got != tt.want {
				t.Errorf("IsCEPValid() = %v, want %v", got, tt.want)
			}
		})
	}
}
