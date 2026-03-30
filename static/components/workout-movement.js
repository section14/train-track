export const debounce = (waitTime) => {
    let timeout = null;

    /**
     * @param {() => void} func
     */
    const caller = (func) => {
        const later = () => {
            clearTimeout(timeout);
            func();
        };

        clearTimeout(timeout);
        timeout = setTimeout(later, waitTime);
    };

    caller.cancel = () => {
        clearTimeout(timeout);
    };

    return caller;
};

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
    };

    index = -1;
    saving = false;
    debouncer = debounce(2000);
    movementForm;
    exercise;
    sets;
    reps;
    weight;

    submitDone = () => {
        this.saving = false;

        const event = new CustomEvent("change-status", {
            detail: { saving: false },
            bubbles: true,
            cancelable: true,
        });

        const mEvent = new CustomEvent("movement-changed", {
            detail: { id: this.dataObject.workoutId, saving: false },
            bubbles: true,
            cancelable: true,
        });
        document.dispatchEvent(mEvent);

        this.saver.dispatchEvent(event);
        this.deleteBtn.disabled = false;
    };

    submitForm = (e) => {
        if (!e) return;
        if (e) e.preventDefault();
        this.debouncer.cancel();
        const formData = new FormData(this.movementForm);

        let entries = {};
        for (const [key, value] of formData.entries()) {
            entries[key] = value;
        }

        const mObj = {
            id: parseInt(this.dataObject.id),
            workoutId: parseInt(this.dataObject.workoutId),
            exerciseId: parseInt(entries["exercise"]),
            sets: parseInt(entries["sets"]),
            reps: parseInt(entries["reps"]),
            weight: parseInt(entries["weight"]),
        };

        const event = new CustomEvent("update-workout", {
            detail: { movement: mObj, callback: () => this.submitDone() },
            bubbles: true,
            cancelable: true,
        });
        document.dispatchEvent(event);
    };

    deleteMovement = () => {
        const event = new CustomEvent("delete-movement", {
            detail: { id: this.dataObject.id, workoutId: this.dataObject.workoutId },
            bubbles: true,
            cancelable: true,
        });
        document.dispatchEvent(event);
    };

    movementUpdate(e) {
        //start a saving indicator here
        if (!this.saving) {
            const event = new CustomEvent("change-status", {
                detail: { saving: true },
                bubbles: true,
                cancelable: true,
            });

            const mEvent = new CustomEvent("movement-changed", {
                detail: { id: this.dataObject.workoutId, saving: true },
                bubbles: true,
                cancelable: true,
            });
            document.dispatchEvent(mEvent);

            this.deleteBtn.disabled = true;
            this.saver.dispatchEvent(event);
        }
        this.saving = true;
        this.debouncer(() => this.submitForm(e));
    }

    mobileLayout(e) {
        const buttons = this.shadowRoot.getElementById("button-group")

        if (e.matches) {
            buttons.style.setProperty("margin-top", "0");
            buttons.style.setProperty("padding-bottom", "1.0rem");
            buttons.style.setProperty("border-bottom", "1px solid var(--med)");
            buttons.style.setProperty("justify-content", "end");
            this.exercise.setAttribute("show-label", "yes");
            this.sets.setAttribute("show-label", "yes");
            this.reps.setAttribute("show-label", "yes");
            this.weight.setAttribute("show-label", "yes");
        } else {
            if (this.index === "0") {
                buttons.style.setProperty("margin-top", "2.0rem");
                this.exercise.setAttribute("show-label", "yes");
                this.sets.setAttribute("show-label", "yes");
                this.reps.setAttribute("show-label", "yes");
                this.weight.setAttribute("show-label", "yes");
            } else {
                buttons.style.setProperty("margin-top", "0");
                this.exercise.setAttribute("show-label", "no");
                this.sets.setAttribute("show-label", "no");
                this.reps.setAttribute("show-label", "no");
                this.weight.setAttribute("show-label", "no");
            }

            buttons.style.setProperty("padding-bottom", "0");
            buttons.style.setProperty("border-bottom", "0");
            buttons.style.setProperty("justify-content", "start");
        }
    }

    connectedCallback() {
        // mobile check
        this.mediaMobile = window.matchMedia("(width < 760px)");
        this.mediaMobile.addEventListener("change", this.mobileLayout.bind(this))

        // form submit listener
        this.movementForm = this.shadowRoot.getElementById("movement-form");
        this.movementForm.addEventListener("submit", this.submitForm.bind(this));
        this.movementForm.addEventListener("movement-update", this.movementUpdate.bind(this));

        // delete listener
        this.deleteBtn = this.shadowRoot.getElementById("delete-mvmt");
        this.deleteBtn.addEventListener("click", this.deleteMovement.bind(this));

        // save indicator
        this.saver = this.shadowRoot.getElementById("saver");

        // child element refs -- set defaults
        this.exercise = this.shadowRoot.getElementById("exercise");
        this.sets = this.shadowRoot.getElementById("sets");
        this.reps = this.shadowRoot.getElementById("reps");
        this.weight = this.shadowRoot.getElementById("weight");

        this.exercise.setAttribute("value", this.dataObject.exerciseId);
        this.sets.setAttribute("value", this.dataObject.sets);
        this.reps.setAttribute("value", this.dataObject.reps);
        this.weight.setAttribute("value", this.dataObject.weight);

        // check size on first render
        this.mobileLayout(this.mediaMobile)

        // register current data with parent
        this.submitForm(null);
    }

    disconnectedCallback() {
        this.mediaMobile.removeEventListener("change", this.mobileLayout.bind(this))
        this.movementForm.removeEventListener("submit", this.submitForm.bind(this));
        this.movementForm.removeEventListener("movement-update", this.movementUpdate.bind(this));
        this.deleteBtn.removeEventListener("click", this.deleteMovement.bind(this));
    }

    static get observedAttributes() {
        return ["id", "index", "workout-id", "exercise-id", "sets", "reps", "weight"];
    }

    attributeChangedCallback(name, oldValue, newValue) {
        switch (name) {
            case "index":
                this.index = newValue;
                break
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
