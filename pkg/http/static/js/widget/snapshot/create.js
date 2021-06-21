import {FormModal} from "../form/modal.js";


export class SnapshotCreate extends FormModal {
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
                <h7 class="modal-title" id="">{{'create snapshot' | i}}</h7>
            </div>
            <div id="" class="modal-body">
                <form name="snapshot-new">
                    <div class="form-group">
                        <label for="bus" class="col-form-label-sm ">{{'name' | i}}</label>
                          <div class="input-group">
                            <input type="text" class="form-control form-control-sm" name="name" value="snapshot-"/>
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
