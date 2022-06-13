// Code generated by httptest; DO NOT EDIT.
// github.com/protogodev/httptest

package usersvc_test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/RussellLuo/structool"
	"github.com/protogodev/httptest/examples/usersvc"
)

var clientCodec = structool.New().TagName("httptest").
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

type RoundTripFunc func(req *http.Request) (*http.Response, error)

func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

// NewTestClient creates a new *http.Client with Transport mocked.
func NewTestHTTPClient(fn RoundTripFunc) *http.Client {
	return &http.Client{Transport: fn}
}

type ClientRequest struct {
	Method string
	Path   string
	Header http.Header
	Body   []byte
}

func (want ClientRequest) Equal(req *http.Request) error {
	if req.Method != want.Method {
		return fmt.Errorf("Method: got (%s), want (%s)", req.Method, want.Method)
	}
	if req.URL.Path != want.Path {
		return fmt.Errorf("Path: got (%s), want (%s)", req.URL.Path, want.Path)
	}

	var gotHeader http.Header
	if len(want.Header) > 0 {
		gotHeader = http.Header{}
	}
	for key := range want.Header {
		gotHeader[key] = req.Header.Values(key)
	}
	if !reflect.DeepEqual(gotHeader, want.Header) {
		return fmt.Errorf("Header: got (%#v), want (%#v)", gotHeader, want.Header)
	}

	if len(want.Body) == 0 {
		return nil
	}

	gotBody, _ := ioutil.ReadAll(req.Body)
	gotContentType := req.Header.Get("Content-Type")
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

type ClientResponse struct {
	StatusCode int
	Header     http.Header
	Body       []byte
}

func (resp ClientResponse) HTTPResponse() *http.Response {
	statusCode := resp.StatusCode
	if statusCode == 0 {
		statusCode = http.StatusOK
	}

	return &http.Response{
		StatusCode: statusCode,
		Header:     resp.Header,
		Body:       io.NopCloser(bytes.NewReader(resp.Body)),
	}
}

func TestHTTPClient_GetUser(t *testing.T) {
	// in contains all the input parameters (except ctx) of GetUser.
	type in struct {
		Ctx  context.Context
		Name string `httptest:"name"`
	}

	// out contains all the output parameters of GetUser.
	type out struct {
		User *usersvc.User `httptest:"user"`
		Err  error         `httptest:"err"`
	}

	tests := []struct {
		name        string
		in          map[string]interface{}
		wantRequest ClientRequest
		response    ClientResponse
		wantOut     map[string]interface{}
	}{
		{
			name: "ok",
			in:   map[string]interface{}{"name": "foo"},
			wantRequest: ClientRequest{
				Method: "GET",
				Path:   "/users/foo",
			},
			response: ClientResponse{
				StatusCode: 200,
				Body:       []byte(`{"name":"foo","sex":"male","birth":"2022-01-01T00:00:00Z"}`),
			},
			wantOut: map[string]interface{}{"err": "", "user": map[interface{}]interface{}{"birth": "2022-01-01T00:00:00Z", "name": "foo", "sex": "male"}},
		},
		{
			name: "err",
			in:   map[string]interface{}{"name": "foo"},
			wantRequest: ClientRequest{
				Method: "GET",
				Path:   "/users/foo",
			},
			response: ClientResponse{
				StatusCode: 400,
				Body:       []byte(`{"error":"not found"}`),
			},
			wantOut: map[string]interface{}{"err": "not found", "user": interface{}(nil)},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var gotRequest *http.Request
			httpClient := NewTestHTTPClient(func(req *http.Request) (*http.Response, error) {
				gotRequest = req
				return tt.response.HTTPResponse(), nil
			})
			sut, err := usersvc.NewHTTPClient(httpClient, "http://localhost:8080")
			if err != nil {
				t.Errorf("err when creating Client: %v", err)
			}

			var in in
			if err := serverCodec.Decode(tt.in, &in); err != nil {
				t.Errorf("err when decoding In: %v", err)
			}
			in.Ctx = context.Background()

			var gotOut out
			gotOut.User, gotOut.Err = sut.GetUser(in.Ctx, in.Name)

			if err := tt.wantRequest.Equal(gotRequest); err != nil {
				t.Error(err.Error())
			}

			encodedGotOut, err := clientCodec.Encode(gotOut)
			if err != nil {
				t.Errorf("err when encoding gotOut: %v", err)
			}

			// Using "%+v" instead of "%#v" as a workaround for https://github.com/go-yaml/yaml/issues/139.
			if fmt.Sprintf("%+v", encodedGotOut) != fmt.Sprintf("%+v", tt.wantOut) {
				t.Errorf("Out: Got (%+v) != Want (%+v)", encodedGotOut, tt.wantOut)
			}
		})
	}
}

