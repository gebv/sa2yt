package grifts

import (
	"github.com/gobuffalo/buffalo"
	"github.com/sayto/actions"
)

func init() {
	buffalo.Grifts(actions.App())
}
