{{define "page-home"}}
{{- $lang := .lang }}
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{ $lang.Loc "Title" "TTT" }}</title>
    <link rel="icon" type="image/png" href="/static/share/icons/dice-5.svg" />
    <!-- htmx -->
    <script src="https://unpkg.com/htmx.org@1.9.10"></script>
    <script src="https://unpkg.com/htmx.org/dist/ext/ws.js"></script>
    <script src="https://unpkg.com/hyperscript.org@0.9.12"></script>
    <script src="/static/share/ws.js" defer></script>
    <script src="/static/share/dice.js" defer></script>
    <script src="/static/tictactoe/game.js" defer></script>
    <!-- css -->
    <link rel="stylesheet" href="/static/share/luciole.css"/>
    <link rel="stylesheet" href="/static/share/game.css"/>
    <link rel="stylesheet" href="/static/share/dice.css"/>
    <link rel="stylesheet" href="/static/share/avatar.css"/>
    <link rel="stylesheet" href="/static/tictactoe/game.css"/>
</head>
<body>
	
	<!-- user -->
    <div id="user" hx-get="/htmx/user" hx-target="this" hx-swap="outerHTML" hx-trigger="load">...</div>
	
	<!-- websocket status -->
	<div id="ws-status" class="cloud-on">
		<svg id="cloud-off" color="#777" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24"><title>cloud-off</title><path d="M19.8 22.6L17.15 20H6.5Q4.2 20 2.6 18.4T1 14.5Q1 12.58 2.19 11.08 3.38 9.57 5.25 9.15 5.33 8.95 5.4 8.76 5.5 8.57 5.55 8.35L1.4 4.2L2.8 2.8L21.2 21.2M21.6 18.75L8.05 5.23Q8.93 4.63 9.91 4.31 10.9 4 12 4 14.93 4 16.96 6.04 19 8.07 19 11 20.73 11.2 21.86 12.5 23 13.78 23 15.5 23 16.5 22.63 17.31 22.25 18.15 21.6 18.75Z"></path></svg>
		<svg id="cloud-connecting" color="#f2ae16" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24"><title>cloud-refresh-variant</title><path d="M21.86 12.5C21.1 11.63 20.15 11.13 19 11C19 9.05 18.32 7.4 16.96 6.04C15.6 4.68 13.95 4 12 4C10.42 4 9 4.47 7.75 5.43S5.67 7.62 5.25 9.15C4 9.43 2.96 10.08 2.17 11.1S1 13.28 1 14.58C1 16.09 1.54 17.38 2.61 18.43C3.69 19.5 5 20 6.5 20H18.5C19.75 20 20.81 19.56 21.69 18.69C22.56 17.81 23 16.75 23 15.5C23 14.35 22.62 13.35 21.86 12.5M16 13H12L13.77 11.23C13.32 10.78 12.69 10.5 12 10.5C10.62 10.5 9.5 11.62 9.5 13S10.62 15.5 12 15.5C12.82 15.5 13.54 15.11 14 14.5H15.71C15.12 15.97 13.68 17 12 17C9.79 17 8 15.21 8 13S9.79 9 12 9C13.11 9 14.11 9.45 14.83 10.17L16 9V13Z"></path></svg>
		<svg id="cloud-on" color="#00ad94" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24"><title>cloud-check-variant</title><path d="M10.35 17L16 11.35L14.55 9.9L10.33 14.13L8.23 12.03L6.8 13.45M6.5 20Q4.22 20 2.61 18.43 1 16.85 1 14.58 1 12.63 2.17 11.1 3.35 9.57 5.25 9.15 5.88 6.85 7.75 5.43 9.63 4 12 4 14.93 4 16.96 6.04 19 8.07 19 11 20.73 11.2 21.86 12.5 23 13.78 23 15.5 23 17.38 21.69 18.69 20.38 20 18.5 20Z"></path></svg>
		<svg id="cloud-error" fill="#b71348" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24"><title>cloud-alert</title><path d="M21.86 12.5C21.1 11.63 20.15 11.13 19 11C19 9.05 18.32 7.4 16.96 6.04C15.6 4.68 13.95 4 12 4C10.42 4 9 4.47 7.75 5.43S5.67 7.62 5.25 9.15C4 9.43 2.96 10.08 2.17 11.1S1 13.28 1 14.58C1 16.09 1.54 17.38 2.61 18.43C3.69 19.5 5 20 6.5 20H18.5C19.75 20 20.81 19.56 21.69 18.69C22.56 17.81 23 16.75 23 15.5C23 14.35 22.62 13.35 21.86 12.5M13 17H11V15H13V17M13 13H11V7H13V13Z"></path></svg>
	</div>
	
	<!-- websocket -->
    <div id="main" hx-ext="ws" ws-connect="/ttt/htmx/connect" hx-trigger="load">
	
	    <!-- header -->        
		<div id="header">
            <div class="title">{{ $lang.Loc "Title" "TTT" }}</div>
        </div>
        
        <!-- content -->  
        <div id="content">
            <div class="center">
                <div class="loading-dot"><div></div><div></div><div></div><div></div></div>
            </div>
        </div>
        
        <!-- notifications --> 
        <div id="notifications"></div>
        
    </div>
</body>
</html>
{{end}}