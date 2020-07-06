import {Api} from "./api.js"


export class NetworkApi extends Api {
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
            return super.url(`/network/${uuid}`);
        }
        return super.url('/network');
    }

    create(data) {
        let range = data.range || "";
        data.range = [];
        if (range !== "") {
            let lines = range.split(/\r?\n/);
            for (let line of lines) {
                if (line.indexOf(',') > 0) {
                    let [start, end] = line.split(',', 2);
                    data.range.push({start, end});
                }
            }
        }
        super.create(data);
    }
}
