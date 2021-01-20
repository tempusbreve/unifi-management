package unifi

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"sort"
	"time"
)

// Errors.
var (
	ErrUnifi         = errors.New("unifi error")
	ErrUnInitialized = errors.New("uninitialized session")
)

// NewSession prepares a session for use.
func NewSession(opts ...Option) (*Session, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, fmt.Errorf("%v\n%w", err, ErrUnifi)
	}

	s := &Session{
		endpoint: "http://unifi",
		username: "ubnt",
		password: "ubnt",
		client: &http.Client{
			Jar:           jar,
			Timeout:       time.Minute * 1,
			Transport:     UserAgent("unifibot/2.0")(nil),
			CheckRedirect: nil,
		},
		err: nil,
	}

	s.login = s.webLogin

	for _, opt := range opts {
		if opt != nil {
			opt(s)
		}
	}

	return s, nil
}

// Option describes a function that modifies Session config.
type Option func(*Session)

// Endpoint sets the UniFi base URL.
func Endpoint(url string) func(*Session) { return func(s *Session) { s.endpoint = url } }

// Credentials sets the UniFi base URL.
func Credentials(username, password string) func(*Session) {
	return func(s *Session) {
		s.username = username
		s.password = password
		s.login = s.webLogin
	}
}

// Session wraps metadata to manage session state.
type Session struct {
	endpoint string
	username string
	password string
	csrf     string
	client   *http.Client
	login    func() (string, error)
	err      error
}

// Login performs authentication with the UniFi server, and stores the
// http credentials.
func (s *Session) Login() (string, error) {
	if s.login == nil {
		s.login = func() (string, error) {
			return "", ErrUnInitialized
		}
	}

	return s.login()
}

// ListDevices describes the known UniFi clients.
func (s *Session) ListDevices() ([]Device, error) {
	if s.err != nil {
		return nil, s.err
	}

	u, err := url.Parse(fmt.Sprintf("%s/proxy/network/api/s/default/rest/user", s.endpoint))
	if err != nil {
		return nil, s.setError(err)
	}

	data, err := s.get(u)
	if err != nil {
		return nil, s.setError(err)
	}

	resp := &Response{}
	if err := json.Unmarshal([]byte(data), resp); err != nil {
		return nil, s.setError(err)
	}

	devices := resp.Data
	sort.Slice(devices, func(i, j int) bool {
		return devices[i].LastSeen < devices[j].LastSeen
	})

	return devices, nil
}

// Kick disconnects a connected client, identified by MAC address.
func (s *Session) Kick(mac string) (string, error) {
	return s.macAction("kick-sta", mac)
}

// Block prevents a specific client (identified by MAC) from connecting
// to the UniFi network.
func (s *Session) Block(mac string) (string, error) {
	return s.macAction("block-sta", mac)
}

// Unblock re-enables a specific client.
func (s *Session) Unblock(mac string) (string, error) {
	return s.macAction("unblock-sta", mac)
}

func (s *Session) macAction(action, mac string) (string, error) {
	if b, err := s.login(); err != nil {
		return b, err
	}

	u, err := url.Parse(fmt.Sprintf("%s/proxy/network/api/s/default/cmd/stamgr", s.endpoint))
	if err != nil {
		return "", s.setError(err)
	}

	r := bytes.NewBufferString(fmt.Sprintf(`{"cmd":%q,"mac":%q}`, action, mac))

	return s.post(u, r)
}

func (s *Session) get(u fmt.Stringer) (string, error) {
	return s.verb("GET", u, nil)
}

func (s *Session) post(u fmt.Stringer, body io.Reader) (string, error) {
	return s.verb("POST", u, body)
}

func (s *Session) verb(verb string, u fmt.Stringer, body io.Reader) (string, error) {
	req, err := http.NewRequestWithContext(context.Background(), verb, u.String(), body)
	if err != nil {
		return "", s.setError(err)
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Origin", s.endpoint)

	if s.csrf != "" {
		req.Header.Set("x-csrf-token", s.csrf)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return "", s.setError(err)
	}

	defer resp.Body.Close()

	if tok := resp.Header.Get("x-csrf-token"); tok != "" {
		s.csrf = tok
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", s.setError(err)
	}

	if resp.StatusCode < http.StatusOK || http.StatusBadRequest <= resp.StatusCode {
		_ = s.setErrorString(http.StatusText(resp.StatusCode))
	}

	return string(respBody), s.err
}

func (s *Session) webLogin() (string, error) {
	if s.err != nil {
		return "", s.err
	}

	u, err := url.Parse(fmt.Sprintf("%s/api/auth/login", s.endpoint))
	if err != nil {
		return "", s.setError(err)
	}

	r := bytes.NewBufferString(
		fmt.Sprintf(
			`{"username":%q,"password":%q,"strict":"true","remember":"true"}`,
			s.username, s.password))

	respBody, err := s.post(u, r)
	if err == nil {
		s.login = func() (string, error) { return respBody, nil }
	}

	return respBody, err
}

func (s *Session) setError(e error) error {
	if e == nil {
		return nil
	}

	if s.err == nil {
		s.err = e
	} else {
		s.err = fmt.Errorf("%s\n%w", e, s.err)
	}

	return s.err
}

func (s *Session) setErrorString(e string) error {
	if len(e) == 0 {
		return nil
	}

	return s.setError(fmt.Errorf("%s\n%w", e, ErrUnifi))
}

// UserAgent returns an http.RoundTripper factory that sets the user-agent.
func UserAgent(agent string) func(http.RoundTripper) http.RoundTripper {
	return func(rt http.RoundTripper) http.RoundTripper {
		if len(agent) == 0 {
			return rt
		}

		return &customHeaderTransport{key: "User-Agent", value: agent, next: rt}
	}
}

type customHeaderTransport struct {
	key   string
	value string
	next  http.RoundTripper
}

func (t *customHeaderTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	r.Header.Set(t.key, t.value)

	rt := t.next
	if rt == nil {
		rt = http.DefaultTransport
	}

	return rt.RoundTrip(r)
}
