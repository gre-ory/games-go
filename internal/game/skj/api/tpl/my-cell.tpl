{{- define "my-cell" }}
    {{- if .CanFlip }}
        <div class="my {{ .Labels }} select" ws-send data-action="flip" data-column="{{ .Column }}" data-row="{{ .Row }}">
    {{- else if .IsEmpty }}
        <div class="my {{ .Labels }}">
    {{- else }}
        {{- if .IsVisible }}
            {{ .Card }}
        {{- else }}
            &nbsp;
        {{- end }}
    </div>
{{- end }}