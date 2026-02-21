class WorkoutMovement extends HTMLElement {
    constructor() {
        super();

        let template = document.getElementById("workout-movement-template");
        let templateContent = template.content;

        const shadowRoot = this.attachShadow({ mode: "open" });
        shadowRoot.appendChild(templateContent.cloneNode(true));
    }

    dataObject = {
        id: 0,
        workoutId: 0,
        exerciseId: 0,
        sets: 0,
        reps: 0,
        weight: 0,
    }

    movementForm;
    exercise;
    sets;
    reps;
    weight;

    submitForm = (e) => {
        if (!e) return
        if (e) e.preventDefault();
        const formData = new FormData(this.movementForm);

        let entries = {};
        for (var [key, value] of formData.entries()) {
            entries[key] = value;
        }

        const mObj = {
            id: parseInt(this.dataObject.id),
            workoutId: parseInt(this.dataObject.workoutId),
            exerciseId: parseInt(entries["exercise"]),
            sets: parseInt(entries["sets"]),
            reps: parseInt(entries["reps"]),
            weight: parseInt(entries["weight"]),
        }

        const event = new CustomEvent("update-workout", {
            detail: mObj,
            bubbles: true,
            cancelable: true,
        })
        document.dispatchEvent(event)
    }

    deleteMovement = () => {
        const event = new CustomEvent("delete-movement", {
            detail: { id: this.dataObject.id, workoutId: this.dataObject.workoutId },
            bubbles: true,
            cancelable: true,
        })
        document.dispatchEvent(event)
    }

    connectedCallback() {
        //form submit listener
        this.movementForm = this.shadowRoot.getElementById("movement-form")
        this.movementForm.addEventListener("submit", this.submitForm.bind(this))

        //delete listener
        this.deleteBtn = this.shadowRoot.getElementById("delete-mvmt")
        this.deleteBtn.addEventListener("click", this.deleteMovement.bind(this))

        //child element refs -- set defaults
        this.exercise = this.shadowRoot.getElementById("exercise")
        this.sets = this.shadowRoot.getElementById("sets")
        this.reps = this.shadowRoot.getElementById("reps")
        this.weight = this.shadowRoot.getElementById("weight")

        this.exercise.setAttribute("value", this.dataObject.exerciseId)
        this.sets.setAttribute("value", this.dataObject.sets)
        this.reps.setAttribute("value", this.dataObject.reps)
        this.weight.setAttribute("value", this.dataObject.weight)

        //register current data with parent
        this.submitForm(null)
    }

    disconnectedCallback() {
        this.movementForm.removeEventListener("submit", this.submitForm.bind(this))
        this.deleteBtn.removeEventListener("click", this.deleteMovement.bind(this))
    }

    static get observedAttributes() {
        return ["id", "workout-id", "exercise-id", "sets", "reps", "weight"];
    }

    attributeChangedCallback(name, oldValue, newValue) {
        switch (name) {
            case "id":
                this.dataObject.id = newValue;
                break;
            case "workout-id":
                this.dataObject.workoutId = newValue;
                break;
            case "exercise-id":
                this.dataObject.exerciseId = newValue;
                break;
            case "sets":
                this.dataObject.sets = newValue;
                break;
            case "reps":
                this.dataObject.reps = newValue;
                break;
            case "weight":
                this.dataObject.weight = newValue;
                break;
        }
    }
}

customElements.define("workout-movement", WorkoutMovement);
