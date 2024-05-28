
// //////////////////////////////////////////////////
// ws helpers

updateWsStatus( 'off' )

registerDefaultWsHelpers()

function clearNotifications( event ) {
    let notifications = document.getElementById('notifications');
    notifications.innerText = "";
}

// htmx.on( "htmx:wsBeforeMessage", clearNotifications )

// rollOnClick( 'D1-1', 'D1', 1 )
// rollOnClick( 'D1-2', 'D1', 2 )
// rollOnClick( 'D1-3', 'D1', 3 )
// rollOnClick( 'D1-4', 'D1', 4 )
// rollOnClick( 'D1-5', 'D1', 5 )
// rollOnClick( 'D1-6', 'D1', 6 )

// rotateOnClick( 'D1-1', 'D1', 45, 45 )
// rotateOnClick( 'D1-2', 'D1', 40, 50 )
// rotateOnClick( 'D1-3', 'D1', 35, 55 )
// rotateOnClick( 'D1-4', 'D1', 30, 60 )

// fillDices()

// const d1 = document.getElementById( 'D1' )
// const d2 = document.getElementById( 'D2' )
// const d3 = document.getElementById( 'D3' )
// const d4 = document.getElementById( 'D4' )
// const d5 = document.getElementById( 'D5' )
// const d6 = document.getElementById( 'D6' )

// onClick( 'roll-D', () => {
//     rollDiceToFace( d1, randomFace() )
//     rollDiceToFace( d2, randomFace() )
//     rollDiceToFace( d3, randomFace() )
//     rollDiceToFace( d4, randomFace() )
//     rollDiceToFace( d5, randomFace() )
//     rollDiceToFace( d6, randomFace() )
// } )

// const p1 = document.getElementById( 'P1' )
// const p2 = document.getElementById( 'P2' )
// const p3 = document.getElementById( 'P3' )
// const p4 = document.getElementById( 'P4' )
// const p5 = document.getElementById( 'P5' )

// onClick( 'roll-P', () => {
//     rollDiceToFace( p1, randomFace() )
//     rollDiceToFace( p2, randomFace() )
//     rollDiceToFace( p3, randomFace() )
//     rollDiceToFace( p4, randomFace() )
//     rollDiceToFace( p5, randomFace() )
// } )

// const pm1 = document.getElementById( 'PM1' )
// const pm2 = document.getElementById( 'PM2' )
// const pm3 = document.getElementById( 'PM3' )
// const pm4 = document.getElementById( 'PM4' )
// const pm5 = document.getElementById( 'PM5' )
// const pm6 = document.getElementById( 'PM6' )
// const pm7 = document.getElementById( 'PM7' )
// const pm8 = document.getElementById( 'PM8' )

// onClick( 'roll-PM', () => {
//     rollDiceToFace( pm1, randomFace() )
//     rollDiceToFace( pm2, randomFace() )
//     rollDiceToFace( pm3, randomFace() )
//     rollDiceToFace( pm4, randomFace() )
//     rollDiceToFace( pm5, randomFace() )
//     rollDiceToFace( pm6, randomFace() )
//     rollDiceToFace( pm7, randomFace() )
//     rollDiceToFace( pm8, randomFace() )
// } )
