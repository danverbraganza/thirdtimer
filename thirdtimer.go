package main

import (
	"github.com/danverbraganza/go-mithril/moria"
	"honnef.co/go/js/dom"

	"thirdtimer/components"
)

func main() {
	myTimer := &components.Timer{Ratio: 3}

	moria.Route(
		dom.GetWindow().Document().QuerySelector("body"), "/",
		map[string]moria.Component{
			"/": myTimer,
		})
}
