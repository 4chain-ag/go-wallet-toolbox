package {{ .Package }}

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

{{ end }}

