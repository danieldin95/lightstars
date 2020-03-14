import {Api} from "./api.js"
import {Alert} from "../com/alert.js";


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
            return `/api/network/${uuid}`
        }
        return '/api/network'
    }

    list(data, func) {
        $.GET(this.url(), {format: 'schema'}, (resp, status) => {
            func({data, resp});
        }).fail((e) => {
            $(this.tasks).append(Alert.danger((`GET ${this.url()}: ${e.responseText}`)));
        });
    }

    create(data) {
        if (data.range && data.range !== "") {
            let lines = data.range.split(/\r?\n/);

            data.range = [];
            for (let i = 0; i < lines.length; i ++) {
                if (lines[i].indexOf(',') > 0) {
                    let [start, end] = lines[i].split(',', 2);
                    data.range.push({start, end});
                }
            }
        }
        $.POST(this.url(), JSON.stringify(data), (resp, status) => {
            $(this.tasks).append(Alert.success(`create ${resp.message} success`));
        }).fail((e) => {
            $(this.tasks).append(Alert.danger((`POST ${this.url()}: ${e.responseText}`)));
        });
    }

    delete() {
        this.uuids.forEach((uuid, index, err) => {
            $.DELETE(this.url(uuid), (resp, status) => {
                $(this.tasks).append(Alert.success(`remove '${uuid}' ${resp.message} success`));
            }).fail((e) => {
                $(this.tasks).append(Alert.danger((`DELETE ${this.url(uuid)}: ${e.responseText}`)));
            });
        });
    }
}