import {FormModal} from "../form/modal.js";


export class InstanceRemove extends FormModal {
    //
    constructor (props) {
        super(props);
        this.name = props.name;
        this.uuid = props.uuid;
        this.render();
        this.loading();
    }

    template() {
        return (`
        <div class="modal-dialog modal-dialog-centered model-md" role="document">
        <div class="modal-content">
            <div class="modal-header">
                <h6 class="modal-title text-dark" id="">Warning</h6>
            </div>
            <div id="" class="modal-body">
                <p class="text-center">
                  Are you sure to remove guest <span class="text-danger font-weight-bold">${this.name}</span>
                </p>
                <p class="text-center font-weight-bold">
                  If you confirm to remove it, all data of this guest will be clear.
                </p>
            </div>
            <div id="" class="modal-footer">
                <button name="cancel-btn" class="btn btn-outline-dark btn-sm">Cancel</button>
                <button name="finish-btn" class="btn btn-outline-danger btn-sm">Confirm</button>
            </div>
        </div>
        </div>`);
    }
}
