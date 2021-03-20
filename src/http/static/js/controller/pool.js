import {Ctl} from "./ctl.js"
import {VolumeCtl} from "./volume.js";
import {FileUpload} from "../widget/common/upload.js";
import {UploadApi} from "../api/upload.js";


export class PoolCtl extends Ctl {

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
    }
}
