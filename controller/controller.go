package controller

import (
	"github.com/cookiengineer/systempanel/model"
	"github.com/cookiengineer/systempanel/view"
)

type Controller interface {
	View() view.View
	Model() model.Model
	Handle(action string, args map[string]any) error
}
