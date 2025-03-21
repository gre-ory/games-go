{{- define "board" }}
{{- $lang := .Lang }}
{{- if .Game.WasStarted }}
{{- $game := .Game }}
{{- $current_player := .Player }}
{{- $playing := .Player.IsPlaying }}
<div id="board" hx-swap-oob="outerHTML">
    <div class="board">

        <!-- players board -->

        <div class="other-players {{ $game.NbPlayer }}-players cols-{{ $game.NbOtherPlayer }}">
            {{- range $game.Players }}
                {{- if ne .Id $current_player.Id }}
                    <div class="other-player col-1">
                        {{ template "player-board" . }}
                    </div>
                {{- end }}
            {{- end }}
        </div>

        <!-- center board -->

        <div class="cols-4">

            <!-- draw deck -->

            <div id="draw-deck col-1">
                {{- if .Game.DrawDeck.IsEmpty }}
                    {{ template "empty-deck-off" }}
                {{- else }}
                    {{- if .Game.CanDrawCard .Player }}
                        <div class="deck" ws-send data-action="draw-card">
                        <div class="nb-card">{{ .Game.DrawDeck.Size }}</div>
                        {{ template "card-back-select" }}
                    </div>
                    {{- else }}
                        <div class="deck">
                            <div class="nb-card">{{ .Game.DrawDeck.Size }}</div>
                            {{ template "card-back-off" }}
                        </div>
                    {{- end }}
                {{- end }}
            </div>

            <!-- discard deck -->

            <div id="discard-deck col-1">
                {{- if .Game.DiscardDeck.IsEmpty }}
                    {{ template "empty-deck-off" }}
                {{- else }}
                    {{- $select := false }}        
                    {{- if .Game.CanDrawDiscardCard .Player }}
                        <div class="discard-deck deck" ws-send data-action="draw-discard-card">
                            <div class="nb-card">{{ .Game.DiscardDeck.Size }}</div>
                            {{ template "card-front-select" .Game.DiscardDeck.TopCard }}
                        </div>
                    {{- else if .Game.CanDiscardCard .Player }}
                        <div class="discard-deck deck" ws-send data-action="discard-card">
                            <div class="nb-card">{{ .Game.DiscardDeck.Size }}</div>
                            {{ template "card-front-drop" .Game.DiscardDeck.TopCard }}
                        </div>
                    {{- else }}
                        <div class="discard-deck deck">
                            <div class="nb-card">{{ .Game.DiscardDeck.Size }}</div>
                            {{ template "card-front-off" .Game.DiscardDeck.TopCard }}
                        </div>
                    {{- end }}
                            
                {{- end }}
            </div>

            <!-- space -->

            <div class="col-1">&nbsp;</div>

            <div id="selected-card col-1">
                {{- if .Game.HasSelectedCard }}
                    {{ template "card-front-off" .Game.SelectedCard }}
                {{- else }}
                    {{ template "card-spot-off" }}
                {{- end }}
            </div>

        </div>

        <!-- my board -->

        <div class="my-player-board">
            {{ template "my-board" dict "Game" $game "Player" $current_player }}
        </div>

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