        const addExerciseForm = document.getElementById("add-exercise-form")
        const addExercise = (e) => {
            e.preventDefault()
            PostForm("/api/exercises", addExerciseForm).then((res) => {
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
        }


        const dummyFunc = (e) => {
            console.log("event: ", e.detail)
        }

        const page = document.getElementById("exercise-page")
        document.addEventListener("custom-page-event", dummyFunc)

        swapContent("/api/partials/exercises", "exercise-list-container")

    const deleteExercise = (id) => {
        Delete(`/api/exercises/${id}`).then((res) => {
            swapContent("/api/partials/exercises", "exercise-list-container")
        }).catch((err) => {
            console.log("delete err: ", err)
        })
    }

    const editExercise = (id) => {
        const itemRow = document.getElementById(`item-row-${id}`)
        const itemName = document.getElementById(`item-name-${id}`)
        const itemNameEdit = document.getElementById(`item-name-edit-${id}`)
        const itemButtons = document.getElementById(`item-buttons-${id}`)
        const itemEditButtons = document.getElementById(`item-edit-buttons-${id}`)

        itemRow.classList.add("highlighted")

        itemName.classList.remove("show")
        itemName.classList.add("hide")

        itemNameEdit.classList.remove("hide")
        itemNameEdit.classList.add("show")

        itemButtons.classList.remove("show")
        itemButtons.classList.add("hide")

        itemEditButtons.classList.remove("hide")
        itemEditButtons.classList.add("show")
    }

    const editExerciseCancel = (id) => {
        const itemRow = document.getElementById(`item-row-${id}`)
        const itemName = document.getElementById(`item-name-${id}`)
        const itemNameEdit = document.getElementById(`item-name-edit-${id}`)
        const itemButtons = document.getElementById(`item-buttons-${id}`)
        const itemEditButtons = document.getElementById(`item-edit-buttons-${id}`)

        itemRow.classList.remove("highlighted")

        itemName.classList.remove("hide")
        itemName.classList.add("show")

        itemNameEdit.classList.remove("show")
        itemNameEdit.classList.add("hide")

        itemButtons.classList.remove("hide")
        itemButtons.classList.add("show")

        itemEditButtons.classList.remove("show")
        itemEditButtons.classList.add("hide")
    }

    const patchListener = (id) => {
        const editExerciseForm = document.getElementById(`edit-exercise-form-${id}`)
        PatchForm("/api/exercises", id, editExerciseForm).then((res) => {
            //update list
            swapContent("/api/partials/exercises", "exercise-list-container")
        }).catch((err) => {
            console.log("post err", err)
        })
    }


