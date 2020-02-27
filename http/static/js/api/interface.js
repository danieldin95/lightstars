import {Api} from "./api.js"
import {AlertDanger, AlertSuccess} from "../com/alert.js";


export class InterfaceApi extends Api {
    // {
    //   uuids: [],
    //   tasks: 'tasks',
    //   name: ''
    // }
    constructor(props) {
        super(props);

        this.instance = props.instance;
    }

    url(instance, uuid) {
        if (uuid) {
            return `/api/instance/${instance}/interface/${uuid}`
        }
        return `/api/instance/${instance}/interface`
    }

    create(data) {
        let your = this;

        $.post(your.url(this.instance), JSON.stringify(data), function (data, status) {
            $(your.tasks).append(AlertSuccess(`start interface ${data.message}`));
        }).fail(function (e) {
            $(your.tasks).append(AlertDanger((`${this.type} ${this.url}: ${e.responseText}`)));
        });
    }

    delete() {
        let your = this;

        this.uuids.forEach(function (uuid, index, err) {
            $.delete(your.url(your.instance, uuid), function (data, status) {
                $(your.tasks).append(AlertSuccess(`remove interface '${uuid}' ${data.message}`));
            }).fail(function (e) {
                $(your.tasks).append(AlertDanger((`${this.type} ${this.url}: ${e.responseText}`)));
            });
        });
    }

    edit(data) {
        let your = this;

        $.post(your.url(your.instance, your.uuids[0]), JSON.stringify(data), function (data, status) {
            $(your.tasks).append(AlertSuccess(`edit interface '${data.name}' ${data.message}`));
        }).fail(function (e) {
            $(your.tasks).append(AlertDanger((`${this.type} ${this.url}: ${e.responseText}`)));
        });
    }
}