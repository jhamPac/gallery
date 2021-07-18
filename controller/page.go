package controller

import "github.com/jhampac/gallery/view"

type Static struct {
	Home    *view.View
	Contact *view.View
}

func NewStaticPage() *Static {
	return &Static{
		Home:    view.New("index", "page/home"),
		Contact: view.New("index", "page/contact"),
	}
}
