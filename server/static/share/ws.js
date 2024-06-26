
// //////////////////////////////////////////////////
// ws-status helpers

function updateWsStatus( value ) {
    let status = document.getElementById('ws-status');
    if ( status ) {
        status.className = `cloud-${value}`;
    } else {
        console.log(`element #ws-status NOT found!`)
    }
}

// //////////////////////////////////////////////////
// attach data helpers

const dataPrefix = 'data-'

function attachDataToRequest( event ) {
    if ( event.detail.elt.hasAttributes() ) {
        modified = false
        for ( const attr of event.detail.elt.attributes ) {
            if ( attr.name.startsWith( dataPrefix ) ) {
                parameter = attr.name.substring(dataPrefix.length)
                event.detail.parameters[parameter] = attr.value
                modified = true
            }
        }
        if ( modified ) {
            console.log(`parameters attached to ws request`, event.detail.parameters )
        }
    }
}

// //////////////////////////////////////////////////
// events

function defaultOnWsConnecting( event ) {
    updateWsStatus( 'connecting' )
}

function defaultOnWsOpen( event ) {
    updateWsStatus( 'on' )
}

function defaultOnWsClose( event ) {
    updateWsStatus( 'off' )
}

function defaultOnWsError( event ) {
    updateWsStatus( 'error' )
}

function defaultOnWsConfigSend( event ) {
    attachDataToRequest( event )
}

function registerDefaultWsHelpers() {
    htmx.on( "htmx:wsConnecting", defaultOnWsConnecting )
    htmx.on( "htmx:wsOpen", defaultOnWsOpen )
    htmx.on( "htmx:wsClose", defaultOnWsClose )
    htmx.on( "htmx:wsError", defaultOnWsError )
    htmx.on( "htmx:wsConfigSend", defaultOnWsConfigSend )
}