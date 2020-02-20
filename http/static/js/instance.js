import Alert from "./widget/alert.js";

export default class Instance {
    constructor() {
        this.instanes = [];

        let instanceDom = $("on-instance input");
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
        // register click handle.
        $("btn-create").on("click", this, function (e) {
            e.data.create(this);
        });
        $("btn-start, btn-more-start").on("click", this, function (e) {
            e.data.start(this);
        });
        $("btn-shutdown, btn-more-shutdown").on("click", this, function (e) {
            e.data.shutdown(this);
        });
        $("btn-more-reset").on("click", this, function (e) {
            e.data.reset(this);
        });
    }
    start(on) {
        this.instanes.forEach(function (item, index, err) {
            $.post("instance/"+item, {action: 'start'}, function (data, status) {
                console.log("success", status, data);
            }).fail(function (e) {
                $("errors").append(new Alert((`${this.type} ${this.url}: ${e.responseText}`)).danger());
            });
        });
    }
    shutdown(on) {
        this.instanes.forEach(function (item, index, err) {
            $.post("instance/"+item, {action: 'shutdown'}, function (data, status) {
                console.log("success", status, data);
            }).fail(function (e) {
                $("errors").append(new Alert((`${this.type} ${this.url}: ${e.responseText}`)).danger());
            });
        });
    }
    reset(on) {
        this.instanes.forEach(function (item, index, err) {
            $.post("instance/"+item, {action: 'reset'}, function (data, status) {
                console.log("success", status, data);
            }).fail(function (e) {
                $("errors").append(new Alert((`${this.type} ${this.url}: ${e.responseText}`)).danger());
            });
        });
    }
    suspend(on) {
        this.instanes.forEach(function (item, index, err) {
            $.post("instance/"+item, {action: 'suspend'}, function (data, status) {
                console.log("success", status, data);
            }).fail(function (e) {
                $("errors").append(new Alert((`${this.type} ${this.url}: ${e.responseText}`)).danger());
            });
        });
    }
    resume(on) {
        this.instanes.forEach(function (item, index, err) {
            $.post("instance/"+item, {action: 'resume'}, function (data, status) {
                console.log("success", status, data);
            }).fail(function (e) {
                $("errors").append(new Alert((`${this.type} ${this.url}: ${e.responseText}`)).danger());
            })
        });
    }
    create (on) {
        this.instanes.forEach(function (item, index, err) {
            $.post("instance/"+item, {action: 'create'}, function (data, status) {
                console.log("success", status, data);
            }).fail(function (e) {
                $("errors").append(new Alert((`${this.type} ${this.url}: ${e.responseText}`)).danger());
            });
        });
    }
}