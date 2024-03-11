{{- define "error" }}
    <div id="errors" class="htmx-content" hx-swap-oob="innerHTML">
        error: {{ .error }}
    </div>
{{- end }}