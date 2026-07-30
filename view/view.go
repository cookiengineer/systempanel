package view

import "github.com/cookiengineer/systempanel/bindings/gtk4"

type View interface {
	Widget() *gtk4.Widget
	Name() string
	Title() string
	IconName() string
	OnShow()
	OnHide()
	Destroy()
}

type ViewDescriptor struct {
	Name     string
	Title    string
	IconName string
	DetectFn func() bool
	Factory  func() View
}

var Registry []ViewDescriptor
