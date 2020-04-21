package helpers

import "errors"

// ParamsHelper help to determin which controller to use
type ParamsHelper struct {
	params []string
}

func (p ParamsHelper) getController() (string, error) {
	if len(p.params[0]) > 0 {
		return p.params[0], nil
	}

	return "", errors.New("controller does not exist")
}

func (p ParamsHelper) getMethod() (string, error) {
	if len(p.params) <= 0 || len(p.params[0]) <= 0 {
		return "", errors.New("controller can not be empty")
	}

	// default method should be index
	if len(p.params) == 1 && len(p.params[0]) > 0 {
		return "index", nil
	}

	if len(p.params) >= 2 && len(p.params[1]) > 0 {
		return p.params[1], nil
	}

	return "", errors.New("unknown error")
}
