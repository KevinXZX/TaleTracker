package model

type Category struct {
	ID       int       `json:"id"`
	Name     string    `json:"name"`
	ParentID int       `json:"parent_id"`
	Parent   *Category `json:"parent,omitempty"`
}
