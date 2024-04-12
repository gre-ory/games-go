
// //////////////////////////////////////////////////
// ws helpers

updateWsStatus( 'off' )

registerDefaultWsHelpers()

function clearNotifications( event ) {
    let notifications = document.getElementById('notifications');
    notifications.innerText = "";
}
