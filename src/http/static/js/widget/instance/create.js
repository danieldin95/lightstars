import {Option} from "../option.js";
import {FormModal} from "../form/modal.js";
import {FormWizard} from "../form/wizard.js";
import {BridgeApi} from "../../api/bridge.js";
import {IsoApi} from "../../api/iso.js"
import {DataStoreApi} from "../../api/datastores.js";


export class InstanceCreate extends FormModal {
    // {
    //   id: "#InstanceCrateModal",
    //   wizardId: "",
    // }
    constructor (props) {
        super(props);

        this.render();
        this.loading();

        $(this.id).off('show.bs.modal');
        $(this.id).on('show.bs.modal', (e) => {
            this.fetch();
        });
    }

    fetch() {
        let iso = {
            selector: this.view.find("select[name='disk0File']"),
            fresh: function (datastore) {
                let selector = this.selector;

                new IsoApi().list(datastore, (data) => {
                    selector.find("option").remove();
                    for (let ele of data.resp) {
                        selector.append(Option(ele['path'], ele['path']));
                    }
                    selector.append(Option('CDROM device:/sr0', '/dev/sr0'));
                });
            },
        };

        let store = {
            selector: this.view.find("select[name='datastore']"),
            refresh: function () {
                let selector = this.selector;
                new DataStoreApi().list(this,  (data) => {
                    let resp = data.resp;
                    selector.find("option").remove();
                    for (let ele of resp.items) {
                        selector.append(Option(ele['name'], ele['name']));
                    }
                    if (resp.items.length > 0) {
                        iso.fresh(resp.items[0]['name']);
                    }
                });
            },
        };

        let iface = {
            fresh: function (){
                let selector = this.selector;

                new BridgeApi().list(this, (data) => {
                    selector.find("option").remove();
                    for (let ele of data.resp) {
                        if (ele['type'] === 'bridge') {
                            selector.append(Option(`Linux Bridge #${ele['name']}`, ele['name']));
                        } else if (ele['type'] === 'openvswitch') {
                            selector.append(Option(`Open vSwitch #${ele['name']}`, ele['name']));
                        }
                    }
                });
            },
            selector: this.view.find("select[name='interface0Source']"),
        };

        iface.fresh();
        store.refresh();
        store.selector.on("change", this, function (e) {
            iso.fresh($(this).val());
        });
    }

    render() {
        super.render();
        this.view.find("select[name='cpu'] option").remove();
        for (let i = 1; i < 17; i++) {
            this.view.find("select[name='cpu']").append(new Option(i, i));
        }
    }

    loading() {
        new FormWizard({
            id: this.id,
            default: '#head1',
            navigation: '#nav-tabs li a',
            form: '#form',
            buttons: {
                prev: '#btn-prev',
                next: '#btn-next',
                submit: '#btn-submit',
                cancel: '#btn-cancel',
            },
        }).load({
            submit: (e) => {
                this.submit(e);
                $(this.id).modal('hide');
            },
            cancel: (e) => {
                $(this.id).modal('hide');
            },
        });
    }

    template(props) {
        return this.compile(`
        <div class="modal-dialog modal-dialog-centered modal-lg" role="document">
        <div class="modal-content ">
            <div class="modal-header">
                <h7>{{'create a instance' | i}}</h7>
            </div>
            <div class="modal-body">
            <div class="row">
                <div class="col-3 border-right">
                    <div class="">
                        <ul id="nav-tabs" class="nav flex-column nav-pills">
                            <li class="nav-item">
                                <a id="head1" class="nav-link" data-target="#guest">{{'configure guest' | i}}</a>
                            </li>
                            <li class="nav-item">
                                <a id="head2" class="nav-link" data-target="#datastore">{{'select datastore' | i}}</a>
                            </li>
                            <li class="nav-item">
                                <a id="head3" class="nav-link" data-target="#custom">{{'custom setting' | i}}</a>
                            </li>
                        </ul>
                    </div>
                </div>
                <div class="col-9">
                <form>
                <div id="guest" class="d-none">
                    <div class="form-group">
                        <label for="name" class="col-form-label-sm">{{'guest name' | i}}</label>
                        <div class="input-group">
                            <input type="text" class="form-control form-control-sm" name="name" value="guest.01"/>
                        </div>
                    </div>
                    <div class="form-group">
                        <label for="family" class="col-form-label-sm">{{'operating system' | i}}</label>
                        <div class="input-group">
                            <select class="select-lg" name="family">
                                <option value="linux" selected>linux</option>
                                <option value="windows">windows</option>
                                <option value="other">other</option>
                            </select>
                        </div>
                    </div>
                </div>
                <div id="datastore" class="d-none">
                    <div class="form-group">
                        <label for="datastore" class="col-form-label-sm">{{'datastore location' | i}}</label>
                        <div class="input-group">
                            <select class="select-lg" name="datastore">
                                <option value="datastore/01" selected>datastore01</option>
                                <option value="datastore/02">datastore02</option>
                            </select>
                        </div>
                    </div>
                </div>
                <div id="custom" class="d-none">
                    <div class="form-group">
                        <label for="cpu" class="col-form-label-sm">{{'processors' | i}}</label>
                        <div class="input-group">
                            <select class="form-control form-control-sm" name="cpuMode">
                                <option value="" selected>Default</option>
                                <option value="host-passthrough">Intel VT-x or AMD-V</option>
                            </select>
                            <select class="select-twice-md" name="cpu">
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
                            <input type="text" class="form-control form-control-sm"
                                   name="memSize" value="2048"/>
                            <select class="select-unit-right" name="memUnit">
                                <option value="MiB" selected>MiB</option>
                                <option value="GiB">GiB</option>
                            </select>
                        </div>
                    </div>
                    <div class="form-group">
                        <label for="diskSize" class="col-form-label-sm">{{'hardware disk' | i}}</label>
                        <div class="input-group">
                            <input type="text" class="form-control form-control-sm"
                                   name="disk1Size" value="10"/>
                            <select class="select-unit-right" name="disk1Unit">
                                <option value="Mib">MiB</option>
                                <option value="GiB" selected>GiB</option>
                                <option value="TiB">TiB</option>
                            </select>
                        </div>
                    </div>
                    <div class="form-group">
                        <label for="disk0File" class="col-form-label-sm">{{'select ISO or image' | i}}</label>
                        <div class="input-group">
                            <select class="form-control form-control-sm" name="disk0File">
                                <option value="/dev/sr0">sr0</option>
                            </select>
                        </div>
                    </div>
                    <div class="form-group">
                        <label for="interface0Source" class="col-form-label-sm">{{'network interface' | i}}</label>
                        <div class="input-group">
                            <select class="form-control form-control-sm" name="interface0Source">
                                <option value="virbr0" selected>Linux Bridge #virbr0</option>
                                <option value="virbr1">Linux Bridge #virbr1</option>
                                <option value="virbr2">Linux Bridge #virbr2</option>
                                <option value="virbr3">Linux Bridge #virbr3</option>
                            </select>
                        </div>
                        </div>
                    </div>
                </div>
                </form>
                </div>
            </div>
            <div class="modal-footer text-right">
                <button id="btn-prev" class="btn btn-outline-dark btn-sm">{{'previous' | i}}</button>
                <button id="btn-next" class="btn btn-outline-info btn-sm">{{'next' | i}}</button>
                <button id="btn-cancel" class="btn btn-outline-dark btn-sm">{{'cancel' | i}}</button>
                <button id="btn-submit" class="btn btn-outline-success btn-sm">{{'submit' | i}}</button>
            </div>
        </div>
        </div>`);
    }
}
