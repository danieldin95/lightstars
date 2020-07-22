import {Api} from "./api.js"


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
        super.edit(data);
    }
}
