import {AlertDanger} from "./widget/alert.js";
import {AlertWarn} from "./widget/alert.js";

export default class Instance {
    constructor() {
        this.instanes = [];

        let instanceDom = $("instance-on input");
        for (let i = 0; i < instanceDom.length; i++) {
            instanceDom.eq(i).on("change", this, function(e) {
                let uuid = $(this).attr("data");
                if ($(this).prop("checked")) {
                    e.data.instanes.push(uuid)
                } else {
                    e.data.instanes = e.data.instanes.filter(v => v != uuid);
                }
            });
        }
        // Register click handle.
        $("instance-create").on("click", this, function (e) {
            e.data.create(this);
        });
        $("instance-start, instance-more-start").on("click", this, function (e) {
            e.data.start(this);
        });
        $("instance-shutdown, instance-more-shutdown").on("click", this, function (e) {
            e.data.shutdown(this);
        });
        $("instance-more-reset").on("click", this, function (e) {
            e.data.reset(this);
        });
    }
    start(on) {
        this.instanes.forEach(function (item, index, err) {
            $.post("instance/"+item, {action: 'start'}, function (data, status) {
                console.log("success", status, data);
            }).fail(function (e) {
                $("errors").append(AlertDanger((`${this.type} ${this.url}: ${e.responseText}`)));
            });
        });
    }
    shutdown(on) {
        this.instanes.forEach(function (item, index, err) {
            $.post("instance/"+item, {action: 'shutdown'}, function (data, status) {
                console.log("success", status, data);
            }).fail(function (e) {
                $("errors").append(AlertWarn((`${this.type} ${this.url}: ${e.responseText}`)));
            });
        });
    }
    reset(on) {
        this.instanes.forEach(function (item, index, err) {
            $.post("instance/"+item, {action: 'reset'}, function (data, status) {
                console.log("success", status, data);
            }).fail(function (e) {
                $("errors").append(AlertDanger((`${this.type} ${this.url}: ${e.responseText}`)));
            });
        });
    }
    suspend(on) {
        this.instanes.forEach(function (item, index, err) {
            $.post("instance/"+item, {action: 'suspend'}, function (data, status) {
                console.log("success", status, data);
            }).fail(function (e) {
                $("errors").append(AlertDanger((`${this.type} ${this.url}: ${e.responseText}`)));
            });
        });
    }
    resume(on) {
        this.instanes.forEach(function (item, index, err) {
            $.post("instance/"+item, {action: 'resume'}, function (data, status) {
                console.log("success", status, data);
            }).fail(function (e) {
                $("errors").append(AlertDanger((`${this.type} ${this.url}: ${e.responseText}`)));
            })
        });
    }
    create (on) {
        console.log("TODO diag create wizard.")
    }
}