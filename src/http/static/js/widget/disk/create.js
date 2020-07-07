import {FormModal} from "../form/modal.js";
import {Option} from "../option.js";


export class DiskCreate extends FormModal {
    //
    constructor (props) {
        super(props);

        this.render();
        this.loading();
    }

    render() {
        super.render();
        this.view.find("select[name='seq'] option").remove();
        for (let i = 1; i < 16; i++) {
            this.view.find("select[name='seq']").append(new Option(i, i));
        }
    }

    template() {
        return this.compile(`
        <div class="modal-dialog modal-dialog-centered model-md" role="document">
        <div class="modal-content">
            <div class="modal-header">
                <h5 class="modal-title" id="">{{'add disk' | i}}</h5>
            </div>
            <div id="" class="modal-body">
                <form name="disk-new">
                    <div class="form-group">
                        <label for="bus" class="col-form-label-sm ">{{'target bus' | i}}</label>
                        <div class="input-group">
                            <select class="select-lg" name="bus">
                                <option value="virtio" selected>Linux Virtual IO</option>
                                <option value="scsi">Logical SCSI</option>
                                <option value="ide">Logical IDE</option>
                            </select>  
                        </div>            
                    </div>
                    <div class="form-group">
                        <label for="slot" class="col-form-label-sm ">{{'sequence number' | i}}</label>
                        <div class="input-group">
                            <select class="select-lg" name="seq">
                                <option value="0" selected>0</option>
                                <option value="1">1</option>
                                <option value="2">2</option>
                            </select>
                        </div>              
                    </div>
                    <div class="form-group">
                        <label for="size" class="col-form-label-sm">{{'disk size' | i}}</label>
                        <div class="input-group">
                            <input type="text" class="form-control form-control-sm" name="size" value="10"/>
                            <select class="select-unit-right" name="sizeUnit">
                                <option value="Mib">MiB</option>
                                <option value="GiB" selected>GiB</option>
                                <option value="TiB">TiB</option>
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
