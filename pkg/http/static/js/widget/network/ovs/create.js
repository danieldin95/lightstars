import {FormModal} from "../../form/modal.js";


export class OVSCreate extends FormModal {
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
                <h7 class="modal-title" id="">{{'create open vswitch network' | i}}</h7>
            </div>
            <div id="" class="modal-body">
                <form>
                    <input type="text" class="d-none" name="mode" value="bridge"/>
                    <input type="text" class="d-none" name="dhcp" value="no"/>
                    <input type="text" class="d-none" name="type" value="openvswitch"/>
                    <div class="form-group">
                        <label for="name" class="col-form-label-sm ">{{'network name' | i}}</label>
                        <div class="input-group">
                            <input type="text" class="form-control form-control-sm" name="name" value=""/>
                        </div>
                    </div>
                    <div class="form-group">
                        <label for="name" class="col-form-label-sm ">{{'existing ovs bridge' | i}}</label>
                        <div class="input-group">
                            <input type="text" class="form-control form-control-sm input-lg" name="bridge" value="br-int"/>                                    
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
