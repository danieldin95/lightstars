import {FormModal} from "../form/modal.js";
import {Option} from "../option.js";


export class InstanceSet extends FormModal {
    //
    constructor (props) {
        super(props);

        this.cpu = props.cpu ? props.cpu : 1;
        this.mem = props.mem ? props.mem / 1024 : 1;
        this.render();
        this.loading();
    }

    render() {
        super.render();
        let cpu = this.view.find("select[name='cpu']");
        cpu.find("option").remove();
        for (let i = 1; i < 17; i++) {
            cpu.append(new Option(i, i));
        }
        cpu.val(this.cpu).change();
    }

    template() {
        return (`
        <div class="modal-dialog modal-dialog-centered model-md" role="document">
        <div class="modal-content">
            <div class="modal-header">
                <h5 class="modal-title" id="">Instance setting</h5>
            </div>
            <form name="instance-set">
                <input type="text" class="d-none" name="mode" value=""/>
                <div id="" class="modal-body">
                    <div class="form-group row">
                        <label for="cpu" class="col-sm-2 col-md-4 col-form-label-sm">
                            Processors
                        </label>
                        <div class="col-sm-10 col-md-6">
                            <div class="input-group">
                                <select class="select-md" name="cpu" value="${this.cpu}">
                                    <option value="1">1</option>
                                    <option value="2" selected>2</option>
                                    <option value="3">3</option>
                                    <option value="4">4</option>
                                </select>
                            </div>
                        </div>
                    </div>
                    <div class="form-group row">
                        <label for="MaxMem" class="col-sm-4 col-md-4 col-form-label-sm">Memory size</label>
                        <div class="col-sm-10 col-md-6">
                            <div class="input-group">
                                <input type="text" class="form-control form-control-sm input-number-lg"
                                       name="memSize" value="${this.mem}"/>
                                <select class="select-unit-right" name="memUnit">
                                    <option value="MiB" selected>MiB</option>
                                    <option value="GiB">GiB</option>
                                </select>
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