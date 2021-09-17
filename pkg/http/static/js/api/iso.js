import {Api} from "./api.js"
import {Alert} from "../lib/alert.js";


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
        if (typeof datastore == "function") {
            func = datastore;
        }
        $.GET(this.url(), {datastore}, (resp, status) => {
            func({datastore, resp});
        }).fail((e) => {
            Alert.danger(this.tasks,`GET ${this.url()}: ${e.responseText}`);
        });
    }
}
