import {Api} from "./api.js"


export class VolumeApi extends Api {
    // {
    //   uuids: [],
    //   tasks: 'tasks',
    //   name: ''
    // }
    constructor(props) {
        super(props);
        this.pool = this.props.pool
    }

    url(uuid) {

        if (uuid) {
            return super.url(`/datastore/${this.pool}/volume/${uuid}`);
        }
        return super.url(`/datastore/${this.pool}/volume`);
    }
}