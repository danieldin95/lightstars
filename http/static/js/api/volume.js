import {Api} from "./api.js"


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
}
