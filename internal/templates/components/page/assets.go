package page

// AssetURL resolves a logical asset name to its URL (including content hash in production).
// Set this before rendering any templates, e.g. in cmd/generate/main.go or cmd/server/main.go.
// Defaults to serving from /assets/ with no hashing.
var AssetURL func(string) string = func(name string) string {
	return "/assets/" + name
}
