import {DataStoreApi} from "./api/datastore.js";
import {CheckBoxTop} from "./com/utils.js";


export class DataStore {
    // {
    //   id: "#datastores"
    // }
    constructor(props) {
        this.id = props.id;
        this.props = props;
        this.datastoreOn = new DataStoreOn(props);
        this.datastores = this.datastoreOn.uuids;

        // register buttons's  click.
        $(`${this.id} #delete`).on("click", this.datastores, function (e) {
            new DataStoreApi({uuids: e.data.store}).delete();
        });
    }

    create(data) {
        new DataStoreApi().create(data);
    }
}


export class DataStoreOn {
    // {
    //   id: "#datastores"
    // }
    constructor(props) {
        this.id = props.id;
        this.props = props;
        this.uuids = {store: [], id: this.id};

        let change = this.change;
        let record = this.uuids;

        new CheckBoxTop({
            one: `${this.id} #on-one`,
            all: `${this.id} #on-all`,
            change: function (e) {
                change(record, e);
            }
        });

        // disabled firstly.
        change(record, this.uuids);
    }

    change(record, from) {
        record.store = from.store;
        console.log("DataStoreOn.change", record.store);

        if (from.store.length == 0) {
            $(`${record.id} #edit`).addClass('disabled');
            $(`${record.id} #delete`).addClass('disabled');
        } else {
            $(`${record.id} #edit`).removeClass('disabled');
            $(`${record.id} #delete`).removeClass('disabled');
        }
        if (from.store.length != 1) {
            $(`${record.id} #edit`).addClass('disabled');
        } else {
            $(`${record.id} #edit`).removeClass('disabled');
        }
    }
}