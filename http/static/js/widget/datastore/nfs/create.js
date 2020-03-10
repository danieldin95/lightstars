import {FormModal} from "../../form/modal.js";
import {Option} from "../../option.js";
import {Utils} from "../../../com/utils.js";

export class NFSCreate extends FormModal {
    //
    constructor (props) {
        super(props);

        this.render();
        this.loading();
    }

    render() {
        super.render();
        let name = {
            fresh: function() {
                this.selector.find('option').remove();
                for (let i = 10; i <= 20; i++) {
                    let alias = "datastore@"+Utils.iton(i, 2);
                    this.selector.append(new Option(alias, alias));
                }
            },
            selector: this.view.find("select[name='name']"),
        };

        name.fresh();
    }

    template() {
        return (`
        <div class="modal-dialog modal-dialog-centered model-md" role="document">
            <div class="modal-content">
                <div class="modal-header">
                    <h5 class="modal-title" id="">New datastore by NFS</h5>
                </div>
                <form name="datastore-new">
                <input type="text" class="d-none" name="type" value="netfs"/>
                <input type="text" class="d-none" name="format" value="nfs"/>
                <div id="" class="modal-body">
                    <div class="form-group row">
                        <label for="name" class="col-sm-4 col-form-label-sm ">DataStore</label>
                        <div class="col-sm-6">
                            <div class="input-group">
                                <select class="select-lg" name="name">
                                    <option value="datastore@01" selected>datastore@01</option>
                                </select>
                            </div>
                        </div>
                    </div>
                    <div class="form-group row">
                        <label for="host" class="col-sm-4 col-form-label-sm">Host address</label>
                        <div class="col-sm-6">
                            <div class="input-group">
                                <input type="text" class="form-control form-control-sm" name="host" value=""/>                                         
                            </div>
                        </div>
                    </div>  
                    <div class="form-group row">
                        <label for="host" class="col-sm-4 col-form-label-sm">Directory path</label>
                        <div class="col-sm-6">
                            <div class="input-group">
                                <input type="text" class="form-control form-control-sm" name="path" value="/"/>                                         
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