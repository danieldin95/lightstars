import {Api} from "./api.js"
import {Alert} from "../lib/alert.js";


export class PortApi extends Api {
    // {
    //   uuids: [],
    //   tasks: 'tasks',
    //   query: {uuid: '00-00-00'}
    //   name: ''
    // }
    constructor(props) {
        super(props);

        this.bridge = props.bridge;
        if (props.query) {
            this.query = props.query;
        } else {
            this.query = {format: 'json'};
        }
    }

    url(uuid) {
        if (uuid) {
            return super.url(`/network/${this.bridge}/interface/${uuid}`);
        }
        return super.url(`/network/${this.bridge}/interface`);
    }

    list(data, func) {
        if (typeof data == "function") {
            func = data;
            data = {};
        }
        $.GET(this.url(), this.query, (resp, status) => {
            func({data, resp});
        }).fail((e) => {
            Alert.danger(this.tasks, `GET ${this.url()}: ${e.responseText}`);
        });
    }
}
