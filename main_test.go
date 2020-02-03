package main

import (
	"testing"

	resty "github.com/go-resty/resty/v2"
)

func Test_printResponseOrError(t *testing.T) {
	type args struct {
		resp *resty.Response
		err  error
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			printResponseOrError(tt.args.resp, tt.args.err)
		})
	}
}
