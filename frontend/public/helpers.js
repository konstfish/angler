
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

function testAuth(){
    var response = fetch('/api/backend/v1/alive', {
        method: "POST",
        headers: {
            "Authorization": "Bearer " + localStorage.getItem("token")
        }
    })

    console.log("testAuth", response)
}