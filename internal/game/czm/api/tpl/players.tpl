{{- define "players" }}
{{- $lang := .lang }}
<div id="players" hx-swap-oob="outerHTML">
    <div class="players cols-2">
    {{- $game := .game }}
    {{- $current_id := .player.Id }}
    {{- range $id, $player := .game.Players }}
        {{- if eq $id $current_id }}
            <div class="{{ $game.PlayerLabels $id }} current col-1 item">
        {{- else }}
            <div class="{{ $game.PlayerLabels $id }} col-1 item">
        {{- end }}
                <div class="title center">
                    {{ $player.ExtraSmallAvatarHtml }}
                    <div class="name truncate">{{ or $player.Name $player.Id }}</div>
                </div>
                <div class="content center">
                    {{- if eq $id $current_id }}
                        {{ $game.YourPlayerMessage $lang $id }}
                    {{- else }}
                        {{ $game.PlayerMessage $lang $id }}
                    {{- end }}
                </div>
            </div>
    {{- end }}
    </div>
</div>
{{- end }}