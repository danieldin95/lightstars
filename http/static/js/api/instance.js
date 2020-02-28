import {Api} from "./api.js";
import {AlertDanger, AlertSuccess, AlertWarn} from "../com/alert.js";


export class InstanceApi extends Api {
    // {
    //   uuids: [],
    //   tasks: 'tasks',
    //   name: ''
    // }
    constructor(props) {
        super(props)
    }

    url(uuid) {
        if (uuid) {
            return `/api/instance/${uuid}`
        }
        return 'api/instance'
    }

    list(data, func) {
        let your = this;

        $.get(your.url(), function (resp, status) {
            func({data, resp});
        }).fail(function (e) {
            $(your.tasks).append(AlertDanger((`${this.type} ${this.url}: ${e.responseText}`)));
        });
    }

    start() {
        let your = this;
        let data = JSON.stringify({action: 'start'});

        this.uuids.forEach(function (uuid, index, err) {
            $.put(your.url(uuid), data, function (data, status) {
                $(your.tasks).append(AlertSuccess(`start instance '${uuid}' success`));
            }).fail(function (e) {
                $(your.tasks).append(AlertDanger((`${this.type} ${this.url}: ${e.responseText}`)));
            });
        });
    }

    shutdown() {
        let your = this;
        let data = JSON.stringify({action: 'shutdown'});

        this.uuids.forEach(function (uuid, index, err) {
            $.put(your.url(uuid), data, function (data, status) {
                $(your.tasks).append(AlertSuccess(`shutdown instance '${uuid}' success`));
            }).fail(function (e) {
                $(your.tasks).append(AlertWarn((`${this.type} ${this.url}: ${e.responseText}`)));
            });
        });
    }

    reset() {
        let your = this;
        let data = JSON.stringify({action: 'reset'});

        this.uuids.forEach(function (uuid, index, err) {
            $.put(your.url(uuid), data, function (data, status) {
                $(your.tasks).append(AlertSuccess(`reset instance '${uuid}' success`));
            }).fail(function (e) {
                $(your.tasks).append(AlertDanger((`${this.type} ${this.url}: ${e.responseText}`)));
            });
        });
    }

    suspend() {
        let your = this;
        let data = JSON.stringify({action: 'suspend'});

        this.uuids.forEach(function (uuid, index, err) {
            $.put(your.url(uuid), data, function (data, status) {
                $(your.tasks).append(AlertSuccess(`suspend instance '${uuid}' success`));
            }).fail(function (e) {
                $(your.tasks).append(AlertDanger((`${this.type} ${this.url}: ${e.responseText}`)));
            });
        });
    }

    resume() {
        let your = this;
        let data = JSON.stringify({action: 'resume'});

        this.uuids.forEach(function (uuid, index, err) {
            $.put(your.url(uuid), data, function (data, status) {
                $(your.tasks).append(AlertSuccess(`resume instance '${uuid}' success`));
            }).fail(function (e) {
                $(your.tasks).append(AlertDanger((`${this.type} ${this.url}: ${e.responseText}`)));
            });
        });
    }

    destroy() {
        let your = this;
        let data = JSON.stringify({action: 'destroy'});

        this.uuids.forEach(function (uuid, index, err) {
            $.put(your.url(uuid), data, function (data, status) {
                $(your.tasks).append(AlertSuccess(`destroy instance '${uuid}' success`));
            }).fail(function (e) {
                $(your.tasks).append(AlertDanger((`${this.type} ${this.url}: ${e.responseText}`)));
            });
        });
    }

    remove() {
        let your = this;

        this.uuids.forEach(function (uuid, index, err) {
            $.delete(your.url(uuid), function (data, status) {
                $(your.tasks).append(AlertSuccess(`remove instance '${uuid}' success`));
            }).fail(function (e) {
                $(your.tasks).append(AlertDanger((`${this.type} ${this.url}: ${e.responseText}`)));
            });
        });
    }

    console() {
        let password = this.props.passwd;
        this.uuids.forEach(function (uuid, index, err) {
            window.open("/ui/console?instance="+uuid+"&password="+password[uuid]);
        });
    }

    create (data) {
        let your = this;

        $.post(your.url(), JSON.stringify(data), function (data, status) {
            $(your.tasks).append(AlertSuccess(`create instance '${data.name}' success`));
        }).fail(function (e) {
            $(your.tasks).append(AlertDanger((`${this.type} ${this.url}: ${e.responseText}`)));
        });
    }
}