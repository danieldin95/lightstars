import {Controller} from "./controller.js"
import {VolumeCtl} from "./volume.js";
import {FileUpload} from "../widget/common/upload.js";
import {UploadApi} from "../api/upload.js";
import {DataStoreApi} from "../api/datastores.js";


export class PoolCtl extends Controller {

    constructor(props) {
        super(props);

        let name = $(this.id).attr("name");
        let uuid = $(this.id).attr("data");
        this.uuid = uuid;
        this.name = name;
        this.tasks = props.tasks || "tasks";
        this.volumes = new VolumeCtl({
            ...props.volumes, uuid, name,
            upload: props.volumes.upload,
        });
        this.upload = new FileUpload({
            id: props.upload
        });
        this.upload.onsubmit((e) => {
            new UploadApi({
                uuids: this.uuid,
                id: '#process'
            }).upload(e.form);
        });

        let root = this.child('#header');
        let api = new DataStoreApi({tasks: this.tasks});
        let state = ($(this.id).attr("state") || "").toLowerCase();
        let autostart = ($(this.id).attr("autostart") || "").toLowerCase() === "true";

        $(root + " #destroy").on("click", () => {
            api.destroy(this.uuid);
        });
        if (state === "inactive") {
            $(root + " #destroy").text("start");
            $(root + " #destroy").off("click").on("click", () => {
                api.start(this.uuid);
            });
        }
        $(root + " #autostart").on("click", () => {
            api.autostart(this.uuid, !autostart);
            autostart = !autostart;
            $(root + " #autostart").text($.i18n(autostart ? "disable autostart" : "enable autostart"));
        });
        $(root + " #clean").on("click", () => {
            api.clean(this.uuid);
        });
        $(root + " #remove").on("click", () => {
            api.removeAction(this.uuid);
        });
    }
}
