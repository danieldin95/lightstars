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
            <div id="" class="modal-body">
                <form>
                    <input type="text" class="d-none" name="mode" value="bridge"/>
                    <input type="text" class="d-none" name="dhcp" value="no"/>
                    <div class="form-group">
                        <label for="name" class="col-form-label-sm ">Network Name</label>
                        <div class="input-group">
                            <input type="text" class="form-control form-control-sm" name="name" value=""/>
                        </div>
                    </div>
                    <div class="form-group">
                        <label for="name" class="col-form-label-sm ">Existing Bridge</label>
                        <div class="input-group">
                            <input type="text" class="form-control form-control-sm input-lg" name="bridge" value="br0"/>                                    
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
