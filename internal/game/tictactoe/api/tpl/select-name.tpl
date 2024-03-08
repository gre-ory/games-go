{{define "select-name"}}
    <div class="htmx-content" hx-swap-oob="innerHTML:#content">
        Welcome {{ .player.Id }}!
        <br/>
        Please select your name:
        <form>
            <input type="hidden" name="action" value="set-name"/>
            <input type="text" name="player_name" value="{{ .player.Name }}"/>
            <button ws-send type="submit">Submit</button>
        </form>
    </div>
{{end}}