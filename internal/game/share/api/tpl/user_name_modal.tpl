{{- define "user-name-modal" }}
<div id="user-name-modal" class="modal" _="on closeModal add .closing then wait for animationend then remove me">
	<div class="modal-underlay" _="on click trigger closeModal"></div>
	<div class="modal-content">
        <form hx-put="/htmx/user" hx-target="#user" hx-swap="outerHTML">
            <div class="left">
                <label>Name</label>
                <input type="text" name="user_name" value="{{ or .user.Name .user.Id }}">
            </div>
            <div class="right">
                <button _="on click trigger closeModal">Submit</button>
                <button hx-get="/htmx/user" hx-target="#user" hx-swap="outerHTML" _="on click trigger closeModal">Cancel</button>
            </div>
        </form>
	</div>
</div>
{{- end }}