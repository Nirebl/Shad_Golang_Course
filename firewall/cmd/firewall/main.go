//go:build !solution

package main

import (
	"flag"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"bytes"
	"regexp"
	"strings"
)

var (
	confPath     = flag.String("conf", "./configs/example.yaml", "path to config file")
	firewallAddr = flag.String("addr", ":8081", "address of firewall")
	serviceAddr  = flag.String("service-addr", "http://eu.httpbin.org/", "address of protected service")
)

var _ http.RoundTripper = (*firewall)(nil)

type firewall struct {
	config *Rules
}

func NewFirewall(config *Rules) http.RoundTripper {
	f := &firewall{
		config: config,
	}
	return f
}

func (f *firewall) RoundTrip(request *http.Request) (*http.Response, error) {
	rule := f.getRule(request.RequestURI)

	if !rule.reqAllowed(request) {
		return &http.Response{
			Body:       io.NopCloser(bytes.NewBufferString("Forbidden")),
			StatusCode: 403,
		}, nil
	}

	res, err := http.DefaultTransport.RoundTrip(request)
	if err != nil {
		return nil, err
	}

	if !rule.resAllowed(res) {
		return &http.Response{
			Body:       io.NopCloser(bytes.NewBufferString("Forbidden")),
			StatusCode: 403,
		}, nil
	}

	return res, nil
}

func (f *firewall) getRule(reqURI string) *Rule {
	for _, rule := range f.config.Rules {
		if rule.Endpoint == reqURI {
			return &rule
		}
	}
	return nil
}

func (r *Rule) reqAllowed(req *http.Request) bool {
	if r == nil {
		return true
	}

	UserAgent := req.UserAgent()
	for _, v := range r.ForbiddenUserAgents {
		matchString, err := regexp.MatchString(v, UserAgent)
		if matchString || err != nil {
			return false
		}
	}

	for _, v := range r.RequiredHeaders {
		if req.Header.Get(v) == "" {
			return false
		}
	}

	for _, v := range r.ForbiddenHeaders {
		fhPair := strings.SplitN(v, ": ", 2)
		fieldName, fieldRegex := fhPair[0], fhPair[1]
		matchString, err := regexp.MatchString(fieldRegex, req.Header.Get(fieldName))
		if matchString || err != nil {
			return false
		}
	}

	if req.Body != nil {
		body, err := io.ReadAll(req.Body)
		if err != nil {
			return false
		}

		req.Body = io.NopCloser(bytes.NewReader(body))

		if BEL(body, r.MaxRequestLengthBytes) || BForbidden(string(body), r.ForbiddenRequestRe) {
			return false
		}
	}

	return true
}

func (r *Rule) resAllowed(response *http.Response) bool {
	if r == nil {
		return true
	}

	for _, v := range r.ForbiddenResponseCodes {
		if response.StatusCode == v {
			return false
		}
	}

	if response.Body != nil {
		body, err := io.ReadAll(response.Body)
		if err != nil {
			return false
		}

		response.Body = io.NopCloser(bytes.NewBuffer(body))

		if BEL(body, r.MaxResponseLengthBytes) || BForbidden(string(body), r.ForbiddenResponseRe) {
			return false
		}
	}

	return true
}

func BEL(body []byte, limit int) bool {
	if limit <= 0 {
		return false
	}

	return len(body) > limit
}

func BForbidden(body string, ForbiddenRes []string) bool {
	if len(body) <= 0 {
		return false
	}
	for _, v := range ForbiddenRes {
		matchString, err := regexp.MatchString(v, body)
		if matchString || err != nil {
			return true
		}
	}

	return false
}

func main() {
	flag.Parse()
	confData, err := os.ReadFile(*confPath)
	if err != nil {
		log.Fatal(err)
	}

	rules, err := ParseR(confData)
	if err != nil {
		log.Fatal(err)
	}

	parsedURL, err := url.Parse(*serviceAddr)
	if err != nil {
		log.Fatal(err)
	}

	reverseProxy := httputil.NewSingleHostReverseProxy(parsedURL)
	reverseProxy.Transport = NewFirewall(rules)

	http.HandleFunc("/", reverseProxy.ServeHTTP)
	err = http.ListenAndServe(*firewallAddr, nil)
	if err != nil {
		log.Fatal(err)
	}
}
