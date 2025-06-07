package azki

import (
	"testing"
	"strings"
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
// dummy Handle for testing; adjust if Handle has a specific signature
var dummyHandle Handle = nil

// -- existing handleWildCard tests omitted for brevity --

// Test inserting purely static paths
func TestInsertChild_Static(t *testing.T) {
	root := &Node{}
	h := dummyHandle

	root.insertChild("/home", h)

	if len(root.children) != 1 {
		t.Fatalf("expected 1 child under root, got %d", len(root.children))
	}
	home := root.children[0]

	// allow either "home" or "/home"
	actual := strings.TrimPrefix(home.path, "/")
	if actual != "home" {
		t.Errorf("expected path 'home', got '%s'", home.path)
	}
	if home.wildCard {
		t.Errorf("expected wildCard=false for static segment")
	}
	if home.wildcardType != static {
		t.Errorf("expected wildcardType=static, got %v", home.wildcardType)
	}
}

// Test inserting a path with a :param segment
func TestInsertChild_Param(t *testing.T) {
	root := &Node{}
	h := dummyHandle

	root.insertChild("/user/:id", h)

	// root -> "user"
	if len(root.children) != 1 {
		t.Fatalf("expected 1 child under root, got %d", len(root.children))
	}
	user := root.children[0]
	name := strings.Trim(user.path, "/")
	if name != "user" {
		t.Fatalf("expected static child 'user', got '%s'", user.path)
	}
	if user.wildCard {
		t.Errorf("expected user.wildCard=false")
	}

	// user -> ":id"
	if len(user.children) != 1 {
		t.Fatalf("expected 1 child under 'user', got %d", len(user.children))
	}
	idNode := user.children[0]
	paramName := strings.Trim(idNode.path, "/")
	if paramName != ":id" {
		t.Errorf("expected param child ':id', got '%s'", idNode.path)
	}
	if !idNode.wildCard || idNode.wildcardType != param {
		t.Errorf("expected param wildcardType, got wildCard=%v, type=%v", idNode.wildCard, idNode.wildcardType)
	}
}


// Test inserting a path with a *catchAll segment
func TestInsertChild_CatchAll(t *testing.T) {
	root := &Node{}
	h := dummyHandle

	root.insertChild("/static/*filepath", h)

	// root -> "static"
	if len(root.children) != 1 {
		t.Fatalf("expected 1 child under root, got %d", len(root.children))
	}
	s := root.children[0]
	staticName := strings.Trim(s.path, "/")
	if staticName != "static" {
		t.Fatalf("expected static child 'static', got '%s'", s.path)
	}

	// static -> "*filepath"
	if len(s.children) != 1 {
		t.Fatalf("expected 1 child under 'static', got %d", len(s.children))
	}
	catch := s.children[0]
	catchName := strings.Trim(catch.path, "/")
	if catchName != "*filepath" {
		t.Errorf("expected catchAll child '*filepath', got '%s'", catch.path)
	}
	if !catch.wildCard || catch.wildcardType != catchAll {
		t.Errorf("expected catchAll wildcardType, got wildCard=%v, type=%v", catch.wildCard, catch.wildcardType)
	}
}


