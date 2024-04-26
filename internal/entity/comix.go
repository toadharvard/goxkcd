package entity

type Comix struct {
	ID       int      `json:"id"`
	URL      string   `json:"url"`
	Keywords []string `json:"keywords"`
}

func NewComix(id int, url string, keywords []string) *Comix {
	return &Comix{
		ID:       id,
		URL:      url,
		Keywords: keywords,
	}
}
