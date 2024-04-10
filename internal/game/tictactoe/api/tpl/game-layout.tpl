{{- define "game-layout" }}
    <div id="content" hx-swap-oob="innerHTML">
        <div id="players">
            <div class="center">
                <div class="loading-dot"><div></div><div></div><div></div><div></div></div>
            </div>
        </div>
        <div id="board">
            <div class="center">
                <div class="loading-dot"><div></div><div></div><div></div><div></div></div>
            </div>
        </div>
    </div>
{{- end }}

