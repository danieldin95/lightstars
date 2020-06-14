import {Ctl} from "./ctl.js"
import {DataStoreApi} from "../api/datastore.js";
import {DataStoreTable} from "../widget/datastore/table.js";
import {FileUpload} from "../widget/datastore/upload.js";
import {UploadApi} from "../api/upload.js";
import {CheckBox} from "../widget/checkbox/checkbox.js";


class CheckBoxCtl extends CheckBox {
    change(from) {
        super.change(from);
        if (from.store.length !== 1) {
            $(this.child('#upload')).addClass('disabled');
        } else {
            $(this.child('#upload')).removeClass('disabled');
        }
    }
}


export class DataStoresCtl extends Ctl {
    // {
    //   id: "#datastores"
    // }
    constructor(props) {
        super(props);
        this.checkbox = new CheckBoxCtl(props);
        this.uuids = this.checkbox.uuids;
        this.table = new DataStoreTable({id: this.child('#display-table')});
        this.upload = new FileUpload({id: props.upload});

        this.upload.onsubmit(this.uuids, function (e) {
            new UploadApi({uuids: e.data.store, id: '#Process'}).upload(e.form);
        });
        // register buttons's  click.
        $(this.child('#delete')).on("click", this.uuids, function (e) {
            new DataStoreApi({uuids: e.data.store}).delete();
        });

        // refresh table and register refresh click.
        $(this.child('#refresh')).on("click", (e) => {
            this.table.refresh((e) => {
                this.checkbox.refresh();
            });
        });
        this.table.refresh((e) => {
            this.checkbox.refresh();
        });
    }

    create(data) {
        new DataStoreApi().create(data);
    }
}
