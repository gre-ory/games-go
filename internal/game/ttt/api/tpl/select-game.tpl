{{- define "select-game" }}
{{- $lang := .Lang }}
    <div id="content" hx-swap-oob="innerHTML">
        {{- if .HasWaitingUsers }}
        <div class="cols-1">
            <div class="new-game col-1 item">
                <div class="title center">{{ $lang.Loc "Lobby" }}</div>
                <div class="content left">
                    {{- range .WaitingUsers }}
                    <div class="badge user player waiting">
                        {{ .User.Avatar.XS }}
                        <div class="name truncate">{{ .User.Name }}</div>
                    </div>
                    {{- end }}
                </div>
            </div>
        </div>
        {{- end }}
        {{- range .NewGames }}
        {{- $game := . }}
            <div class="cols-1">
                <div class="join-game col-1 item {{ .Labels }}"> 
                    <div class="title center">{{ $lang.Loc "GameTitle" .Id }}</div>
                    <div class="content">
                        <div class="left">
                        {{- range .Players }}
                            <div class="badge user {{ $game.PlayerLabels .Id }}">
                                {{ .User.Avatar.XS }}
                                <div class="name truncate">{{ .User.Name }}</div>
                            </div>
                        {{- end }}
                        </div>
                        <div class="right">
                            <button ws-send data-action="join-game" data-game="{{ .Id }}">
                                {{ $lang.Loc "JoinAction" }}
                            </button>
                        </div>
                    </div>
                </div>
            </div>
        {{- end }}
        <div class="cols-1">
            <div class="new-game col-1 item">
                <div class="title center">{{ $lang.Loc "NewGame" }}</div>
                <div class="content right">
                    <button ws-send data-action="create-game">
                        {{ $lang.Loc "CreateAction" }}
                    </button>
                </div>
            </div>
        </div>
        {{- range .OtherGames }}
        {{- $game := . }}
            <div class="cols-1">
                <div class="join-game col-1 item {{ .Labels }}"> 
                    <div class="title center">{{ $lang.Loc "GameTitle" .Id }}</div>
                    <div class="content">
                        <div class="left">
                        {{- range .Players }}
                            <div class="badge user {{ $game.PlayerLabels .Id }}">
                                {{ .User.Avatar.XS }}
                                <div class="name truncate">{{ .User.Name }}</div>
                            </div>
                        {{- end }}
                        </div>
                    </div>
                </div>
            </div>
        {{- end }}
    </div>
{{- end }}