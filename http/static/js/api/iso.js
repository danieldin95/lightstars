import {Api} from "./api.js"
import {Alert} from "../com/alert.js";


export class IsoApi extends Api {
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
            return super.url(`/iso/${uuid}`);
        }
        return super.url('/iso');
    }

    list(datastore, func) {
        if (typeof data == "function") {
            func = data;
        }
        $.GET(this.url(), {datastore}, (resp, status) => {
            func({datastore, resp});
        }).fail((e) => {
            $(this.tasks).append(Alert.danger(`GET ${this.url()}: ${e.responseText}`));
        });
    }
}