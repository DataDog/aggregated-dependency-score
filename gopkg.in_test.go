package aggregdepscore

import "testing"

func TestGetGopkginRepository(t *testing.T) {
	for _, each := range []struct {
		packageName string
		expected    string
		expectError bool
	}{
		{
			packageName: "gopkg.in/yaml.v2",
			expected:    "github.com/go-yaml/yaml",
			expectError: false,
		},
		{
			packageName: "gopkg.in/DataDog/dd-trace-go.v1",
			expected:    "github.com/DataDog/dd-trace-go",
			expectError: false,
		},
		{
			packageName: "wrong",
			expectError: true,
		},
		{
			packageName: "gopkg.in/noVersion",
			expectError: true,
		},
		{
			packageName: "gopkg.in/one/two/three",
			expectError: true,
		},
	} {
		t.Run(each.packageName, func(t *testing.T) {
			actual, err := getGopkginRepository(each.packageName)
			if each.expectError {
				if err == nil {
					t.Errorf("Expected error, but got nil")
				}
				return
			}
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			if actual != each.expected {
				t.Errorf("Expected %q, but got %q", each.expected, actual)
			}
		})
	}
}
