import {Api} from "./api.js"


export class ImageApi extends Api {
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
            return super.url(`/image/${uuid}`);
        }
        return super.url('/image');
    }
}
