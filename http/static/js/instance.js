import {InstanceApi} from './api/instance.js';
import {Disk} from "./disk.js";
import {Interface} from "./interface.js"

export class Instance {
    // {
    //   id: '#instance'
    //   header: {
    //     id: '#'
    //  }
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

        this.head = props.header.id;
        // register buttons's click.
        $(`${this.head} #start, ${this.id} #more-start`).on("click", this, function (e) {
            new InstanceApi({uuids: uuid}).start();
        });
        $(`${this.head} #shutdown`).on("click", this, function (e) {
            new InstanceApi({uuids: uuid}).shutdown();
        });
        $(`${this.head} #reset`).on("click", this, function (e) {
            new InstanceApi({uuids: uuid}).reset();
        });
        $(`${this.head} #suspend`).on("click", this, function (e) {
            new InstanceApi({uuids: uuid}).suspend();
        });
        $(`${this.head} #resume`).on("click", this, function (e) {
            new InstanceApi({uuids: uuid}).resume();
        });
        $(`${this.head} #destroy`).on("click", this, function (e) {
            new InstanceApi({uuids: uuid}).destroy();
        });
        $(`${this.head} #remove`).on("click", this, function (e) {
            new InstanceApi({uuids: uuid}).remove();
        });
    }
}