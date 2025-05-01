package zepto

// CleanPath cleans the path by removing trailing slashes.
func cleanPath(path string) string {
	if len(path) == 0 {
		return "/"
	}
	if path[len(path)-1] == '/' {
		path = path[:len(path)-1]
	}
	return path
}
