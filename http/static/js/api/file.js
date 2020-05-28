import {Api} from "./api.js"
import {Alert} from "../com/alert.js";


export class FileApi extends Api {

    constructor(props) {
        super(props);
    }

    url() {

        // return '/ext/files/01/'
        return window.location.pathname
    }

    list(data, func) {
        if (typeof data == 'function') {
            func = data;
            data = {};
        }
        $.POST(this.url(), (resp, status) => {
            func({data, resp})
        }).fail((e) => {
            $(this.tasks).append(Alert.danger(`GET ${this.url()}: ${e.responseText}`));
        });
    }
}