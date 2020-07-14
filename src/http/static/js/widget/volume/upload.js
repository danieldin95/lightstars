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
        return this.compile(`
        <div class="modal-dialog modal-dialog-centered model-md" role="document">
        <div class="modal-content">
            <div class="modal-header">
                <h7 class="modal-title" id="">{{'upload a file' | i}}</h7>
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
                <button name="cancel-btn" class="btn btn-outline-dark btn-sm">{{'cancel' | i}}</button>
                <button name="finish-btn" class="btn btn-outline-success btn-sm">{{'finish' | i}}</button>
            </div>
        </div>
        </div>`);
    }
}
