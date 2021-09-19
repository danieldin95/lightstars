import {Api} from "./api.js"


export class HistoryApi extends Api {
    // {
    //   uuids: [],
    //   tasks: 'tasks',
    //   name: ''
    // }
    constructor(props) {
        super(props);
        this._host = ""
    }

    url(action) {
        if (action) {
            return super.url(`/history/${action}`);
        }
        return super.url('/history');
    }
}
