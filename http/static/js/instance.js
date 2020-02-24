import {DiskApi} from './api/disk.js';
import {InterfaceApi} from './api/interface.js';
import {InstanceApi} from './api/instance.js';
import {ListenChangeAll} from "./com/utils.js";


export class Instance {
    // null
    constructor() {
        let uuid = $('instance').attr("data");
        this.uuid = uuid;

        this.disk = new Disk(uuid);
        this.interface = new Interface(uuid);

        // Register click handle.
        $("instance-start, instance-more-start").on("click", this, function (e) {
            new InstanceApi({uuids: uuid}).start();
        });
        $("instance-more-shutdown").on("click", this, function (e) {
            new InstanceApi({uuids: uuid}).shutdown();
        });
        $("instance-more-reset").on("click", this, function (e) {
            new InstanceApi({uuids: uuid}).reset();
        });
        $("instance-more-suspend").on("click", this, function (e) {
            new InstanceApi({uuids: uuid}).suspend();
        });
        $("instance-more-resume").on("click", this, function (e) {
            new InstanceApi({uuids: uuid}).resume();
        });
        $("instance-more-destroy").on("click", this, function (e) {
            new InstanceApi({uuids: uuid}).destroy();
        });
        $("instance-more-remove").on("click", this, function (e) {
            new InstanceApi({uuids: uuid}).remove();
        });
    }
}


export class Disk {
    //uuid
    constructor(instance, name) {
        this.name = name;
        this.instance = instance;

        this.diskOn = new DiskOn();
        this.disks = this.diskOn.disks;

        $("disk-remove").on("click", this, function (e) {
            new DiskApi({instance: e.data.instance, uuids: e.data.disks.store}).delete();
        });
    }

    create(data) {
        new DiskApi({instance: this.instance}).create(data);
    }

    edit(data) {
        new DiskApi({instance: this.instance}).edit(data);
    }
}


export class DiskOn {

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
    //
    constructor(instance, name) {
        this.name = name;
        this.instance = instance;

        this.interfaceOn = new InterfaceOn();
        this.interfaces = this.interfaceOn.interfaces;

        $("interface-remove").on("click", this, function (e) {
            new InterfaceApi({instance: e.data.instance, uuids: e.data.interfaces.store}).delete();
        });
    }

    create(data) {
        new InterfaceApi({instance: this.instance}).create(data);
    }

    edit(data) {
        new InterfaceApi({instance: this.instance}).edit(data);
    }
}


export class InterfaceOn {
    // nil
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