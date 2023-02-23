package fetch

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

var EmptyBody = bytes.NewBuffer([]byte{})

type fetch struct {
	req  *http.Request
	resp *http.Response
	err  error
}

func Get(uri string) *fetch {
	req, err := http.NewRequest(http.MethodGet, uri, nil)
	return &fetch{req: req, err: err}
}

func Post(uri string, body io.Reader) *fetch {
	req, err := http.NewRequest(http.MethodPost, uri, body)
	return &fetch{req: req, err: err}
}

func (f *fetch) Header(k string, v string) *fetch {
	if f.err == nil {
		f.req.Header.Add(k, v)
	}
	return f
}

func (f *fetch) JSON() *fetch {
	return f.ContentType("application/json")
}

func (f *fetch) ContentType(v string) *fetch {
	if f.err == nil {
		f.req.Header.Add("Content-Type", v)
	}
	return f
}

func (f *fetch) LastErr() error {
	return f.err
}

func (f *fetch) RespCode() int {
	if f.resp != nil {
		return f.resp.StatusCode
	}
	return 0
}

func (f *fetch) Apply(v ...interface{}) *fetch {
	if f.err != nil {
		return f
	}

	f.resp, f.err = http.DefaultClient.Do(f.req)
	if f.err != nil {
		return f
	}

	defer f.resp.Body.Close()
	if len(v) > 0 {
		f.err = json.NewDecoder(f.resp.Body).Decode(v[0])
		if f.err != nil {
			return f
		}
	}
	return f
}
