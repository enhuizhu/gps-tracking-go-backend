package controllers

import "fmt"

// HomeController to deal with root request
type HomeController struct {
}

// Index default method
func (h *HomeController) Index() {
	fmt.Println("hello home controller")
}
