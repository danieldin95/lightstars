import {DataStoreApi} from "../api/datastore.js";
import {DataStoreTable} from "../widget/datastore/table.js";
import {FileUpload} from "../widget/datastore/upload.js";
import {UploadApi} from "../api/upload.js";
import {CheckBoxTab} from "../widget/checkbox/checkbox.js";


export class CheckBox extends CheckBoxTab {
    change(from) {
        super.change(from);
        if (from.store.length !== 1) {
            $(`${this.uuids.id} #upload`).addClass('disabled');
        } else {
            $(`${this.uuids.id} #upload`).removeClass('disabled');
        }
    }
}


export class DataStoresCtl {
    // {
    //   id: "#datastores"
    // }
    constructor(props) {
        this.id = props.id;
        this.props = props;
        this.checkbox = new CheckBox(props);
        this.uuids = this.checkbox.uuids;
        this.table = new DataStoreTable({id: `${this.id} #display-table`});
        this.upload = new FileUpload({id: props.upload});

        this.upload.onsubmit(this.uuids, function (e) {
            new UploadApi({uuids: e.data.store, id: '#Process'}).upload(e.form);
        });
        // register buttons's  click.
        $(`${this.id} #delete`).on("click", this.uuids, function (e) {
            new DataStoreApi({uuids: e.data.store}).delete();
        });

        // refresh table and register refresh click.
        $(`${this.id} #refresh`).on("click", (e) => {
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
