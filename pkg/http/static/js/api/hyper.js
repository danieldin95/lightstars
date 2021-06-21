import {Api} from "./api.js"
import {Alert} from "../lib/alert.js";


export class HyperApi extends Api {
    // {
    //   uuids: [],
    //   tasks: 'tasks',
    //   name: ''
    // }
    constructor(props) {
        super(props);
    }

    url(uuid, action) {
        if (uuid) {
            if (action) {
                return super.url(`/hyper/${uuid}/${action}`);
            }
            return super.url(`/hyper/${uuid}`);
        }
        return super.url('/hyper');
    }

    statics(data, func, fail) {
        if (typeof data == "function") {
            fail = func;
            func = data;
            data = {};
        }
        let url = this.url('statics');
        if (this.uuids.length > 0) {
            url = this.url(this.uuids[0]);
        }
        $.GET(url, {format: 'schema'}, (resp, status) => {
            func({data, resp});
        }).fail((e) => {
            if (fail) {
                fail(e)
            }
            console.log(`GET ${this.url()}: ${e.responseText}`);
        });
    }
}
