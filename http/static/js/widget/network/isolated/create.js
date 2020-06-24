import {FormModal} from "../../form/modal.js";
import {Option} from "../../option.js";


export class IsolatedCreate extends FormModal {
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
        prefix.fresh();
    }

    template() {
        return (`
        <div class="modal-dialog modal-dialog-centered model-md" role="document">
        <div class="modal-content">
            <div class="modal-header">
                <h5 class="modal-title" id="">Create isolated network</h5>
            </div>
            <form name="network-new">
                <input type="text" class="d-none" name="mode" value=""/>
                <div id="" class="modal-body">
                    <div class="form-group row">
                        <label for="name" class="col-sm-4 col-form-label-sm ">Name</label>
                        <div class="col-sm-6">
                            <div class="input-group">
                                <input type="text" class="form-control form-control-sm" name="name" value=""/>
                            </div>
                        </div>
                    </div>
                    <div class="form-group row">
                        <label for="address" class="col-sm-4 col-form-label-sm">IP Address</label>
                        <div class="col-sm-6">
                            <div class="input-group">
                                <input type="text" class="form-control form-control-sm input-number-lg"
                                       name="address" value="192.168.30.1"/>
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
                                <select class="select-lg" name="dhcp">
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
                                    name="range" rows="3">192.168.30.100,192.168.30.200</textarea>                                   
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
