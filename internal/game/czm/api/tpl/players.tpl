{{- define "players" }}
{{- $lang := .lang }}
<div id="players" hx-swap-oob="outerHTML">
    <div class="players cols-2">
    {{- $game := .game }}
    {{- $current_player := .player }}
    {{- range .game.Players }}
        {{- if eq .Id $current_player.Id }}
            <div class="{{ $game.PlayerLabels .Id }} current col-1 item">
        {{- else }}
            <div class="{{ $game.PlayerLabels .Id }} col-1 item">
        {{- end }}
                <div class="title center">
                    {{ .Avatar.ExtraSmallHtml }}
                    <div class="name truncate">{{ or .Name .Id }}</div>
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