package client

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"gopkg.in/yaml.v2"
)

type request struct {
	Method string      `yaml:"method"`
	Path   string      `yaml:"path"`
	Header http.Header `yaml:"header"`
	Body   string      `yaml:"body"`
}

type response struct {
	StatusCode string      `yaml:"statusCode"`
	Header     http.Header `yaml:"header"`
	Body       string      `yaml:"body"`
}

type Subtest struct {
	Name        string                 `yaml:"name"`
	In          map[string]interface{} `yaml:"in"`
	WantRequest request                `yaml:"wantRequest"`
	Response    response               `yaml:"response"`
	WantOut     map[string]interface{} `yaml:"wantOut"`
}

type Test struct {
	Name     string    `yaml:"name"`
	Subtests []Subtest `yaml:"subtests"`
}

type Import struct {
	Path  string `yaml:"path"`
	Alias string `yaml:"alias"`
}

func (i Import) ImportString() string {
	s := fmt.Sprintf("%q", i.Path)
	if i.Alias != "" {
		s = i.Alias + " " + s
	}
	return s
}

type Spec struct {
	RawImports []string `yaml:"imports"`
	Imports    []Import `yaml:"-"`
	Client     string   `yaml:"handler"`
	Codec      string   `yaml:"codec"`
	Tests      []Test   `yaml:"tests"`
}

func NewSpec(testFilename string) (*Spec, error) {
	b, err := ioutil.ReadFile(testFilename)
	if err != nil {
		return nil, err
	}

	spec := &Spec{}
	err = yaml.Unmarshal(b, spec)
	if err != nil {
		return nil, err
	}

	imports, err := getImports(spec.RawImports)
	if err != nil {
		return nil, err
	}
	spec.Imports = append(spec.Imports, imports...)

	return spec, nil
}

func getImports(rawImports []string) (imports []Import, err error) {
	var path, alias string

	for i, str := range rawImports {
		fields := strings.Fields(str)
		switch len(fields) {
		case 1:
			alias, path = "", fields[0]
		case 2:
			alias, path = fields[0], fields[1]
		default:
			return nil, fmt.Errorf("invalid path in imports[%d]: %s", i, str)
		}

		if !strings.HasPrefix(path, `"`) {
			path = fmt.Sprintf("%q", path)
		}
		imports = append(imports, Import{
			Path:  path,
			Alias: alias,
		})
	}

	return
}
