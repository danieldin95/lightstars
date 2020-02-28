import {DiskApi} from './api/disk.js';
import {InterfaceApi} from './api/interface.js';
import {InstanceApi} from './api/instance.js';
import {CheckBoxTop} from "./com/utils.js";


export class Instance {
    // {
    //   id: '#instance'
    //   disks: {
    //     id: '#disks'
    //   },
    //   interfaces: {
    //     id: "#interfaces"
    //   },
    // }
    constructor(props) {
        this.id = props.id;
        this.props = props;
        let name = $(this.id).attr("name");
        let uuid = $(this.id).attr("data");
        this.uuid = uuid;
        this.name = name;

        this.disk = new Disk({id: props.disks.id, uuid, name});
        this.interface = new Interface({id: props.interfaces.id, uuid, name});

        // register buttons's click.
        $(`${this.id} #start, ${this.id} #more-start`).on("click", this, function (e) {
            new InstanceApi({uuids: uuid}).start();
        });
        $(`${this.id} #shutdown`).on("click", this, function (e) {
            new InstanceApi({uuids: uuid}).shutdown();
        });
        $(`${this.id} #reset`).on("click", this, function (e) {
            new InstanceApi({uuids: uuid}).reset();
        });
        $(`${this.id} #suspend`).on("click", this, function (e) {
            new InstanceApi({uuids: uuid}).suspend();
        });
        $(`${this.id} #resume`).on("click", this, function (e) {
            new InstanceApi({uuids: uuid}).resume();
        });
        $(`${this.id} #destroy`).on("click", this, function (e) {
            new InstanceApi({uuids: uuid}).destroy();
        });
        $(`${this.id} #remove`).on("click", this, function (e) {
            new InstanceApi({uuids: uuid}).remove();
        });
    }
}


export class Disk {
    // {
    //   id: '#instance #disk',
    //   uuid: uuid of instance,
    //   name: name of instance,
    // }
    constructor(props) {
        this.id = props.id;
        this.name = props.name;
        this.instance = props.uuid;

        this.diskOn = new DiskOn(props);
        this.disks = this.diskOn.disks;

        // register button's click.
        $(`${this.id} #remove`).on("click", this, function (e) {
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
    // {
    //   id: '#instane #disk'
    // }
    constructor(props) {
        this.id = props.id;
        this.disks = {store: [], id: this.id};

        let record = this.disks;
        let change = this.change;

        new CheckBoxTop({
            one: `${this.id} #on-one`,
            all: `${this.id} #on-all`,
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
            $(`${record.id} #remove`).addClass('disabled');
        } else {
            $(`${record.id} #remove`).removeClass('disabled');
        }
        if (from.store.length != 1) {
            $(`${record.id} #edit`).addClass('disabled');
        } else {
            $(`${record.id} #edit`).removeClass('disabled');
        }
    }
}


export class Interface {
    // {
    //   id: '#instance #interface',
    //   uuid: uuid of instance,
    //   name: name of instance,
    // }
    constructor(props) {
        this.id = props.id;
        this.props = props;
        this.name = props.name;
        this.instance = props.uuid;

        this.interfaceOn = new InterfaceOn(props);
        this.interfaces = this.interfaceOn.interfaces;

        // register buttons's click
        $(`${this.id} #remove`).on("click", this, function (e) {
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
    // {
    //   id: '#instane #disk'
    // }
    constructor(props) {
        this.id = props.id;
        this.props = props;
        this.interfaces = {store: [], id: this.id};

        let record = this.interfaces;
        let change = this.change;

        new CheckBoxTop({
            one: `${this.id} #on-one`,
            all: `${this.id} #on-all`,
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
            $(`${record.id} #remove`).addClass('disabled');
        } else {
            $(`${record.id} #remove`).removeClass('disabled');
        }
        if (from.store.length != 1) {
            $(`${record.id} #edit`).addClass('disabled');
        }  else {
            $(`${record.id} #edit`).removeClass('disabled');
        }
    }
}