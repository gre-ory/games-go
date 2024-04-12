{{- define "user" }}
<div id="user" class="user">
    {{template "user-content" . }}
</div>
{{- end }}

{{- define "user-oob" }}
<div id="user" class="user" hx-target="#user" hx-swap-oob="outerHTML">
    {{template "user-content" . }}
</div>
{{- end }}

{{- define "user-content" }}
    <div class="avatar-{{ .user.Avatar }} s click" hx-get="/htmx/user-avatar-modal" hx-trigger="click" hx-target="body" hx-swap="beforeend"></div>  
    <div class="id">{{ .user.Id }}</div>
    <div class="name s click" hx-get="/htmx/user-name-modal" hx-trigger="click" hx-target="body" hx-swap="beforeend">{{ .user.Name }}</div>
    <div class="language-{{ .user.Language }} click" hx-get="/htmx/user-language-modal" hx-trigger="click" hx-target="body" hx-swap="beforeend"></div>
{{- end }}