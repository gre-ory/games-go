{{define "game-layout"}}
    <div id="content" class="htmx-content" hx-swap-oob="innerHTML">
        Welcome {{ .player.Name }}!
        <br/>
        <div id="players">
            <span class="loading"></span>
        </div>
        <div id="board">
            <span class="loading"></span>
        </div>
        <div id="messages">
            <span class="loading"></span>
        </div>
    </div>
{{end}}

