{{- define "select-game" }}
    <div class="htmx-content" hx-swap-oob="innerHTML:#content">
        Welcome {{ .player.Name }}!
        <br/>
        Please select game:
        {{- range .games }}
            <button 
                class="join-game" 
                ws-send
                data-action="join-game"
                data-game="{{ .Id }}">
                Join game {{ .Id }}...
            </button>
        {{- end }}
        <button 
            class="new-game"
            ws-send
            data-action="create-game">
            New game...
        </button>
    </div>
{{- end }}