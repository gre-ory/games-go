{{define "board"}}
    <div id="board" class="htmx-content" hx-swap-oob="innerHTML">
        {{- $playing := .player.Playing }}
        <div class="{{ .player.Labels }}">
            {{- range $y, $row := .game.Rows }}
                <div class="row">
                    {{- range $x, $cell := $row.Cells }}
                        <div class="cell-container">
                        {{- if and $playing .IsEmpty }}
                            <form>
                                <input type="hidden" name="action" value="play"/>
                                <input type="hidden" name="play_x" value="{{ $x }}"/>
                                <input type="hidden" name="play_y" value="{{ $y }}"/>
                                <div ws-send class="{{ .Labels }}">&nbsp;</div>
                            </form>
                        {{- else if .IsEmpty }}
                            <div class="{{ .Labels }}">&nbsp;</div>
                        {{- else }}
                            <div class="{{ .Labels }}">{{ . }}</div>
                        {{- end }}
                        </div>
                    {{- end }}
                </div>
            {{- end }}
        </div>
    </div>
{{end}}