{{- define "my-board" }}
    {{- $game := .Game }}
    {{- $player := .Player }}
    {{- $playing := $player.IsPlaying }}
    {{- $board := $player.Board }}
    <div class="my {{ $player.Labels }}">
        {{- if $board.IsReady }}
            <div class="score">{{ $board.Total }}</div>
            {{- range $rowIndex, $row := $board.Rows }}
                <div class="cols-{{ $board.NbColumn }}">
                    {{- range $columnIndex, $cell := $row }}
                        {{- if $game.CanPutCard $player $cell }}
                            <div class="col-1 center my {{ $cell.Labels }} can-put-card" ws-send data-action="put-card" data-column="{{ $cell.Column }}" data-row="{{ $cell.Row }}">
                            {{- if $cell.IsFlipped }}
                                {{ template "card-front-drop" $cell.Card }}
                            {{- else }}
                                {{ template "card-back-drop" }}
                            {{- end }}
                            </div>
                        {{- else if $game.CanFlipCard $player $cell }}
                            <div class="col-1 center my {{ $cell.Labels }} can-flip-card" ws-send data-action="flip-card" data-column="{{ $cell.Column }}" data-row="{{ $cell.Row }}">
                            {{- if $cell.IsFlipped }}
                                {{ template "card-front-select" $cell.Card }}
                            {{- else }}
                                {{ template "card-back-select" }}
                            {{- end }}
                            </div>
                        {{- else }}
                            <div class="col-1 center my {{ $cell.Labels }} off">
                                {{- if $cell.IsFlipped }}
                                    {{ template "card-front-off" $cell.Card }}
                                {{- else }}
                                    {{ template "card-back-off" }}
                                {{- end }}
                            </div>
                        {{- end }}
                    {{- end }}
                </div>
            {{- end }}
        {{- else }}
            &nbsp;
        {{- end }}
    </div>
{{- end }}