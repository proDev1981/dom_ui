package main

import (
	"app/dom"
	"app/ui"
)

func main() {

	window := dom.NewWindow().
		SetIcon("./logo.png").
		SetTitle("new-app").
		SetPosition(dom.Center()).
		SetSize(200,200)

	dom.New( ui.App , window )
	dom.OnWait()
}
