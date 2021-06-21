import {Api} from "./api.js"


export class GraphicsApi extends Api {
    // {
    //   inst: 'uuid',
    //   uuids: [uuid],
    //   tasks: 'tasks',
    //   name: ''
    // }
    constructor(props) {
        super(props);
        this.inst = props.inst;
    }

    url(inst, uuid) {
        if (uuid) {
            return super.url(`/instance/${this.inst}/graphics/${uuid}`);
        }
        return super.url(`/instance/${this.inst}/graphics`);
    }
}
