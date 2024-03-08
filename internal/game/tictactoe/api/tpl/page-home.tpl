{{define "page-home"}}
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{ .title }}</title>
    <link rel="icon" type="image/png" href="/static/share/logo.png" />
    <!-- htmx -->
    <script src="https://unpkg.com/htmx.org@1.9.10"></script>
    <script src="https://unpkg.com/htmx.org/dist/ext/ws.js"></script>
    <script src="/static/tictactoe/game.js"></script>
    <!-- css -->
    <link rel="stylesheet" href="/static/share/pico.min.css"/>
    <link rel="stylesheet" href="/static/share/game.css"/>
    <link rel="stylesheet" href="/static/tictactoe/game.css"/>
</head>
<body>
    <div id="main" hx-ext="ws" ws-connect="/htmx/tictactoe/connect">
        <div id="header">{{ .title }} / {{ .cookie.PlayerId }}</div>
        <div id="content"></div>
        <div id="errors"></div>
        <div id="notifications"></div>
        <span id="ws-status"/>
    </div>
</body>
</html>
{{end}}