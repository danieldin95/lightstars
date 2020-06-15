import {Api} from "./api.js"


export class VolumeApi extends Api {
    // {
    //   uuids: [],
    //   tasks: 'Tasks',
    //   name: ''
    // }
    constructor(props) {
        super(props);
        this.uuid = this.props.uuid
    }

    url(uuid) {
        if (uuid) {
            return super.url(`/volume/${this.uuid}`);
        }
        return super.url(`/volume/${this.uuid}`);
    }
}