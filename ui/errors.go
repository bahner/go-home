package ui

import "errors"

var (
	ErrNoEntitySelected = errors.New("no entity selected")
	ErrYouDontExist     = errors.New("you don't exist")
)