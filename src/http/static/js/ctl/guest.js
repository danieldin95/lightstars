import {Ctl} from "./ctl.js";
import {InstanceApi} from "../api/instance.js";
import {DiskCtl} from "./disk.js";
import {InterfaceCtl} from "./interface.js"
import {GraphicsCtl} from "./graphics.js";


class HeaderCtl extends Ctl {
    // {
    //   id: "#xx"
    // }
    constructor(props) {
        super(props);
        this.api = new InstanceApi({uuids: props.uuid});

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
        $(this.child('#start')).on("click", this, (e) => {
            this.api.start();
        });
        $(this.child('#shutdown')).on("click", (e) => {
            this.api.shutdown();
        });
        $(this.child('#reset')).on("click", (e) => {
            this.api.reset();
        });
        $(this.child('#suspend')).on("click", this, function (e) {
            if ($(this).hasClass('disabled')) {
                return
            }
            e.data.api.suspend();
        });
        $(this.child('#resume')).on("click", this, function (e) {
            if ($(this).hasClass('disabled')) {
                return
            }
            e.data.api.resume();
        });
        $(this.child('#destroy')).on("click", (e) => {
            this.api.destroy();
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
}

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
        this.api = new InstanceApi({uuids: uuid});
        this.header = new HeaderCtl({id: props.header.id, uuid, name});
        this.disk = new DiskCtl({id: props.disks.id, uuid, name});
        this.interface = new InterfaceCtl({id: props.interfaces.id, uuid, name});
        this.graphics = new GraphicsCtl({id: props.graphics.id, uuid, name});
    }

    edit(data) {
        this.api.edit(data);
    }

    remove() {
        this.api.remove();
    }
}
