import {Api} from "./api.js"
import {Alert} from "../com/alert.js";


export class GraphicsApi extends Api {
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

    url(inst, uuid) {
        if (uuid) {
            return `/api/instance/${this.inst}/graphics/${uuid}`
        }
        return `/api/instance/${this.inst}/graphics`
    }

    list(data, func) {
        $.get(this.url(), {format: 'schema'}, (resp, status) => {
            func({data, resp});
        }).fail((e) => {
            $(this.tasks).append(Alert.danger(`GET ${this.url()}: ${e.responseText}`));
        });
    }

    create(data) {
        $.post(this.url(), JSON.stringify(data), (resp, status) => {
            $(this.tasks).append(Alert.success(`create ${resp.message}`));
        }).fail((e) => {
            $(this.tasks).append(Alert.danger(`POST ${this.url()}: ${e.responseText}`));
        });
    }

    edit(data) {
        let url = this.url(this.uuids[0]);
        
        $.put(url, JSON.stringify(data), (resp, status) => {
            $(this.tasks).append(Alert.success(`edit '${resp.name}' ${resp.message}`));
        }).fail((e) => {
            $(this.tasks).append(Alert.danger(`PUT ${url}: ${e.responseText}`));
        });
    }
}