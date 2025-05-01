package zepto

type Node struct {
	Segment    string           // literal segment or “:name”
	Children   map[string]*Node // literal->child
	ParamChild *Node            // child for exactly one “:param” at this level
	ParamName  string           // name without the “:”
	Handler    HandlerFunc
}
