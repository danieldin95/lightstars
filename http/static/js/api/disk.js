import {Api} from "./api.js"
import {Alert} from "../com/alert.js";


export class DiskApi extends Api {
    // {
    //   inst: 'uuid',
    //   uuids: [uuid],
    //   tasks: 'Tasks',
    //   name: ''
    // }
    constructor(props) {
        super(props);

        this.inst = props.inst;
    }

    url(uuid) {
        if (uuid) {
            return super.url(`/instance/${this.inst}/disk/${uuid}`);
        }
        return super.url(`/instance/${this.inst}/disk`);
    }

    list(data, func) {
        $.GET(this.url(), {format: 'schema'}, (resp, status) => {
            func({data, resp});
        }).fail((e) => {
            $(this.tasks).append(Alert.danger(`GET ${this.url()}: ${e.responseText}`));
        });
    }

    create(data) {
        $.POST(this.url(), JSON.stringify(data), (resp, status) => {
            $(this.tasks).append(Alert.success(`create ${resp.message}`));
        }).fail((e) => {
            $(this.tasks).append(Alert.danger(`POST ${this.url()}: ${e.responseText}`));
        });
    }

    delete() {
        this.uuids.forEach((uuid, index, err) => {
            $.DELETE(this.url(uuid), (resp, success) => {
                $(this.tasks).append(Alert.success(`remove ${uuid} ${resp.message}`));
            }).fail((e) => {
                $(this.tasks).append(Alert.danger(`DELETE ${this.url(uuid)}: ${e.responseText}`));
            });
        });
    }

    edit(data) {
        let url = this.url(this.uuids[0]);

        $.PUT(url, JSON.stringify(data), (resp, success) => {
            $(this.tasks).append(Alert.success(`edit ${resp.name} ${resp.message}`));
        }).fail((e) => {
            $(this.tasks).append(Alert.danger(`PUT ${url}: ${e.responseText}`));
        });
    }
}
