import {FormModalWid} from "../form/formmodal.js";


export class BridgeCreateWid extends FormModalWid {
    //
    constructor (props) {
        super(props);

        this.render();
        this.loading();
    }

    template() {
        return this.compile(`
        <div class="modal-dialog modal-dialog-centered model-md" role="document">
        <div class="modal-content">
            <div class="modal-header">
                <h7 class="modal-title" id="">{{'create linux bridge network' | i}}</h7>
            </div>
            <div id="" class="modal-body">
                <form>
                    <input type="text" class="d-none" name="mode" value="bridge"/>
                    <input type="text" class="d-none" name="dhcp" value="no"/>
                    <div class="form-group">
                        <label for="bridge-name" class="col-form-label-sm ">{{'network name' | i}}</label>
                        <div class="input-group">
                            <input type="text" class="form-control form-control-sm" id="bridge-name" name="name" value=""/>
                        </div>
                    </div>
                    <div class="form-group">
                        <label for="bridge-host" class="col-form-label-sm ">{{'existing bridge on host' | i}}</label>
                        <div class="input-group">
                            <input type="text" class="form-control form-control-sm input-lg" id="bridge-host" name="bridge" value="" placeholder="br0"/>
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
