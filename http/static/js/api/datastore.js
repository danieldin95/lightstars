import {Api} from "./api.js"
import {Alert} from "../com/alert.js";


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
            return `/api/datastore/${uuid}`
        }
        return '/api/datastore'
    }

    list(data, func) {
        $.get(this.url(), {format: 'schema'}, (resp, status) => {
            func({data, resp});
        }).fail((e) => {
            $(this.tasks).append(Alert.danger((`GET ${this.url()}: ${e.responseText}`)));
        });
    }

    create(data) {
        if (data.format == 'nfs') {
            data.nfs = { host: data.host, path: data.path, format: 'nfs' };
        }
        $.post(this.url(), JSON.stringify(data), (resp, status) => {
            $(this.tasks).append(Alert.success(`create datastore ${resp.message}`));
        }).fail((e) => {
            $(this.tasks).append(Alert.danger(`POST ${this.url()}: ${e.responseText}`));
        });
    }

    delete() {
        this.uuids.forEach((uuid, index, err) => {
            let url = this.url(uuid);

            $.delete(url, (resp, status) => {
                $(this.tasks).append(Alert.success(`remove datastore '${uuid}' ${resp.message}`));
            }).fail((e) => {
                $(this.tasks).append(Alert.danger(`DELETE ${url}: ${e.responseText}`));
            });
        });
    }
}