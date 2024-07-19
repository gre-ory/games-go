{{- define "players" }}
{{- $lang := .Lang }}
<div id="players" hx-swap-oob="outerHTML">
    <div class="players cols-2">
    {{- $game := .Game }}
    {{- $current_player := .Player }}
    {{- range .Game.Players }}
        {{- if eq .Id $current_player.Id }}
            <div class="{{ $game.PlayerLabels .Id }} current col-1 item">
        {{- else }}
            <div class="{{ $game.PlayerLabels .Id }} col-1 item">
        {{- end }}
                <div class="title center">
                    {{ .User.Avatar.XS }}
                    <div class="name truncate">{{ .User.Name }}</div>
                </div>
                <div class="content center">
                    {{- if eq .Id $current_player.Id }}
                        {{ $game.YourPlayerMessage $lang .Id }}
                    {{- else }}
                        {{ $game.PlayerMessage $lang .Id }}
                    {{- end }}
                </div>
            </div>
    {{- end }}
    </div>
</div>
{{- end }}