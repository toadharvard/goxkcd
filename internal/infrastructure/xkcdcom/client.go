package xkcdcom

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/toadharvard/goxkcd/internal/entity"
	"github.com/toadharvard/goxkcd/pkg/iso6391"
)

type XKCDClient struct {
	URL      string
	Language iso6391.ISOCode6391
	client   *http.Client
}

func NewClient(URL string, language iso6391.ISOCode6391, timeout time.Duration) *XKCDClient {
	client := &XKCDClient{
		URL:      URL,
		Language: language,
		client:   &http.Client{Timeout: timeout},
	}
	return client
}

func (c *XKCDClient) GetComixByID(ctx context.Context, id int) (comix entity.Comixer, err error) {
	comix = NewXKCDComix(c.Language)
	URL := fmt.Sprintf("%s/%d/info.0.json", c.URL, id)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, URL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&comix)

	return comix, err
}

func (c *XKCDClient) GetLastComix(ctx context.Context) (comix entity.Comixer, err error) {
	comix = NewXKCDComix(c.Language)
	URL := fmt.Sprintf("%s/info.0.json", c.URL)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, URL, nil)
	if err != nil {
		return
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("status code is not OK, recieved: %d", resp.StatusCode)
		return
	}

	err = json.NewDecoder(resp.Body).Decode(comix)
	return
}
