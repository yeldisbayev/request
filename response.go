package request

import (
	"encoding/json"
	"encoding/xml"
	"net/http"
	"strings"
)

type Response struct {
	*http.Response
}

// IsSuccess checks response status code for success.
func (res *Response) IsSuccess() bool {
	return res.StatusCode >= 200 && res.StatusCode < 300
}

// Decoder returns JSON or XML decoder depending on content type.
func (res *Response) Decoder() Decoder {
	contentType := res.Header.Get(ContentType)
	if strings.Contains(contentType, ApplicationXML) {
		return xml.NewDecoder(res.Body)
	} else {
		return json.NewDecoder(res.Body)
	}

}

// JSONDecoder returns JSON decoder.
func (res *Response) JSONDecoder() Decoder {
	return json.NewDecoder(res.Body)
}

// XMLDecoder returns XML decoder.
func (res *Response) XMLDecoder() Decoder {
	return xml.NewDecoder(res.Body)
}
