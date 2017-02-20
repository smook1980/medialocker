package assets

//go:generate ./gen_assets.sh

import (
	rice "github.com/GeertJohan/go.rice"
	"net/http"
)

func StaticAsseetServer() http.Handler {
	assetHandler := http.FileServer(rice.MustFindBox("web").HTTPBox())

	return assetHandler
}
