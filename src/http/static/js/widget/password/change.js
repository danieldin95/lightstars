import {FormModal} from "../form/modal.js";

export class ChangePassword extends FormModal {

    constructor (props) {
        super(props);

        this.render();
        this.loading();
    }

    render() {
        super.render();
    }

    template(props) {
        return (`
        <div class="modal-dialog modal-dialog-centered model-md" role="document">
            <div class="modal-content ">
                <div class="modal-header">
                        Change user password
                </div>
                <div class="modal-body">
                    <form id="form">
                        <div class="form-group">
                            <label for="name" class="col-form-label-sm">Old Password</label>
                            <div class="input-group">
                                <input type="password" class="form-control form-control-sm" name="old" value=""/>
                            </div>
                        </div>
                        <div class="form-group">
                            <label for="name" class=" col-form-label-sm">New Password</label>
                            <div class="input-group">
                                <input type="password" class="form-control form-control-sm" name="new" value=""/>
                            </div>
                        </div>
                        <div class="form-group">
                            <label for="name" class="col-form-label-sm">Repeat New Password</label>
                            <div class="input-group">
                                <input type="password" class="form-control form-control-sm" name="repeat" value=""/>
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
