message {{.Name}} {
{{- range .OneOfs }}
    oneof {{.Name}} {
        {{- range .Fields }}
            {{.Type}} {{.Name}} = {{.Index}};
        {{- end }}
    }
{{- end }}
{{- range .Fields }}
    {{ .Qualifier }} {{.Type}} {{.Name}} = {{.Index}};
{{- end }}
{{ range .Enums }}
    {{ template "enum.proto.tpl" . }}
{{ end }}
{{- range .SubMessages }}
    {{ template "message.proto.tpl" . }}
{{- end }}
}