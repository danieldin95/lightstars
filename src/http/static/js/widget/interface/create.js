import {FormModal} from "../form/modal.js";
import {Option} from "../option.js";
import {BridgeApi} from "../../api/bridge.js";


export class InterfaceCreate extends FormModal {
    //
    constructor (props) {
        super(props);

        this.render();
        this.loading();
    }

    render() {
        super.render();
        let iface = {
            fresh: function () {
                let selector = this.selector;

                new BridgeApi().list((data) => {
                    selector.find("option").remove();
                    for (let e of data.resp) {
                        if (e['type'] === 'bridge') {
                            selector.append(Option(`Linux Bridge #${e['name']}`, e['name']));
                        } else if (e['type'] === 'openvswitch') {
                            selector.append(Option(`Open vSwitch #${e['name']}`, e['name']));
                        }
                    }
                });
            },
            selector: this.view.find("select[name='source']"),
        };
        let cpu = {
            fresh: function() {
                this.selector.find('option').remove();
                for (let i = 1; i < 17; i++) {
                    this.selector.append(new Option(i, i));
                }
            },
            selector: this.view.find("select[name='seq']"),
        };

        cpu.fresh();
        iface.fresh();
    }

    template() {
        return this.compile(`
        <div class="modal-dialog modal-dialog-centered model-md" role="document">
        <div class="modal-content">
            <div class="modal-header">
                <h7 class="modal-title" id="">{{'add interface' | i}}</h7>
            </div>
            <div id="" class="modal-body">
            <form>
                <div class="form-group">
                    <label for="model" class="col-form-label-sm">{{'target model' | i}}</label>
                    <div class="input-group">
                        <select class="select-lg" name="model">
                            <option value="virtio" selected>Linux Virtual IO</option>
                            <option value="rtl8139">Realtek rtl8139</option>
                            <option value="e1000">Intel e1000</option>
                        </select>
                    </div>
                </div>
                <div class="form-group">
                    <label for="source" class="col-form-label-sm">{{'bridge name' | i}}</label>
                    <div class="input-group">
                        <select class="select-lg" name="source">
                            <option value="ovs-br1">Open vSwitch #ovs-br1</option>
                            <option value="ovs-br2">Open vSwitch #ovs-br2</option>
                            <option value="br-mgt">Open vSwitch #br-mgt</option>
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
