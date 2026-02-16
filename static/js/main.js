globalThis.docLoaded = false;

window.onload = function() {
    console.log("loaded");
    globalThis.docLoaded = true;
}

const divOnLoad = (e) => {
    console.log("div loaded", e);
}

const el = document.querySelectorAll("div[data-div-load]");
document.addEventListener("change", divOnLoad);
