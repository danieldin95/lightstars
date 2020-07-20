import {FormModal} from "../form/modal.js";


export class GraphicsCreate extends FormModal {
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
                <h7 class="modal-title" id="">{{'add graphics' | i}}</h7>
            </div>
            <div id="" class="modal-body">
                <form name="disk-new">
                    <input type="text" class="d-none" name="autoport" value="yes"/>
                    <div class="form-group">
                        <label for="bus" class="col-form-label-sm ">{{'type' | i}}</label>
                        <div class="input-group">
                            <select class="select-lg" name="type">
                                <option value="vnc" selected>VNC</option>
                                <option value="spice">Spice</option>
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
