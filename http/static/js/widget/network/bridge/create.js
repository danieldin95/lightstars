import {FormModal} from "../../form/modal.js";


export class BridgeCreate extends FormModal {
    //
    constructor (props) {
        super(props);

        this.render();
        this.loading();
    }

    template() {
        return (`
        <div class="modal-dialog modal-dialog-centered model-md" role="document">
        <div class="modal-content">
            <div class="modal-header">
                <h5 class="modal-title" id="">Create Network</h5>
            </div>
            <form name="network-new">
                <input type="text" class="d-none" name="mode" value="bridge"/>
                <input type="text" class="d-none" name="dhcp" value="no"/>
                <div id="" class="modal-body">
                    <div class="form-group row">
                        <label for="name" class="col-sm-4 col-form-label-sm ">Bridge's name</label>
                        <div class="col-sm-6">
                            <div class="input-group">
                                <input type="text" class="form-control form-control-sm input-lg"
                                        name="name" value="br0"/>                                    
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