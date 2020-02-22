import {Option} from "./option.js";
import {AlertDanger} from "../alert.js";

export class InstanceCreateModal {
    constructor () {
        this.view = this.render();

        this.view.find("select[name='cpu'] option").remove();
        for (let i = 1; i <= 16; i++) {
            this.view.find("select[name='cpu']").append(new Option(i, i));
        }
        this.container = $("#instanceCreateModal");
        this.container.html(this.view);

        this.wizard = $("#instanceCreateWizard");
        this.forms = $("#instanceCreateWizard form");
        this.prevbtn = "#instanceCreateWizard #prev-btn";
        this.nextbtn = "#instanceCreateWizard #next-btn";
        
        this.load();
        this.fetch();
        this.events = {
            submit: {
                func: function (e) {
                },
                data: undefined,
            }
        }
    }

    fetch() {
        let cpu_selector = this.view.find("select[name='isoFile']");
        let cpu_refresh = function(datastore) {
            $.getJSON("/api/iso", {datastore: datastore}, function (data) {
                cpu_selector.find("option").remove();
                data.forEach(function (ele, index) {
                    cpu_selector.append(Option(ele['name'], ele['path']));
                })
            }).fail(function (e) {
                $("errors").append(AlertDanger((`${this.type} ${this.url}: ${e.responseText}`)));
            });
        };

        let store_selector = this.view.find("select[name='datastore']");
        let store_refresh = function () {
            $.getJSON("/api/datastore", function (data) {
                store_selector.find("option").remove();
                data.forEach(function (ele, index) {
                    store_selector.append(Option(ele['name'], ele['path']));
                })
            }).fail(function (e) {
                $("errors").append(AlertDanger((`${this.type} ${this.url}: ${e.responseText}`)));
            });
        };

        store_refresh();
        cpu_refresh("datastore/01");

        store_selector.on("change", this, function (e) {
            cpu_refresh($(this).val());
        });
    }

    submit() {
        if (this.events.submit.func) {
            this.events.submit.func({
                data: this.events.submit.data,
                array: this.forms.serializeArray()
            });
        }
    }

    onSubmit(data, fn) {
        this.events.submit.data = data;
        this.events.submit.func = fn;
    }

    render() {
        return $(`
        <div class="modal-dialog modal-lg modal-dialog-centered sw-modal overflow-auto" role="document">
            <div class="modal-content">
                <div class="modal-header">
                    <h5 class="modal-title" id="instanceCreateModalLabel">Create Instance</h5>
                </div>
                <!-- Smart Wizard HTML -->
                <div id="instanceCreateWizard" class="modal-body">
                    <ul>
                        <li><a href="#step-os">Select name<br /><small>Configure name and guest OS</small></a></li>
                        <li><a href="#step-storage">Select storage<br /><small>Select one datastore for storage</small></a></li>
                        <li><a href="#step-custom">Custom configuration<br /><small>Configure VM's hardware disk, interface and others</small></a></li>
                    </ul>
                    <div>
                        <div id="step-os" class="">
                            <form name="os-config">
                                <div class="form-group row">
                                    <label for="name" class="col-sm-4 col-form-label-sm">Name</label>
                                    <div class="col-sm-6">
                                        <div class="input-group">
                                            <input type="text" class="form-control form-control-sm" name="name" value="centos.xx">
                                        </div>
                                    </div>
                                </div>
                            </form>
                        </div>
                        <div id="step-storage" class="">
                            <form name="storage-config">
                                <div class="form-group row">
                                    <label for="storage" class="col-sm-4 col-form-label-sm">Storage</label>
                                    <div class="col-sm-6">
                                        <div class="input-group">
                                            <select class="select-simple" name="datastore">
                                                <option value="datastore/01" selected>datastore01</option>
                                                <option value="datastore/02">datastore02</option>
                                            </select>
                                        </div>
                                    </div>
                                </div>
                            </form>
                        </div>
                        <div id="step-custom" class="">
                            <form name="custom-config">
                                <div class="form-group row">
                                    <label for="staticEmail" class="col-sm-4 col-form-label-sm">CPU</label>
                                    <div class="col-sm-6">
                                        <div class="input-group">
                                            <select class="select-simple select-unit" name="cpu">
                                                <option value="1">1</option>
                                                <option value="2" selected>2</option>
                                                <option value="3">3</option>
                                                <option value="4">4</option>
                                            </select>   
                                        </div>
                                    </div>
                                </div>
                                <div class="form-group row">
                                    <label for="inputPassword" class="col-sm-4 col-form-label-sm">Memory</label>
                                    <div class="col-sm-6">
                                        <div class="input-group">
                                            <input type="text" class="form-control form-control-sm input-number-lg" name="memorySize" value="2048">
                                            <select class="select-simple select-unit-right" name="memoryUnit">
                                                <option value="Mib" selected>MiB</option>
                                                <option value="GiB">GiB</option>
                                            </select>       
                                        </div>
                                    </div>
                                </div>                                
                                <div class="form-group row">
                                    <label for="staticEmail" class="col-sm-4 col-form-label-sm">Hardware disk01</label>
                                    <div class="col-sm-6">
                                        <div class="input-group">
                                            <input type="text" class="form-control form-control-sm input-number-lg" name="diskSize" value="10">
                                            <select class="select-simple select-unit-right" name="diskUnit">
                                                <option value="Mib">MiB</option>
                                                <option value="GiB" selected>GiB</option>
                                                <option value="TiB">TiB</option>
                                            </select>                                                                                     
                                        </div>
                                    </div>
                                </div>
                                <div class="form-group row">
                                    <label for="inputPassword" class="col-sm-4 col-form-label-sm">Datastore ISO file</label>
                                    <div class="col-sm-6">
                                        <div class="input-group">
                                            <select class="select-simple" name="isoFile">
                                                <option value="/dev/sr0">sr0</option>
                                            </select>   
                                        </div>
                                    </div>
                                </div>
                                <div class="form-group row">
                                    <label for="staticEmail" class="col-sm-4 col-form-label-sm">Network interface01</label>
                                    <div class="col-sm-6">
                                        <div class="input-group">
                                            <select class="select-simple" name="interface">
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
        </div>`)
    }

    load() {
        let prevbtn = this.prevbtn;
        let nextbtn = this.nextbtn;

        // Step show event
        this.wizard.on("showStep", function(e, anchorObject, stepNumber, stepDirection, stepPosition) {
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
        let btnFinish = $('<button id="finish-btn"></button>').text('Finish')
            .addClass('btn btn-outline-success btn-sm')
            .on('click', this, function(e) {
                e.data.submit();
                e.data.container.modal("hide");
            });
        let btnCancel = $('<button id="cancel-btn"></button>').text('Cancel')
            .addClass('btn btn-outline-dark btn-sm')
            .on('click', this, function(e) {
                e.data.container.modal("hide");
            });

        // Smart Wizard
        this.wizard.smartWizard({
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
    }
}