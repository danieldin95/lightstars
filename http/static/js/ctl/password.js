import {PasswordApi} from "../api/password.js";

export class Password {

    constructor(props) {
        this.id = props.id;
        this.props = props;
        this.tasks = props.tasks || "tasks";
        
    }

    set(data) {
        new PasswordApi().set(data);
    }
}
