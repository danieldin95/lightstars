import {Api} from "./api.js"


export class FileApi extends Api {

    constructor(props) {
        super(props);
    }

    url() {
        return window.location.pathname
    }
}
