
class Angler {
    constructor(domain, target) {
        this.domain = domain
        this.target = target
        
        console.log(this.domain, '->', this.target)

        this.register().then(data => {
            console.log(data)
            this.setSession(data)
        });
    }

    async register() {
        // todo; split these into initial information & per page/event

        // send over initial info
        var ua = window.navigator.userAgent
        console.log(ua)
        var loc = document.location
        // pathname, protocol, hash, 
        console.log(loc)
        // load times

        // meh
        // window.screen
        //var height = window.innerHeight
        //var width = window.innerWidth

        // timings
        var start = window.performance.timeOrigin
        console.log(start)
        // defer this for time spent on pages

        // load time
        console.log(window.performance.now())

        let response = await fetch(this.target + "/get/" + this.domain)
        let data = await response.json()
        return data
    }

    setSession(data) {
        sessionStorage.setItem("anglerKey", data.key);
    }

}

var domain = document.currentScript.getAttribute("domain")
// rewrite so target is fetched from script source
var target = document.currentScript.getAttribute("target")  || "http://localhost:8084/ping" // || "https://angler.konst.fish/ingress"

const angler = new Angler(domain, target)
