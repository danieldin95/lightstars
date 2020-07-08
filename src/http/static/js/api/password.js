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

    url(uuid) {
        if (uuid) {
            return `/api/user/password/${uuid}`;
        }
        return '/api/user/password';
    }

    set(data) {
        super.create(data);
    }
}
