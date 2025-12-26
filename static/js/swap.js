import { Get } from './request.js'

export const swapContent = (url, container, autoload) => {
    Get(url).then((res) => {
        const cont = document.getElementById(container)

        document.startViewTransition(() => {
            cont.innerHTML = res
            if (autoload) autoload()
        })
    }).catch((err) => {
        console.log("error fetching content", err)
    })
}