func TestHTTPClient_ListUsers(t *testing.T) {
	// in contains all the input parameters (except ctx) of ListUsers.
	type in struct {
		Ctx context.Context
	}

	// out contains all the output parameters of ListUsers.
	type out struct {
		Users []*usersvc.User `httptest:"users"`
		Err   error           `httptest:"err"`
	}

	tests := []struct {
		name        string
		in          map[string]interface{}
		wantRequest ClientRequest
		response    ClientResponse
		wantOut     map[string]interface{}
	}{
		{
			name: "ok",
			in:   map[string]interface{}(nil),
			wantRequest: ClientRequest{
				Method: "GET",
				Path:   "/users",
			},
			response: ClientResponse{
				StatusCode: 200,
				Body:       []byte(`{"users":[{"name":"foo","sex":"male","birth":"2022-01-01T00:00:00Z"}]}`),
			},
			wantOut: map[string]interface{}{"err": "", "users": []interface{}{map[interface{}]interface{}{"birth": "2022-01-01T00:00:00Z", "name": "foo", "sex": "male"}}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var gotRequest *http.Request
			httpClient := NewTestHTTPClient(func(req *http.Request) (*http.Response, error) {
				gotRequest = req
				return tt.response.HTTPResponse(), nil
			})
			sut, err := usersvc.NewHTTPClient(httpClient, "http://localhost:8080")
			if err != nil {
				t.Errorf("err when creating Client: %v", err)
			}

			var in in
			if err := serverCodec.Decode(tt.in, &in); err != nil {
				t.Errorf("err when decoding In: %v", err)
			}
			in.Ctx = context.Background()

			var gotOut out
			gotOut.Users, gotOut.Err = sut.ListUsers(in.Ctx)

			if err := tt.wantRequest.Equal(gotRequest); err != nil {
				t.Error(err.Error())
			}

			encodedGotOut, err := clientCodec.Encode(gotOut)
			if err != nil {
				t.Errorf("err when encoding gotOut: %v", err)
			}

			// Using "%+v" instead of "%#v" as a workaround for https://github.com/go-yaml/yaml/issues/139.
			if fmt.Sprintf("%+v", encodedGotOut) != fmt.Sprintf("%+v", tt.wantOut) {
				t.Errorf("Out: Got (%+v) != Want (%+v)", encodedGotOut, tt.wantOut)
			}
		})
	}
}

