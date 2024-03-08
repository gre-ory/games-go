{{define "game-list"}}

    <div id="header">
        <b>Tic-Tac-Toe</b>
    </div>

    {{ if .player.Name }}
        Welcome {{ .player.Name }}!
    {{ else }}
        Welcome {{ .player.Id }}!
    {{ end }}

    <div id="content">

        {{ range .games }}
            <button 
                class="join-game" 
                hx-put="/htmx/tictactoe/join-game/{{ .Id }}"
                hx-target="#main"
                hx-swap="innerHTML">
                Join game {{ .Id }}...
            </button>
        {{ end }}

        <button 
            class="new-game"
            hx-put="/htmx/tictactoe/create-game"
            hx-target="#main"
            hx-swap="innerHTML">
            New game...
        </button>
        
    </div>

    <div id="footer">
    </div>

{{end}}