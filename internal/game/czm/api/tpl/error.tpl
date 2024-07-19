{{- define "error" }}
<div id="notifications" hx-swap-oob="innerHTML">
    <div class="error">
        <div class="icon-error"></div>
        <div class="message">{{ .error }}</div>
    </div>
</div>
{{- end }}