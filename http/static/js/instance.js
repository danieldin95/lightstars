import {DiskApi} from './api/disk.js';
import {InterfaceApi} from './api/interface.js';
import {InstanceApi} from './api/instance.js';
import {CheckBoxTop} from "./com/utils.js";


export class Instance {
    // nil
    constructor() {
        let name = $('instance').attr("name");
        let uuid = $('instance').attr("data");
        this.uuid = uuid;
        this.name = name;
        this.disk = new Disk(uuid, name);
        this.interface = new Interface(uuid, name);

        // register buttons's click.
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

        // register button's click.
        $("disk-remove").on("click", this, function (e) {
            new DiskApi({
                instance: e.data.instance,
                uuids: e.data.disks.store,
                name: e.data.name}).delete();
        });
    }

    create(data) {
        new DiskApi({instance: this.instance, name: this.name}).create(data);
    }

    edit(data) {
        new DiskApi({instance: this.instance, name: this.name}).edit(data);
    }
}


export class DiskOn {
    // nil
    constructor() {
        this.disks = {store: []};

        let record = this.disks;
        let change = this.change;

        new CheckBoxTop({
            one: "disk-on-one input",
            all: "disk-on-all input",
            change: function (e) {
                change(record, e);
            },
        });

        // disabled firstly.
        change(record, this.disks);
    }

    change(record, from) {
        record.store = from.store;

        if (from.store.length == 0) {
            $("disk-remove button").addClass('disabled');
        } else {
            $("disk-remove button").removeClass('disabled');
        }
        if (from.store.length != 1) {
            $("disk-edit button").addClass('disabled');
        }
        else {
            $("disk-edit button").removeClass('disabled');
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

        // register buttons's click
        $("interface-remove").on("click", this, function (e) {
            new InterfaceApi({
                instance: e.data.instance,
                uuids: e.data.interfaces.store,
                name: this.name}).delete();
        });
    }

    create(data) {
        new InterfaceApi({instance: this.instance, name: this.name}).create(data);
    }

    edit(data) {
        new InterfaceApi({instance: this.instance, name: this.name}).edit(data);
    }
}


export class InterfaceOn {
    // nil
    constructor() {
        this.interfaces = {store: []};

        let record = this.interfaces;
        let change = this.change;

        new CheckBoxTop({
            one: "interface-on-one input",
            all: "interface-on-all input",
            change: function (e) {
                change(record, e);
            }
        });

        // disabled firstly.
        change(record, this.interfaces);
    }

    change(record, from) {
        record.store = from.store;

        if (from.store.length == 0) {
            $("interface-remove button").addClass('disabled');
        } else {

            $("interface-remove button").removeClass('disabled');
        }
        if (from.store.length != 1) {
            $("interface-edit button").addClass('disabled');
        }
        else {
            $("interface-edit button").removeClass('disabled');
        }
    }
}