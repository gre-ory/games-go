{{define "players"}}
    <div id="players" class="htmx-content" hx-swap-oob="innerHTML">
        <span class="{{ .player.Labels }} current"> {{ .player.Name }}</span>
        {{ $current_id := .player.GetId }}
        {{ range $id, $player := .game.Players }}
            {{ if ne $id $current_id }}
                <span class="{{ $player.Labels }}"> {{ $player.Name }}</span>
            {{ end }}
        {{ end }}
    </div>
{{end}}