func TestHTTPClient_CreateUser(t *testing.T) {
	// in contains all the input parameters (except ctx) of CreateUser.
	type in struct {
		Ctx  context.Context
		User *usersvc.User `httptest:"user"`
	}

	// out contains all the output parameters of CreateUser.
	type out struct {
		Err error `httptest:"err"`
	}

	tests := []struct {
		name        string
		in          map[string]interface{}
		wantRequest ClientRequest
		response    ClientResponse
		wantOut     map[string]interface{}
	}{
		{
			name: "ok",
			in:   map[string]interface{}{"user": map[interface{}]interface{}{"birth": "2022-01-01T00:00:00Z", "name": "foo", "sex": "male"}},
			wantRequest: ClientRequest{
				Method: "POST",
				Path:   "/users",
				Body:   []byte(`{"name":"foo","sex":"male","birth":"2022-01-01T00:00:00Z"}`),
			},
			response: ClientResponse{
				StatusCode: 204,
			},
			wantOut: map[string]interface{}{"err": ""},
		},
		{
			name: "err",
			in:   map[string]interface{}{"user": map[interface{}]interface{}{"birth": "2022-01-01T00:00:00Z", "name": "foo", "sex": "male"}},
			wantRequest: ClientRequest{
				Method: "POST",
				Path:   "/users",
				Body:   []byte(`{"name":"foo","sex":"male","birth":"2022-01-01T00:00:00Z"}`),
			},
			response: ClientResponse{
				StatusCode: 400,
				Body:       []byte(`{"error":"already exists"}`),
			},
			wantOut: map[string]interface{}{"err": "already exists"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var gotRequest *http.Request
			httpClient := NewTestHTTPClient(func(req *http.Request) (*http.Response, error) {
				gotRequest = req
				return tt.response.HTTPResponse(), nil
			})
			sut, err := usersvc.NewHTTPClient(httpClient, "http://localhost:8080")
			if err != nil {
				t.Errorf("err when creating Client: %v", err)
			}

			var in in
			if err := serverCodec.Decode(tt.in, &in); err != nil {
				t.Errorf("err when decoding In: %v", err)
			}
			in.Ctx = context.Background()

			var gotOut out
			gotOut.Err = sut.CreateUser(in.Ctx, in.User)

			if err := tt.wantRequest.Equal(gotRequest); err != nil {
				t.Error(err.Error())
			}

			encodedGotOut, err := clientCodec.Encode(gotOut)
			if err != nil {
				t.Errorf("err when encoding gotOut: %v", err)
			}

			// Using "%+v" instead of "%#v" as a workaround for https://github.com/go-yaml/yaml/issues/139.
			if fmt.Sprintf("%+v", encodedGotOut) != fmt.Sprintf("%+v", tt.wantOut) {
				t.Errorf("Out: Got (%+v) != Want (%+v)", encodedGotOut, tt.wantOut)
			}
		})
	}
}

func TestHTTPClient_UpdateUser(t *testing.T) {
	// in contains all the input parameters (except ctx) of UpdateUser.
	type in struct {
		Ctx  context.Context
		Name string        `httptest:"name"`
		User *usersvc.User `httptest:"user"`
	}

	// out contains all the output parameters of UpdateUser.
	type out struct {
		Err error `httptest:"err"`
	}

	tests := []struct {
		name        string
		in          map[string]interface{}
		wantRequest ClientRequest
		response    ClientResponse
		wantOut     map[string]interface{}
	}{
		{
			name: "ok",
			in:   map[string]interface{}{"name": "foo", "user": map[interface{}]interface{}{"birth": "2022-01-01T00:00:00Z", "sex": "male"}},
			wantRequest: ClientRequest{
				Method: "PATCH",
				Path:   "/users/foo",
				Body:   []byte(`{"sex":"male","birth":"2022-01-01T00:00:00Z"}`),
			},
			response: ClientResponse{
				StatusCode: 204,
			},
			wantOut: map[string]interface{}{"err": ""},
		},
		{
			name: "err",
			in:   map[string]interface{}{"name": "foo", "user": map[interface{}]interface{}{"birth": "2022-01-01T00:00:00Z", "sex": "male"}},
			wantRequest: ClientRequest{
				Method: "PATCH",
				Path:   "/users/foo",
				Body:   []byte(`{"sex":"male","birth":"2022-01-01T00:00:00Z"}`),
			},
			response: ClientResponse{
				StatusCode: 400,
				Body:       []byte(`{"error":"not found"}`),
			},
			wantOut: map[string]interface{}{"err": "not found"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var gotRequest *http.Request
			httpClient := NewTestHTTPClient(func(req *http.Request) (*http.Response, error) {
				gotRequest = req
				return tt.response.HTTPResponse(), nil
			})
			sut, err := usersvc.NewHTTPClient(httpClient, "http://localhost:8080")
			if err != nil {
				t.Errorf("err when creating Client: %v", err)
			}

			var in in
			if err := serverCodec.Decode(tt.in, &in); err != nil {
				t.Errorf("err when decoding In: %v", err)
			}
			in.Ctx = context.Background()

			var gotOut out
			gotOut.Err = sut.UpdateUser(in.Ctx, in.Name, in.User)

			if err := tt.wantRequest.Equal(gotRequest); err != nil {
				t.Error(err.Error())
			}

			encodedGotOut, err := clientCodec.Encode(gotOut)
			if err != nil {
				t.Errorf("err when encoding gotOut: %v", err)
			}

			// Using "%+v" instead of "%#v" as a workaround for https://github.com/go-yaml/yaml/issues/139.
			if fmt.Sprintf("%+v", encodedGotOut) != fmt.Sprintf("%+v", tt.wantOut) {
				t.Errorf("Out: Got (%+v) != Want (%+v)", encodedGotOut, tt.wantOut)
			}
		})
	}
}

