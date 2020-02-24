import {ModalFormBase} from "../form/modal.js";

export class DiskCreate extends ModalFormBase {
    //
    constructor (props) {
        super(props);

        this.render();
        this.loading();
    }

    loading() {
        this.container().find('#finish-btn').on('click', this, function(e) {
            e.data.submit();
            e.data.container().modal("hide");
        });
        this.container().find('#cancel-btn').on('click', this, function(e) {
            e.data.container().modal("hide");
        });
    }

    template() {
        return `
    <div class="modal-dialog modal-dialog-centered model-md" role="document">
        <div class="modal-content">
            <div class="modal-header">
                <h5 class="modal-title" id="">Create Disk</h5>
            </div>
            <form name="disk-new">
            <div id="" class="modal-body">
                <div class="form-group row">
                    <label for="bus" class="col-sm-4 col-form-label-sm ">Target bus</label>
                    <div class="col-sm-6">
                        <div class="input-group">
                            <select class="select-simple select-middle" name="bus">
                                <option value="virtio" selected>Linux Virtual IO</option>
                                <option value="scsi">Logical SCSI</option>
                                <option value="ide">Logial IDE</option>
                            </select>  
                        </div>
                    </div>              
                </div>
                <div class="form-group row">
                    <label for="size" class="col-sm-4 col-form-label-sm">Virtual disk size</label>
                    <div class="col-sm-6">
                        <div class="input-group">
                            <input type="text" class="form-control form-control-sm input-number-lg" name="size" value="10">
                            <select class="select-simple select-unit-right" name="unit">
                                <option value="Mib">MiB</option>
                                <option value="GiB" selected>GiB</option>
                                <option value="TiB">TiB</option>
                            </select>
                        </div>
                    </div>
                </div>
            </div>
            <div id="" class="modal-footer">
                <div class="mr-0" rol="group">
                    <button id="finish-btn" class="btn btn-outline-success btn-sm">Finish</button>
                    <button id="reset-btn" class="btn btn-outline-dark btn-sm" type="reset">Reset</button>
                    <button id="cancel-btn" class="btn btn-outline-dark btn-sm">Cancel</button>
                </div>
            </div>
            </form>
        </div>
    </div>`
    }
}