import {FormModal} from "../form/modal.js";
import {Option} from "../option.js";


export class InstanceSet extends FormModal {
    //
    constructor (props) {
        super(props);

        this.cpu = props.data.maxCpu || 0;
        this.mem = props.data.maxMem ?  props.data.maxMem / 1024 : 0;
        this.cpuMode = props.data.cpuMode || '';
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
        let cpuMode = this.view.find("select[name='cpuMode']");
        cpuMode.val(this.cpuMode).change();
    }

    template() {
        return this.compile(`
        <div class="modal-dialog modal-dialog-centered model-md" role="document">
        <div class="modal-content">
            <div class="modal-header">
                <h7 class="modal-title" id="">{{'instance setting' | i}}</h7>
            </div>
            <div id="" class="modal-body">
                <form>
                    <input type="text" class="d-none" name="mode" value=""/>
                    <div class="form-group">
                        <label for="cpu" class="col-form-label-sm">{{'processors' | i}}</label>
                        <div class="input-group">
                           <select class="form-control form-control-sm" name="cpuMode" value="${this.cpuMode}">
                                <option value="" selected>Default</option>
                                <option value="host-passthrough">Host passthrough</option>
                            </select>
                            <select class="select-twice-md" name="cpu" value="${this.cpu}">
                                <option value="1">1</option>
                                <option value="2" selected>2</option>
                                <option value="3">3</option>
                                <option value="4">4</option>
                            </select>
                        </div>
                    </div>
                    <div class="form-group">
                        <label for="MaxMem" class="col-form-label-sm">{{'memory size' | i}}</label>
                        <div class="input-group">
                            <input type="text" class="form-control form-control-sm input-md"
                                   name="memSize" value="${this.mem}"/>
                            <select class="select-twice-md" name="memUnit">
                                <option value="MiB" selected>MiB</option>
                                <option value="GiB">GiB</option>
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
