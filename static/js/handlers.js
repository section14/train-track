export const toggleHandler = (id) => {
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
