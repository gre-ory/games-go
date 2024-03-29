{{- define "players" }}
    <div id="players" hx-swap-oob="outerHTML">
        <div class="players cols-2">
        {{- $current_id := .player.Id }}
        {{- range $id, $player := .game.Players }}
            {{- if eq $id $current_id }}
                <div class="{{ $player.Labels }} current col-1 item">
            {{- else }}
                <div class="{{ $player.Labels }} col-1 item">
            {{- end }}
                    <div class="title center">
                        {{ $player.ExtraSmallAvatarHtml }}
                        <div class="name truncate">{{ or $player.Name $player.Id }}</div>
                    </div>
                    <div class="content center">
                        {{- if eq $id $current_id }}
                            {{ $player.YourMessage }}
                        {{- else }}
                            {{ $player.Message }}
                        {{- end }}
                    </div>
                </div>
        {{- end }}
        </div>
    </div>
{{- end }}