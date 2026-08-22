package request

import (
	"testing"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func TestCreateTargetRequest_IPv4Binding(t *testing.T) {
	v, ok := binding.Validator.Engine().(*validator.Validate)
	if !ok {
		t.Fatal("binding validator is not *validator.Validate")
	}

	tests := []struct {
		ip    string
		valid bool
	}{
		{"10.10.10.10", true},
		{"1.2.3.4", true},
		{"1111111111", false},
		{"localhost", false},
		{"256.1.1.1", false},
	}

	for _, tc := range tests {
		err := v.Struct(CreateTargetRequest{
			IPAddress: tc.ip,
			Port:      9100,
			Labels:    []byte(`{}`),
		})
		if tc.valid && err != nil {
			t.Errorf("ip %q should be valid: %v", tc.ip, err)
		}
		if !tc.valid && err == nil {
			t.Errorf("ip %q should be invalid", tc.ip)
		}
	}
}
