import {AlertDanger, AlertSuccess, AlertWarn} from "../com/alert.js";

export class InstanceApi {

    constructor(uuids, tasks, name) {
        console.log(uuids, typeof uuids);
        if (tasks) {
            this.tasks = tasks;
        } else {
            this.tasks = "tasks";
        }
        this.name = name;
        if (typeof uuids == "string") {
            console.log("...")
            this.uuids = [uuids];
        } else {
            this.uuids = uuids;
        }
        console.log(this.uuids);
    }

    start() {
        let tasks = $(this.tasks);
        let data = {action: 'start'};

        this.uuids.forEach(function (item, index, err) {
            $.put("/api/instance/"+item, JSON.stringify(data), function (data, status) {
                tasks.append(AlertSuccess(`start instance '${item}' success`));
            }).fail(function (e) {
                tasks.append(AlertDanger((`${this.type} ${this.url}: ${e.responseText}`)));
            });
        });
    }

    shutdown() {
        let tasks = $(this.tasks);
        let data = {action: 'shutdown'};

        this.uuids.forEach(function (item, index, err) {
            $.put("/api/instance/"+item, JSON.stringify(data), function (data, status) {
                tasks.append(AlertSuccess(`shutdown instance '${item}' success`));
            }).fail(function (e) {
                tasks.append(AlertWarn((`${this.type} ${this.url}: ${e.responseText}`)));
            });
        });
    }

    reset() {
        let tasks = $(this.tasks);
        let data = {action: 'reset'};

        this.uuids.forEach(function (item, index, err) {
            $.put("/api/instance/"+item, JSON.stringify(data), function (data, status) {
                tasks.append(AlertSuccess(`reset instance '${item}' success`));
            }).fail(function (e) {
                tasks.append(AlertDanger((`${this.type} ${this.url}: ${e.responseText}`)));
            });
        });
    }

    suspend() {
        let tasks = $(this.tasks);
        let data = {action: 'suspend'};

        this.uuids.forEach(function (item, index, err) {
            $.put("/api/instance/"+item, JSON.stringify(data), function (data, status) {
                tasks.append(AlertSuccess(`suspend instance '${item}' success`));
            }).fail(function (e) {
                tasks.append(AlertDanger((`${this.type} ${this.url}: ${e.responseText}`)));
            });
        });
    }

    resume() {
        let tasks = $(this.tasks);
        let data = {action: 'resume'};

        this.uuids.forEach(function (item, index, err) {
            $.put("/api/instance/" + item, JSON.stringify(data), function (data, status) {
                tasks.append(AlertSuccess(`resume instance '${item}' success`));
            }).fail(function (e) {
                tasks.append(AlertDanger((`${this.type} ${this.url}: ${e.responseText}`)));
            });
        });
    }

    destroy() {
        let tasks = $(this.tasks);
        let data = {action: 'destroy'};

        this.uuids.forEach(function (item, index, err) {
            $.put("/api/instance/" + item, JSON.stringify(data), function (data, status) {
                tasks.append(AlertSuccess(`destroy instance '${item}' success`));
            }).fail(function (e) {
                tasks.append(AlertDanger((`${this.type} ${this.url}: ${e.responseText}`)));
            });
        });
    }

    remove() {
        let tasks = $(this.tasks);

        this.uuids.forEach(function (item, index, err) {
            $.delete("/api/instance/" + item, function (data, status) {
                tasks.append(AlertSuccess(`remove instance '${item}' success`));
            }).fail(function (e) {
                tasks.append(AlertDanger((`${this.type} ${this.url}: ${e.responseText}`)));
            });
        });
    }

    console() {
        this.uuids.forEach(function (item, index, err) {
            window.open("/ui/console?instance="+item);
        });
    }

    create (data) {
        let tasks = $(this.tasks);

        $.post("/api/instance", JSON.stringify(data), function (data, status) {
            tasks.append(AlertSuccess(`create instance '${data.name}' success`));
        }).fail(function (e) {
            tasks.append(AlertDanger((`${this.type} ${this.url}: ${e.responseText}`)));
        });
    }
}