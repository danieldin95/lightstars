import {Controller} from "./controller.js";
import {InstanceApi} from "../api/instance.js";
import {DiskCtl} from "./disk.js";
import {InterfaceCtl} from "./interface.js"
import {GraphicsCtl} from "./graphics.js";
import {SnapshotCtl} from "./snapshot.js";


class HeaderCtl extends Controller {
    // {
    //   id: "#xx"
    // }
    constructor(props) {
        super(props);
        this.api = new InstanceApi({uuids: props.uuid});

        this.console();
        // register buttons's click.
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
            e.data.api.suspend();
        });
        $(this.child('#resume')).on("click", this, function (e) {
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

    console() {
        let url = $(this.child('#console')).attr('data');
        let view = $(this.child('#console')).attr('data-target');
        $(view).on('show.bs.modal', function (e) {
            $(this).find(".modal-body").html(
                `<iframe width="800px" height="600px" src="${url}" frameborder="0"></iframe>`
            );
        });
    }
}

export class GuestCtl extends Controller {
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
    //   data: {
    //   }
    // }
    constructor(props) {
        super(props);
        let name = props.name;
        let uuid = props.uuid;
        this.uuid = uuid;
        this.name = name;
        this.cpu = props.data.maxCpu || 0;
        this.mem = props.data.maxMem || 0;
        this.tasks = props.tasks || "tasks";
        this.api = new InstanceApi({uuids: uuid});
        this.header = new HeaderCtl({...props.header, uuid, name});
        this.disk = new DiskCtl({...props.disks, uuid, name});
        this.interface = new InterfaceCtl({...props.interfaces, uuid, name});
        this.graphics = new GraphicsCtl({...props.graphics, uuid, name});
        this.snapshot = new SnapshotCtl({...props.snapshot, uuid, name})
    }

    edit(data) {
        this.api.edit(data);
    }

    remove() {
        this.api.remove();
    }

    title(data) {
        this.api.title(data);
    }
}
