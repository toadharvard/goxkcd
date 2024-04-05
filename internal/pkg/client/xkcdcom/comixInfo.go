package xkcdcom

type ComixInfo struct {
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
	Language   string
	Id         int
}

func NewComixInfo(id int, language string) ComixInfo {
	return ComixInfo{
		Id:       id,
		Language: language,
	}
}
