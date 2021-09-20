import {Api} from "./api.js"
import {Alert} from "../lib/alert.js";


export class PasswordApi extends Api {
    // {
    //   uuids: [],
    //   tasks: 'tasks',
    //   name: ''
    // }
    constructor(props) {
        super(props);
    }

    url(name) {
        if (name) {
            return `/api/user/${name}/password`;
        }
        return '/api/user/password';
    }

    set(data) {
        if (data.new === data.repeat) {
            super.edit(data);
        } else {
            Alert.danger(this.tasks, `SET ${this.url()}: The passwords are inconsistent`);
        }
    }
}
