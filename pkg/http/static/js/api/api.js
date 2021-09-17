import {Alert} from "../lib/alert.js";

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

        this.props = props;
        this.name = props.name;
        this.tasks = props.tasks || "tasks";
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
            Alert.danger(this.tasks, `GET ${this.url()}: ${e.responseText}`);
        });
    }

    get(data, func, fail) {
        if (typeof data == "function") {
            fail = func;
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
            if (fail) {
                fail(e)
            }
            Alert.danger(this.tasks,`GET ${this.url()}: ${e.responseText}`);
        });
    }

    create(data) {
        $.POST(this.url(), JSON.stringify(data), (resp, status) => {
            Alert.success(this.tasks, `create ${resp.message}`);
        }).fail((e) => {
            Alert.danger(this.tasks,`POST ${this.url()}: ${e.responseText}`);
        });
    }

    delete() {
        this.uuids.forEach((uuid, index, err) => {
            $.DELETE(this.url(uuid), (resp, success) => {
                Alert.success(this.tasks, `remove ${uuid} ${resp.message}`);
            }).fail((e) => {
                Alert.danger(this.tasks,`DELETE ${this.url(uuid)}: ${e.responseText}`);
            });
        });
    }

    edit(data) {
        let uuid = this.uuids[0];
        let url = this.url(uuid);
        $.PUT(url, JSON.stringify(data), (resp, success) => {
            Alert.success(this.tasks, `edit ${uuid} ${resp.message}`);
        }).fail((e) => {
            Alert.danger(this.tasks,`PUT ${url}: ${e.responseText}`);
        });
    }

    // static func.
    static path(url) {
        url = url || "";
        if (this._host !== "") {
            return `/host/${this._host}/${url}`;
        }
        return url;
    }

    static host(name) {
        if (name !== undefined) {
            this._host = name
        }
        return this._host
    }
}
