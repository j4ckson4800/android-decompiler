enum {{.Name}} {
    {{- range .Values }}
    {{.Name}} = {{.Value}};
    {{- end }}
}