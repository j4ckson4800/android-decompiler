syntax = "proto3";
package {{.PackageName}};
option go_package = "{{.GoPackageName}}";

{{ range .Imports }}
import "{{.}}.proto";
{{ end }}


{{- range .Messages }}
{{ if .IsGlobal }}{{ include "message.proto.tpl" . }}{{else}}{{ end }}
{{- end }}

{{- range .Enums }}
{{ template "enum.proto.tpl" . }}
{{- end }}