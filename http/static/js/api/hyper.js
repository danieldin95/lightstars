import {Api} from "./api.js"
import {AlertDanger, AlertSuccess} from "../com/alert.js";


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
        let your = this;

        $.get(your.url(), {format: 'schema'}, function (resp, status) {
            func({data, resp});
        }).fail(function (e) {
            $(your.tasks).append(AlertDanger((`${this.type} ${this.url}: ${e.responseText}`)));
        });
    }
}