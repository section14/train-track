class FormInput extends HTMLElement {
    constructor() {
        super();

        /**
         * @returns {HTMLTemplateElement}
         */

        let template = document.getElementById("form-input-template")
        console.log("template: ", template)
        let templateContent = template.content

        const shadowRoot = this.attachShadow({mode: "open"})
        shadowRoot.appendChild(templateContent.cloneNode(true))
    }

    connectedCallback() {
        console.log("I've been a connect")

    }

    static get observedAttributes() {
        return ["value", "label", "placeholder", "error-message"];
    }

    valueChanged(value) {

    }

    labelChanged(label) {

    }

    placeholderChanged(placeholder) {

    }

    errorChanged(error) {

    }

    /**
     * Runs when the value of an attribute is changed on the component
     * @param  {String} name     The attribute name
     * @param  {String} oldValue The old attribute value
     * @param  {String} newValue The new attribute value
     */
    attributeChangedCallback(name, oldValue, newValue) {
        console.log(`Attribute ${name} has changed.`);
        switch (name) {
            case "value":
                this.valueChanged(newValue);
                break;
            case "label":
                this.labelChanged(newValue);
                break;
            case "placeholder":
                this.placeholderChanged(newValue);
                break;
            case "error-message":
                this.errorChanged(newValue);
                break;
            default:
                console.log("unknown attribute: ", name);
        }
    }
}

customElements.define("form-input", FormInput);
