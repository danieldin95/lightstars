import {Api} from "./api.js"
import {Alert} from "../com/alert.js";


export class LeaseApi extends Api {
    // {
    //   net: 'uuid',
    //   uuids: [uuid],
    //   tasks: 'tasks',
    //   name: ''
    // }
    constructor(props) {
        super(props);

        this.net = props.net;
    }

    url(uuid) {
        if (uuid) {
            return super.url(`/network/${this.net}/lease/${uuid}`);
        }
        return super.url(`/network/${this.net}/lease`);
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
