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
            return super.url(`/network/${uuid}`);
        }
        return super.url('/network');
    }

    list(data, func) {
        $.GET(this.url(), {format: 'schema'}, (resp, status) => {
            func({data, resp});
        }).fail((e) => {
            $(this.tasks).append(Alert.danger((`GET ${this.url()}: ${e.responseText}`)));
        });
    }

    get(data, func) {
        let url = this.url(this.uuids[0]);

        $.GET(url, {format: 'schema'}, (resp, status) => {
            func({data, resp});
        }).fail((e) => {
            $(this.tasks).append(Alert.danger(`GET ${url}: ${e.responseText}`));
        });
    }

    create(data) {
        if (data.range && data.range !== "") {
            let lines = data.range.split(/\r?\n/);

            data.range = [];
            for (let line of lines) {
                if (line.indexOf(',') > 0) {
                    let [start, end] = line.split(',', 2);
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
