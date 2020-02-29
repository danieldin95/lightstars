import {FormModal} from "../form/modal.js";
import {Option} from "../option.js";
import {Utils} from "../../com/utils.js";

export class DataStoreCreate extends FormModal {
    //
    constructor (props) {
        super(props);

        this.render();
        this.loading();
    }

    render() {
        this.view = $(this.template());
        this.container().html(this.view);

        let name = {
            fresh: function() {
                this.selector.find('option').remove();
                for (let i = 1; i <= 16; i++) {
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
                    <h5 class="modal-title" id="">Create DataStore</h5>
                </div>
                <form name="datastore-new">
                <div id="" class="modal-body">
                    <div class="form-group row">
                        <label for="name" class="col-sm-4 col-form-label-sm ">name</label>
                        <div class="col-sm-6">
                            <div class="input-group">
                                <select class="" name="name">
                                    <option value="datastore@01" selected>datastore@01</option>
                                </select>
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