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
                    let alias = "datastore@"+Utils.a2n(i, 2);
                    this.selector.append(new Option(alias, alias));
                }
            },
            selector: this.view.find("select[name='name']"),
        };

        name.fresh();
    }

    template() {
        return this.compile(`
        <div class="modal-dialog modal-dialog-centered model-md" role="document">
        <div class="modal-content">
            <div class="modal-header">
                <h7 class="modal-title" id="">{{'new datastore by nfs' | i}}</h7>
            </div>
            <div id="" class="modal-body">
            <form name="datastore-new">
                <input type="text" class="d-none" name="type" value="netfs"/>
                <input type="text" class="d-none" name="format" value="nfs"/>
                <div class="form-group">
                    <label for="name" class="col-form-label-sm ">{{'datastore' | i}}</label>
                    <div class="input-group">
                        <select class="select-lg" name="name">
                            <option value="datastore@01" selected>datastore@01</option>
                        </select>
                    </div>
                </div>
                <div class="form-group">
                    <label for="host" class="col-form-label-sm">{{'host address' | i}}</label>
                    <div class="input-group">
                        <input type="text" class="form-control form-control-sm" name="host" value=""/>                                         
                    </div>
                </div>  
                <div class="form-group">
                    <label for="host" class="col-form-label-sm">{{'remote directory' | i}}</label>
                    <div class="input-group">
                        <input type="text" class="form-control form-control-sm" name="path" value="/public"/>                                         
                    </div>
                </div>  
            </form>                                          
            </div>
            <div id="" class="modal-footer">
                <button name="cancel-btn" class="btn btn-outline-dark btn-sm">{{'cancel' | i}}</button>
                <button name="finish-btn" class="btn btn-outline-success btn-sm">{{'finish' | i}}</button>
            </div>
        </div>
        </div>`);
    }
}
