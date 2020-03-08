import {Api} from "./api.js"
import {AlertDanger, AlertSuccess, AlertWarn} from "../com/alert.js";


export class GraphicsApi extends Api {
    // {
    //   instance: 'uuid',
    //   uuids: [uuid],
    //   tasks: 'tasks',
    //   name: ''
    // }
    constructor(props) {
        super(props);
        this.instance = props.instance;
    }

    url(instance, uuid) {
        if (uuid) {
            return `/api/instance/${instance}/graphics/${uuid}`
        }
        return `/api/instance/${instance}/graphics`
    }

    list(data, func) {
        let your = this;

        $.get(your.url(this.instance), {format: 'schema'}, function (resp, status) {
            func({data, resp});
        }).fail(function (e) {
            $(your.tasks).append(AlertDanger((`${this.type} ${this.url}: ${e.responseText}`)));
        });
    }

    create(data) {
        let your = this;

        $.post(your.url(this.instance), JSON.stringify(data), function (data, status) {
            $(your.tasks).append(AlertSuccess(`create graphics ${data.message}`));
        }).fail(function (e) {
            $(your.tasks).append(AlertDanger((`${this.type} ${this.url}: ${e.responseText}`)));
        });
    }

    edit(data) {
        let your = this;

        $.post(your.url(your.instance, your.uuids[0]), JSON.stringify(data), function (data, status) {
            $(your.tasks).append(AlertSuccess(`edit graphics '${data.name}' ${data.message}`));
        }).fail(function (e) {
            $(your.tasks).append(AlertDanger((`${this.type} ${this.url}: ${e.responseText}`)));
        });
    }
}