const getHTML = async (url) => {
    try {
        const response = await fetch(url)
        if (!response.ok) {
            throw new Error(`Error getting HTML: ${response.status}`)
        }

        const res = await response.text()

        console.log("res: ", res)
    } catch (err) {
        console.log("error caught getting HTML: ", err)
    }
}

const Get = (url) => {
    return new Promise((resolve, reject) => {
        fetch(url).then((response) => {
            if (!response.ok) {
                throw new Error(`Error getting HTML: ${response.status}`)
            }

            return resolve(response.text())
        }).catch((err) => {
            reject(err)
        })
    })
}

const PostForm = (url, data) => {
    return new Promise((resolve, reject) => {
        fetch(url, { method: "post", body: new FormData(data) }).then((response) => {
            if (!response.ok) {
                throw new Error(`Error POST'ing: ${response.status}`)
            }

            return resolve(response.text())
        }).catch((err) => {
            reject(err)
        })
    })
}

const Delete = (url) => {
    return new Promise((resolve, reject) => {
        fetch(url, {method: "delete"}).then((response) => {
            if (!response.ok) {
                throw new Error(`Error deleting: ${response.status}`)
            }

            return resolve(response.text())
        }).catch((err) => {
            reject(err)
        })
    })
}

const swapContent = (url, container) => {
    Get(url).then((res) => {
        const cont = document.getElementById(container)

        document.startViewTransition(() => {
            cont.innerHTML = res
        })
    }).catch((err) => {
        console.log("error fetching widgets", err)
    })
}
