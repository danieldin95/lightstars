import {Api} from "./api.js"


export class HyperApi extends Api {
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
            return super.url(`/hyper/${uuid}`);
        }
        return super.url('/hyper');
    }
}
