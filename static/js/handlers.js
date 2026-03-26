export const toggleWorkoutHandler = (id) => {
    const collapse = document.getElementById(`collapse-${id}`)
    collapse.classList.toggle("show")

    const addMovement = document.getElementById(`add-movement-${id}`)
    const delMovement = document.getElementById(`delete-movement-${id}`)

    if (collapse.classList.contains("show")) {
        addMovement.disabled = false;
        delMovement.disabled = false;
    } else {
        addMovement.disabled = true;
        delMovement.disabled = true;
    }
}

export const toggleExerciseHandler = (show) => {
    const collapse = document.getElementById("exercise-collapse")
    if (!collapse) return;

    if (typeof show !== "boolean") {
        collapse.classList.toggle("show")
        return
    }

    show ? collapse.classList.add("show") : collapse.classList.remove("show")
}

const toggleExercises = (e) => {
    const {show} = e.detail
    toggleExerciseHandler(show)
}

document.addEventListener("toggle-exercises", toggleExercises)
