import {Controller} from "./controller.js"
import {LeasesCtl} from "./leases.js";
import {PortCtl} from "./port.js";
import {NetworkApi} from "../api/network.js";
import {ConfirmAction} from "../widget/common/confirm.js";

export class NetworkCtl extends Controller {
    // {
    //   id: '#network'
    //   header: {
    //     id: '#'
    //  }
    //   leases: {
    //     id: '#leases'
    //   },
    // }
    constructor(props) {
        super(props);
        let name = $(this.id).attr("name");
        let uuid = $(this.id).attr("data");
        this.uuid = uuid;
        this.name = name;
        this.tasks = props.tasks || "tasks";
        this.confirm = props.confirm;
        this.state = ($(this.id).attr("state") || "").toLowerCase();
        this.autostart = ($(this.id).attr("autostart") || "").toLowerCase() === "true";

        this.leases = new LeasesCtl({...props.leases, uuid, name});
        this.ports = new PortCtl({...props.ports, uuid, name, confirm: this.confirm});

        let api = new NetworkApi({tasks: this.tasks});
        let root = this.child('#header');
        let autoBtn = $(root + " #autostart");
        let autoEnabled = this.autostart;

        $(root + " #destroy").on("click", () => {
            new ConfirmAction({
                id: this.confirm,
                action: "destroy",
                name: this.name,
                message: "destroy",
            }).onsubmit(() => {
                api.destroy(this.uuid);
            });
            $(this.confirm).modal("show");
        });
        $(root + " #remove").on("click", () => {
            new ConfirmAction({
                id: this.confirm,
                action: "remove",
                name: this.name,
                message: "remove",
            }).onsubmit(() => {
                api.removeAction(this.uuid);
            });
            $(this.confirm).modal("show");
        });
        $(root + " #autostart").on("click", () => {
            api.autostart(this.uuid, !autoEnabled);
            autoEnabled = !autoEnabled;
            $(root + " #autostart").text($.i18n(autoEnabled ? "disable autostart" : "enable autostart"));
        });
        if (this.state === "inactive") {
            $(root + " #destroy").text("start");
            $(root + " #destroy").off("click").on("click", () => {
                api.start(this.uuid);
            });
        }
    }
}
