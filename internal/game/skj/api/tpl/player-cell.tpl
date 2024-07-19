{{- define "player-cell" }}
    {{- if .CanFlip }}
        <div class="{{ .Labels }} select" ws-send data-action="flip" data-column="{{ .Column }}" data-row="{{ .Row }}">
    {{- else if .IsEmpty }}
        <div class="{{ .Labels }}">
    {{- end }}
        {{- if .IsVisible }}
            {{ .Card }}
        {{- else }}
            &nbsp;
        {{- end }}
    </div>
{{- end }}