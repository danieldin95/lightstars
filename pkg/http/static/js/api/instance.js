import {Api} from "./api.js";
import {Alert} from "../lib/alert.js";
import {Utils} from "../lib/utils.js";


export class InstanceApi extends Api {
    // {
    //   uuids: [],
    //   tasks: 'tasks',
    //   name: ''
    // }
    constructor(props) {
        super(props);
    }

    url(uuid, action) {
        if (uuid) {
            if (action) {
                return super.url(`/instance/${uuid}/${action}`);
            }
            return super.url(`/instance/${uuid}`);
        }
        return super.url(`/instance`);
    }

    start() {
        this.uuids.forEach((uuid) => {
            let url = this.url(uuid, 'start');

            $.PUT(url, (resp, status) => {
                Alert.success(this.tasks, `start '${uuid}' success`);
            }).fail((e) => {
                Alert.danger(this.tasks,`PUT ${url}: ${e.responseText}`);
            });
        });
    }

    shutdown() {
        this.uuids.forEach((uuid) => {
            let url = this.url(uuid, 'shutdown');

            $.PUT(url, (resp, status) => {
                Alert.success(this.tasks, `shutdown '${uuid}' success`);
            }).fail((e) => {
                Alert.warn(this.tasks,`PUT ${url}: ${e.responseText}`);
            });
        });
    }

    reset() {
        this.uuids.forEach((uuid, index, err) => {
            let url = this.url(uuid, 'reset');

            $.PUT(url, (resp, status) => {
                Alert.success(this.tasks, `reset '${uuid}' success`);
            }).fail((e) => {
                Alert.danger(this.tasks,`PUT ${url}: ${e.responseText}`);
            });
        });
    }

    suspend() {
        this.uuids.forEach((uuid, index, err) => {
            let url = this.url(uuid, 'suspend');

            $.PUT(url, (resp, status) => {
                Alert.success(this.tasks, `suspend '${uuid}' success`);
            }).fail((e) => {
                Alert.danger(this.tasks,`PUT ${url}: ${e.responseText}`);
            });
        });
    }

    resume() {
        this.uuids.forEach((uuid, index, err) => {
            let url = this.url(uuid, 'resume');

            $.PUT(url, (resp, status) => {
                Alert.success(this.tasks, `resume '${uuid}' success`);
            }).fail((e) => {
                Alert.danger(this.tasks,`PUT ${url}: ${e.responseText}`);
            });
        });
    }

    destroy() {
        this.uuids.forEach((uuid, index, err) => {
            let url = this.url(uuid, 'destroy');

            $.PUT(url, (resp, status) => {
                Alert.success(this.tasks, `destroy '${uuid}' success`);
            }).fail((e) => {
                Alert.danger(this.tasks, `PUT ${url}: ${e.responseText}`);
            });
        });
    }

    remove() {
        this.uuids.forEach((uuid, index, err) => {
            let url = this.url(uuid);

            $.DELETE(url, (resp, status) => {
                Alert.success(this.tasks, `remove '${uuid}' success`);
            }).fail((e) => {
                Alert.danger(this.tasks, `DELETE ${url}: ${e.responseText}`);
            });
        });
    }

    edit(data) {
        let uuid = this.uuids[0];

        if (data.cpu !== "") {
            let url = this.url(uuid, 'processor');
            let cpu = {cpu: data.cpu, mode: data.cpuMode};

            $.PUT(url, JSON.stringify(cpu), (resp, status) => {
                Alert.success(this.tasks, `set processor for '${uuid}' success`);
            }).fail((e) => {
                Alert.warn(this.tasks,`PUT ${url}: ${e.responseText}`);
            });
        }
        if (data.memSize !== "") {
            let url = this.url(uuid, 'memory');
            let mem = {size: data.memSize, unit: data.memUnit};

            $.PUT(url, JSON.stringify(mem), (resp, status) => {
                Alert.success(this.tasks, `set memory for '${uuid}' success`);
            }).fail((e) => {
                Alert.warn(this.tasks,`PUT ${url}: ${e.responseText}`);
            });
        }
    }

    title(data) {
        let uuid = this.uuids[0];
        let url = this.url(uuid, 'title');
        let params = JSON.stringify({title: data.title});

        $.PUT(url, params, (resp, status) => {
            Alert.success(this.tasks, `destroy '${uuid}' success`);
        }).fail((e) => {
            Alert.danger(this.tasks, `PUT ${url}: ${e.responseText}`);
        });
    }

    console() {
        this.uuids.forEach((uuid, index, err) => {
            let password = '', name = '';
            let host = Api.host() || "";
            if (this.props.passwd) {
                password = this.props.passwd[uuid] || "";
            }
            if (this.props.name) {
                name = this.props.name[uuid] || "";
            }
            let url = "/ui/console?id="+uuid+"&password="+password+"&title="+name;
            if (host !== "") {
                url += "&node=" + host
            }
            window.open(url);
        });
    }

    create (data) {
        let schema = {};
        schema.name = data.name;
        schema.family = data.family;
        schema.datastore = data.datastore;
        schema.boots = 'hd,cdrom,network';
        schema.maxCpu = parseInt(data.cpu);
        schema.cpuMode = data.cpuMode || '';
        schema.maxMem = Utils.toKiB(data.memSize, data.memUnit);
        schema.disks = [
            {
                source: data.disk0File,
            },
            {
                size: data.disk1Size,
                sizeUnit: data.disk1Unit || 'MiB',
                source: data.disk1File || '',
            },
         ];
        schema.interfaces = [
            {
                source: data.interface0Source,
                Modal: data.interface0Model || 'virtio',
            },
        ];
        $.POST(this.url(), JSON.stringify(schema), (resp, status) => {
            Alert.success(this.tasks, `create ${resp.name} success`);
        }).fail((e) => {
            Alert.danger(this.tasks,`POST ${this.url()}: ${e.responseText}`);
        });
    }

    stats(data, func, fail) {
        if (typeof data == 'function') {
            func = data;
            data = {};
        }
        let url = this.url('stats');

        $.GET(url, {format: 'schema'}, (resp, status) => {
            func({data, resp});
        }).fail((e) => {
            if (fail) {
                fail(e);
            }
            Alert.danger(this.tasks,`GET ${url}: ${e.responseText}`);
        });
    }
}
