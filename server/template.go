package server

var (
	Template = `// Code generated by httptest; DO NOT EDIT.
// github.com/protogodev/httptest

package {{$.DstPkgName}}_test

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"

	{{- range $.Data.Imports}}
	{{.ImportString}}
	{{- end}}

	{{- range $.Spec.Imports}}
	{{.ImportString}}
	{{- end}}
)

{{- $srcPkgPrefix := $.Data.SrcPkgQualifier}}
{{- $interfaceName := $.Data.InterfaceName}}
{{- $mockInterfaceName := printf "%sMock" $interfaceName}}

{{- if $.Spec.Codec}}
var serverCodec = {{$.Spec.Codec}}
{{- else}}
var serverCodec = structool.New().TagName("httptest").
	DecodeHook(
		structool.DecodeStringToError,
		structool.DecodeStringToTime(time.RFC3339),
		structool.DecodeStringToDuration,
	).
	EncodeHook(
		structool.EncodeErrorToString,
		structool.EncodeTimeToString(time.RFC3339),
		structool.EncodeDurationToString,
	)
{{- end}} {{/* if $.Spec.Codec */}}

// Ensure that {{$mockInterfaceName}} does implement {{$srcPkgPrefix}}{{$interfaceName}}.
var _ {{$srcPkgPrefix}}{{$interfaceName}} = &{{$mockInterfaceName}}{}

type {{$mockInterfaceName}} struct {
{{- range $.Data.Methods}}
	{{.Name}}Func func({{.ArgList}}) {{.ReturnArgNamedValueList}}
{{- end}}
}

{{- range $.Data.Methods}}

func (mock *{{$mockInterfaceName}}) {{.Name}}({{.ArgList}}) {{.ReturnArgNamedValueList}} {
	if mock.{{.Name}}Func == nil {
		panic("{{$mockInterfaceName}}.{{.Name}}Func: not implemented")
	}
	return mock.{{.Name}}Func({{.CallArgList}})
}
{{- end}}

type ServerRequest struct {
	Method string
	Path   string
	Header http.Header
	Body   []byte
}

func (r ServerRequest) ServedBy(handler http.Handler) *http.Response {
	req := httptest.NewRequest(r.Method, r.Path, nil)
	if len(r.Body) > 0 {
		req = httptest.NewRequest(r.Method, r.Path, bytes.NewReader(r.Body))
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
	}

	for key, values := range r.Header {
		for _, v := range values {
			req.Header.Add(key, v)
		}
	}

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	return w.Result()
}

type ServerResponse struct {
	StatusCode int
	Header     http.Header
	Body       []byte
}

func (want ServerResponse) Equal(resp *http.Response) error {
	gotStatusCode := resp.StatusCode
	if gotStatusCode != want.StatusCode {
		return fmt.Errorf("StatusCode: got (%d), want (%d)", gotStatusCode, want.StatusCode)
	}
	
	if gotStatusCode == http.StatusNoContent {
		return nil
	}

	var gotHeader http.Header
	if len(want.Header) > 0 {
		gotHeader = http.Header{}
	}
	for key := range want.Header {
		gotHeader[key] = resp.Header.Values(key)
	}
	if !reflect.DeepEqual(gotHeader, want.Header) {
		return fmt.Errorf("Header: got (%#v), want (%#v)", gotHeader, want.Header)
	}

	gotBody, _ := ioutil.ReadAll(resp.Body)
	gotContentType := resp.Header.Get("Content-Type")
	if strings.HasPrefix(gotContentType, "application/json") {
		// Remove the trailing newline from the JSON bytes encoded by Go.
		// See https://github.com/golang/go/issues/37083.
		gotBody = bytes.TrimSuffix(gotBody, []byte("\n"))
	}

	if !bytes.Equal(gotBody, want.Body) {
		return fmt.Errorf("Body: got (%q), want (%q)", gotBody, want.Body)
	}

	return nil
}

{{- range $.Spec.Tests}}

{{$method := interfaceMethod .Name}}
{{$params := $method.Params}}
{{$nonCtxParams := nonCtxParams $params}}

func TestHTTPServer_{{.Name}}(t *testing.T) {
	// in contains all the input parameters (except ctx) of {{.Name}}.
	type in struct {
		{{- range $nonCtxParams}}
		{{title .Name}} {{.TypeString}} ` + "`httptest:\"{{.Name}}\"`" + `
		{{- end}}
	}

	// out contains all the output parameters of {{.Name}}.
	type out struct {
		{{- range $method.Returns}}
		{{title .Name}} {{.TypeString}} ` + "`httptest:\"{{.Name}}\"`" + `
		{{- end}}
	}

	tests := []struct {
		name         string
		request      ServerRequest
		wantIn       map[string]interface{}
		out          map[string]interface{}
		wantResponse ServerResponse
	}{
		{{- range .Subtests}}
		{
			name: "{{.Name}}",
			request: ServerRequest{
				Method: "{{.Request.Method}}",
				Path:   "{{.Request.Path}}",
				{{- if .Request.Header}}
				Header: {{goString .Request.Header}},
				{{- end}}
				{{- if .Request.Body}}
				Body: {{bodyToBytes .Request.Body}},
				{{- end}}
			},
			wantIn: {{goString .WantIn}},
			out: {{goString .Out}},
			wantResponse: ServerResponse{
				StatusCode: {{.WantResponse.StatusCode}},
				{{- if .WantResponse.Header}}
				Header: {{goString .WantResponse.Header}},
				{{- end}}
				{{- if .WantResponse.Body}}
				Body: {{bodyToBytes .WantResponse.Body}},
				{{- end}}
			},
		},
		{{- end}}
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			newServer := {{$.Spec.Server}}

			var out out
			if err := serverCodec.Decode(tt.out, &out); err != nil {
				t.Errorf("err when decoding Out: %v", err)
			}

			var gotIn in
			resp := tt.request.ServedBy(newServer(&{{$mockInterfaceName}}{
				{{.Name}}Func: func({{$method.ArgList}}) {{$method.ReturnArgNamedValueList}} {
					gotIn = in{
						{{- range $nonCtxParams}}
						{{title .Name}}: {{.Name}},
						{{- end}}
					}
					return {{fmtArgCSV $method.ReturnArgValueList "out.>Name"}}
				},
			}))

			encodedGotIn, err := serverCodec.Encode(gotIn)
			if err != nil {
				t.Errorf("err when encoding gotIn: %v", err)
			}

			// Using "%+v" instead of "%#v" as a workaround for https://github.com/go-yaml/yaml/issues/139.
			if fmt.Sprintf("%+v", encodedGotIn) != fmt.Sprintf("%+v", tt.wantIn) {
				t.Errorf("In: Got (%+v) != Want (%+v)", encodedGotIn, tt.wantIn)
			}

			if err := tt.wantResponse.Equal(resp); err != nil {
				t.Error(err.Error())
			}
		})
	}
}
{{- end}}
`
)
