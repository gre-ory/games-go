{{- define "board" }}
{{- if or .game.Started .game.Stopped }}
{{- $playing := .player.Playing }}
<div id="board" class="{{ .player.Labels }}" hx-swap-oob="outerHTML">
    <div class="board">
        {{- range $y, $row := .game.Rows }}
        <div class="row">
        {{- range $x, $cell := $row.Cells }}
        {{- if and $playing .IsEmpty }}
            <div class="cell {{ .Labels }} select" ws-send data-action="play" data-x="{{ $x }}" data-y="{{ $y }}">&nbsp;</div>
        {{- else if .IsEmpty }}
            <div class="cell {{ .Labels }}">&nbsp;</div>
        {{- else }}
            <div class="cell {{ .Labels }}">{{ .IconHtml }}</div>
        {{- end }}
        {{- end }}
        </div>
        {{- end }}
    </div>
    <div class="center">
        {{- if .game.Stopped }}
            <button ws-send data-action="create-game">New game</button>
        {{- end }}
        <button ws-send data-action="leave-game">Leave</button>
    </div>
</div>
{{- else }}
<div id="board" class="center" hx-swap-oob="outerHTML">
    {{- if .game.CanStart }}
    <button ws-send data-action="start-game">Start</button>
    {{- else }}
    <button class="off">Start</button>
    {{- end }}
    <button ws-send data-action="leave-game">Leave</button>
</div>
{{- end }}
{{- end }}