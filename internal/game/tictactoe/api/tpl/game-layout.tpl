{{- define "game-layout" }}
    <div id="content" class="htmx-content" hx-swap-oob="innerHTML">
        <div id="board">
            <span class="loading"></span>
        </div>
        <div id="players">
            <span class="loading"></span>
        </div>
    </div>
{{- end }}

