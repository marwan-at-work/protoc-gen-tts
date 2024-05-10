package tts

import (
	"text/template"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var title = cases.Title(language.English, cases.NoLower)

var tmpl = template.Must(template.New("tts").Funcs(template.FuncMap{
	"title": func(s string) string {
		return title.String(s)
	},
}).Parse(serviceTemplate))

const serviceTemplate = `/* tslint:disable */

// This file has been generated by marwan.io/protoc-gen-tss.
// Do not edit.

{{ if .Services }}
import {
  createTwirpRequest,
  Fetch,
  throwTwirpError
} from './twirp'
{{ end }}

{{ range .Imports }}
import {
{{ range .Declarations }}
  {{ . }},
{{ end }}
} from './{{ .Name }}'
{{ end }}

{{ range .Services }}
{{ .Doc }}
export interface {{.Name}} {
	{{ range .Methods }}
	{{ .Doc }}
	{{ .Name }}: (request: {{.Input}}Properties, headers?: object) => Promise<{{.Output}}>
	{{ end }}
}

{{ .Doc }}
export class {{ .Name }}Client implements {{ .Name }} {
  private hostname: string
  private fetch: Fetch
  private path = '{{ .PathPrefix }}'
  private opts: RequestInit

  constructor(hostname: string, fetch: Fetch, opts: RequestInit = {}) {
    this.hostname = hostname
    this.fetch = fetch
    this.opts = opts
  }

  private url(name: string): string {
    return this.hostname + this.path + name
  }

  {{ range .Methods }}
  {{ .Doc }}
  public {{ .Name }}({{ if .InputHasFields }}request{{ else }}_{{ end }}: {{ .Input }}Properties, headers: object = {}): Promise<{{ .Output }}> {
  {{ if .InputHasFields }}
	const body = new {{ .Input }}(request);
  {{ else }}
  const body = new {{ .Input }}();
  {{ end -}}
    return this.fetch(
      this.url('{{ .Name | title }}'),
      createTwirpRequest(body, headers, this.opts)
    ).then((res) => {
      if (!res.ok) {
        return throwTwirpError(res)
      }
      return res.json().then((props) => { return {{ .Output }}.fromJSON(props)})
    })
  }
  {{ end }} 
}
{{ end }}

{{ range .Enums }}
export enum {{ .Name }} {
	{{- range .Values }}
	{{ . }} = '{{ . }}',
	{{- end -}}
}
{{ end }}

{{ range .Messages }}
{{ .Doc }}
export interface {{ .Name }}Properties {
	{{- range .Fields }}
	{{ .Name }}?: {{ .PrintTypeProperties }}
	{{- end }}
	toJSON?(): object
}

interface {{ .Name }}JSON {
	{{- range .Fields }}
	{{ .JSONName }}?: {{ .PrintType }}
	{{- end -}}
}

{{ .Doc }}
export class {{ .Name }} implements {{ .Name }}Properties {
  {{ $msg := . }}
  {{ range .Fields -}}
  {{ .Name }}{{ if $msg.Optional }}?{{ end }}: {{ .PrintType }}
  {{ end }}

  {{ if (gt (len .Fields) 0) }}
	constructor(props: {{ .Name }}Properties) {
    {{- range .Fields }}
    this.{{ .Name }} = {{ .SetConstructorProp $msg.Optional }}
    {{- end -}}
	}
  {{ end }}

  {{ if (gt (len .Fields) 0) }}
  static fromJSON(props: {{ .Name }}JSON): {{ .Name }} {
    if (!props) {
      props = {};
    }
    return new {{ .Name }}({
      {{- range .Fields }}
      {{ .Name }}: {{ .ResolveType }},
      {{- end -}}
    })
  }
  {{ else }}
  static fromJSON(_: {{ .Name }}JSON): {{ .Name }} {
    return new {{ .Name }}();
  }
  {{ end }}

	public toJSON(): {{ .Name }}JSON {
    return {
      {{- range .Fields }}
      '{{ .JSONName }}': this.{{ .Name }},
      {{- end -}}
    }
  }
  
  public toObject(): {{ .Name }}Properties {
    return {
      {{- range .Fields }}
      {{ .Name }}: {{ .SetToObjectProp $msg.Optional }},
      {{- end }}
    }
  }
}
{{ end }}
`
