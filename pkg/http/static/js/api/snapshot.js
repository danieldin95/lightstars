import {Api} from "./api.js"
import {Alert} from "../lib/alert.js";


export class SnapshotApi extends Api {
    // {
    //   inst: 'uuid',
    //   uuids: [uuid],
    //   tasks: 'tasks',
    //   name: ''
    // }
    constructor(props) {
        super(props);
        this.inst = props.inst;
    }

    url(uuid, action) {
        if (uuid) {
            if (action) {
                return super.url(`/instance/${this.inst}/snapshot/${uuid}/${action}`);
            }
            return super.url(`/instance/${this.inst}/snapshot/${uuid}`);
        }
        return super.url(`/instance/${this.inst}/snapshot`);
    }

    revert() {
        let uuid = this.uuids[0];
        let url = this.url(uuid, 'revert');
        $.PUT(url, (resp, status) => {
            //$(this.tasks).append(Alert.success(`destroy '${uuid}' success`));
        }).fail((e) => {
            Alert.danger(this.tasks,`PUT ${url}: ${e.responseText}`);
        });
    }
}
