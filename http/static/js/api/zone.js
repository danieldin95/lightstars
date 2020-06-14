import {Api} from "./api.js"
import {Alert} from "../com/alert.js";


export class ZoneApi extends Api {
    // {
    //   uuids: [],
    //   tasks: 'Tasks',
    //   name: ''
    // }
    constructor(props) {
        super(props);
    }

    url(uuid) {
        if (uuid) {
            return `/api/zone/${uuid}`;
        }
        return '/api/zone';
    }

    list(data, func) {
        if (typeof data == "function") {
            func = data;
        }
        $.GET(this.url(), (resp, status) => {
            func({data, resp});
        }).fail((e) => {
            $(this.tasks).append(Alert.danger(`GET ${this.url()}: ${e.responseText}`));
        });
    }
}
