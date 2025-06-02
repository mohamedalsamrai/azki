package azki

import (
	"testing"
)

func TestHandleWildCard_ValidPaths(t *testing.T) {
	validPaths := []string{
		"/users/:id",
		"/assets/*filepath",
		"/product/:pid/review",
		"/api/v1/resource",
		"/files/:name/info",
		"/download/*file",
	}

	for _, path := range validPaths {
		t.Run("Valid "+path, func(t *testing.T) {
			result := handleWildCard(path)
			if result != path {
				t.Errorf("expected '%s', got '%s'", path, result)
			}
		})
	}
}

func TestHandleWildCard_InvalidPaths(t *testing.T) {
	invalidPaths := []string{
		"user/:id",                 // does not start with "/"
		"/user//name",              // contains "//"
		"/file/*/edit",             // "*" followed by "/"
		"/file/:*id",               // ":" followed by "*"
		"/assets/*filepath/images", // catch-all not at the end
		"/products/:/details",      // ":" without a variable name
		"/search/:q uery",          // space in parameter name
		"/path/**double",           // consecutive "*"
		"/path/:id:extra",          // consecutive ":"
		"//", 					// starts with "//"
		
	}

	for _, path := range invalidPaths {
		t.Run("Invalid "+path, func(t *testing.T) {
			defer func() {
				if r := recover(); r == nil {
					t.Errorf("expected panic for path '%s', but no panic occurred", path)
				}
			}()
			_ = handleWildCard(path)
		})
	}
}
