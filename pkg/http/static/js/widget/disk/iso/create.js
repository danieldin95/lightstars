import {FormModal} from "../../form/modal.js";
import {Option} from "../../option.js";
import {DataStoreApi} from "../../../api/datastores.js";
import {IsoApi} from "../../../api/iso.js";


export class IsoCreate extends FormModal {
    //
    constructor (props) {
        super(props);
        this.usedDriveSeqs = new Set();

        this.render();
        this.loading();
    }

    setUsedSeqs({driveSeqs}) {
        this.usedDriveSeqs = new Set(driveSeqs || []);
        this.render();
        this.loading();
    }

    render() {
        super.render();

        let seqSelector = this.view.find("select[name='seq']");
        seqSelector.find("option").remove();
        for (let i = 1; i < 16; i++) {
            if (!this.usedDriveSeqs.has(i)) {
                seqSelector.append(Option(i, i));
            }
        }

        let iso = {
            selector: this.view.find("select[name='source']"),
            fresh: function (datastore) {
                let selector = this.selector;

                new IsoApi().list(datastore, (data) => {
                    selector.find("option").remove();
                    for (let ele of data.resp) {
                        selector.append(Option(ele['path'], ele['path']));
                    }
                    selector.append(Option('CDROM device:/sr0', '/dev/sr0'));
                });
            },
        };

        let store = {
            selector: this.view.find("select[name='datastore']"),
            refresh: function () {
                let selector = this.selector;
                new DataStoreApi().list(this,  (data) => {
                    let resp = data.resp;
                    selector.find("option").remove();
                    for (let ele of resp.items) {
                        selector.append(Option(ele['name'], ele['name']));
                    }
                    if (resp.items.length > 0) {
                        iso.fresh(resp.items[0]['name']);
                    }
                });
            },
        };

        store.refresh();
        store.selector.on("change", this, function (e) {
            iso.fresh($(this).val());
        });
    }

    template() {
        return this.compile(`
        <div class="modal-dialog modal-dialog-centered model-md" role="document">
        <div class="modal-content">
            <div class="modal-header">
                <h7 class="modal-title" id="">{{'add cdrom' | i}}</h7>
            </div>
            <div id="" class="modal-body">
                <form name="disk-new">
                    <div class="form-group">
                        <label for="datastore" class="col-form-label-sm">{{'datastore location' | i}}</label>
                        <div class="input-group">
                            <select class="select-lg" id="datastore" name="datastore">
                                <option value="datastore/01" selected>datastore01</option>
                                <option value="datastore/02">datastore02</option>
                            </select>
                        </div>
                    </div>
                    <div class="form-group">
                        <label for="source" class="col-form-label-sm">{{'select ISO file' | i}}</label>
                        <div class="input-group">
                            <select class="form-control form-control-sm" id="source" name=source>
                                <option value="/dev/sr0">sr0</option>
                            </select>
                        </div>
                    </div>
                    <div class="form-group">
                        <label for="seq" class="col-form-label-sm ">{{'sequence number' | i}}</label>
                        <div class="input-group">
                            <select class="select-lg" id="seq" name="seq">
                                <option value="1" selected>1</option>
                                <option value="2">2</option>
                                <option value="3">3</option>
                                <option value="4">4</option>
                            </select>
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
