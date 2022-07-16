package ui

import "app/dom"

var Btn *dom.State
var Botonera = dom.NewComponent(
		//Action
		func(){
			Btn = dom.NewState("press")

			dom.Bind("handleBtn",func(){
				Message.Set("Bye world")
				Btn.Set("Apretado")
			})
		},
		//Model
		func()string{
			return `
				<button class='botonera' onclick=handleBtn()>$Btn</button>
			`
		},
)