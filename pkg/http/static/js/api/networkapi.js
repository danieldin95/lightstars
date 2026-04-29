import {Api} from "./api.js"
import {Alert} from "../lib/alert.js";


export class NetworkApi extends Api {
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
            return super.url(`/network/${uuid}`);
        }
        return super.url('/network');
    }

    create(data) {
        let range = data.range || "";
        data.range = [];
        if (range !== "") {
            let lines = range.split(/\r?\n/);
            for (let line of lines) {
                if (line.indexOf(',') > 0) {
                    let [start, end] = line.split(',', 2);
                    data.range.push({start, end});
                }
            }
        }
        super.create(data);
    }

    putAction(uuid, action) {
        let url = this.url(uuid) + "/" + action;
        $.PUT(url, (resp, status) => {
            Alert.success(this.tasks, `${action} '${uuid}' success`);
        }).fail((e) => {
            Alert.danger(this.tasks, `PUT ${url}: ${e.responseText}`);
        });
    }

    start(uuid) {
        this.putAction(uuid, "start");
    }

    destroy(uuid) {
        this.putAction(uuid, "destroy");
    }

    autostart(uuid, enable) {
        let url = this.url(uuid) + "/autostart?enable=" + (enable ? "true" : "false");
        $.PUT(url, (resp, status) => {
            Alert.success(this.tasks, `autostart '${uuid}' success`);
        }).fail((e) => {
            Alert.danger(this.tasks, `PUT ${url}: ${e.responseText}`);
        });
    }

    removeAction(uuid) {
        this.putAction(uuid, "remove");
    }

    remove(uuid) {
        let url = this.url(uuid);
        $.DELETE(url, (resp, status) => {
            Alert.success(this.tasks, `remove '${uuid}' success`);
        }).fail((e) => {
            Alert.danger(this.tasks, `DELETE ${url}: ${e.responseText}`);
        });
    }
}
