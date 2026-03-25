globalThis.docLoaded = false;

window.onload = function() {
    globalThis.docLoaded = true;
}

const divOnLoad = (e) => {
    console.log("div loaded", e);
}

const movementChanged = (e) => {
    const { id, saving } = e.detail

    const addMovement = document.getElementById(`add-movement-${id}`)
    const delMovement = document.getElementById(`delete-movement-${id}`)

    if (saving) {
        addMovement.disabled = true;
        delMovement.disabled = true;
    } else {
        addMovement.disabled = false;
        delMovement.disabled = false;
    }
}

//const el = document.querySelectorAll("div[data-div-load]");
//document.addEventListener("change", divOnLoad);
document.addEventListener("movement-changed", movementChanged)
