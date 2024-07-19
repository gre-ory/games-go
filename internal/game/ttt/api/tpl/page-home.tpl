{{define "page-home"}}
{{- $lang := .Lang }}
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{ $lang.Loc "Title" }}</title>
    <link rel="icon" type="image/png" href="/static/share/icons/dice-5.svg" />
    <!-- htmx -->
    <script src="https://unpkg.com/htmx.org@1.9.10"></script>
    <script src="https://unpkg.com/htmx.org/dist/ext/ws.js"></script>
    <script src="https://unpkg.com/hyperscript.org@0.9.12"></script>
    <script src="/static/share/ws.js" defer></script>
    <script src="/static/ttt/game.js" defer></script>
    <!-- css -->
    <link rel="stylesheet" href="/static/share/luciole.css"/>
    <link rel="stylesheet" href="/static/share/game.css"/>
    <link rel="stylesheet" href="/static/share/avatar.css"/>
    <link rel="stylesheet" href="/static/{{ .AppId }}/game.css"/>
</head>
<body>
    {{ .Share.UserBadge }}
    {{ .Share.WsStatusBadge }}

	<!-- websocket -->
    <div id="main" hx-ext="ws" ws-connect="{{ .ConnectUrl }}" hx-trigger="load">

	    <!-- header -->        
		<div id="header">
            <div class="title">{{ $lang.Loc "Title" }}</div>
        </div>
        
        <!-- content -->  
        <div id="content">
            {{ .Share.LoadingDot }}
        </div>
        
        <!-- notifications --> 
        <div id="notifications"></div>
        
    </div>
</body>
</html>
{{end}}