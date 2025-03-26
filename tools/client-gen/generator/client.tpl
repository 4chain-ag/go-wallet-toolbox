package {{ .Package }}

import (
	"context"
	"github.com/filecoin-project/go-jsonrpc"
)

{{ range .Interfaces }}
{{ $clientName := printf "%sClient" .Name }}

type {{ $clientName }} struct {
	client *rpc{{ .Name }}
}

	{{ range .Methods }}
func (c *{{ $clientName }}) {{ .Name }}({{ range .Arguments }}{{ .Name }} {{ .Type }}, {{ end }}) ({{ range .Results }}{{ .Type }},{{ end }}) {
	{{ if gt (len .Results) 0 }} return {{ end }} c.client.{{ .Name }}({{ range .Arguments }}{{ .Name }}, {{ end }})
}
	{{ end }}

type rpc{{ .Name }} struct {
	{{ range .Methods }}
	{{ .Name }} func({{ range .Arguments }}{{ .Type }}, {{ end }}) ({{ range .Results }}{{ .Type }},{{ end }})
	{{ end }}
}


func NewClient(addr string, overrideOptions ...InternalOverrides) (*{{ $clientName}}, func(), error) {
    opts := defaultClientOptions()
    client := &{{ $clientName}}{
        client: &rpc{{ .Name }}{},
    }

    for _, opt := range overrideOptions {
      opt(&opts)
    }

    cleanup, err := jsonrpc.NewMergeClient(
      context.Background(),
      addr,
      "remote_storage",
      []any{client.client},
      nil,
      opts.options...
    )

    return client, cleanup, err
}

{{ end }}

// ===== TODO: REMOVE BELOW LINES FROM TEMPLATE =====
// ===== THIS IS JUST FOR SHOWING THE MODEL =====

// {{ .FullPackage }}
// Is in the same package: {{ .IsSamePackage }}

type InterfaceInfo struct {
	Name    string
	Methods []MethodInfo
}

type MethodInfo struct {
	Name      string
	Arguments []ParamInfo
	Results   []TypeInfo
}

type ParamInfo struct {
	Name string
	Type string
}

type TypeInfo struct {
	Type string
}

var Interfaces = []InterfaceInfo{
{{ range .Interfaces }}	{
		Name: "{{ .Name }}",
		Methods: []MethodInfo{
		{{ range .Methods }}	{
				Name: "{{ .Name }}",
				Arguments: []ParamInfo{
				{{ range .Arguments }}	{Name: "{{ .Name }}", Type: "{{ .Type }}"},
				{{ end }}},
				Results: []TypeInfo{
				{{ range .Results }}	{Type: "{{ .Type }}"},
				{{ end }}},
			},
		{{ end }}},
	},
{{ end }}}
