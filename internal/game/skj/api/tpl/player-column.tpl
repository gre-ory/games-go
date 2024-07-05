{{- define "player-column" }}
    <div class="{{ .Labels }}">
        {{- range $rowIndex, $cell := .Cells }}
            {{ template "player-cell" $cell }}
        {{- end }}
    </div>
{{- end }}