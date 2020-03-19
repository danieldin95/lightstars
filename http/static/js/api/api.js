

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

        this.name = props.name;
        this.props = props;
        this.tasks = props.tasks || "tasks";
        if (typeof props.uuids == "string") {
            this.uuids = [props.uuids];
        } else {
            this.uuids = props.uuids;
        }
        this.host = Api.host || ""
    }

    url(suffix) {
        suffix = suffix || "";
        if (this.host !== "") {
            return `/host/${this.host}/api${suffix}`
        }
        return `/api${suffix}`
    }

    static prefix() {
        if (this.host) {
            return `/host/${this.host}`
        }
        return ""
    }

    static Host(name) {
        if (typeof name !== 'undefined') {
            this.host = name
        }
        return this.host
    }
}