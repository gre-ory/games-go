{{- define "user-badge" }}
<!-- user -->
<div id="user" class="user" hx-get="/htmx/user" hx-target="this" hx-swap="outerHTML" hx-trigger="load"></div>
{{- end }}