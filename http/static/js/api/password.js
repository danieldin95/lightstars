import {Api} from "./api.js"
import {Alert} from "../com/alert.js";


export class PasswordApi extends Api {
    // {
    //   uuids: [],
    //   tasks: 'Tasks',
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
    set(data, func) {
        
        $.POST(this.url(), JSON.stringify(data), (resp, status) => {
            if (resp.message == "success") {
                $(this.tasks).append(Alert.success(`update password success`));
            } else {
                $(this.tasks).append(Alert.danger((`please enter right old password`)));
            }
            
        }).fail((e) => {
            $(this.tasks).append(Alert.danger((`POST ${this.url()}: ${e.responseText}`)));
        });

    }
}
