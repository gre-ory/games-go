

const randomFace = () => {
    return Math.floor(Math.random() * 6) + 1
}

const rollDiceToFace = ( dice, face ) => {

    dice.style.transition = 'none';
    dice.style.animation = 'none';
    dice.style.transform = `rotateX(-720deg) rotateY(-720deg)`;

    setTimeout(() => {

        dice.style.transition = '2s ease';
        switch (face) {
            case 1:
                dice.style.transform = 'rotateX(720deg) rotateY(720deg)';
                break;

            case 6:
                dice.style.transform = 'rotateX(900deg) rotateY(720deg)';
                break;

            case 2:
                dice.style.transform = 'rotateX(630deg) rotateY(720deg)';
                break;

            case 5:
                dice.style.transform = 'rotateX(810deg) rotateY(720deg)';
                break;

            case 3:
                dice.style.transform = 'rotateX(720deg) rotateY(810deg)';
                break;

            case 4:
                dice.style.transform = 'rotateX(720deg) rotateY(630deg)';
                break;

            default:
                break;
        }

    }, 50);

}

const rollDiceToFace2 = ( dice, face ) => {

    const extraX = Math.floor(Math.random() * 200) + 300
    const extraY = Math.floor(Math.random() * 200) + 300

    dice.style.transform = `rotateX(${extraX}deg) rotateY(${extraY}deg)`;

    setTimeout(() => {

        switch (face) {
            case 1:
                dice.style.transform = 'rotateX(0deg) rotateY(0deg)';
                break;

            case 6:
                dice.style.transform = 'rotateX(180deg) rotateY(0deg)';
                break;

            case 2:
                dice.style.transform = 'rotateX(-90deg) rotateY(0deg)';
                break;

            case 5:
                dice.style.transform = 'rotateX(90deg) rotateY(0deg)';
                break;

            case 3:
                dice.style.transform = 'rotateX(0deg) rotateY(90deg)';
                break;

            case 4:
                dice.style.transform = 'rotateX(0deg) rotateY(-90deg)';
                break;

            default:
                break;
        }

        dice.style.animation = 'none';

    }, 2050);

}

const onClick = ( eltId, callback ) => {
    document.getElementById( eltId ).addEventListener( 'click', callback );
}

const rollOnClick = ( eltId, diceId, face ) => {
    const elt = document.getElementById( eltId );
    const dice = document.getElementById( diceId );
    if ( elt && dice ) {
        console.log( `[roll] click event set on #${eltId} for dice #${diceId} ( face ${face} )`, elt, dice )
        elt.addEventListener('click', ( event ) => {
            console.log( `[roll] elt ${eltId} >>> roll dice ${diceId} to face ${face}`, event )
            rollDiceToFace( dice, face )
        });
    } else {
        console.log( `[roll] missing #${eltId} or #${diceId}!`, elt, dice )
    }
}

const fillDice = ( diceId ) => {
    const diceElt = document.getElementById( diceId );
    if ( diceElt ) {
        for ( let face = 1; face <= 6 ; face++ ) {
            let faceElt = document.createElement( 'div' );
            faceElt.className = `face-${face}`
            for ( let dot = 1; dot <= face ; dot++ ) {
                let dotElt = document.createElement( 'div' );
                dotElt.className = `dot-${dot}`
                faceElt.appendChild( dotElt )
            }
            diceElt.appendChild( faceElt )
        }
    } else {
        console.log( `[fill] missing dice #${diceId}!` )
    }
}

const fillDices = () => {
    document.querySelectorAll(".dice.d6").forEach( ( d6 ) => {
        fillD6( d6, 6 )
    } );
}

const fillD6 = ( diceElt ) => {
    for ( let face = 1; face <= 6 ; face++ ) {
        let faceElt = document.createElement( 'div' );
        faceElt.className = `face-${face}`
        for ( let dot = 1; dot <= face ; dot++ ) {
            let dotElt = document.createElement( 'div' );
            dotElt.className = `dot-${dot}`
            faceElt.appendChild( dotElt )
        }
        diceElt.appendChild( faceElt )    
    }
}

