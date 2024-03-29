{{- define "info" }}
<div id="notifications" hx-swap-oob="innerHTML">
    <div class="info">
        <div class="icon-info"></div>
        <div class="message">{{ .info }}</div>
    </div>
</div>
{{- end }}