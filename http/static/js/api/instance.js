import {Api} from "./api.js";
import {Alert} from "../com/alert.js";


export class InstanceApi extends Api {
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
            return super.url(`/instance/${uuid}`);
        }
        return super.url(`/instance`);
    }

    list(data, func) {
        $.GET(this.url(), {format: 'schema'}, (resp, status) => {
            func({data, resp});
        }).fail((e) => {
            $(this.tasks).append(Alert.danger(`GET ${this.url()}: ${e.responseText}`));
        });
    }

    get(data, func) {
        let url = this.url(this.uuids[0]);

        $.GET(url, {format: 'schema'}, (resp, status) => {
            func({data, resp});
        }).fail((e) => {
            $(this.tasks).append(Alert.danger(`GET ${url}: ${e.responseText}`));
        });
    }

    start() {
        this.uuids.forEach((uuid) => {
            let url = this.url(uuid);
            let data = JSON.stringify({action: 'start'});

            $.PUT(url, data, (resp, status) => {
                $(this.tasks).append(Alert.success(`start '${uuid}' success`));
            }).fail((e) => {
                $(this.tasks).append(Alert.danger(`PUT ${url}: ${e.responseText}`));
            });
        });
    }

    shutdown() {
        this.uuids.forEach((uuid) => {
            let url = this.url(uuid);
            let data = JSON.stringify({action: 'shutdown'});

            $.PUT(url, data, (resp, status) => {
                $(this.tasks).append(Alert.success(`shutdown '${uuid}' success`));
            }).fail((e) => {
                $(this.tasks).append(Alert.warn(`PUT ${url}: ${e.responseText}`));
            });
        });
    }

    reset() {
        this.uuids.forEach((uuid, index, err) => {
            let url = this.url(uuid);
            let data = JSON.stringify({action: 'reset'});

            $.PUT(url, data, (resp, status) => {
                $(this.tasks).append(Alert.success(`reset '${uuid}' success`));
            }).fail((e) => {
                $(this.tasks).append(Alert.danger(`PUT ${url}: ${e.responseText}`));
            });
        });
    }

    suspend() {
        this.uuids.forEach((uuid, index, err) => {
            let url = this.url(uuid);
            let data = JSON.stringify({action: 'suspend'});

            $.PUT(this.url(uuid), data, (resp, status) => {
                $(this.tasks).append(Alert.success(`suspend '${uuid}' success`));
            }).fail((e) => {
                $(this.tasks).append(Alert.danger(`PUT ${url}: ${e.responseText}`));
            });
        });
    }

    resume() {
        this.uuids.forEach((uuid, index, err) => {
            let url = this.url(uuid);
            let data = JSON.stringify({action: 'resume'});

            $.PUT(url, data, (resp, status) => {
                $(this.tasks).append(Alert.success(`resume '${uuid}' success`));
            }).fail((e) => {
                $(this.tasks).append(Alert.danger(`PUT ${url}: ${e.responseText}`));
            });
        });
    }

    destroy() {
        this.uuids.forEach((uuid, index, err) => {
            let url = this.url(uuid);
            let data = JSON.stringify({action: 'destroy'});

            $.PUT(url, data, (resp, status) => {
                $(this.tasks).append(Alert.success(`destroy '${uuid}' success`));
            }).fail((e) => {
                $(this.tasks).append(Alert.danger((`PUT ${url}: ${e.responseText}`)));
            });
        });
    }

    remove() {
        this.uuids.forEach((uuid, index, err) => {
            let url = this.url(uuid);

            $.DELETE(url, (resp, status) => {
                $(this.tasks).append(Alert.success(`remove '${uuid}' success`));
            }).fail((e) => {
                $(this.tasks).append(Alert.danger((`DELETE ${url}: ${e.responseText}`)));
            });
        });
    }

    edit(data) {
        let uuid = this.uuids[0];
        let url = this.url(uuid);

        if (data.cpu !== "") {
            let api = url+"/processor";
            let cpu = {cpu: data.cpu};

            $.PUT(api, JSON.stringify(cpu), (resp, status) => {
                $(this.tasks).append(Alert.success(`set processor for '${uuid}' success`));
            }).fail((e) => {
                $(this.tasks).append(Alert.warn((`PUT ${api}: ${e.responseText}`)));
            });
        }
        if (data.memSize !== "") {
            let api = url+"/memory";
            let mem = {size: data.memSize, unit: data.memUnit};

            $.PUT(api, JSON.stringify(mem), (resp, status) => {
                $(this.tasks).append(Alert.success(`set memory for '${uuid}' success`));
            }).fail((e) => {
                $(this.tasks).append(Alert.warn((`PUT ${api}: ${e.responseText}`)));
            });
        }
    }

    console() {
        this.uuids.forEach((uuid, index, err) => {
            let password = '';
            let host = Api.host || "";
            if (this.props.passwd) {
                password = this.props.passwd[uuid];
            }
            let url = "/ui/console?id=" + uuid + "&password=" + password;
            if (host !== "") {
                url += "&node=" + host
            }
            window.open(url);
        });
    }

    create (data) {
        $.POST(this.url(), JSON.stringify(data), (resp, status) => {
            $(this.tasks).append(Alert.success(`create '${resp.name}' success`));
        }).fail((e) => {
            $(this.tasks).append(Alert.danger(`POST ${this.url()}: ${e.responseText}`));
        });
    }
}