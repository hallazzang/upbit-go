package upbit

import (
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt"
)

var DefaultBaseURL = "https://api.upbit.com/v1/"

type Client struct {
	baseURL    *url.URL
	httpClient *http.Client
	accessKey  string
	secretKey  string
}

func NewClient(accessKey string, opts *ClientOptions) (*Client, error) {
	if accessKey == "" {
		return nil, fmt.Errorf("access key must be provided")
	}
	c := &Client{accessKey: accessKey}

	baseURL := DefaultBaseURL
	if opts.BaseURL != nil {
		baseURL = *opts.BaseURL
	}
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("parse base url: %w", err)
	}
	c.baseURL = u

	if opts.SecretKey != nil {
		c.secretKey = *opts.SecretKey
	}

	if opts.HTTPClient != nil {
		c.httpClient = opts.HTTPClient
	} else {
		c.httpClient = &http.Client{}
	}

	return c, nil
}

func (c *Client) AuthToken(query url.Values) (string, error) {
	nonce, err := uuid.NewV4()
	if err != nil {
		return "", fmt.Errorf("new uuid: %w", err)
	}
	claims := jwt.MapClaims{
		"access_key": c.accessKey,
		"nonce":      nonce.String(),
	}
	if len(query) != 0 {
		h := sha512.Sum512([]byte(query.Encode()))
		claims["query_hash"] = hex.EncodeToString(h[:])
		claims["query_hash_alg"] = "SHA512"
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	s, err := token.SignedString([]byte(c.secretKey))
	if err != nil {
		return "", fmt.Errorf("sign token: %w", err)
	}
	return "Bearer " + s, nil
}

func (c *Client) GET(path string, query url.Values) (*http.Response, error) {
	path = strings.TrimPrefix(path, "/") // Trim leading slash
	u, err := c.baseURL.Parse(path)
	if err != nil {
		return nil, fmt.Errorf("resolve path: %w", err)
	}
	u.RawQuery = query.Encode()
	req, _ := http.NewRequest("GET", u.String(), nil)
	token, err := c.AuthToken(query)
	if err != nil {
		return nil, fmt.Errorf("get auth token: %w", err)
	}
	req.Header.Set("Authorization", token)
	return c.httpClient.Do(req)
}

type ClientOptions struct {
	SecretKey  *string
	BaseURL    *string
	HTTPClient *http.Client
}

func NewClientOptions() *ClientOptions {
	return &ClientOptions{}
}

func (opts *ClientOptions) SetSecretKey(secretKey string) *ClientOptions {
	opts.SecretKey = &secretKey
	return opts
}

func (opts *ClientOptions) SetBaseURL(baseURL string) *ClientOptions {
	opts.BaseURL = &baseURL
	return opts
}

func (opts *ClientOptions) SetHTTPClient(httpClient *http.Client) *ClientOptions {
	opts.HTTPClient = httpClient
	return opts
}
