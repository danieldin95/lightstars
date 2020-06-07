import {InstanceApi} from "../api/instance.js";
import {DiskCtl} from "./disk.js";
import {InterfaceCtl} from "./interface.js"
import {GraphicsCtl} from "./graphics.js";


export class GuestCtl {
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

        this.disk = new DiskCtl({id: props.disks.id, uuid, name});
        this.interface = new InterfaceCtl({id: props.interfaces.id, uuid, name});
        this.graphics = new GraphicsCtl({id: props.graphics.id, uuid, name});

        // register buttons's click.
        $(`${this.id} #console`).on("click", this, function (e) {
            if ($(this).hasClass('disabled')) {
                return
            }
            let url = $(this).attr('data');
            let target = $(this).attr('data-target');
            let iframe = `<iframe width="800px" height="600px" src="${url}" frameborder="0"></iframe>`;

            $(target).modal('show');
            $(`${target} .modal-body`).html(iframe);
            $(target).on('hidden.bs.modal', function (e) {
                $(target).find(".modal-body").empty();
            });
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
