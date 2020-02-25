import {ModalFormBase} from "../form/modal.js";
import {Option} from "../option.js";


export class DiskCreate extends ModalFormBase {
    //
    constructor (props) {
        super(props);

        this.render();
        this.loading();
    }

    render() {
        this.view = $(this.template());
        this.view.find("select[name='slot'] option").remove();
        for (let i = 1; i < 9; i++) {
            this.view.find("select[name='slot']").append(new Option(i, i));
        }
        this.container().html(this.view);
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
                            <select class="select-lg" name="bus">
                                <option value="virtio" selected>Linux Virtual IO</option>
                                <option value="scsi">Logical SCSI</option>
                                <option value="ide">Logical IDE</option>
                            </select>  
                        </div>
                    </div>              
                </div>
                <div class="form-group row">
                    <label for="slot" class="col-sm-4 col-form-label-sm ">Slot address</label>
                    <div class="col-sm-6">
                        <div class="input-group">
                            <select class="select-lg" name="slot">
                                <option value="0" selected>0</option>
                                <option value="1">1</option>
                                <option value="2">2</option>
                            </select>  
                        </div>
                    </div>              
                </div>
                <div class="form-group row">
                    <label for="size" class="col-sm-4 col-form-label-sm">Virtual disk size</label>
                    <div class="col-sm-6">
                        <div class="input-group">
                            <input type="text" class="form-control form-control-sm input-number-lg" name="size" value="10">
                            <select class="select-unit-right" name="unit">
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
                    <button name="finish-btn" class="btn btn-outline-success btn-sm">Finish</button>
                    <button name="reset-btn" class="btn btn-outline-dark btn-sm" type="reset">Reset</button>
                    <button name="cancel-btn" class="btn btn-outline-dark btn-sm">Cancel</button>
                </div>
            </div>
            </form>
        </div>
    </div>`
    }
}