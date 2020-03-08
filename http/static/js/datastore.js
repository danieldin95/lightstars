import {DataStoreApi} from "./api/datastore.js";
import {CheckBoxTop} from "./com/utils.js";
import {DataStoreTable} from "./widget/datastore/table.js";
import {FileUpload} from "./widget/datastore/upload.js";
import {UploadApi} from "./api/upload.js";

export class DataStore {
    // {
    //   id: "#datastores"
    // }
    constructor(props) {
        this.id = props.id;
        this.props = props;
        this.checkbox = new Checkbox(props);
        this.uuids = this.checkbox.uuids;
        this.table = new DataStoreTable({id: `${this.id} #display-table`});
        this.upload = new FileUpload({id: props.upload});

        this.upload.onsubmit(this.uuids, function (e) {
            new UploadApi({uuids: e.data.store, id: '#process'}).upload(e.form);
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

export class Checkbox {
    // {
    //   id: "#datastores"
    // }
    constructor(props) {
        this.id = props.id;
        this.props = props;
        this.uuids = {store: [], id: this.id};

        this.top = new CheckBoxTop({
            one: `${this.id} #on-one`,
            all: `${this.id} #on-all`,
            change: (e) => {
                this.change(this.uuids, e);
            }
        });

        // disabled firstly.
        this.change(this.uuids, this.uuids);
    }

    refresh() {
        this.top.refresh();
    }

    change(record, from) {
        record.store = from.store;

        if (from.store.length === 0) {
            $(`${record.id} #edit`).addClass('disabled');
            $(`${record.id} #delete`).addClass('disabled');
        } else {
            $(`${record.id} #edit`).removeClass('disabled');
            $(`${record.id} #delete`).removeClass('disabled');
        }
        if (from.store.length !== 1) {
            $(`${record.id} #edit`).addClass('disabled');
            $(`${record.id} #upload`).addClass('disabled');
        } else {
            $(`${record.id} #edit`).removeClass('disabled');
            $(`${record.id} #upload`).removeClass('disabled');
        }
    }
}