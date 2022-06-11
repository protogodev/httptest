package httptest

import (
	"encoding/hex"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/protogodev/httptest/client"
	"github.com/protogodev/httptest/server"
	protogocmd "github.com/protogodev/protogo/cmd"
	"github.com/protogodev/protogo/generator"
	"github.com/protogodev/protogo/parser"
	"github.com/protogodev/protogo/parser/ifacetool"
)

func init() {
	protogocmd.MustRegister(&protogocmd.Plugin{
		Name: "httptest",
		Cmd:  protogocmd.NewGen(&Generator{}),
	})
}

type Generator struct {
	OutFileName      string `name:"out" help:"output filename (default \"./<srcPkgName>_<mode>_test.go\")"`
	Formatted        bool   `name:"fmt" default:"true" help:"whether to make the test code formatted"`
	TestSpecFileName string `name:"spec" required:"" help:"the test specification in YAML"`
	Mode             string `name:"mode" required:"" enum:"server,client" help:"generation mode (server or client)"`
}

func (g *Generator) Generate(data *ifacetool.Data) (*generator.File, error) {
	if g.OutFileName == "" {
		g.OutFileName = fmt.Sprintf("./%s_%s_test.go", data.SrcPkgName, g.Mode)
	}

	var template string
	var tmplData interface{}

	switch g.Mode {
	case "server":
		spec, err := server.NewSpec(g.TestSpecFileName)
		if err != nil {
			return nil, err
		}
		tmplData = struct {
			DstPkgName string
			Data       *ifacetool.Data
			Spec       *server.Spec
		}{
			DstPkgName: parser.PkgNameFromDir(filepath.Dir(g.OutFileName)),
			Data:       data,
			Spec:       spec,
		}

		template = server.Template

	case "client":
		spec, err := client.NewSpec(g.TestSpecFileName)
		if err != nil {
			return nil, err
		}
		tmplData = struct {
			DstPkgName string
			Data       *ifacetool.Data
			Spec       *client.Spec
		}{
			DstPkgName: parser.PkgNameFromDir(filepath.Dir(g.OutFileName)),
			Data:       data,
			Spec:       spec,
		}

		template = client.Template

	default:
		panic(fmt.Errorf("bad mode: %s", g.Mode))
	}

	methodMap := make(map[string]*ifacetool.Method)
	for _, method := range data.Methods {
		methodMap[method.Name] = method
	}

	return generator.Generate(template, tmplData, generator.Options{
		Funcs: map[string]interface{}{
			"title": strings.Title,
			"fmtArgCSV": func(csv string, format string) string {
				if csv == "" {
					return ""
				}

				sep := ", "
				args := strings.Split(csv, sep)

				var results []string
				for _, a := range args {
					r := strings.NewReplacer("$Name", a, ">Name", strings.Title(a))
					results = append(results, r.Replace(format))
				}

				return strings.Join(results, sep)
			},
			"interfaceMethod": func(name string) *ifacetool.Method {
				method, ok := methodMap[name]
				if !ok {
					return nil
				}
				return method
			},
			"goString": func(v interface{}) string {
				return fmt.Sprintf("%#v", v)
			},
			"ctxParam": func(params []*ifacetool.Param) *ifacetool.Param {
				for _, p := range params {
					if p.TypeString == "context.Context" {
						return p
					}
				}
				return nil
			},
			"nonCtxParams": func(params []*ifacetool.Param) (out []*ifacetool.Param) {
				for _, p := range params {
					if p.TypeString != "context.Context" {
						out = append(out, p)
					}
				}
				return
			},
			"bodyToBytes": func(s string) string {
				if s == "" {
					// An empty string indicates a nil byte slice.
					return "[]byte(nil)"
				}

				if strings.HasPrefix(s, "0x") {
					// This is a hexadecimal string, decode it into bytes.
					//
					// Note that kun borrows the idea from eth2.0 to represent binary data
					// as hex encoded strings, see https://github.com/ethereum/eth2.0-spec-tests/issues/5.
					decoded, err := hex.DecodeString(s[2:])
					if err != nil {
						panic(err)
					}

					var hexes []string
					for _, b := range decoded {
						hexes = append(hexes, fmt.Sprintf("0x%x", b))
					}
					return fmt.Sprintf("[]byte{%s}", strings.Join(hexes, ", "))
				}

				// This is a normal string, leave it as is.
				return fmt.Sprintf("[]byte(`%s`)", s)
			},
		},
		Formatted:      g.Formatted,
		TargetFileName: g.OutFileName,
	})
}
