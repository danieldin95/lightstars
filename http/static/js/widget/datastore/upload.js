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
                <form name="upload-new">
                <div id="" class="modal-body">
                    <div class="form-group row">
                        <label for="name" class="col-sm-4 col-form-label-sm ">Select file</label>
                        <div class="col-sm-6">
                            <div class="input-group">
                                <input type="file" class="form-control-file" name="file" value="Upload">
                            </div>
                        </div>
                    </div>
                </div>
                <div id="" class="modal-footer">
                    <button name="reset-btn" class="btn btn-outline-dark btn-sm" type="reset">Reset</button>
                    <button name="cancel-btn" class="btn btn-outline-dark btn-sm">Cancel</button>
                    <button name="finish-btn" class="btn btn-outline-success btn-sm">Finish</button>
                </div>
                </form>
            </div>
        </div>`);
    }
}