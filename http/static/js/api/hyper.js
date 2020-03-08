import {Api} from "./api.js"
import {Alert} from "../com/alert.js";


export class HyperApi extends Api {
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
            return `/api/hyper/${uuid}`
        }
        return '/api/hyper'
    }

    get(data, func) {
        $.get(this.url(), {format: 'schema'}, (resp, status) => {
            func({data, resp});
        }).fail((e) => {
            $(this.tasks).append(Alert.danger(`GET ${this.url()}: ${e.responseText}`));
        });
    }
}