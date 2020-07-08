import {Api} from "./api.js"


export class LeaseApi extends Api {
    // {
    //   net: 'uuid',
    //   uuids: [uuid],
    //   tasks: 'tasks',
    //   name: ''
    // }
    constructor(props) {
        super(props);

        this.net = props.net;
    }

    url(uuid) {
        if (uuid) {
            return super.url(`/network/${this.net}/lease/${uuid}`);
        }
        return super.url(`/network/${this.net}/lease`);
    }
}
