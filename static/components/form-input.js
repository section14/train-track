class FormInput extends HTMLElement {
    constructor() {
        super();

        /**
         * @returns {HTMLTemplateElement}
         */

        //expose input element change functions
        this.internals = this.attachInternals();

        let template = document.getElementById("form-input-template");
        let templateContent = template.content;

        const shadowRoot = this.attachShadow({ mode: "open" });
        shadowRoot.appendChild(templateContent.cloneNode(true));
    }

    //input field refs
    inputField;
    labelElem;
    showLabel = "no";
    style = "";

    //component state
    dataObject = {
        name: "",
        type: "text",
        value: "",
        placeholder: "",
        label: "",
    };

    //enable values to make it into a form
    static formAssociated = true;

    updateForm(e) {
        this.internals.setFormValue(e.target.value);

        const event = new CustomEvent("movement-update", {
            detail: "",
            bubbles: true,
            cancelable: true,
        });
        this.internals.form.dispatchEvent(event);
    }

    connectedCallback() {
        //input
        this.inputField = this.shadowRoot.querySelector("input");
        this.inputField.name = this.dataObject.name;
        this.inputField.type = this.dataObject.type;
        this.inputField.placeholder = this.dataObject.placeholder;
        this.inputField.setAttribute("style", this.style);

        this.internals.setFormValue(this.dataObject.value);

        this.renderLabel();

        //listen for changes
        this.inputField.addEventListener("input", this.updateForm.bind(this));
    }

    renderLabel() {
        //label
        this.labelElem = this.shadowRoot.querySelector("label");
        if (this.dataObject.label && this.showLabel === "yes") {
            this.labelElem.style.setProperty("display", "block");
            this.labelElem.innerHTML = this.dataObject.label;
        } else {
            this.labelElem.style.setProperty("display", "none");
        }
    }

    disconnectedCallback() {
        this.inputField.removeEventListener("change", this.updateForm.bind(this));
    }

    static get observedAttributes() {
        return [
            "name",
            "type",
            "value",
            "label",
            "show-label",
            "placeholder",
            "style",
            "error-message",
        ];
    }

    //these methods update when the component attributes are changed from a parent
    typeChanged(value) {
        this.dataObject.type = value;
    }

    nameChanged(value) {
        this.dataObject.name = value;
    }

    valueChanged(value) {
        this.dataObject.value = value;
        if (this.inputField) {
            this.inputField.setAttribute("value", value);
        }
    }

    labelChanged(value) {
        this.dataObject.label = value;
    }

    placeholderChanged(value) {
        this.dataObject.placeholder = value;
    }

    styleChanged(value) {
        this.style = value;
        if (this.inputField) {
            this.inputField.setAttribute("style", value);
        }
    }

    errorChanged(error) {}

    /**
     * Runs when the value of an attribute is changed on the component
     * @param  {string} name     The attribute name
     * @param  {string} oldValue The old attribute value
     * @param  {string} newValue The new attribute value
     */
    attributeChangedCallback(name, oldValue, newValue) {
        switch (name) {
            case "type":
                this.typeChanged(newValue);
                break;
            case "name":
                this.nameChanged(newValue);
                break;
            case "value":
                this.valueChanged(newValue);
                if (this.internals) {
                    this.internals.setFormValue(newValue);
                }
                break;
            case "label":
                this.labelChanged(newValue);
                break;
            case "show-label":
                this.showLabel = newValue;
                this.renderLabel();
                break;
            case "placeholder":
                this.placeholderChanged(newValue);
                break;
            case "style":
                this.styleChanged(newValue);
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
