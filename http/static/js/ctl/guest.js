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
        this.header = props.header;
        this.tasks = props.tasks || "Tasks";

        this.disk = new DiskCtl({id: props.disks.id, uuid, name});
        this.interface = new InterfaceCtl({id: props.interfaces.id, uuid, name});
        this.graphics = new GraphicsCtl({id: props.graphics.id, uuid, name});

        // register buttons's click.
        $(this.headerId('#console')).on("click", this, function (e) {
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
        $(this.headerId('#start'), this.headerId('#more-start')).on("click", this, function (e) {
            new InstanceApi({uuids: uuid}).start();
        });
        $(this.headerId('#shutdown')).on("click", this, function (e) {
            new InstanceApi({uuids: uuid}).shutdown();
        });
        $(this.headerId('#reset')).on("click", this, function (e) {
            new InstanceApi({uuids: uuid}).reset();
        });
        $(this.headerId('#suspend')).on("click", this, function (e) {
            if ($(this).hasClass('disabled')) {
                return
            }
            new InstanceApi({uuids: uuid}).suspend();
        });
        $(this.headerId('#resume')).on("click", this, function (e) {
            if ($(this).hasClass('disabled')) {
                return
            }
            new InstanceApi({uuids: uuid}).resume();
        });
        $(this.headerId('#destroy')).on("click", this, function (e) {
            new InstanceApi({uuids: uuid}).destroy();
        });
        $(this.headerId('#remove')).on("click", this, function (e) {
            new InstanceApi({uuids: uuid}).remove();
        });

        // console
        $(this.headerId('#console-self')).on('click', this, function (e) {
            let url = $(this).attr('data');
            window.open(url, '_self');
        });
        $(this.headerId('#console-blank')).on('click', this, function (e) {
            let url = $(this).attr('data');
            window.open(url, '_blank');
        });
        $(this.headerId('#console-window')).on('click', this, function (e) {
            let url = $(this).attr('data');
            window.open(url, e.data.name,'width=800,height=600');
        });
    }

    edit(data) {
        new InstanceApi({uuids: this.uuid}).edit(data);
    }

    headerId(id) {
        return [this.id, this.header.id, id].join(" ")
    }
}
