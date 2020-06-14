import {Api} from "./api.js"


export class DiskApi extends Api {
    // {
    //   inst: 'uuid',
    //   uuids: [uuid],
    //   tasks: 'Tasks',
    //   name: ''
    // }
    constructor(props) {
        super(props);

        this.inst = props.inst;
    }

    url(uuid) {
        if (uuid) {
            return super.url(`/instance/${this.inst}/disk/${uuid}`);
        }
        return super.url(`/instance/${this.inst}/disk`);
    }
}
