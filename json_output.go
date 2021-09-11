package headlysis

type Output struct {
	Headlysis []UrlAnalysisOutput `json:"headlysis"`
}

type UrlAnalysisOutput struct {
	Url               string             `json:"target_url"`
	PresentHeaders    []PresentHeader    `json:"present_security_headers"`
	NotPresentHeaders []NotPresentHeader `json:"not_present_security_header"`
}

type PresentHeader struct {
	Name  string `json:"header_name"`
	Value string `json:"header_value"`
}

type NotPresentHeader struct {
	Name    string `json:"header_name"`
	InfoUrl string `json:"info_url"`
}
