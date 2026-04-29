import {FormModalWid} from "../form/formmodal.js";


export class PreferencesWid extends FormModalWid {
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
                    <h7>{{'Preferences' | i}}</h7>
                </div>
                <div class="modal-body">
                    <form id="form-preferences">
                        <div class="form-group">
                            <label for="lang" class="col-form-label-sm">{{'language and region' | i}}</label>
                            <div class="input-group">
                                <select class="select-lg" id="lang" name="lang">
                                    <option value="en-US" selected>English/US</option>
                                    <option value="zh-CN">简体中文/中国</option>
                                </select>
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
