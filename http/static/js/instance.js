import {InstanceApi} from './api/instance.js'
import {ListenChangeAll} from "./com/utils.js";

export class Instance {

    constructor() {
        this.disk = new Disk();
        this.interface = new Interface();

        // Register click handle.
        $("instance-start, instance-more-start").on("click", this, function (e) {
            new InstanceApi($(this).attr("data")).start();
        });
        $("instance-more-shutdown").on("click", this, function (e) {
            new InstanceApi($(this).attr("data")).shutdown();
        });
        $("instance-more-reset").on("click", this, function (e) {
            new InstanceApi($(this).attr("data")).reset();
        });
        $("instance-more-suspend").on("click", this, function (e) {
            new InstanceApi($(this).attr("data")).suspend();
        });
        $("instance-more-resume").on("click", this, function (e) {
            new InstanceApi($(this).attr("data")).resume();
        });
        $("instance-more-destroy").on("click", this, function (e) {
            new InstanceApi($(this).attr("data")).destroy();
        });
        $("instance-more-remove").on("click", this, function (e) {
            new InstanceApi($(this).attr("data")).remove();
        });
    }
}


export class Disk {

    constructor() {
        this.disks = {store: []};

        let record = this.disks;
        let change = this.change;

        ListenChangeAll("disk-on-one input", "disk-on-all input", function (e) {
           change(record, e);
        });

        // Disabled firstly.
        change(record, this.disks);
    }

    change(record, from) {
        record.store = from.store;

        if (from.store.length == 0) {
            $("disk-edit button").addClass('disabled');
            $("disk-remove button").addClass('disabled');
        } else {
            $("disk-edit button").removeClass('disabled');
            $("disk-remove button").removeClass('disabled');
        }
    }
}

export class Interface {

    constructor() {
        this.interfaces = {store: []};

        let record = this.interfaces;
        let change = this.change;

        ListenChangeAll("interface-on-one input", "interface-on-all input", function (e) {
            change(record, e);
        });

        // Disabled firstly.
        change(record, this.interfaces);
    }

    change(record, from) {
        record.store = from.store;

        if (from.store == 0) {
            $("interface-edit button").addClass('disabled');
            $("interface-remove button").addClass('disabled');
        } else {
            $("interface-edit button").removeClass('disabled');
            $("interface-remove button").removeClass('disabled');
        }
    }
}