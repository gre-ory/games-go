{{- define "my-board" }}
    <div class="my {{ .Labels }}">
        <div class="score">{{ .Total }}</div>
        {{- range $columnIndex, $column := .Columns }}
            {{ template "my-column" $column }}
        {{- end }}
    </div>
{{- end }}