package assets

//go:generate rice embed-go

import (
	"github.com/GeertJohan/go.rice"
)

func UiBox() *rice.Box {
	templateBox, err := rice.FindBox("public")
	return templateBox
}
