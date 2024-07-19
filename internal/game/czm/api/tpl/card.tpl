{{- define "card" }}
{{- if .Selectable }}
<div class="{{ .Card.Labels }} selectable" ws-send data-action="select-card" data-card="{{ .CardIndex }}">
{{- else if .Selected }}
<div class="{{ .Card.Labels }} selected">
{{- else }}
<div class="{{ .Card.Labels }}">
{{- end }}
    <div class="symbol-top">{{ .Card.Value }}</div>
    <div class="value">{{ .Card.Value }}</div>
    <div class="symbol-bottom">{{ .Card.Value }}</div>
</div>
{{- end }}