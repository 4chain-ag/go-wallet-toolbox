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


func NewClient(addr string, overrideOptions ...StorageClientOverrides) (*{{ $clientName}}, func(), error) {
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

