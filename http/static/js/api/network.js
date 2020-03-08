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
        $.get(this.url(), {format: 'schema'}, (resp, status) => {
            func({data, resp});
        }).fail((e) => {
            $(this.tasks).append(Alert.danger((`GET ${this.url()}: ${e.responseText}`)));
        });
    }

    create(data) {
        if (data.range && data.range !== "") {
            let range = data.range.split(',', 2);
            if (range.length == 2) {
                data.rangeStart = range[0];
                data.rangeEnd = range[1];
            }
        }
        $.post(this.url(), JSON.stringify(data), (resp, status) => {
            $(this.tasks).append(Alert.success(`create ${resp.message} success`));
        }).fail((e) => {
            $(this.tasks).append(Alert.danger((`GET ${this.url()}: ${e.responseText}`)));
        });
    }

    delete() {
        this.uuids.forEach((uuid, index, err) => {
            $.delete(this.url(uuid), (resp, status) => {
                $(this.tasks).append(Alert.success(`remove '${uuid}' ${resp.message} success`));
            }).fail((e) => {
                $(this.tasks).append(Alert.danger((`GET ${this.url(uuid)}: ${e.responseText}`)));
            });
        });
    }
}