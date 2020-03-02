import {Option} from "../option.js";
import {AlertDanger} from "../../com/alert.js";
import {FormModal} from "../form/modal.js";


export class InstanceCreate extends FormModal {
    // {
    //   id: "#instanceCrateModal",
    //   wizardId: "",
    // }
    constructor (props) {
        super(props);

        this.props.wizardId = props.wizardId || 'instanceCreateWizard';
        this.wizardId = this.props.wizardId; // id
        console.log('InstanceCreate.constructor', this.props, this.wizardId);

        this.forms = `#${this.wizardId} form`;
        this.prevbtn = `#${this.wizardId} #prev-btn`;
        this.nextbtn = `#${this.wizardId} #next-btn`;
        console.log('forms', this.forms, this.prevbtn, this.nextbtn);

        this.render();
        this.loading();

        $(this.id).on('show.bs.modal', this, function (e) {
            e.data.fetch();
        });
    }

    fetch() {
        let iso = {
            selector: this.view.find("select[name='disk0File']"),
            fresh: function (datastore) {
                let selector = this.selector;

                $.getJSON("/api/iso", {datastore: datastore}, function (data) {
                    selector.find("option").remove();
                    data.forEach(function (ele, index) {
                        selector.append(Option(ele['path'], ele['path']));
                    });
                    selector.append(Option('CDROM device:/sr0', '/dev/sr0'));
                }).fail(function (e) {
                    $("tasks").append(AlertDanger((`${this.type} ${this.url}: ${e.responseText}`)));
                });
            },
        };

        let store = {
            selector: this.view.find("select[name='datastore']"),
            refresh: function () {
                let selector = this.selector;

                $.getJSON("/api/datastore", function (data) {
                    selector.find("option").remove();
                    data.forEach(function (ele, index) {
                        selector.append(Option(ele['name'], ele['path']));
                    });
                    if (data.length > 0) {
                        iso.fresh(data[0]['name']);
                    }
                }).fail(function (e) {
                    $("tasks").append(AlertDanger((`${this.type} ${this.url}: ${e.responseText}`)));
                });
            },
        };

        let iface = {
            fresh: function (){
                let selector = this.selector;

                $.getJSON("/api/bridge", function (data) {
                    selector.find("option").remove();
                    data.forEach(function (e, i) {
                        if (e['type'] == 'bridge') {
                            selector.append(Option(`Linux Bridge #${e['name']}`, e['name']));
                        } else if (e['type'] == 'openvswith') {
                            selector.append(Option(`Open vSwitch #${e['name']}`, e['name']));
                        }
                    });
                }).fail(function (e) {
                    $("tasks").append(AlertDanger((`${this.type} ${this.url}: ${e.responseText}`)));
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
        this.view = $(this.template());

        this.view.find("select[name='cpu'] option").remove();
        for (let i = 1; i < 17; i++) {
            this.view.find("select[name='cpu']").append(new Option(i, i));
        }
        this.container().html(this.view);
    }


    wizard() {
        return $(`#${this.wizardId}`);
    }

    loading() {
        let prevbtn = this.prevbtn;
        let nextbtn = this.nextbtn;

        // Step show event
        this.wizard().on("showStep", function(e, anchorObject, stepNumber, stepDirection, stepPosition) {
            if (stepPosition === 'first') {
                $(prevbtn).addClass('disabled');
            } else if (stepPosition === 'final') {
                $(nextbtn).addClass('disabled');
            } else {
                $(prevbtn).removeClass('disabled');
                $(nextbtn).removeClass('disabled');
            }
        });

        // Toolbar extra buttons
        let btnFinish = $('<button name="finish-btn"></button>').text('Finish')
            .addClass('btn btn-outline-success btn-sm');
        let btnCancel = $('<button name="cancel-btn"></button>').text('Cancel')
            .addClass('btn btn-outline-dark btn-sm');

        // Smart wizard
        this.wizard().smartWizard({
            selected: 0,
            theme: 'dots',
            transitionEffect: 'fade',
            showStepURLhash: false,
            autoAdjustHeight: false,
            toolbarSettings: {
                toolbarPosition: 'bottom',
                toolbarExtraButtons: [btnFinish, btnCancel],
            }
        });

        // loading super for finish and cancel buttons.
        super.loading();
    }

    template(props) {
        return (`
        <div class="modal-dialog modal-lg modal-dialog-centered" role="document">
            <div class="modal-content">
                <div class="modal-header">
                    <h5 class="modal-title" id="">Create Instance</h5>
                </div>
                <div id="${this.wizardId}" class="modal-body">
                    <!-- Wizard navigations -->
                    <ul class="wizard-navs">
                        <li>
                            <a href="#step-os">Select name<br />
                                <small>Configure name and guest OS</small>
                            </a>
                        </li>
                        <li>
                            <a href="#step-storage">Select Datastore<br />
                                <small>Select one datastore for storage</small>
                            </a>
                        </li>
                        <li>
                            <a href="#step-custom">Custom configuration<br />
                                <small>Configure VM's hardware disk, interface and others</small>
                            </a>
                        </li>
                    </ul>
                    <!-- Wizard content -->
                    <div class="wizard-main">
                    <!-- Gust OS -->
                    <div id="step-os" class="">
                        <form name="os-config">
                        <div class="form-group row">
                            <label for="name" class="col-sm-4 col-form-label-sm">Name</label>
                            <div class="col-sm-6">
                                <div class="input-group">
                                    <input type="text" class="form-control form-control-sm"
                                           name="name" value="guest.01"/>
                                </div>
                            </div>
                        </div>
                        <div class="form-group row">
                            <label for="family" class="col-sm-4 col-form-label-sm">Guest OS</label>
                            <div class="col-sm-6">
                                <div class="input-group">
                                    <select class="select-lg" name="family">
                                        <option value="linux" selected>Linux</option>
                                        <option value="windows">Windows</option>
                                        <option value="other">Other</option>
                                    </select>
                                </div>
                            </div>
                        </div>
                        </form>
                    </div>
                    <!-- DataStore -->
                    <div id="step-storage" class="">
                        <form name="storage-config">
                        <div class="form-group row">
                            <label for="datastore" class="col-sm-4 col-form-label-sm">
                                Datastore
                            </label>
                            <div class="col-sm-6">
                                <div class="input-group">
                                    <select class="select-lg" name="datastore">
                                        <option value="datastore/01" selected>datastore01</option>
                                        <option value="datastore/02">datastore02</option>
                                    </select>
                                </div>
                            </div>
                        </div>
                        </form>
                    </div>
                    <!-- Custom instance -->
                    <div id="step-custom" class="">
                    <form name="custom-config">
                        <div class="form-group row">
                            <label for="cpu" class="col-sm-2 col-md-4 col-form-label-sm">
                                Processors
                            </label>
                            <div class="col-sm-10 col-md-6">
                                <div class="input-group">
                                    <select class="" name="cpu">
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
                            <label for="MaxMem" class="col-sm-4 col-md-4 col-form-label-sm">Memory</label>
                            <div class="col-sm-10 col-md-6">
                                <div class="input-group">
                                    <input type="text" class="form-control form-control-sm input-number-lg"
                                           name="memSize" value="2048"/>
                                    <select class="select-unit-right" name="memUnit">
                                        <option value="Mib" selected>MiB</option>
                                        <option value="GiB">GiB</option>
                                    </select>
                                </div>
                            </div>
                        </div>
                        <div class="form-group row">
                            <label for="diskSize" class="col-sm-4 col-md-4 col-form-label-sm">
                                Hardware disk01
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
                                Datastore ISO file
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
                                Network interface01
                            </label>
                            <div class="col-sm-10 col-md-6">
                                <div class="input-group">
                                    <select class="" name="interface0Source">
                                        <option value="virbr0" selected>Linux Bridge #virbr0</option>
                                        <option value="virbr1">Linux Bridge #virbr1</option>
                                        <option value="virbr2">Linux Bridge #virbr2</option>
                                        <option value="virbr3">Linux Bridge #virbr3</option>
                                    </select>
                                </div>
                            </div>
                        </div>
                        </form>
                    </div>
                    </div>
                </div>
            </div>
        </div>`);
    }
}