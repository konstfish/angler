
class Angler {
    constructor(domain, target) {
        this.domain = domain
        this.target = target
        
        console.log(this.domain, '->', this.target)

        window.addEventListener('locationchange', () => this.update('location'));
        window.addEventListener('hashchange', () => this.update('hash'), false);

        // to make spa navigation events work
        var c
        window.history.pushState && (c = window.history.pushState, window.history.pushState = function() {
            c.apply(this, arguments),
            window.angler.update('location')
        });
        
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
            },
            ev: event
        }
    }

    async push() {
        if(!sessionStorage.getItem("angler_key")){
            var res = await this.register()
        }else{
            var res = await this.update("unknown")
        }

        console.log(res)
    }

    async register(){
        var data = this.getInit()

        console.log(data)

        var response = await fetch(this.target + "/session/" + this.domain, {
            method: "POST",
            body: JSON.stringify(data),
            headers: {
                "Content-Type": "application/json"
            }
        })

        if(response.ok){
            var res_body = await response.json()
            console.log(res_body)
            this.setSession(res_body.id)

            this.update("init")

            return res_body
        }
    }

    async update(event){
        // add to data; event type etc
        if(!this.getSession()){
            return {}
        }

        if(event == "location" || event == "hash" || event == "unknown"){
            var event = this.getNavigationType()
        }
        
        var data = this.getState(event)

        console.log(data)

        var response = await fetch(this.target + "/event/" + this.domain + "/session/" + this.getSession(), {
            method: "POST",
            body: JSON.stringify(data),
            headers: {
                "Content-Type": "application/json"
            }
        })

        return response.status
    }

    getSession() {
        return sessionStorage.getItem("angler_key");
    }

    setSession(key) {
        sessionStorage.setItem("angler_key", key);
    }

    getNavigationType() {
        let performanceEntries = performance.getEntriesByType("navigation");
        
        if (performanceEntries && performanceEntries.length > 0) {
            return performanceEntries[0].type;
        }
    
        return 'unknown';
    }
}

var domain = document.currentScript.getAttribute("domain")
var target = document.currentScript.getAttribute("target")  || "https://angler.konst.fish/ingress/v1"

window.angler = new Angler(domain, target)
