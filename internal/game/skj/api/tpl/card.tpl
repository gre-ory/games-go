
{{- define "card-front" }}
    <div class="{{ .Labels }}">
        <div class="card-front">{{ .Value }}</div>
    </div>
{{- end }}

{{- define "card-front-select" }}
    <div class="{{ .Labels }} select">
        <div class="card-front">{{ .Value }}</div>
    </div>
{{- end }}

{{- define "card-front-selected" }}
    <div class="{{ .Labels }} selected">
        <div class="card-front">{{ .Value }}</div>
    </div>
{{- end }}

{{- define "card-front-drop" }}
    <div class="{{ .Labels }} drop">
        <div class="card-front">{{ .Value }}</div>
    </div>
{{- end }}

{{- define "card-front-dropped" }}
    <div class="{{ .Labels }} dropped">
        <div class="card-front">{{ .Value }}</div>
    </div>
{{- end }}

{{- define "card-front-off" }}
    <div class="{{ .Labels }} off">
        <div class="card-front">{{ .Value }}</div>
    </div>
{{- end }}

{{- define "card-back" }}
    <div class="card">
        <div class="card-back">&nbsp;</div>
    </div>
{{- end }}

{{- define "card-back-select" }}
    <div class="card select">
        <div class="card-back">&nbsp;</div>
    </div>
{{- end }}

{{- define "card-back-selected" }}
    <div class="card selected">
        <div class="card-back">&nbsp;</div>
    </div>
{{- end }}

{{- define "card-back-drop" }}
    <div class="card drop">
        <div class="card-back">&nbsp;</div>
    </div>
{{- end }}

{{- define "card-back-dropped" }}
    <div class="card dropped">
        <div class="card-back">&nbsp;</div>
    </div>
{{- end }}

{{- define "card-back-off" }}
    <div class="card off">
        <div class="card-back">&nbsp;</div>
    </div>
{{- end }}

{{- define "card-spot" }}
    <div class="card-spot">&nbsp;</div>
{{- end }}

{{- define "card-spot-off" }}
    <div class="card-spot off">&nbsp;</div>
{{- end }}

{{- define "empty-deck" }}
    <div class="deck">
        <div class="nb-card">0</div>
        <div class="card-spot">&nbsp;</div>
    </div>
{{- end }}

{{- define "empty-deck-off" }}
    <div class="deck">
        <div class="nb-card">0</div>
        <div class="card-spot off">&nbsp;</div>
    </div>
{{- end }}
