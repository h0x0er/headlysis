package headlysis

import "net/http"

type SecureHeader struct {
	Name string
}

func (s SecureHeader) GetUrl() string {
	return "https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/" + s.Name
}

//https://www.researchgate.net/publication/326280169_HTTP_security_headers_analysis_of_top_one_million_websites

var secureHeaders = []SecureHeader{
	{"X-Frame-Options"},
	{"Content-Security-Policy"},
	{"X-Xss-Protection"},
	{"X-Content-Type-Options"},
	{"Strict-Transport-Security"},
	{"Cache-Control"},
	{"Clear-Site-Data"},
	{"Referrer-Policy"},
	{"Expect-CT"},
}

func isPresent(what string, in []string) bool {

	var out = false

	for _, k := range in {
		if k == what {
			out = true
			break
		}
	}
	return out

}

func GetMissingHeaders(reqHeaders http.Header) ([]SecureHeader, []SecureHeader) {

	var responseHeadersNames []string
	for key, _ := range reqHeaders {
		responseHeadersNames = append(responseHeadersNames, key)
	}

	var notPresent []SecureHeader
	var present []SecureHeader
	for _, sh := range secureHeaders {
		if !isPresent(sh.Name, responseHeadersNames) {
			notPresent = append(notPresent, sh)
		} else {
			present = append(present, sh)
		}
	}

	return notPresent, present

}
