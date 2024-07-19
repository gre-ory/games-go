{{- define "error" }}
<div id="notifications" hx-swap-oob="innerHTML">
    <div class="error">
        <div class="icon-error"></div>
        <div class="message">{{ .Error }}</div>
    </div>
</div>
{{- end }}