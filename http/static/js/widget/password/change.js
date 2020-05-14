import {FormModal} from "../form/modal.js";
import {FormWizard} from "../form/wizard.js";

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
                        Change User Password
                </div>
               
                <div class="modal-body">
                    <form id="form">
                        <div class="form-group row">
                            <label for="name" class="col-sm-4 col-md-4 col-form-label-sm">Old Password</label>
                            <div class="col-sm-10 col-md-6"">
                                <div class="input-group">
                                    <input type="password" class="form-control form-control-sm"
                                           name="oldPassword" value=""/>
                                </div>
                            </div>
                        </div>
                        <div class="form-group row">
                            <label for="name" class="col-sm-4 col-md-4 col-form-label-sm">New Password</label>
                            <div class="col-sm-10 col-md-6"">
                                <div class="input-group">
                                    <input type="password" class="form-control form-control-sm"
                                           name="newPassword" value=""/>
                                </div>
                            </div>
                        </div>
                        <div class="form-group row">
                            <label for="name" class="col-sm-4 col-md-4 col-form-label-sm">Repeat New Password</label>
                            <div class="col-sm-10 col-md-6"">
                                <div class="input-group">
                                    <input type="password" class="form-control form-control-sm"
                                           name="repeatPassword" value=""/>
                                </div>
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
