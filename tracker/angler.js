
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

    getInit() {
        return {
            ua: window.navigator.userAgent,
            to: window.performance.timeOrigin,
            st: this.getState()
        }
    }

    getState(){
        return {
            loc: {
                path: document.location.pathname,
                hash: document.location.hash,
                protocol: document.location.protocol
            },
            timing: {
                now: window.performance.now()
            },
        }
    }

    async register() {
        var push = this.getInit()

        console.log(push)

        var response = await fetch(this.target + "/session", { //+ this.domain)
            method: "POST",
            body: JSON.stringify(push),
            headers: {
                "Content-Type": "application/json"
            }
        })

        return response.json()
    }

    setSession(data) {
        sessionStorage.setItem("anglerKey", data.key);
    }

}

var domain = document.currentScript.getAttribute("domain")
// rewrite so target is fetched from script source
var target = document.currentScript.getAttribute("target")  || "http://localhost:8084/v1" // || "https://angler.konst.fish/ingress"

const angler = new Angler(domain, target)
