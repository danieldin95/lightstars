import {FormModal} from "../form/modal.js";
import {Option} from "../option.js";
import {Utils} from "../../lib/utils.js";
import {DataStoreApi} from "../../api/datastores.js";


export class DirCreate extends FormModal {
    //
    constructor (props) {
        super(props);
        this.usedNames = new Set();

        this.render();
        this.loading();
    }

    setUsedNames(names) {
        this.usedNames = new Set(names || []);
        this.refreshNameOptions();
    }

    refreshUsedNames() {
        new DataStoreApi().list(this, (resp) => {
            let items = (resp && resp.resp && resp.resp.items) ? resp.resp.items : [];
            let names = items.map((item) => item.name).filter((name) => !!name);
            this.setUsedNames(names);
        });
    }

    refreshNameOptions() {
        let selector = this.view.find("select[name='name']");
        if (selector.length === 0) {
            return;
        }
        selector.find('option').remove();
        for (let i = 1; i <= 9; i++) {
            let alias = "datastore@" + Utils.a2n(i, 2);
            if (!this.usedNames.has(alias)) {
                selector.append(new Option(alias, alias));
            }
        }
        if (selector.find('option').length === 0) {
            selector.append(new Option('N/A', ''));
        }
    }

    render() {
        super.render();
        this.refreshNameOptions();
    }

    template() {
        return this.compile(`
        <div class="modal-dialog modal-dialog-centered model-md" role="document">
        <div class="modal-content">
            <form name="datastore-new">
            <div class="modal-header">
                <h7 class="modal-title" id="">{{'new a datastore' | i}}</h7>
            </div>
            <div id="" class="modal-body">
                <input type="text" class="d-none" name="type" value="dir"/>
                <div class="form-group">
                    <label for="name" class="col-form-label-sm ">{{'select datastore' | i}}</label>
                    <div class="input-group">
                        <select class="select-lg" id="name" name="name">
                            <option value="datastore@01" selected>datastore@01</option>
                        </select>
                    </div>
                </div>
            </div>
            <div id="" class="modal-footer">
                <button type="button" name="cancel-btn" data-dismiss="modal" class="btn btn-outline-dark btn-sm">{{'cancel' | i}}</button>
                <button type="button" name="finish-btn" class="btn btn-outline-success btn-sm">{{'finish' | i}}</button>
            </div>
            </form>
        </div>
        </div>`);
    }
}
