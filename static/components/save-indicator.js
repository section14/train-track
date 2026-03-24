class SaveIndicator extends HTMLElement {
    constructor() {
        super();

        let template = document.getElementById("save-indicator-template");
        let templateContent = template.content;

        const shadowRoot = this.attachShadow({ mode: "open" });
        shadowRoot.appendChild(templateContent.cloneNode(true));
    }

    changeStatus(e) {
        const { saving } = e.detail;

        if (saving) {
            this.spinner.classList.add("show");
            this.spinner.classList.remove("hide");
            this.check.classList.add("hide");
            this.check.classList.remove("show");
        } else {
            this.spinner.classList.add("hide");
            this.spinner.classList.remove("show");
            this.check.classList.add("show");
            this.check.classList.remove("hide");
        }
    }

    connectedCallback() {
        this.spinner = this.shadowRoot.getElementById("spinner");
        this.check = this.shadowRoot.getElementById("check");
        this.addEventListener("change-status", this.changeStatus.bind(this));
    }

    disconnectedCallback() {
        this.removeEventListener("change-status", this.changeStatus.bind(this));
    }
}

customElements.define("save-indicator", SaveIndicator);
