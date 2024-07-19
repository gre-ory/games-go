{{- define "players" }}
{{- $lang := .Lang }}
<div id="players" hx-swap-oob="outerHTML">
    <div class="players cols-2">
    {{- $game := .Game }}
    {{- $current_player := .player }}
    {{- range $index, $player := .Game.Players }}
        {{- if eq $player.Id $current_player.Id }}
            <div class="{{ $game.PlayerLabels $player.Id }} current col-1 item">
        {{- else }}
            <div class="{{ $game.PlayerLabels $player.Id }} col-1 item">
        {{- end }}
                <div class="title center">
                    {{ $player.User.Avatar.XS }}
                    <div class="name truncate">{{ $player.User.Name }}</div>
                </div>
                <div class="content center">
                    {{- if eq $player.Id $current_player.Id }}
                        {{ $game.YourPlayerMessage $lang $player.Id }}
                    {{- else }}
                        {{ $game.PlayerMessage $lang $player.Id }}
                    {{- end }}
                </div>
            </div>
    {{- end }}
    </div>
</div>
{{- end }}