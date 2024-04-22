package core_test

import (
	"testing"

	"github.com/copito/goscaffold/core"
)

func TestDeltaRelativePath(t *testing.T) {
	testCases := []struct {
		basePath     string
		currentPath  string
		expectedPath string
	}{
		{
			basePath:     "/home/user1/Documents/",
			currentPath:  "/home/user1/Documents/example1/example3",
			expectedPath: "example1/example3",
		},
		{
			basePath:     "/home/user1/Documents/",
			currentPath:  "/home/user1/Documents/",
			expectedPath: "",
		},
		{
			basePath:     "/home/user1/Documents/",
			currentPath:  "/home/user1/Documents",
			expectedPath: "",
		},
		{
			basePath:     "/home/user1/Documents/",
			currentPath:  "/home/user1/Documents/example1",
			expectedPath: "example1",
		},
		{
			basePath:     "/home/user1/Documents/",
			currentPath:  "/home/user1/Documents/example1/example3/",
			expectedPath: "example1/example3",
		},
		{
			basePath:     "/home/user1/Documents/",
			currentPath:  "/home/user1/Documents/example1/example3/lemon.go",
			expectedPath: "example1/example3/lemon.go",
		},
		{
			basePath:     "/home/user1/",
			currentPath:  "/home/user1/Documents/example1/{{scaffold.project_name}}/lemon.go",
			expectedPath: "Documents/example1/{{scaffold.project_name}}/lemon.go",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.currentPath, func(t *testing.T) {
			actual := core.DeltaRelativePath(tc.basePath, tc.currentPath)
			if actual != tc.expectedPath {
				t.Errorf("DeltaRelativePath(%q, %q) = %q, expected %q", tc.basePath, tc.currentPath, actual, tc.expectedPath)
			}
		})
	}
}
