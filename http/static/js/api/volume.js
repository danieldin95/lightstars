import {Api} from "./api.js"
import {Alert} from "../com/alert.js";


export class VolumeApi extends Api {
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
            return super.url(`/volume/${uuid}`);
        }
        return super.url('/volume');
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
