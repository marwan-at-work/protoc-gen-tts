package tts

import (
	"strings"
	"text/template"
)

var tmpl = template.Must(template.New("tts").Funcs(template.FuncMap{
	"title": strings.Title,
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

  constructor(hostname: string, fetch: Fetch) {
    this.hostname = hostname
    this.fetch = fetch
  }

  private url(name: string): string {
    return this.hostname + this.path + name
  }

  {{ range .Methods }}
  {{ .Doc }}
  public {{ .Name }}(request: {{ .Input }}Properties, headers: object = {}): Promise<{{ .Output }}> {
	const body = new {{ .Input }}(request)
    return this.fetch(
      this.url('{{ .Name | title }}'),
      createTwirpRequest(body, headers)
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
	constructor(props?: {{ .Name }}Properties) {
		if (props) {
			{{- range .Fields }}
			this.{{ .Name }} = {{ .SetConstructorProp $msg.Optional }}
			{{- end -}}
		}
	}

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
