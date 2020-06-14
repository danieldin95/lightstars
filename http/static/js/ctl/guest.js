import {Ctl} from "./ctl.js";
import {InstanceApi} from "../api/instance.js";
import {DiskCtl} from "./disk.js";
import {InterfaceCtl} from "./interface.js"
import {GraphicsCtl} from "./graphics.js";


export class GuestCtl extends Ctl {
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
        super(props);
        let name = $(this.id).attr("name");
        let uuid = $(this.id).attr("data");
        this.cpu = $(this.id).attr("cpu");
        this.mem = $(this.id).attr("memory");
        this.uuid = uuid;
        this.name = name;
        this.tasks = props.tasks || "Tasks";

        this.disk = new DiskCtl({id: props.disks.id, uuid, name});
        this.interface = new InterfaceCtl({id: props.interfaces.id, uuid, name});
        this.graphics = new GraphicsCtl({id: props.graphics.id, uuid, name});

        // register buttons's click.
        $(this.child('#console')).on("click", this, function (e) {
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
        $(this.child('#start'), this.child('#more-start')).on("click", this, function (e) {
            new InstanceApi({uuids: uuid}).start();
        });
        $(this.child('#shutdown')).on("click", this, function (e) {
            new InstanceApi({uuids: uuid}).shutdown();
        });
        $(this.child('#reset')).on("click", this, function (e) {
            new InstanceApi({uuids: uuid}).reset();
        });
        $(this.child('#suspend')).on("click", this, function (e) {
            if ($(this).hasClass('disabled')) {
                return
            }
            new InstanceApi({uuids: uuid}).suspend();
        });
        $(this.child('#resume')).on("click", this, function (e) {
            if ($(this).hasClass('disabled')) {
                return
            }
            new InstanceApi({uuids: uuid}).resume();
        });
        $(this.child('#destroy')).on("click", this, function (e) {
            new InstanceApi({uuids: uuid}).destroy();
        });
        $(this.child('#remove')).on("click", this, function (e) {
            new InstanceApi({uuids: uuid}).remove();
        });

        // console
        $(this.child('#console-self')).on('click', this, function (e) {
            let url = $(this).attr('data');
            window.open(url, '_self');
        });
        $(this.child('#console-blank')).on('click', this, function (e) {
            let url = $(this).attr('data');
            window.open(url, '_blank');
        });
        $(this.child('#console-window')).on('click', this, function (e) {
            let url = $(this).attr('data');
            window.open(url, e.data.name,'width=800,height=600');
        });
    }

    edit(data) {
        new InstanceApi({uuids: this.uuid}).edit(data);
    }
}
