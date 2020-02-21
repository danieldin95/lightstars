
export class InstanceCreateModal {
    constructor () {
        let view = this.render();
        $("#instanceCreateModal").html(view);
        this.load();
    }

    render() {
        return `
        <div class="modal-dialog modal-lg modal-dialog-centered" role="document">
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
                                        <input type="text" class="form-control form-control-sm" name="name" value="centos.xx">
                                    </div>
                                </div>
                            </form>
                        </div>
                        <div id="step-storage" class="">
                            <form name="storage-config">
                                <div class="form-group row">
                                    <label for="storage" class="col-sm-4 col-form-label-sm">Storage</label>
                                    <div class="col-sm-6">
                                        <select class="custom-select input-select-sm" name="datastore">
                                            <option value="datastore01" selected>datastore01</option>
                                            <option value="datastore02">datastore02</option>
                                        </select>
                                    </div>
                                </div>
                            </form>
                        </div>
                        <div id="step-custom" class="">
                            <form name="custom-config">
                                <div class="form-group row">
                                    <label for="staticEmail" class="col-sm-4 col-form-label-sm">CPU</label>
                                    <div class="col-sm-6">
                                        <input type="text" class="form-control form-control-sm input-number" name="cpuNum" value="2">
                                    </div>
                                </div>
                                <div class="form-group row">
                                    <label for="staticEmail" class="col-sm-4 col-form-label-sm">Disk</label>
                                    <div class="col-sm-6">
                                        <div class="input-group">
                                            <input type="text" class="form-control form-control-sm input-number-lg" name="diskSize" value="10">
                                            <select class="" name="diskUnit">
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
                                        <input type="text" class="form-control form-control-sm" name="isoFile" value="/dev/sr0">
                                    </div>
                                </div>
                                <div class="form-group row">
                                    <label for="inputPassword" class="col-sm-4 col-form-label-sm">Memory</label>
                                    <div class="col-sm-6">
                                        <div class="input-group">
                                            <input type="text" class="form-control form-control-sm input-number-lg" name="memorySize" value="2048">
                                            <select class="" name="memoryUnit">
                                                <option value="Mib" selected>MiB</option>
                                                <option value="GiB">GiB</option>
                                            </select>       
                                        </div>
                                    </div>
                                </div>
                                <div class="form-group row">
                                    <label for="staticEmail" class="col-sm-4 col-form-label-sm">Bridge interface</label>
                                    <div class="col-sm-6">
                                        <input type="text" class="form-control form-control-sm input-string-sm" name="bridge" value="virbr0">
                                    </div>
                                </div>
                            </form>
                        </div>
                    </div>
                </div>
            </div>
        </div>`
    }

    load() {
        // Step show event
        $("#instanceCreateWizard").on("showStep", function(e, anchorObject, stepNumber, stepDirection, stepPosition) {
            if(stepPosition === 'first') {
                $("#prev-btn").addClass('disabled');
            }else if(stepPosition === 'final'){
                $("#next-btn").addClass('disabled');
            }else{
                $("#prev-btn").removeClass('disabled');
                $("#next-btn").removeClass('disabled');
            }
        });
        // Toolbar extra buttons
        var btnFinish = $('<button></button>').text('Finish')
            .addClass('btn btn-outline-success btn-sm')
            .on('click', function() {
                instance.create("#instaceCreateWizard");
            });
        var btnCancel = $('<button></button>').text('Cancel')
            .addClass('btn btn-outline-dark btn-sm')
            .on('click', function(){ $("#instanceCreateModal").modal("hide"); });

        // Smart Wizard
        $('#instanceCreateWizard').smartWizard({
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