package reelser

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
)

var (
	ErrReelIDNotFoundInLink = errors.New("reel id not found in link")
)

var reelIDRe = regexp.MustCompile("reel/(.*?)/")

const baseURL = "https://www.instagram.com/graphql/query/"

type Client struct {
	httpClient *http.Client
}

func New() *Client {
	return &Client{
		httpClient: http.DefaultClient,
	}
}

type Response struct {
	Data struct {
		ShortMedia struct {
			VideoURL string `json:"video_url"`
		} `json:"shortcode_media"`
	} `json:"data"`
}

func (c *Client) GetVideoURL(reelURL string) (*url.URL, error) {
	reelID, err := reelIDFromURL(reelURL)
	if err != nil {
		return nil, fmt.Errorf("reel id from url: %w", err)
	}

	req, err := http.NewRequest(http.MethodGet, baseURL, nil)
	if err != nil {
		return nil, fmt.Errorf("http new request: %w", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.3; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.193 Safari/537.36")
	q := req.URL.Query()
	q.Add("hl", "en")
	q.Add("query_hash", "b3055c01b4b222b8a47dc12b090e4e64")
	q.Add("variables", fmt.Sprintf(`{"child_comment_count":1,"fetch_comment_count":1,"has_threaded_comments":true,"parent_comment_count":1,"shortcode":"%s"}`, reelID))

	req.URL.RawQuery = q.Encode()
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("client do: %w", err)
	}

	defer resp.Body.Close()

	var data Response

	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("json decoder decode: %w", err)
	}

	videoURL, err := url.Parse(data.Data.ShortMedia.VideoURL)
	if err != nil {
		return nil, fmt.Errorf("url parse: %w", err)
	}

	return videoURL, nil
}

func (c *Client) DownloadReel(url *url.URL, w io.Writer) error {
	req, err := http.NewRequest(http.MethodGet, url.String(), nil)
	if err != nil {
		return fmt.Errorf("http new request")
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.3; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.193 Safari/537.36")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("client do: %w", err)
	}

	defer resp.Body.Close()

	if _, err = io.Copy(w, resp.Body); err != nil {
		return fmt.Errorf("io copy: %w", err)
	}

	return nil
}

func reelIDFromURL(reelURL string) (string, error) {
	match := reelIDRe.FindStringSubmatch(reelURL)
	if len(match) < 2 {
		return "", ErrReelIDNotFoundInLink
	}

	return match[1], nil
}
