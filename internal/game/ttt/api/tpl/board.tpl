{{- define "board" }}
{{- $lang := .Lang }}
{{- if .Game.WasStarted }}
{{- $playing := .Player.IsPlaying }}
<div id="board" class="{{ .Player.Labels }}" hx-swap-oob="outerHTML">
    <div class="board">
        {{- range $y, $row := .Game.Rows }}
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
        <button ws-send data-action="leave-game">{{ $lang.Loc "LeaveAction" }}</button>
    </div>
</div>
{{- else }}
<div id="board" class="center" hx-swap-oob="outerHTML">
    {{- if .Game.CanStart }}
    <button ws-send data-action="start-game">{{ $lang.Loc "StartAction" }}</button>
    {{- else }}
    <button class="off">{{ $lang.Loc "StartAction" }}</button>
    {{- end }}
    <button ws-send data-action="leave-game">{{ $lang.Loc "LeaveAction" }}</button>
</div>
{{- end }}
{{- end }}