class FormInput extends HTMLElement {
    constructor() {
        super();

        /**
         * @returns {HTMLTemplateElement}
         */

        //expose input element change functions
        this.internals = this.attachInternals();

        let template = document.getElementById("form-input-template");
        console.log("template: ", template);
        let templateContent = template.content;

        const shadowRoot = this.attachShadow({ mode: "open" });
        shadowRoot.appendChild(templateContent.cloneNode(true));
    }

    //input field refs
    inputField;
    labelElem;

    //component state
    dataObject = {
        name: "",
        value: "",
        placeholder: "",
        label: "",
    };

    //enable values to make it into a form
    static formAssociated = true;

    updateForm(e) {
        console.log("changer value: ", e.target.value)
        this.internals.setFormValue(e.target.value)
    }

    connectedCallback() {
        console.log("I've been a connect");
        //input
        this.inputField = this.shadowRoot.querySelector("input");
        this.inputField.name = this.dataObject.name;
        this.inputField.placeholder = this.dataObject.placeholder;

        //this.inputField.value = this.dataObject.value;
        this.internals.setFormValue(this.dataObject.value);

        //listen for changes
        this.inputField.addEventListener("change", this.updateForm.bind(this))

        //label
        this.labelElem = this.shadowRoot.querySelector("label");
        if (this.dataObject.label) {
            this.labelElem.innerHTML = this.dataObject.label;
        }
    }

    disconnectedCallback() {
        this.inputField.removeEventListener("change", this.updateForm.bind(this))
    }

    static get observedAttributes() {
        return ["name", "value", "label", "placeholder", "error-message"];
    }

    //these methods update when the component attributes are changed from a parent
    nameChanged(value) {
        this.dataObject.name = value;
    }

    valueChanged(value) {
        this.dataObject.value = value;
    }

    labelChanged(value) {
        this.dataObject.label = value;
    }

    placeholderChanged(value) {
        this.dataObject.placeholder = value;
    }

    errorChanged(error) {}

    /**
     * Runs when the value of an attribute is changed on the component
     * @param  {String} name     The attribute name
     * @param  {String} oldValue The old attribute value
     * @param  {String} newValue The new attribute value
     */
    attributeChangedCallback(name, oldValue, newValue) {
        console.log(`Attribute ${name} has changed.`);
        if (!this.inputField) {
            console.log("why:", this.inputField);
        }

        switch (name) {
            case "name":
                this.nameChanged(newValue);
                break;
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
