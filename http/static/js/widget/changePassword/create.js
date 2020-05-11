import {FormModal} from "../form/modal.js";
import {FormWizard} from "../form/wizard.js";

export class ChangePassword extends FormModal {

    constructor (props) {
        super(props);

        this.render();
        this.loading();

    }


    render() {
        console.log('??')
        super.render();
    }

    loading() {
        new FormWizard({
            id: this.id,
            default: '#head1',
            navigation: '#nav-tabs li a',
            form: '#form',
            buttons: {
                submit: '#btn-submit',
                cancel: '#btn-cancel',
            },
        }).load({
            submit: (e) => {
                this.submit(e);
                $(this.id).modal('hide');
            },
            cancel: (e) => {
                $(this.id).modal('hide');
            },
        });
    }

    template(props) {
        return (`
        <div class="modal-dialog modal-dialog-centered modal-lg" role="document">
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
                
                <div class="modal-footer text-right">
                    <button id="btn-cancel" class="btn btn-outline-dark btn-sm">Cancel</button>
                    <button id="btn-submit" class="btn btn-outline-success btn-sm">Submit</button>
                </div>
            </div>
        </div>`);
    }
}
