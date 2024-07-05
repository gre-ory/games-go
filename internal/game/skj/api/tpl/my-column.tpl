{{- define "my-column" }}
    <div class="my {{ .Labels }}">
        {{- range $rowIndex, $cell := .Cells }}
            {{ template "my-cell" $cell }}
        {{- end }}
    </div>
{{- end }}