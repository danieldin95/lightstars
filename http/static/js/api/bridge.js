import {Api} from "./api.js"
import {Alert} from "../com/alert.js";


export class BridgeApi extends Api {
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
            return super.url(`/bridge/${uuid}`);
        }
        return super.url('/bridge');
    }

    list(data, func) {
        if (typeof data == "function") {
            func = data;
        }
        $.GET(this.url(), {format: 'schema'}, (resp, status) => {
            func({data, resp});
        }).fail((e) => {
            $(this.tasks).append(Alert.danger(`GET ${this.url()}: ${e.responseText}`));
        });
    }
}