package libretranslate

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
)

// DefaultBaseURL contains the default base url for the LibreTranslate API.
const DefaultBaseURL = "https://translate.argosopentech.com"

// Client handles the interaction with the LibreTranslate API.
type Client struct {
	baseUrl string
	client  *http.Client
}

// NewClient returns a new API client with the given token.
func NewClient() *Client {
	return &Client{
		baseUrl: DefaultBaseURL,
		client:  http.DefaultClient,
	}
}

// Detection represents the result of a dectection query.
type Detection struct {
	// Confidence value
	Confidence float64 `json:"confidence"`
	// Language code
	Language string `json:"language"`
}

// Language represents the result for the languages query.
type Language struct {
	// Language code
	Code string `json:"code"`
	// Human-readable language name (in English)
	Name string `json:"name"`
}

// Language represents the result for a translation query.
type TranslateResult struct {
	TranslatedText string `json:"translatedText"`
}

// Detect makes a request to detects the language of a given text.
func (c *Client) Detect(q string) ([]Detection, error) {
	params := url.Values{}
	params.Set("q", q)

	req, err := c.buildRequest(http.MethodPost, "/detect", params)
	if err != nil {
		return nil, err
	}

	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	responseBody, err := checkForResponseErrors(res)
	if err != nil {
		return nil, err
	}

	defer responseBody.Close()

	result := []Detection{}
	if err := json.NewDecoder(responseBody).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}

// Getlanguages makes a request to retrieve the list of supported languages.
func (c *Client) GetLanguages() ([]Language, error) {
	req, err := c.buildRequest(http.MethodGet, "/languages", url.Values{})

	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	responseBody, err := checkForResponseErrors(res)
	if err != nil {
		return nil, err
	}

	defer responseBody.Close()

	result := []Language{}
	if err := json.NewDecoder(responseBody).Decode(&result); err != nil {
		return result, err
	}

	return result, nil
}

// Translate makes a request to translate a given text from one language to another.
func (c *Client) Translate(query, source, target string) (string, error) {
	params := url.Values{}
	params.Set("q", query)
	params.Set("source", source)
	params.Set("target", target)

	req, err := c.buildRequest(http.MethodPost, "/translate", params)
	if err != nil {
		return "", err
	}

	res, err := c.client.Do(req)
	if err != nil {
		return "", err
	}

	responseBody, err := checkForResponseErrors(res)
	if err != nil {
		return "", err
	}

	defer responseBody.Close()

	result := TranslateResult{}
	if err := json.NewDecoder(responseBody).Decode(&result); err != nil {
		return "", err
	}

	return result.TranslatedText, nil
}

// buildRequest constructs an HTTP request with the specified HTTP method, endpoint, and parameters.
func (c *Client) buildRequest(method, endpoint string, params url.Values) (*http.Request, error) {
	uri, err := url.Parse(c.baseUrl)
	if err != nil {
		return nil, fmt.Errorf("URL parsing error: %s", err)
	}

	uri.Path = path.Join(uri.Path, endpoint)

	req, err := http.NewRequest(method, uri.String(), bytes.NewBufferString(params.Encode()))
	if err != nil {
		return nil, fmt.Errorf("HTTP request creation error: %s", err)
	}

	if method == http.MethodPost {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}

	return req, nil
}

type apiError struct {
	Error string `json:"error"`
}

// checkForResponseErrors checks an HTTP response for errors and returns the response body.
func checkForResponseErrors(res *http.Response) (io.ReadCloser, error) {
	body := res.Body
	status := res.StatusCode

	if status != 200 {
		var result apiError
		if err := json.NewDecoder(body).Decode(&result); err != nil {
			return nil, fmt.Errorf("API error: non-ok response (%d) from the API and failed to decode error message", status)
		}

		return nil, fmt.Errorf("API error: code %d - %s", status, result.Error)
	}

	return body, nil
}