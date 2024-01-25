package model

type Tale struct {
	ID         int        `json:"id"`
	Url        string     `json:"url"`
	Author     string     `json:"author"`
	Title      string     `json:"title"`
	Blurb      string     `json:"blurb"`
	Added      string     `json:"added"`
	Published  string     `json:"published"`
	Updated    string     `json:"updated"`
	Review     Review     `json:"review"`
	Tags       []Tag      `json:"tags"`
	Categories []Category `json:"categories"`
}
