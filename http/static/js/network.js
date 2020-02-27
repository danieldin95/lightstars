import {NetworkApi} from "./api/network.js";
import {ListenChangeAll} from "./com/utils.js";


export class Network {
    // nil
    constructor() {
        this.networkOn = new NetworkOn();
        this.networks = this.networkOn.uuids;

        // register buttons's click.
        $("network-delete").on("click", this.networks, function (e) {
            new NetworkApi({uuids: e.data.store}).remove();
        });
    }

    create(data) {
        new NetworkApi().create(data);
    }
}


export class NetworkOn {

    constructor() {
        this.uuids = {store: []};

        let change = this.change;
        let record = this.uuids;

        ListenChangeAll("network-on-one input", "network-on-all input", function(e) {
           change(record, e);
        });

        // disabled firstly.
        change(record, this.uuids);
    }

    change(record, from) {
        record.store = from.store;
        console.log(record.store);

        if (from.store.length == 0) {
            $("network-edit button").addClass('disabled');
            $("network-delete button").addClass('disabled');
        } else {
            $("network-edit button").removeClass('disabled');
            $("network-delete button").removeClass('disabled');
        }

        if (from.store.length != 1) {
            $("network-edit button").addClass('disabled');
        }
        else {
            $("network-edit button").removeClass('disabled');
        }
    }
}