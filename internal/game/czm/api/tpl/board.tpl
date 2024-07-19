{{- define "board" }}
{{- $lang := .Lang }}
{{- if .Game.WasStarted }}
{{- $playing := .Player.IsPlaying }}
<div id="board" class="{{ .Player.Labels }}" hx-swap-oob="outerHTML">
    <div class="board cols-5">
    {{ template "deck" dict "Deck" .Game.DrawDeck }}
    {{- range $y, $deck := .Game.DiscardDecks }}
    {{ template "card" dict "Card" $deck.TopCard }}
    {{- end }}
    </div>
    <div class="center">
        {{- if .Game.Stopped }}
            <button ws-send data-action="create-game">{{ $lang.Loc "NewGameAction" }}</button>
        {{- end }}
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