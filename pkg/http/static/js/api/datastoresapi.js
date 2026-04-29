import {Api} from "./api.js"
import {Alert} from "../lib/alert.js";


export class DataStoreApi extends Api {
    // {
    //   uuids: [],
    //   tasks: 'tasks',
    //   name: ''
    // }
    constructor(props) {
        super(props);
    }

    url(uuid) {
        if (uuid) {
            return super.url(`/datastore/${uuid}`);
        }
        return super.url('/datastore');
    }
    create(data) {
        if (data.format === 'nfs') {
            data.nfs = { host: data.host, path: data.path, format: 'nfs' };
        }
        super.create(data)
    }

    putAction(uuid, action) {
        let url = this.url(uuid) + "/" + action;
        $.PUT(url, (resp, status) => {
            Alert.success(this.tasks, `${action} '${uuid}' success`);
        }).fail((e) => {
            Alert.danger(this.tasks, `PUT ${url}: ${e.responseText}`);
        });
    }

    start(uuid) { this.putAction(uuid, "start"); }
    destroy(uuid) { this.putAction(uuid, "destroy"); }
    refresh(uuid) { this.putAction(uuid, "refresh"); }
    autostart(uuid, enable) {
        let url = this.url(uuid) + "/autostart?enable=" + (enable ? "true" : "false");
        $.PUT(url, (resp, status) => {
            Alert.success(this.tasks, `autostart '${uuid}' success`);
        }).fail((e) => {
            Alert.danger(this.tasks, `PUT ${url}: ${e.responseText}`);
        });
    }
    clean(uuid) { this.putAction(uuid, "clean"); }
    removeAction(uuid) { this.putAction(uuid, "remove"); }
}
