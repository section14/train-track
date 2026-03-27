class HintButton extends HTMLElement {
    constructor() {
        super();

        let template = document.getElementById("hint-button-template");
        let templateContent = template.content;

        const shadowRoot = this.attachShadow({ mode: "open" });
        shadowRoot.appendChild(templateContent.cloneNode(true));
    }

    mouseover() {
        this.hint.showPopover()
    }

    mouseout() {
        this.hint.hidePopover()
    }

    connectedCallback() {
        this.btn = this.shadowRoot.querySelector('slot[name="button"]');

        // todo: easier way to achieve this? Seems over-the-top
        const slotHint = this.shadowRoot.querySelector('slot[name="hint"]');
        this.hint = slotHint.assignedElements()[0].querySelector("#collapse-hint")

        this.btn.addEventListener("mouseover", this.mouseover.bind(this))
        this.btn.addEventListener("mouseout", this.mouseout.bind(this))
    }


}

customElements.define("hint-button", HintButton);
