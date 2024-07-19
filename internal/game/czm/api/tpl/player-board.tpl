{{- define "player-board" }}
{{- $lang := .Lang }}
<div id="player-board" hx-swap-oob="outerHTML">
    <div class="cols-4">
    {{- $game := .Game }}
    {{- $player := .Player }}
    {{- range $y, $card := $player.Cards }}
    {{- $selectable := $game.IsCardSelectable $player.Id $y }}
    {{- $selected := $game.IsCardSelected $player.Id $y }}
    {{ template "card" dict "Card" $card "Selectable" $selectable "Selected" $selected "CardIndex" $y }}
    {{- end }}
    </div>
</div>
{{- end }}