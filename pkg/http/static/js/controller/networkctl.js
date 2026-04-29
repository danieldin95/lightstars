import {Controller} from "./controller.js"
import {Leases} from "./leasesctl.js";
import {Port} from "./portctl.js";
import {NetworkApi} from "../api/networkapi.js";
import {confirmaction} from "../widget/common/confirmaction.js";

export class Network extends Controller {
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

        this.leases = new Leases({...props.leases, uuid, name});
        this.ports = new Port({...props.ports, uuid, name, confirm: this.confirm});

        let api = new NetworkApi({tasks: this.tasks});
        let root = this.child('#header');
        let autoBtn = $(root + " #autostart");
        let autoEnabled = this.autostart;

        $(root + " #destroy").on("click", () => {
            new confirmaction({
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
            new confirmaction({
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
