import {Api} from "./api.js"

export class NetworkApi extends Api {
    // {uuids: [], tasks: 'tasks', name: ''}
    constructor(props) {
        super(props);
    }

    url(uuid) {
        if (uuid) {
            return `/api/network/${uuid}`
        }
        return 'api/network'
    }
}