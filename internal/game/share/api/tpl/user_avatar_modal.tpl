{{- define "user-avatar-modal" }}
<div id="user-avatar-modal" class="modal" _="on closeModal add .closing then wait for animationend then remove me">
	<div class="modal-underlay" _="on click trigger closeModal"></div>
	<div class="modal-content">
		{{- $current_avatar := .user.Avatar }}
		{{- range $i, $group := .available_avatars }}
		<div class="cols-5">
		{{- range $j, $avatar := $group }}
		{{- if eq $avatar $current_avatar }}
		<div class="col-1 center select on">
			<div
				class="avatar-{{ $avatar }} xl"
			 	_="on click trigger closeModal"></div>
		{{- else }}
		<div class="col-1 center select">
			<div 
				class="avatar-{{ $avatar }} xl" 
				hx-put="/htmx/user?user_avatar={{ $avatar }}" 
				hx-target="#user" 
				hx-swap="outerHTML"
				_="on click trigger closeModal"></div>
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