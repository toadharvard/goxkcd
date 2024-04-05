package comix

import (
	"github.com/toadharvard/goxkcd/internal/pkg/client/xkcdcom"
	"github.com/toadharvard/goxkcd/internal/pkg/stemming"
)

type Comix struct {
	Id       int      `json:"id"`
	URL      string   `json:"url"`
	Keywords []string `json:"keywords"`
}

func New(id int, url string, keywords []string) Comix {
	return Comix{
		Id:       id,
		URL:      url,
		Keywords: keywords,
	}
}

func FromComixInfo(comixInfo xkcdcom.ComixInfo) Comix {
	keywords := []stemming.Token{}
	keywords = append(keywords, stemming.StemString(comixInfo.SafeTitle, comixInfo.Language)...)
	keywords = append(keywords, stemming.StemString(comixInfo.Transcript, comixInfo.Language)...)

	keywords = stemming.RemoveDuplicates(keywords)
	return Comix{
		Id:       comixInfo.Num,
		URL:      comixInfo.Img,
		Keywords: keywords,
	}
}
