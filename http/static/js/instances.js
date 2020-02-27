import {InstanceApi} from "./api/instance.js";
import {ListenChangeAll} from "./com/utils.js";


export class Instances {
    // nil
    constructor() {
        this.instanceOn = new InstanceOn();
        this.instances = this.instanceOn.uuids;

        // Register click handle.
        $("instance-console").on("click", this.instances, function (e) {
            let props = {uuids: e.data.store, passwd: {}};
            e.data.store.forEach(function (v) {
                props.passwd[v] = $(`input[data=${v}]`).attr('passwd');
            });
            new InstanceApi(props).console();
        });
        $("instance-start, instance-more-start").on("click", this.instances, function (e) {
            new InstanceApi({uuids: e.data.store}).start();
        });
        $("instance-more-shutdown").on("click", this.instances, function (e) {
            new InstanceApi({uuids: e.data.store}).shutdown();
        });
        $("instance-more-reset").on("click", this.instances, function (e) {
            new InstanceApi({uuids: e.data.store}).reset();
        });
        $("instance-more-suspend").on("click", this.instances, function (e) {
            new InstanceApi({uuids: e.data.store}).suspend();
        });
        $("instance-more-resume").on("click", this.instances, function (e) {
            new InstanceApi({uuids: e.data.store}).resume();
        });
        $("instance-more-destroy").on("click", this.instances, function (e) {
            new InstanceApi({uuids: e.data.store}).destroy();
        });
        $("instance-more-remove").on("click", this.instances, function (e) {
            new InstanceApi({uuids: e.data.store}).remove();
        });
    }

    create(data) {
        new InstanceApi().create(data);
    }
}


export class InstanceOn {

    constructor() {
        this.uuids = {store: []};

        let change = this.change;
        let record = this.uuids;

        ListenChangeAll("instance-on-one input", "instance-on-all input", function(e) {
           change(record, e);
        });

        // Disabled firstly.
        change(record, this.uuids);
    }

    change(record, from) {
        record.store = from.store;
        console.log(record.store);

        if (from.store.length == 0) {
            $("instance-start button").addClass('disabled');
            $("instance-console button").addClass('disabled');
            $("instance-shutdown button").addClass('disabled');
            $("instance-more button").addClass('disabled');
        } else {
            $("instance-start button").removeClass('disabled');
            $("instance-console button").removeClass('disabled');
            $("instance-shutdown button").removeClass('disabled');
            $("instance-more button").removeClass('disabled');
        }

        if (from.store.length != 1) {
            $("instance-edit button").addClass('disabled');
        }
        else {
            $("instance-edit button").removeClass('disabled');
        }
    }
}