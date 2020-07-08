import {Alert} from "../com/alert.js";

export class Api {
    // static func.
    static prefix() {
        if (this._host) {
            return `/host/${this._host}`
        }
        return ""
    }

    static host(name) {
        if (name !== undefined) {
            this._host = name
        }
        return this._host
    }

    // {
    //   uuids: [],
    //   tasks: 'Tasks',
    //   name: ''
    // }
    constructor(props) {
        if (!props) {
            props = {};
        }

        this.name = props.name;
        this.props = props;
        this.tasks = props.tasks || "Tasks";
        if (typeof props.uuids === "string") {
            this.uuids = [props.uuids];
        } else if (typeof props.uuids === "object") {
            this.uuids = props.uuids;
        } else {
            this.uuids = [];
        }
        this._host = Api._host || ""
    }

    url(suffix) {
        suffix = suffix || "";
        if (this._host !== "") {
            return `/host/${this._host}/api${suffix}`
        }
        return `/api${suffix}`
    }

    list(data, func) {
        if (typeof data == "function") {
            func = data;
            data = {};
        }
        $.GET(this.url(), {format: 'schema'}, (resp, status) => {
            func({data, resp});
        }).fail((e) => {
            $(this.tasks).append(Alert.danger(`GET ${this.url()}: ${e.responseText}`));
        });
    }

    get(data, func) {
        if (typeof data == "function") {
            func = data;
            data = {};
        }
        let url = this.url();
        if (this.uuids.length > 0) {
            url = this.url(this.uuids[0]);
        }
        $.GET(url, {format: 'schema'}, (resp, status) => {
            func({data, resp});
        }).fail((e) => {
            $(this.tasks).append(Alert.danger(`GET ${this.url()}: ${e.responseText}`));
        });
    }

    create(data) {
        $.POST(this.url(), JSON.stringify(data), (resp, status) => {
            $(this.tasks).append(Alert.success(`create ${resp.message}`));
        }).fail((e) => {
            $(this.tasks).append(Alert.danger(`POST ${this.url()}: ${e.responseText}`));
        });
    }

    delete() {
        this.uuids.forEach((uuid, index, err) => {
            $.DELETE(this.url(uuid), (resp, success) => {
                $(this.tasks).append(Alert.success(`remove ${uuid} ${resp.message}`));
            }).fail((e) => {
                $(this.tasks).append(Alert.danger(`DELETE ${this.url(uuid)}: ${e.responseText}`));
            });
        });
    }

    edit(data) {
        let url = this.url(this.uuids[0]);
        $.PUT(url, JSON.stringify(data), (resp, success) => {
            $(this.tasks).append(Alert.success(`edit ${resp.name} ${resp.message}`));
        }).fail((e) => {
            $(this.tasks).append(Alert.danger(`PUT ${url}: ${e.responseText}`));
        });
    }
}
