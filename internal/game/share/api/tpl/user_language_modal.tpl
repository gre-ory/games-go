{{- define "user-language-modal" }}
<div id="user-language-modal" class="modal" _="on closeModal add .closing then wait for animationend then remove me">
	<div class="modal-underlay" _="on click trigger closeModal"></div>
	<div class="modal-content">
		{{- $current_language := .user.Language }}
		{{- range $i, $group := .available_languages }}
		<div class="cols-5">
		{{- range $j, $language := $group }}
		{{- if eq $language $current_language }}
		<div class="col-1 center select on">
			<div
				class="language-{{ $language }} m"
			 	_="on click trigger closeModal">{{ $language }}</div>
		{{- else }}
		<div class="col-1 center select">
			<div 
				class="language-{{ $language }} m" 
				hx-put="/htmx/user?user_language={{ $language }}" 
				hx-target="#user" 
				hx-swap="outerHTML"
				_="on click trigger closeModal">{{ $language }}</div>
		{{- end }}
		</div>
		{{- end }}
		</div>
		{{- end }}
		<div class="right">
			<button _="on click trigger closeModal">Cancel</button>
		</div>
	</div>
</div>
{{- end }}