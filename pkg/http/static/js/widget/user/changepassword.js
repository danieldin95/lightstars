import {formmodal} from "../form/formmodal.js";


export class changepassword extends formmodal {
    // {
    //   id: ""
    // }
    constructor (props) {
        super(props);

        this.render();
        this.loading();
    }

    render() {
        super.render();
    }

    template(props) {
        return this.compile(`
        <div class="modal-dialog modal-dialog-centered model-md" role="document">
            <div class="modal-content ">
                <div class="modal-header">
                    <h7>{{'change user password' | i}}</h7>
                </div>
                <div class="modal-body">
                    <form id="form-change">
                        <div class="form-group">
                            <label for="old" class="col-form-label-sm">{{'old password' | i}}</label>
                            <div class="input-group">
                                <input type="password" class="form-control form-control-sm" id="old" name="old" value="" autofocus/>
                            </div>
                        </div>
                        <div class="form-group">
                            <label for="new" class=" col-form-label-sm">{{'new password' | i}}</label>
                            <div class="input-group">
                                <input type="password" class="form-control form-control-sm" id="new" name="new" value=""/>
                            </div>
                        </div>
                        <div class="form-group">
                            <label for="repeat" class="col-form-label-sm">{{'repeat new password' | i}}</label>
                            <div class="input-group">
                                <input type="password" class="form-control form-control-sm" id="repeat" name="repeat" value=""/>
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
