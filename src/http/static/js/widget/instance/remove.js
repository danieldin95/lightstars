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
        return this.compile(`
        <div class="modal-dialog modal-dialog-centered model-md" role="document">
        <div class="modal-content">
            <div class="modal-header">
                <h6 class="modal-title text-danger" id="">{{'danger' | i}}</h6>
            </div>
            <div id="" class="modal-body">
                <p class="text-center font-weight-bold">
                   {{'are you sure you want to remove' | i}} <span class="text-danger">${this.name}</span>{{'?' | i}}
                </p>
                <p class="text-center font-weight-normal font-italic">
                  {{'if you confirm to remove it, all data of this guest will be clear.' | i}}
                </p>
            </div>
            <div id="" class="modal-footer">
                <button name="cancel-btn" class="btn btn-outline-dark btn-sm">{{'cancel' | i}}</button>
                <button name="finish-btn" class="btn btn-outline-danger btn-sm">{{'confirm' | i}}</button>
            </div>
        </div>
        </div>`);
    }
}
