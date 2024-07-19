{{- define "game-layout" }}
    <div id="content" hx-swap-oob="innerHTML">
        <div id="players">
            {{ .Share.LoadingDot }}
        </div>
        <div id="board">
            {{ .Share.LoadingDot }}
        </div>
        <div id="board-player">
            {{ .Share.LoadingDot }}
        </div>
    </div>
{{- end }}

