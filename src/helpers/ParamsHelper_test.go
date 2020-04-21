package helpers

import (
	"fmt"
	"testing"
)

func TestGetController(t *testing.T) {
	strArr := []string{"a", "b"}
	pp := ParamsHelper{params: strArr}
	controller, err := pp.getController()

	if err != nil {
		t.Error(err)
	}

	if controller != "a" {
		t.Errorf("controller is %s, want a", controller)
	}
}

func TestGetMethodWhenParamsIsEmpty(t *testing.T) {
	strArr := []string{}
	pp := ParamsHelper{params: strArr}
	method, err := pp.getMethod()

	if err == nil || method != "" {
		t.Errorf("it should return error 'controller can not be empty'")
	} else {
		fmt.Println("test getMethod under the scenario that it has 0 parmas")
	}
}

func TestGetMethodWhenParamsHasOneElement(t *testing.T) {
	strArr := []string{"Hello"}
	pp := ParamsHelper{params: strArr}
	method, err := pp.getMethod()

	if err != nil || method != "index" {
		t.Error("method should be index")
	}
}

func TestGetMethodWhenParamsHasMoreThenTwoElements(t *testing.T) {
	strArr := []string{"Hello", "echo"}
	pp := ParamsHelper{params: strArr}
	method, err := pp.getMethod()

	if err != nil || method != "echo" {
		t.Error("method should be echo")
	}
}
