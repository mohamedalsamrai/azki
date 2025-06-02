package azki

import (
	"strings"
)

// HandleWildCard validates the given path for correct wildcard usage.
// It panics if the path is invalid, otherwise returns the same path.
func handleWildCard(path string) string {
	// Ensure the path starts with '/'
	if len(path) == 0 || path[0] != '/' {
		panic("path should start with /")
	}

	// Ensure there is no '//' anywhere in the path
	for i := 1; i < len(path); i++ {
		if path[i] == '/' && path[i-1] == '/' {
			panic("path should not contain //")
		}
	}

	// Iterate over each byte in the path to check wildcards
	for i := range path {
		switch path[i] {
		case '*':
			// '*' must be followed by at least one character
			if i+1 >= len(path) {
				panic("'*' should be followed by characters")
			}
			// '*' must not be directly followed by '/' or ':'
			if path[i+1] == '/' || path[i+1] == ':' {
				panic("'*' should not be followed by '/' or ':'")
			}
			// A catch-all wildcard must be the last segment
			//  i.e., no '/' may appear after '*' in the same segment
			if idx := strings.IndexByte(path[i+1:], '/'); idx != -1 {
				// Found a '/' after "*..." → catch-all is not at the end
				panic("catch-all must be the last segment")
			}
			// Inside the variable name after '*', do not allow another '*', ':' '/', or space
			for j := i + 1; j < len(path); j++ {
				switch path[j] {
				case '*':
					panic("path should not contain '**'")
				case ':', '/', ' ':
					panic("*filepath should not contain ':', '/' or spaces")
				}
			}

		case ':':
			// ':' must be followed by at least one character (letter, digit, or underscore)
			if i+1 >= len(path) {
				panic("':' should be followed by characters")
			}
			// ':' must not be directly followed by '/' or '*'
			if path[i+1] == '/' || path[i+1] == '*' {
				panic("':' should not be followed by '/' or '*'")
			}
			// Inside the variable name after ':', do not allow ':', '*' or space
			for j := i + 1; j < len(path); j++ {
				if path[j] == '/' {
					// End of the parameter name at the first '/'
					break
				}
				if path[j] == ':' {
					panic("path should not contain '::'")
				}
				if path[j] == '*' {
					panic(":param should not contain '*'")
				}
				if path[j] == ' ' {
					panic("parameter name should not contain spaces")
				}
			}
		}
	}

	return path
}
