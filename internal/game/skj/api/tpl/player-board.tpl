{{- define "player-board" }}
    {{- $player := . }}
    {{- $board := .Board }}
    <div class="other {{ $player.Labels }}">
        {{- if $board.IsReady }}
            <div class="score">{{ $board.Total }}</div>
            {{- range $rowIndex, $row := $board.Rows }}
                <div class="cols-{{ $board.NbColumn }}">
                    {{- range $columnIndex, $cell := $row }}
                        <div class="col-1 center {{ $cell.Labels }}">
                            {{- if $cell.IsFlipped }}
                                {{ template "card-front-off" $cell.Card }}
                            {{- else }}
                                {{ template "card-back-off" }}
                            {{- end }}
                        </div>
                    {{- end }}
                </div>
            {{- end }}
        {{- else }}
            &nbsp;
        {{- end }}
    </div>
{{- end }}