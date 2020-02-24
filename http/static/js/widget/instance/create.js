import {Option} from "../option.js";
import {AlertDanger} from "../../com/alert.js";
import {ModalFormBase} from "../form/modal.js";


export class InstanceCreate extends ModalFormBase {
    // {containerId: "", wizardId: ""}
    constructor (props) {
        super(props);

        this.props.wizardId = props.wizardId || 'instanceCreateWizard';
        this.wizardId = this.props.wizardId; // id
        console.log('aa', this.props, this.wizardId);

        this.forms = `#${this.wizardId} form`;
        this.prevbtn = `#${this.wizardId} #prev-btn`;
        this.nextbtn = `#${this.wizardId} #next-btn`;
        console.log('forms', this.forms, this.prevbtn, this.nextbtn);

        this.render();
        this.loading();
        this.fetch();
    }

    fetch() {
        let iso_sel = this.view.find("select[name='isoFile']");
        let iso_refresh = function(datastore) {
            $.getJSON("/api/iso", {datastore: datastore}, function (data) {
                iso_sel.find("option").remove();
                data.forEach(function (ele, index) {
                    iso_sel.append(Option(ele['path'], ele['path']));
                })
            }).fail(function (e) {
                $("tasks").append(AlertDanger((`${this.type} ${this.url}: ${e.responseText}`)));
            });
        };

        let store_sel = this.view.find("select[name='datastore']");
        let store_refresh = function () {
            $.getJSON("/api/datastore", function (data) {
                store_sel.find("option").remove();
                data.forEach(function (ele, index) {
                    store_sel.append(Option(ele['name'], ele['path']));
                });
                if (data.length > 0) {
                    iso_refresh(data[0]['name']);
                }
            }).fail(function (e) {
                $("tasks").append(AlertDanger((`${this.type} ${this.url}: ${e.responseText}`)));
            });
        };

        store_refresh();
        store_sel.on("change", this, function (e) {
            iso_refresh($(this).val());
        });
    }

    render() {
        this.view = $(this.template());

        this.view.find("select[name='cpu'] option").remove();
        for (let i = 1; i <= 16; i++) {
            this.view.find("select[name='cpu']").append(new Option(i, i));
        }
        this.container().html(this.view);
    }

    template(props) {
        return `
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
                        <a href="#step-storage">Select storage<br />
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
                                <input type="text" class="form-control form-control-sm" name="name" value="guest.01">
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
                        <label for="datastore" class="col-sm-4 col-form-label-sm">Storage</label>
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
                        <label for="cpu" class="col-sm-4 col-form-label-sm">CPU</label>
                        <div class="col-sm-6">
                            <div class="input-group">
                                <select class="" name="cpu">
                                    <option value="1">1</option>
                                    <option value="2" selected>2</option>
                                    <option value="3">3</option>
                                    <option value="4">4</option>
                                </select>   
                            </div>
                        </div>
                    </div>
                    <div class="form-group row">
                        <label for="memorySize" class="col-sm-4 col-form-label-sm">Memory</label>
                        <div class="col-sm-6">
                            <div class="input-group">
                                <input type="text" class="form-control form-control-sm input-number-lg" name="memorySize" value="2048">
                                <select class="select-unit-right" name="memoryUnit">
                                    <option value="Mib" selected>MiB</option>
                                    <option value="GiB">GiB</option>
                                </select>       
                            </div>
                        </div>
                    </div>                                
                    <div class="form-group row">
                        <label for="diskSize" class="col-sm-4 col-form-label-sm">Hardware disk01</label>
                        <div class="col-sm-6">
                            <div class="input-group">
                                <input type="text" class="form-control form-control-sm input-number-lg" name="diskSize" value="10">
                                <select class="select-unit-right" name="diskUnit">
                                    <option value="Mib">MiB</option>
                                    <option value="GiB" selected>GiB</option>
                                    <option value="TiB">TiB</option>
                                </select>                                                                                     
                            </div>
                        </div>
                    </div>
                    <div class="form-group row">
                        <label for="isoFile" class="col-sm-4 col-form-label-sm">Datastore ISO file</label>
                        <div class="col-sm-6">
                            <div class="input-group">
                                <select class="" name="isoFile">
                                    <option value="/dev/sr0">sr0</option>
                                </select>   
                            </div>
                        </div>
                    </div>
                    <div class="form-group row">
                        <label for="interface" class="col-sm-4 col-form-label-sm">Network interface01</label>
                        <div class="col-sm-6">
                            <div class="input-group">
                                <select class="" name="interface">
                                    <option value="virbr0" selected>Virtual Bridge0</option>
                                    <option value="virbr1">Virtual Bridge1</option>
                                    <option value="virbr2">Virtual Bridge2</option>
                                    <option value="virbr3">Virtual Bridge3</option>
                                </select>  
                            </div>
                        </div>
                    </div>
                    </form>
                </div>
                </div>
            </div>
        </div>
    </div>`
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
}