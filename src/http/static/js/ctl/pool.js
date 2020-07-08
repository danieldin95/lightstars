import {Ctl} from "./ctl.js"
import {VolumeCtl} from "./volume.js";


export class PoolCtl extends Ctl {

    constructor(props) {
        super(props);

        let name = $(this.id).attr("name");
        let uuid = $(this.id).attr("data");
        this.uuid = uuid;
        this.name = name;
        this.tasks = props.tasks || "tasks";

        this.volumes = new VolumeCtl({id: props.volumes.id, uuid, name});
    }
}
