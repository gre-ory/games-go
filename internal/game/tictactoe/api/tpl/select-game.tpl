{{- define "select-game" }}
    <div id="content" hx-swap-oob="innerHTML">
        {{- range .games }}
            <div class="cols-1">
                <div class="join-game col-1 item"> 
                    <div class="title center">game {{ .Id }}</div>
                    <div class="content">
                        <div class="left">
                        {{- range $id, $player := .Players }}
                            <div class="badge user">
                                {{- if ne .Avatar 0 }}
                                    <div class="avatar-{{ .Avatar }} s"></div>
                                {{- end }}
                                <div class="name">{{ or .Name .Id }}</div>
                            </div>
                        {{- end }}
                        </div>
                        <div class="right">
                            <button ws-send data-action="join-game" data-game="{{ .Id }}">
                                Join
                            </button>
                        </div>
                    </div>
                </div>
            </div>
        {{- end }}
        <div class="cols-1">
            <div class="new-game col-1 item">
                <div class="title center">new game</div>
                <div class="content right">
                    <button ws-send data-action="create-game">
                        Create
                    </button>
                </div>
            </div>
        </div>
    </div>
{{- end }}