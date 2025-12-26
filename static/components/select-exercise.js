import { GetJson } from '/static/js/request.js'

class SelectExercise extends HTMLElement {
    constructor() {
        super();

        //expose input element change functions
        this.internals = this.attachInternals();

        let template = document.getElementById("select-exercise-template");
        let templateContent = template.content;

        const shadowRoot = this.attachShadow({ mode: "open" });
        shadowRoot.appendChild(templateContent.cloneNode(true));
    }

    selectField;

    dataObject = {
        name: "",
        value: "",
    }

    //enable values to make it into a form
    static formAssociated = true;

    updateForm(e) {
        this.internals.setFormValue(e.target.value);
    }

    connectedCallback() {
        this.selectField = this.shadowRoot.querySelector("select");
        this.selectField.name = this.dataObject.name

        GetJson("/api/json/exercises").then((res) => {
            const selectElem = this.shadowRoot.querySelector("select");

            for (let i = 0; i < res.length; i++) {
                const newOpt = document.createElement("option");
                newOpt.value = res[i].id;
                newOpt.label = res[i].name;
                selectElem.appendChild(newOpt);
            }

            if (this.dataObject.value) {
                selectElem.value = this.dataObject.value;
            }

        }).catch((err) => {
            console.log("error getting exercises for select: ", err);
        })

        //listen for changes
        this.selectField.addEventListener("change", this.updateForm.bind(this));
    }

    disconnectedCallback() {
        this.selectField.removeEventListener("change", this.updateForm.bind(this));
    }

    static get observedAttributes() {
        return ["name", "value"];
    }

    attributeChangedCallback(name, oldValue, newValue) {
        switch (name) {
            case "name":
                this.dataObject.name = newValue;
                break;
            case "value":
                this.dataObject.value = newValue;
                this.internals.setFormValue(newValue);
                this.selectField.setAttribute("value", newValue);
                break;
        }
    }
}

customElements.define("select-exercise", SelectExercise);
