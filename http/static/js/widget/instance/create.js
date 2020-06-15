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
        return (`
        <div class="modal-dialog modal-dialog-centered modal-lg" role="document">
        <div class="modal-content ">
            <div class="modal-header">
                    Create a instance
            </div>
            <div class="modal-body">
            <div class="card">
                <div class="card-header">
                    <ul id="nav-tabs" class="nav nav-pills card-header-pills">
                        <li class="nav-item">
                            <a id="head1" class="nav-link text-center" data-target="#guest">
                                1.Configure Guest</a>
                        </li>
                        <li class="nav-item">
                            <a id="head2" class="nav-link text-center" data-target="#datastore">
                                2.Select DataStore</a>
                        </li>
                        <li class="nav-item">
                            <a id="head3" class="nav-link text-center" data-target="#custom">
                                3.Custom Setting</a>
                        </li>
                    </ul>
                </div>
                <form id="form">
                <div id="guest" class="card-body text-center d-none">
                    <div class="form-group row">
                        <label for="name" class="col-sm-4 col-md-4 col-form-label-sm">Name</label>
                        <div class="col-sm-10 col-md-6"">
                            <div class="input-group">
                                <input type="text" class="form-control form-control-sm"
                                       name="name" value="guest.01"/>
                            </div>
                        </div>
                    </div>
                    <div class="form-group row">
                        <label for="family" class="col-sm-4 col-md-4 col-form-label-sm">Guest OS</label>
                        <div class="col-sm-10 col-md-6">
                            <div class="input-group">
                                <select class="select-lg" name="family">
                                    <option value="linux" selected>Linux</option>
                                    <option value="windows">Windows</option>
                                    <option value="other">Other</option>
                                </select>
                            </div>
                        </div>
                    </div>
                </div>
                <div id="datastore" class="card-body text-center d-none">
                    <div class="form-group row">
                        <label for="datastore" class="col-sm-4 col-md-4 col-form-label-sm">
                            DataStore location
                        </label>
                        <div class="col-sm-10 col-md-6"">
                            <div class="input-group">
                                <select class="select-lg" name="datastore">
                                    <option value="datastore/01" selected>datastore01</option>
                                    <option value="datastore/02">datastore02</option>
                                </select>
                            </div>
                        </div>
                    </div>
                </div>
                <div id="custom" class="card-body text-center d-none">
                    <div class="form-group row">
                        <label for="cpu" class="col-sm-4 col-md-4 col-form-label-sm">
                            Processors
                        </label>
                        <div class="col-sm-10 col-md-6">
                            <div class="input-group">
                                <select class="select-sm" name="cpu">
                                    <option value="1">1</option>
                                    <option value="2" selected>2</option>
                                    <option value="3">3</option>
                                    <option value="4">4</option>
                                </select>
                                <select class="select-twice-md" name="cpuMode">
                                    <option value="" selected>Default</option>
                                    <option value="host-passthrough">Enable Intel VT-x or AMD-V</option>
                                </select>
                            </div>
                        </div>
                    </div>
                    <div class="form-group row">
                        <label for="MaxMem" class="col-sm-4 col-md-4 col-form-label-sm">Memory size</label>
                        <div class="col-sm-10 col-md-6">
                            <div class="input-group">
                                <input type="text" class="form-control form-control-sm input-number-lg"
                                       name="memSize" value="2048"/>
                                <select class="select-unit-right" name="memUnit">
                                    <option value="MiB" selected>MiB</option>
                                    <option value="GiB">GiB</option>
                                </select>
                            </div>
                        </div>
                    </div>
                    <div class="form-group row">
                        <label for="diskSize" class="col-sm-4 col-md-4 col-form-label-sm">
                            Hardware disk-01
                        </label>
                        <div class="col-sm-10 col-md-6">
                            <div class="input-group">
                                <input type="text" class="form-control form-control-sm input-number-lg"
                                       name="disk1Size" value="10"/>
                                <select class="select-unit-right" name="disk1Unit">
                                    <option value="Mib">MiB</option>
                                    <option value="GiB" selected>GiB</option>
                                    <option value="TiB">TiB</option>
                                </select>
                            </div>
                        </div>
                    </div>
                    <div class="form-group row">
                        <label for="disk0File" class="col-sm-4 col-md-4 col-form-label-sm">
                            DataStore ISO file
                        </label>
                        <div class="col-sm-10 col-md-6">
                            <div class="input-group">
                                <select class="" name="disk0File">
                                    <option value="/dev/sr0">sr0</option>
                                </select>
                            </div>
                        </div>
                    </div>
                    <div class="form-group row">
                        <label for="interface0Source" class="col-sm-4 col-md-4 col-form-label-sm">
                            Network interface-01
                        </label>
                        <div class="col-sm-10 col-md-6">
                            <div class="input-group">
                                <select class="select-md" name="interface0Source">
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
                <button id="btn-prev" class="btn btn-outline-dark btn-sm">Previous</button>
                <button id="btn-next" class="btn btn-outline-info btn-sm">Next</button>
                <button id="btn-cancel" class="btn btn-outline-dark btn-sm">Cancel</button>
                <button id="btn-submit" class="btn btn-outline-success btn-sm">Submit</button>
            </div>
        </div>
        </div>`);
    }
}
