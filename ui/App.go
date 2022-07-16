package ui

import (
	"app/dom"
)

var Message *dom.State 
var App = dom.NewComponent(
		//Action
		func(){
			Message = dom.NewState("Hello go!!!")
		},
		//Model
		func()string{
			dom.AddChilds(&Botonera)
			
			return `
				<div class='app'>
					<h1>$Message</h1>
					</Botonera>
				</div>
			`
		},
)