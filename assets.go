package medialocker

import "net/http"

// AssetStore returns a handler to serve the website.
type AssetStore interface {
	Handler() http.Handler
}
