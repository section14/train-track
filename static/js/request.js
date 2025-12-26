export const getHTML = async (url) => {
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

export const Get = (url) => {
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

export const GetJson = (url) => {
    return new Promise((resolve, reject) => {
        fetch(url).then((response) => {
            if (!response.ok) {
                throw new Error(`Error getting HTML: ${response.status}`)
            }

            return resolve(response.json())
        }).catch((err) => {
            reject(err)
        })
    })
}

export const PatchForm = (url, id, data) => {
    return new Promise((resolve, reject) => {
        fetch(`${url}/${id}`, { method: "PATCH", body: new FormData(data) }).then((response) => {
            if (!response.ok) {
                throw new Error(`Error POST'ing: ${response.status}`)
            }

            return resolve(response.text())
        }).catch((err) => {
            reject(err)
        })
    })
}

export const Post = (url) => {
    return new Promise((resolve, reject) => {
        fetch(url, { method: "POST" }).then((response) => {
            if (!response.ok) {
                throw new Error(`Error POST'ing: ${response.status}`)
            }

            return resolve(response.text())
        }).catch((err) => {
            reject(err)
        })
    })
}

export const PostForm = (url, data) => {
    return new Promise((resolve, reject) => {
        fetch(url, { method: "POST", body: new FormData(data) }).then((response) => {
            if (!response.ok) {
                throw new Error(`Error POST'ing: ${response.status}`)
            }

            return resolve(response.text())
        }).catch((err) => {
            reject(err)
        })
    })
}

export const Delete = (url) => {
    return new Promise((resolve, reject) => {
        fetch(url, { method: "DELETE" }).then((response) => {
            if (!response.ok) {
                throw new Error(`Error deleting: ${response.status}`)
            }

            return resolve(response.text())
        }).catch((err) => {
            reject(err)
        })
    })
}
