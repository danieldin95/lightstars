import {FormModal} from "../form/modal.js";

export class ConfirmAction extends FormModal {
    constructor(props) {
        super(props);
        this.action = props.action || "confirm";
        this.name = props.name || "";
        this.render();
        this.loading();
    }

    template() {
        let prompt = this.action === "destroy" ? "are you sure you want to destroy" : "are you sure you want to remove";
        return this.compile(`
        <div class="modal-dialog modal-dialog-centered model-md" role="document">
        <div class="modal-content">
            <div class="modal-header">
                <h7 class="modal-title text-danger">{{'danger' | i}}</h7>
            </div>
            <div class="modal-body">
                <p class="text-center font-weight-normal">
                    {{'${prompt}' | i}}
                    <span class="font-weight-bold">${this.name}</span>{{'?' | i}}
                </p>
            </div>
            <div class="modal-footer">
                <button name="cancel-btn" class="btn btn-outline-dark btn-sm">{{'cancel' | i}}</button>
                <button name="finish-btn" class="btn btn-outline-danger btn-sm">{{'confirm' | i}}</button>
            </div>
        </div>
        </div>`);
    }
}
