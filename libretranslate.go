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
const DefaultBaseURL = "https://libretranslate.com"

// Client handles the interaction with the LibreTranslate API.
type Client struct {
	baseUrl string
	token   string
	client  *http.Client
}

// NewClient returns a new API client with the given token.
func NewClient(token string) *Client {
	return &Client{
		baseUrl: DefaultBaseURL,
		token:   token,
		client:  http.DefaultClient,
	}
}

// NewClientWithBaseURL returns a new API client with the given token.
func NewClientWithBaseURL(baseURL string, token string) *Client {
	return &Client{
		baseUrl: baseURL,
		token:   token,
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

// TranslateResult represents the result for a translation query.
type TranslateResult struct {
	// Detected language information (only for auto detect)
	DetectedLanguage Detection `json:"detectedLanguage"`
	// Translated text
	TranslatedText string `json:"translatedText"`
}

// Detect makes a request to detects the language of a given text.
func (c *Client) Detect(q string) ([]Detection, error) {
	params := url.Values{}
	params.Set("q", q)
	params.Set("api_key", c.token)

	req, err := c.buildRequest(http.MethodPost, "/detect", params)
	if err != nil {
		return nil, err
	}

	responseBody, err := doRequest(c.client, req)
	if err != nil {
		return nil, err
	}

	defer responseBody.Close()

	result := []Detection{}
	err = json.NewDecoder(responseBody).Decode(&result)

	return result, err
}

// Getlanguages makes a request to retrieve the list of supported languages.
func (c *Client) GetLanguages() ([]Language, error) {
	params := url.Values{}
	params.Set("api_key", c.token)

	req, err := c.buildRequest(http.MethodGet, "/languages", params)
	if err != nil {
		return nil, err
	}

	responseBody, err := doRequest(c.client, req)
	if err != nil {
		return nil, err
	}

	defer responseBody.Close()

	result := []Language{}
	err = json.NewDecoder(responseBody).Decode(&result)

	return result, err
}

// Translate makes a request to translate a given text from one language to another.
func (c *Client) Translate(query, source, target string) (string, error) {
	params := url.Values{}
	params.Set("q", query)
	params.Set("source", source)
	params.Set("target", target)
	params.Set("api_key", c.token)

	req, err := c.buildRequest(http.MethodPost, "/translate", params)
	if err != nil {
		return "", err
	}

	responseBody, err := doRequest(c.client, req)
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

// doRequest makes an HTTP request and returns the response body.
func doRequest(client *http.Client, req *http.Request) (io.ReadCloser, error) {
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return checkForResponseErrors(res)
}

type apiError struct {
	Error string `json:"error"`
}

// checkForResponseErrors checks an HTTP response for errors and returns the response body.
func checkForResponseErrors(res *http.Response) (io.ReadCloser, error) {
	if res.StatusCode != http.StatusOK {
		var result apiError
		if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
			return nil, fmt.Errorf(
				"API error: non-ok response (%d) from the API and failed to decode error message",
				res.StatusCode,
			)
		}

		return nil, fmt.Errorf("API error: code %d - %s", res.StatusCode, result.Error)
	}

	return res.Body, nil
}
