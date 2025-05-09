message {{.Name}} {
{{- range .OneOfs }}
    oneof {{.Name}} {
        {{- range .Fields }}
            {{.Type}} {{.Name}} = {{.Index}};
        {{- end }}
    }
{{- end }}
{{- range .Fields }}
    {{with .Qualifier}}{{ . }} {{ end }}{{.Type}} {{.Name}} = {{.Index}};
{{- end }}
{{- range .Enums }}
    {{ template "enum.proto.tpl" . }}
{{- end }}
{{- range .SubMessages }}
{{ include "message.proto.tpl" . | indent 4}}

{{- end }}
}