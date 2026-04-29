import {Api} from "./api.js"


export class BridgeApi extends Api {
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
            return super.url(`/bridge/${uuid}`);
        }
        return super.url('/bridge');
    }
}
