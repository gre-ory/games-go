
htmx.on( "htmx:configRequest", ( event ) => {
    console.log(`[htmx-config-request]`, event)
} )

htmx.on( "htmx:wsConfigSend", ( event ) => {
    console.log(`[htmx-ws-config-send]`, event)
} )

htmx.on( "htmx:wsConnecting", ( event ) => {
    console.log(`[htmx-ws-connecting]`, event)
    updateWsStatus( 'connecting' )
} )

htmx.on( "htmx:wsOpen", ( event ) => {
    console.log(`[htmx-ws-open]`, event)
    updateWsStatus( 'connected' )
} )

htmx.on( "htmx:wsClose", ( event ) => {
    console.log(`[htmx-ws-close]`, event)
    updateWsStatus( 'disconnected' )
} )

htmx.on( "htmx:wsError", ( event ) => {
    console.log(`[htmx-ws-error]`, event)
    updateWsStatus( 'error' )
} )

htmx.on( "htmx:wsBeforeMessage", ( event ) => {
    console.log(`[htmx-ws-before-message]`, event.detail.message)
} )

htmx.on( "htmx:wsAfterMessage", ( event ) => {
    console.log(`[htmx-ws-after-message]`, event)
} )

htmx.on( "htmx:wsBeforeSend", ( event ) => {
    console.log(`[htmx-ws-before-send]`, event)
} )

htmx.on( "htmx:wsAfterSend", ( event ) => {
    console.log(`[htmx-ws-after-send]`, event)
} )

function updateWsStatus( value ) {
    let status = document.getElementById('ws-status');
    status.innerText = value;
	status.setAttribute( 'data-status', value );
}