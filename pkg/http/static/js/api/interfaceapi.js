import {Api} from "./api.js"


export class InterfaceApi extends Api {
    // {
    //   uuids: [],
    //   tasks: 'tasks',
    //   name: ''
    // }
    constructor(props) {
        super(props);

        this.inst = props.inst;
    }

    url(uuid) {
        if (uuid) {
            return super.url(`/instance/${this.inst}/interface/${uuid}`);
        }
        return super.url(`/instance/${this.inst}/interface`);
    }
}
