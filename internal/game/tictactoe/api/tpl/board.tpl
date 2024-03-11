{{- define "board" }}
    <div id="board" class="{{ .player.Labels }}" hx-swap-oob="outerHTML">
        {{- $playing := .player.Playing }}
        <div class="board">
            {{- range $y, $row := .game.Rows }}
                <div class="row">
                    {{- range $x, $cell := $row.Cells }}
                        <div class="cell-container">
                        {{- if and $playing .IsEmpty }}
                            <div class="{{ .Labels }}" 
                                ws-send
                                data-action="play"
                                data-x="{{ $x }}"
                                data-y="{{ $y }}"> </div>
                        {{- else if .IsEmpty }}
                            <div class="{{ .Labels }}"> </div>
                        {{- else }}
                            <div class="{{ .Labels }}">{{ . }}</div>
                        {{- end }}
                        </div>
                    {{- end }}
                </div>
            {{- end }}
        </div>
    </div>
{{- end }}