package medialocker

//go:generate go-bindata --debug -pkg medialocker -tags "!release" --prefix ./ui/dist -o ui-assets-dev.go ./ui/dist/...

//go:generate go-bindata -pkg medialocker -tags release --prefix ./ui/dist -o ui-assets-release.go ./ui/dist/...

import "net/http"

// AssetStore returns a handler to serve the website.
type AssetStore interface {
	Handler() http.Handler
}
