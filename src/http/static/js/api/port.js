import {Api} from "./api.js"


export class PortApi extends Api {
    // {
    //   uuids: [],
    //   tasks: 'tasks',
    //   name: ''
    // }
    constructor(props) {
        super(props);

        this.bridge = props.bridge;
    }

    url(uuid) {
        if (uuid) {
            return super.url(`/network/${this.bridge}/interface/${uuid}`);
        }
        return super.url(`/network/${this.bridge}/interface`);
    }
}
