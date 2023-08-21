class Blitzbase {
    constructor(url) {
        this.url = url
    }

    collection(name) {
        this.collection = name
        return this
    }

    subscribe(action, callback) {
        const sse = new EventSource(`${this.url}/api/realtime`)

        sse.addEventListener(action, callback)


        sse.addEventListener("error", (e) => {
            console.log("Error: ", e)
            sse.close()
        })
    }
}

function main() {

    const elem = document.querySelector("#msg")


    const bb = new Blitzbase("https://blitzbase.onrender.com")
    bb.collection("users").subscribe("create", (e) => {
        const data = JSON.parse(e.data)
        elem.innerHTML = JSON.stringify(data)
        console.log(data)
    })

    // const sse = new EventSource("http://127.0.0.1:3300/realtime/")
    //
    // sse.addEventListener("error", (e) => {
    //     console.log("Error: ", e)
    //     sse.close()
    // })
    //
    // sse.addEventListener("create", (e) => {
    //     const data = JSON.parse(e.data)
    //     elem.innerHTML = JSON.stringify(data)
    // })
    //
    //
    // sse.addEventListener("message", (e) => {
    //     const data = JSON.parse(e.data)
    //     elem.innerHTML = data.data
    // })

}

async function createUser(name) {
    try {
        const options = {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: `{"username":"${name}","email":"${name}@gmail.com","password":"fjYfIf6Ou4dJyKC5"}`
        };

        fetch('https://blitzbase.onrender.com/api/auth/register', options)
            .then(response => response.json())
            .then(response => console.log(response))
            .catch(err => console.error(err));
    } catch (e) {
        console.log(e)
    }
}

window.onload = () => {
    const button = document.querySelector("button")
    const input = document.querySelector("input[name=name]")
    let name = ""

    input.addEventListener("keypress", (e) => {
        name = e.target.value
        if (e.key == 'Enter') {
            button.click()
        }
    })

    button.addEventListener("click", async () => {
        resp = await createUser(name)
        console.log(resp)
    })

    main()
}

