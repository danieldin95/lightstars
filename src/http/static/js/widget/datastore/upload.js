import {FormModal} from "../form/modal.js";

export class FileUpload extends FormModal {
    //
    constructor (props) {
        super(props);

        this.render();
        this.loading();
    }

    submit() {
        if (this.events.submit.func) {
            this.events.submit.func({
                data: this.events.submit.data,
                form: new FormData($(this.forms)[0]),
            });
        }
    }

    template() {
        return (`
        <div class="modal-dialog modal-dialog-centered model-md" role="document">
        <div class="modal-content">
            <div class="modal-header">
                <h5 class="modal-title" id="">Upload a file</h5>
            </div>
            <div id="" class="modal-body">
            <form name="upload-new">
                <div class="form-group">
                    <div class="input-group">
                        <input type="file" class="form-control-file" name="file" value="Upload">
                    </div>
                </div>
            </form>
            </div>
            <div id="" class="modal-footer">
                <button name="cancel-btn" class="btn btn-outline-dark btn-sm">Cancel</button>
                <button name="finish-btn" class="btn btn-outline-success btn-sm">Finish</button>
            </div>
        </div>
        </div>`);
    }
}
