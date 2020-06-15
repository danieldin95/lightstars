import {Api} from "./api.js"


export class FileApi extends Api {

    constructor(props) {
        super(props);
    }

    url(uuid) {
        if (uuid) {
            return super.url(`/api/volume/${uuid}`);
        }
        return super.url(`/api/volume/`);
    }
}
