import {FormModal} from "../form/modal.js";
import {Option} from "../option.js";
import {Alert} from "../../com/alert.js";


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
                $.getJSON("/api/bridge", function (data) {
                    selector.find("option").remove();
                    data.forEach(function (e, i) {
                        if (e['type'] == 'bridge') {
                            selector.append(Option(`Linux Bridge #${e['name']}`, e['name']));
                        } else if (e['type'] == 'openvswitch') {
                            selector.append(Option(`Open vSwitch #${e['name']}`, e['name']));
                        }
                    });
                }).fail(function (e) {
                    $("tasks").append(Alert.danger(`${this.type} ${this.url}: ${e.responseText}`));
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
        return (`
        <div class="modal-dialog modal-dialog-centered model-md" role="document">
        <div class="modal-content">
            <div class="modal-header">
                <h5 class="modal-title" id="">Create Interface</h5>
            </div>
            <form name="interface-new">
            <div id="" class="modal-body">
                <div class="form-group row">
                    <label for="type" class="col-sm-4 col-form-label-sm">Network type</label>
                    <div class="col-sm-6">
                        <div class="input-group">
                            <select class="select-lg" name="type">
                                <option value="bridge" selected>Linux Bridge</option>
                                <option value="openvswitch">Open vSwitch</option>
                            </select>
                        </div>
                    </div>
                </div>
                <div class="form-group row">
                    <label for="model" class="col-sm-4 col-form-label-sm">Target model</label>
                    <div class="col-sm-6">
                        <div class="input-group">
                            <select class="select-lg" name="model">
                                <option value="virtio" selected>Linux Virtual IO</option>
                                <option value="rtl8139">Realtek rtl8139</option>
                                <option value="e1000">Intel e1000</option>
                            </select>
                        </div>
                    </div>
                </div>
                <div class="form-group row">
                    <label for="source" class="col-sm-4 col-form-label-sm">Bridge source</label>
                    <div class="col-sm-6">
                        <div class="input-group">
                            <select class="select-lg" name="source">
                                <option value="ovs-br1">Open vSwitch #ovs-br1</option>
                                <option value="ovs-br2">Open vSwitch #ovs-br2</option>
                                <option value="br-mgt">Open vSwitch #br-mgt</option>
                            </select>
                        </div>
                    </div>
                </div>
                <div class="form-group row">
                    <label for="slot" class="col-sm-4 col-form-label-sm ">Sequence number</label>
                    <div class="col-sm-6">
                        <div class="input-group">
                            <select class="select-lg" name="seq">
                                <option value="0" selected>0</option>
                                <option value="1">1</option>
                                <option value="2">2</option>
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