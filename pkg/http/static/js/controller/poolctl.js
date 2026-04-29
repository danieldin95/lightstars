import {Controller} from "./controller.js"
import {Volume} from "./volumectl.js";
import {fileupload} from "../widget/common/fileupload.js";
import {UploadApi} from "../api/uploadapi.js";
import {DataStoreApi} from "../api/datastoresapi.js";
import {confirmaction} from "../widget/common/confirmaction.js";


export class Pool extends Controller {

    constructor(props) {
        super(props);

        let name = $(this.id).attr("name");
        let uuid = $(this.id).attr("data");
        this.uuid = uuid;
        this.name = name;
        this.tasks = props.tasks || "tasks";
        this.confirm = props.confirm;
        this.volumes = new Volume({
            ...props.volumes, uuid, name,
            upload: props.volumes.upload,
        });
        this.upload = new fileupload({
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
    }
}
