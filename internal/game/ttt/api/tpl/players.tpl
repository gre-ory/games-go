{{- define "players" }}
{{- $lang := .lang }}
<div id="players" hx-swap-oob="outerHTML">
    <div class="players cols-2">
    {{- $game := .game }}
    {{- $current_player := .player }}
    {{- range $index, $player := .game.Players }}
        {{- if eq $player.Id $current_player.Id }}
            <div class="{{ $game.PlayerLabels $player.Id }} current col-1 item">
        {{- else }}
            <div class="{{ $game.PlayerLabels $player.Id }} col-1 item">
        {{- end }}
                <div class="title center">
                    {{ $player.Avatar.ExtraSmallHtml }}
                    <div class="name truncate">{{ or $player.Name $player.Id }}</div>
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