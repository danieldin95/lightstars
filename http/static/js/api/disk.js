import {Api} from "./api.js"
import {AlertDanger, AlertSuccess, AlertWarn} from "../com/alert.js";

export class DiskApi extends Api {
    // {uuids: [], tasks: 'tasks', name: '', instance: ''}
    constructor(props) {
        super(props);

        this.instance = props.instance;
    }

    url(instance, uuid) {
        if (uuid) {
            return `/api/instance/${instance}/disk/${uuid}`
        }
        return `/api/instance/${instance}/disk`
    }

    create(data) {
        let your = this;

        $.post(your.url(this.instance), JSON.stringify(data), function (data, status) {
            $(your.tasks).append(AlertSuccess(`start disk '${data.name}' success`));
        }).fail(function (e) {
            $(your.tasks).append(AlertDanger((`${this.type} ${this.url}: ${e.responseText}`)));
        });
    }

    delete() {
        let your = this;

        this.uuids.forEach(function (uuid, index, err) {
            $.delete(your.url(your.instance, uuid), function (data, status) {
                $(your.tasks).append(AlertSuccess(`remove disk '${item}' success`));
            }).fail(function (e) {
                $(your.tasks).append(AlertDanger((`${this.type} ${this.url}: ${e.responseText}`)));
            });
        });
    }

    edit(data) {
        let your = this;

        $.post(your.url(your.instance, your.uuids[0]), JSON.stringify(data), function (data, status) {
            $(your.tasks).append(AlertSuccess(`edit disk '${data.name}' success`));
        }).fail(function (e) {
            $(your.tasks).append(AlertDanger((`${this.type} ${this.url}: ${e.responseText}`)));
        });
    }
}