func TestHTTPClient_DeleteUser(t *testing.T) {
	// in contains all the input parameters (except ctx) of DeleteUser.
	type in struct {
		Ctx  context.Context
		Name string `httptest:"name"`
	}

	// out contains all the output parameters of DeleteUser.
	type out struct {
		Err error `httptest:"err"`
	}

	tests := []struct {
		name        string
		in          map[string]interface{}
		wantRequest ClientRequest
		response    ClientResponse
		wantOut     map[string]interface{}
	}{
		{
			name: "ok",
			in:   map[string]interface{}{"name": "foo"},
			wantRequest: ClientRequest{
				Method: "DELETE",
				Path:   "/users/foo",
			},
			response: ClientResponse{
				StatusCode: 204,
			},
			wantOut: map[string]interface{}{"err": ""},
		},
		{
			name: "err",
			in:   map[string]interface{}{"name": "foo"},
			wantRequest: ClientRequest{
				Method: "DELETE",
				Path:   "/users/foo",
			},
			response: ClientResponse{
				StatusCode: 400,
				Body:       []byte(`{"error":"not found"}`),
			},
			wantOut: map[string]interface{}{"err": "not found"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var gotRequest *http.Request
			httpClient := NewTestHTTPClient(func(req *http.Request) (*http.Response, error) {
				gotRequest = req
				return tt.response.HTTPResponse(), nil
			})
			sut, err := usersvc.NewHTTPClient(httpClient, "http://localhost:8080")
			if err != nil {
				t.Errorf("err when creating Client: %v", err)
			}

			var in in
			if err := serverCodec.Decode(tt.in, &in); err != nil {
				t.Errorf("err when decoding In: %v", err)
			}
			in.Ctx = context.Background()

			var gotOut out
			gotOut.Err = sut.DeleteUser(in.Ctx, in.Name)

			if err := tt.wantRequest.Equal(gotRequest); err != nil {
				t.Error(err.Error())
			}

			encodedGotOut, err := clientCodec.Encode(gotOut)
			if err != nil {
				t.Errorf("err when encoding gotOut: %v", err)
			}

			// Using "%+v" instead of "%#v" as a workaround for https://github.com/go-yaml/yaml/issues/139.
			if fmt.Sprintf("%+v", encodedGotOut) != fmt.Sprintf("%+v", tt.wantOut) {
				t.Errorf("Out: Got (%+v) != Want (%+v)", encodedGotOut, tt.wantOut)
			}
		})
	}
}
