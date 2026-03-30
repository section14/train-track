globalThis.docLoaded = false;

const divOnLoad = (e) => {
    console.log("div loaded", e);
};

// keep track of queued updates
//todo: I don't love this approach
const movementQueue = [];

const movementChanged = (e) => {
    const { id, saving } = e.detail;

    const addMovement = document.getElementById(`add-movement-${id}`);
    const delMovement = document.getElementById(`delete-movement-${id}`);

    if (saving) {
        movementQueue.push(true);
        addMovement.disabled = true;
        delMovement.disabled = true;
    } else {
        movementQueue.pop();
        if (movementQueue.length === 0) {
            addMovement.disabled = false;
            delMovement.disabled = false;
        }
    }
};

const mediaQuery = window.matchMedia("(width < 1100px)");

const handleSizeChange = (e) => {
    if (e.matches) {
        const button = document.getElementById("toggle-exercise-button");

        if (button) {
            button.style.setProperty("display", "flex");
        }

        const event = new CustomEvent("toggle-exercises", {
            detail: { show: false },
            cancelable: true,
            bubbles: true,
        });
        document.dispatchEvent(event);
    } else {
        const button = document.getElementById("toggle-exercise-button");

        if (button) {
            button.style.setProperty("display", "none");
        }

        const event = new CustomEvent("toggle-exercises", {
            detail: { show: true },
            cancelable: true,
            bubbles: true,
        });
        document.dispatchEvent(event);
    }
};

window.onload = function () {
    const path = window.location.pathname;

    if (path === "/workouts") {
        const btn = document.getElementById("nav-workouts");
        btn.classList.toggle("selected");
    } else {
        const btn = document.getElementById("nav-exercises");
        btn.classList.toggle("selected");
    }

    globalThis.docLoaded = true;
    handleSizeChange(mediaQuery);
};

mediaQuery.addEventListener("change", handleSizeChange);

//const el = document.querySelectorAll("div[data-div-load]");
//document.addEventListener("change", divOnLoad);
document.addEventListener("movement-changed", movementChanged);
