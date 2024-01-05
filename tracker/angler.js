
class Angler {
    constructor(domain, target) {
        this.domain = domain
        this.target = target
        
        console.log(this.domain, '->', this.target)

        window.addEventListener('locationchange', this.push)
        window.addEventListener('hashchange', this.push)

        this.push()
    }

    isMobile(){
        return Number(Math.min(window.screen.width, window.screen.height) < 768 || navigator.userAgent.indexOf("Mobi") > -1)
    }

    getInit() {
        return {
            ua: window.navigator.userAgent,
            to: window.performance.timeOrigin,
            rf: document.referrer,
            dt: this.isMobile(),
        }
    }

    getState(event){
        return {
            loc: {
                pt: document.location.pathname,
                hs: document.location.hash,
                pro: document.location.protocol
            },
            ev: event
        }
    }

    async push() {
        if(!sessionStorage.getItem("angler_key")){
            var data = await this.register()
        }else{
            var data = await this.update()
        }

        console.log(data)
    }

    async register(){
        var data = this.getInit()

        console.log(data)

        var response = await fetch(this.target + "/session/", { //+ this.domain)
            method: "POST",
            body: JSON.stringify(data),
            headers: {
                "Content-Type": "application/json"
            }
        })

        if(response.ok){
            var data = await response.json()
            this.setSession(data.InsertedID)

            this.update()
        }

        return data
    }

    async update(){
        var data = this.getState()

        // add to data; event type etc

        console.log(data)

        var response = await fetch(this.target + "/event/" + this.getSession(), { //+ this.domain + "/session/" + this.getSession() + 
            method: "POST",
            body: JSON.stringify(data),
            headers: {
                "Content-Type": "application/json"
            }
        })

        return response.json()
    }

    getSession() {
        return sessionStorage.getItem("angler_key");
    }

    setSession(key) {
        sessionStorage.setItem("angler_key", key);
    }

}

var domain = document.currentScript.getAttribute("domain")
// rewrite so target is fetched from script source
var target = document.currentScript.getAttribute("target")  || "http://localhost:8084/v1" // || "https://angler.konst.fish/ingress"

const angler = new Angler(domain, target)
