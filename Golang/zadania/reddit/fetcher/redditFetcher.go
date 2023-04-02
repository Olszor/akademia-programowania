package fetcher

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type response struct {
	Data struct {
		Children []struct {
			Data struct {
				Title string `json:"title"`
				URL   string `json:"url"`
			} `json:"data"`
		} `json:"children"`
	} `json:"data"`
}

type RedditFetcher interface {
	Fetch() error
	Save(io.Writer) error
	FetchWithContext(ctx context.Context) error
}

type Fetcher struct {
	Client   http.Client
	Url      string
	response response
}

func (f *Fetcher) Fetch() error {
	req, err := http.NewRequest(http.MethodGet, f.Url, http.NoBody)
	if err != nil {
		return fmt.Errorf("error creating request: %s", err)
	}
	req.Header.Add("Content-Type", "application/json")

	resp, err := f.Client.Do(req)
	if err != nil {
		return fmt.Errorf("error fetching data: %s", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var data response
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return fmt.Errorf("error reading data: %s", err)
	}
	f.response = data

	return nil
}

func (f *Fetcher) Save(writer io.Writer) error {
	for _, child := range f.response.Data.Children {
		_, err := writer.Write([]byte(fmt.Sprintf("%s\n%s\n", child.Data.Title, child.Data.URL)))
		if err != nil {
			return fmt.Errorf("error writing to a file %s", err)
		}
	}
	log.Printf("Saved succesfully!")
	return nil
}

func (f *Fetcher) FetchWithContext(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, f.Url, http.NoBody)
	if err != nil {
		return fmt.Errorf("error creating request: %s", err)
	}
	req.Header.Add("Content-Type", "application/json")

	resp, err := f.Client.Do(req)
	if err != nil {
		return fmt.Errorf("error fetching data: %s", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var data response
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return fmt.Errorf("error reading data: %s", err)
	}
	f.response = data

	return nil
}
