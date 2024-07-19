{{- define "deck" }}
<div class="{{ .Deck.Labels }}">
    <div class="size">{{ .Deck.Size }}</div>
</div>
{{- end }}