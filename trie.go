package zepto

import (
	"strings"
)

// Param holds one URL parameter key and value.
type Param struct {
	Key   string
	Value string
}

// Trie is an HTTP router with zero-alloc lookup for static and param routes.
type Trie struct {
	static map[string]map[string]HandlerFunc // method -> (path -> handler)
	root   map[string]*Node                  // method -> param-trie root
}

// NewTrie creates a new, empty router.
func NewTrie() *Trie {
	return &Trie{
		static: make(map[string]map[string]HandlerFunc),
		root:   make(map[string]*Node),
	}
}

// Handle registers a handler for method+path. Static routes (no ':')
// go into a simple map; param routes go into the trie.
func (t *Trie) Handle(method, path string, handler HandlerFunc) {
	path = cleanPath(path)
	// static route: exact match, zero-alloc at lookup
	if !strings.Contains(path, ":") {
		if t.static[method] == nil {
			t.static[method] = make(map[string]HandlerFunc)
		}
		t.static[method][path] = handler
		return
	}

	// param route: insert into trie
	parts := splitPath(path)
	n, ok := t.root[method]
	if !ok {
		n = &Node{Children: make(map[string]*Node)}
		t.root[method] = n
	}
	for _, part := range parts {
		if strings.HasPrefix(part, ":") {
			// named parameter
			if n.ParamChild == nil {
				n.ParamChild = &Node{
					Segment:   part,
					ParamName: part[1:],
					Children:  make(map[string]*Node),
				}
			}
			n = n.ParamChild
		} else {
			child, ok := n.Children[part]
			if !ok {
				child = &Node{Segment: part, Children: make(map[string]*Node)}
				n.Children[part] = child
			}
			n = child
		}
	}
	// assign handler at leaf
	n.Handler = handler
}

// Find looks up a handler by method+path. Returns handler, params slice, and found.
// Static routes are checked first (zero alloc), then the trie is walked (zero alloc up to cap).
func (t *Trie) Find(method, path string) (HandlerFunc, []Param, bool) {
	// static lookup
	if m, ok := t.static[method]; ok {
		if h, ok2 := m[path]; ok2 {
			return h, nil, true
		}
	}

	// trie lookup for param routes
	root, ok := t.root[method]
	if !ok {
		return nil, nil, false
	}

	// stack-buffer for up to 4 params
	var paramBuf [4]Param
	params := paramBuf[:0]
	n := root
	start := 1 // skip leading '/'
	for start <= len(path) {
		slashIdx := strings.IndexByte(path[start:], '/')
		var seg string
		if slashIdx < 0 {
			seg = path[start:]
			start = len(path) + 1
		} else {
			seg = path[start : start+slashIdx]
			start += slashIdx + 1
		}

		if child, ok := n.Children[seg]; ok {
			n = child
			continue
		}
		if n.ParamChild != nil {
			// record param in stack buffer
			params = append(params, Param{Key: n.ParamChild.ParamName, Value: seg})
			n = n.ParamChild
			continue
		}
		// no match
		return nil, nil, false
	}

	if n.Handler != nil {
		return n.Handler, params, true
	}
	return nil, nil, false
}

// splitPath trims slashes and splits on '/'
func splitPath(p string) []string {
	p = strings.Trim(p, "/")
	if p == "" {
		return nil
	}
	return strings.Split(p, "/")
}
