import {FormModal} from "../form/modal.js";
import {Option} from "../option.js";
import {Utils} from "../../com/utils.js";

export class DirCreate extends FormModal {
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
                for (let i = 1; i <= 9; i++) {
                    let alias = "datastore@"+Utils.a2n(i, 2);
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
                <h5 class="modal-title" id="">New a datastore</h5>
            </div>
            <div id="" class="modal-body">
            <form name="datastore-new">
                <input type="text" class="d-none" name="type" value="dir"/>
                <div class="form-group">
                    <label for="name" class="col-form-label-sm ">Select datastore</label>
                    <div class="input-group">
                        <select class="select-lg" name="name">
                            <option value="datastore@01" selected>datastore@01</option>
                        </select>
                    </div>
                </div>
            </div>
            <div id="" class="modal-footer">
                <button name="cancel-btn" class="btn btn-outline-dark btn-sm">Cancel</button>
                <button name="finish-btn" class="btn btn-outline-success btn-sm">Finish</button>
            </div>
            </form>
        </div>
        </div>`);
    }
}
