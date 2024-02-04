
function clearSession(){
    sessionStorage.clear();
    console.log("cleared session storage")
}

function sendEvent(event_name){
    window.angler.update(event_name);
}

function randomizeHash(){
    window.location.hash = Math.random().toString(36).substring(7);
}
