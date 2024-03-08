{{define "select-game"}}
    <div class="htmx-content" hx-swap-oob="innerHTML:#content">
        Welcome {{ .player.Name }}!
        <br/>
        Please select game:
        {{ range .games }}
            <form>
                <input type="hidden" name="action" value="join-game"/>
                <input type="hidden" name="game_id" value="{{ .Id }}"/>
                <button 
                    class="join-game" 
                    ws-send>
                    Join game {{ .Id }}...
                </button>
            </form>
        {{ end }}

        <form>
            <input type="hidden" name="action" value="create-game"/>
            <button 
                class="new-game"
                ws-send>
                New game...
            </button>
        </form>
    </div>
{{end}}