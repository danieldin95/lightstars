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
        this.cpu = $("#instance").attr("cpu");
        this.mem = $("#instance").attr("mem");
        this.uuid = uuid;
        this.name = name;
        this.tasks = props.tasks || "tasks";

        this.disk = new Disk({id: props.disks.id, uuid, name});
        this.interface = new Interface({id: props.interfaces.id, uuid, name});

        // register buttons's click.
        $(`${this.id} #console`).on("click", this, function (e) {
            if ($(this).hasClass('disabled')) {
                return
            }
            let url = $(this).attr('data');
            let target = $(this).attr('data-target');

            $(target).modal('show');
            $(`${target} iframe`).attr("src", url);
            $(target).on('hidden.bs.modal', function (e) {
                $(target).find("iframe").removeAttr("src");
            });
        });
        $(`${this.id} #refresh`).on("click", this, function (e) {
            window.location.reload();
        });
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
            if ($(this).hasClass('disabled')) {
                return
            }
            new InstanceApi({uuids: uuid}).suspend();
        });
        $(`${this.id} #resume`).on("click", this, function (e) {
            if ($(this).hasClass('disabled')) {
                return
            }
            new InstanceApi({uuids: uuid}).resume();
        });
        $(`${this.id} #destroy`).on("click", this, function (e) {
            new InstanceApi({uuids: uuid}).destroy();
        });
        $(`${this.id} #remove`).on("click", this, function (e) {
            new InstanceApi({uuids: uuid}).remove();
        });


        // console
        $(`${this.id} #console-self`).on('click', this, function (e) {
            let url = $(this).attr('data');
            window.open(url, '_self');
        });
        $(`${this.id} #console-blank`).on('click', this, function (e) {
            let url = $(this).attr('data');
            window.open(url, '_blank');
        });
        $(`${this.id} #console-window`).on('click', this, function (e) {
            let url = $(this).attr('data');
            window.open(url, e.data.name,'width=873,height=655');
        });
        $(`${this.id} #console-spice`).on('click', this, function (e) {
            let url = $(this).attr('data');
            window.open(url, '_blank');
        });
    }

    edit(data) {
        new InstanceApi({uuids: this.uuid}).edit(data);
    }
}