{{- define "game-layout" }}
    <div id="content" hx-swap-oob="innerHTML">
        <div id="players">
            {{ .share.LoadingDot }}
        </div>
        <div id="board">
            {{ .share.LoadingDot }}
        </div>
    </div>
{{- end }}

