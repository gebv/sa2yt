package grifts

import (
	"github.com/gebv/sayto/actions"
	"github.com/gobuffalo/buffalo"
)

func init() {
	buffalo.Grifts(actions.App())
}
