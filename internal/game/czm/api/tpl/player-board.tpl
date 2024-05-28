{{- define "player-board" }}
{{- $lang := .lang }}
<div id="player-board" hx-swap-oob="outerHTML">
    <div class="cols-4">
    {{- $game := .game }}
    {{- $player := .player }}
    {{- range $y, $card := $player.Cards }}
    {{- $selectable := $game.IsCardSelectable $player.Id $y }}
    {{- $selected := $game.IsCardSelected $player.Id $y }}
    {{ template "card" dict "Card" $card "Selectable" $selectable "Selected" $selected "CardIndex" $y }}
    {{- end }}
    </div>
</div>
{{- end }}