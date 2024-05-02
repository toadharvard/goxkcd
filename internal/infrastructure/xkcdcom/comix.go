package xkcdcom

import (
	"github.com/toadharvard/goxkcd/internal/entity"
	"github.com/toadharvard/goxkcd/pkg/iso6391"
	"github.com/toadharvard/goxkcd/pkg/stemming"
)

type XKCDComix struct {
	Month      string `json:"month"`
	Num        int    `json:"num"`
	Link       string `json:"link"`
	Year       string `json:"year"`
	News       string `json:"news"`
	SafeTitle  string `json:"safe_title"`
	Transcript string `json:"transcript"`
	Alt        string `json:"alt"`
	Img        string `json:"img"`
	Title      string `json:"title"`
	Day        string `json:"day"`
	Language   iso6391.ISOCode6391
}

func NewXKCDComix(language iso6391.ISOCode6391) *XKCDComix {
	return &XKCDComix{
		Language: language,
	}
}

func (c *XKCDComix) ToComixEntity() *entity.Comix {
	stemmer := stemming.New()
	keywords := []stemming.Token{}
	keywords = append(keywords, stemmer.StemString(c.SafeTitle, c.Language)...)
	keywords = append(keywords, stemmer.StemString(c.Transcript, c.Language)...)
	keywords = append(keywords, stemmer.StemString(c.Alt, c.Language)...)
	keywords = stemming.RemoveDuplicates(keywords)
	return entity.NewComix(
		c.Num,
		c.Img,
		keywords,
	)
}
