import { swapContent } from '/static/js/swap.js'

class PartialLoader extends HTMLElement {
    constructor() {
        super();
    }

    static get observedAttributes() {
        return ["target-id", "endpoint"];
    }

    targetId = null;
    endpoint = null;

    connectedCallback() {
        console.log("loader connected", this.targetId, " - ", this.endpoint);
        if (this.targetId && this.endpoint) this.loadContent();
    }

    disconnectedCallback() {
        console.log("loader disconnected");
    }

    loadContent() {
        swapContent(this.endpoint, this.targetId);
    }

    attributeChangedCallback(name, oldValue, newValue) {
        switch(name) {
            case "target-id":
                if (oldValue !== newValue) {
                    this.targetId = newValue
                }
                break;
            case "endpoint":
                if (oldValue !== newValue) {
                    this.endpoint = newValue
                }
                break;
        }
    }
};

customElements.define("partial-loader", PartialLoader)
