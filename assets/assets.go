package assets

//go:generate rice embed-go

import (
	"net/http"

	rice "github.com/GeertJohan/go.rice"
)

func StaticAsseetServer() http.Handler {
	assetHandler := http.FileServer(rice.MustFindBox("web").HTTPBox())
	// serves the index.html from rice
	return assetHandler
}
