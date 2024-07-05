{{- define "player-board" }}
    <div class="{{ .Labels }}">
        <div class="score">{{ .Total }}</div>
        {{- range $columnIndex, $column := .Columns }}
            {{ template "player-column" $column }}
        {{- end }}
    </div>
{{- end }}