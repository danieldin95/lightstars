

export class Api {
    // {
    //   uuids: [],
    //   tasks: 'tasks',
    //   name: ''
    // }
    constructor(props) {
        if (!props) {
            props = {};
        }
        if (props.tasks) {
            this.tasks = props.tasks;
        } else {
            this.tasks = "tasks";
        }
        if (typeof props.uuids == "string") {
            this.uuids = [props.uuids];
        } else {
            this.uuids = props.uuids;
        }
        this.name = props.name;
        this.props = props;
    }

    url(data) {
        return '/api/${data}'
    }
}