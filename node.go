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

type WildcardType int

const (
	static WildcardType = iota
	param
	catchAll
)

type Node struct {
	path         string
	children     []*Node
	wildCard     bool
	wildcardType WildcardType
	priority     uint32
	handler      Handle
}

func (n *Node) insertChild(path string, handle Handle) {
    _ = handleWildCard(path) // validate wildcards

    i := 0
    for i < len(path) {
        if path[i] == ':' {
            // Static prefix before :
            if i > 0 {
                staticPart := path[:i]
                child := &Node{
                    path:         staticPart,
                    wildcardType: static,
                }
                n.children = append(n.children, child)
                n = child
            }

            // Param segment
            paramEnd := i + 1
            for paramEnd < len(path) && path[paramEnd] != '/' {
                paramEnd++
            }
            paramPart := path[i:paramEnd]
            child := &Node{
                path:         paramPart,
                wildCard:     true,
                wildcardType: param,
            }
            n.children = append(n.children, child)
            n = child

            path = path[paramEnd:] // continue with remaining path
			i=0
            continue
        }

        if path[i] == '*' {
            if i > 0 {
                staticPart := path[:i]
                child := &Node{
                    path:         staticPart,
                    wildcardType: static,
                }
                n.children = append(n.children, child)
                n = child
            }

            catchAllPart := path[i:]
            child := &Node{
                path:         catchAllPart,
                wildCard:     true,
                wildcardType: catchAll,
            }
            n.children = append(n.children, child)
            n = child
            n.handler = handle
            n.priority++
            return
        }

        i++
    }

    // If no wildcard found, create a static path node
    child := &Node{
        path:         path,
        wildcardType: static,
        handler:      handle,
    }
    n.children = append(n.children, child)
    child.priority++
}
