        const addExerciseForm = document.getElementById("add-exercise-form")
        const addExercise = (e) => {
            e.preventDefault()
            //const data = new FormData(e.target)
            //console.log("e: ", ...data.entries())
            //console.log(e.target.name)
            PostForm("/api/exercises", addExerciseForm).then((res) => {
                console.log("res: ", res)
                //update list
                swapContent("/api/partials/exercises", "exercise-list-container")
            }).catch((err) => {
                console.log("post err", err)
            })
        }

        addExerciseForm.addEventListener("submit", addExercise)

        const clickChange = (id) => {
            //working below with document. Sooo, do you need to get the
            //page element to dispatch event working. I though pushing it
            //on to document would bubble to this div
            const event = new CustomEvent('custom-page-event', {
                detail: "hello",
                bubbles: true,
                cancelable: true,
            })
            document.dispatchEvent(event)

            console.log("click id: ", id)
            //swapContent("/api/partials/widgets/b", "widget-list-container")
        }

        const deleteExercise = (id) => {
            console.log("id: ", id)
            Delete(`/api/exercises/${id}`).then((res) => {
                console.log("delete ok: ", res)
            }).catch((err) => {
                console.log("delete err: ", err)
            })
        }

        const dummyFunc = (e) => {
            console.log("event: ", e.detail)
        }

        const page = document.getElementById("exercise-page")
        document.addEventListener("custom-page-event", dummyFunc)

        swapContent("/api/partials/exercises", "exercise-list-container")

    const deleteTheExercise = (id) => {
        console.log("id of delettion: ", id)
        /*
        Delete(`/api/exercises/${id}`).then((res) => {
            console.log("delete ok: ", res)
        }).catch((err) => {
            console.log("delete err: ", err)
        })
        */
    alert("omg")
    }

