import {FormModal} from "../form/modal.js";
import {Option} from "../option.js";


export class NATCreate extends FormModal {
    //
    constructor (props) {
        super(props);

        this.render();
        this.loading();
    }

    render() {
        super.render();
        let prefix = {
            fresh: function() {
                this.selector.find('option').remove();
                for (let i = 26; i >= 8; i--) {
                    let alias = "/"+i;
                    this.selector.append(new Option(alias, i));
                }
                this.selector.find('option[value=24]').prop('selected', true);
            },
            selector: this.view.find("select[name='prefix']"),
        };

        let name = {
            fresh: function() {
                this.selector.find('option').remove();
                for (let i = 0; i < 8; i++) {
                    let alias = "virbr"+i;
                    this.selector.append(new Option(alias, alias));
                }
            },
            selector: this.view.find("select[name='name']"),
        };

        prefix.fresh();
        name.fresh();
    }

    template() {
        return (`
        <div class="modal-dialog modal-dialog-centered model-md" role="document">
            <div class="modal-content">
                <div class="modal-header">
                    <h5 class="modal-title" id="">Create NAT based network</h5>
                </div>
                <form name="network-new">
                    <input type="text" class="d-none" name="mode" value="nat"/>
                    <div id="" class="modal-body">
                        <div class="form-group row">
                            <label for="name" class="col-sm-4 col-form-label-sm ">Bridge's name</label>
                            <div class="col-sm-6">
                                <div class="input-group">
                                    <select class="select-md" name="name">
                                        <option value="virbr0" selected>virbr0</option>
                                        <option value="virbr1">virbr1</option>
                                        <option value="virbr2">virbr2</option>
                                    </select>
                                </div>
                            </div>
                        </div>
                        <div class="form-group row">
                            <label for="address" class="col-sm-4 col-form-label-sm">Address and prefix</label>
                            <div class="col-sm-6">
                                <div class="input-group">
                                    <input type="text" class="form-control form-control-sm input-number-lg"
                                           name="address" value="192.168.10.1"/>
                                        <select class="select-unit-right" name="prefix">
                                            <option value="24" selected>/24</option>
                                        </select>
                                </div>
                            </div>
                        </div>
                        <div class="form-group row">
                            <label for="dhcp" class="col-sm-4 col-form-label-sm ">DHCP setting</label>
                            <div class="col-sm-6">
                                <div class="input-group">
                                    <select class="select-md" name="dhcp">
                                        <option value="yes" selected>enable</option>
                                        <option value="no">disable</option>
                                    </select>
                                </div>
                            </div>
                        </div>
                        <div class="form-group row">
                            <label for="range" class="col-sm-4 col-form-label-sm">Address range</label>
                            <div class="col-sm-6">
                                <div class="input-group">
                                    <textarea type="text" class="form-control form-control-sm" 
                                        name="range" rows="3">192.168.10.100,192.168.10.200</textarea>                                         
                                </div>
                            </div>
                        </div>
                    </div>
                    <div id="" class="modal-footer">
                        <div class="mr-0" rol="group">
                            <button name="finish-btn" class="btn btn-outline-success btn-sm">Finish</button>
                            <button name="reset-btn" class="btn btn-outline-dark btn-sm" type="reset">Reset</button>
                            <button name="cancel-btn" class="btn btn-outline-dark btn-sm">Cancel</button>
                        </div>
                    </div>
                </form>
            </div>
        </div>`);
    }
}