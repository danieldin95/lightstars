export default class Instance {
    constructor() {
        this.instanes = [];

        let instanceDom = $("tr [name=on-instance]")
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
        $("[name=btn-create]").on("click", this, function (e) {
            e.data.create(this);
        });
        $("[name=btn-start]").on("click", this, function (e) {
            e.data.start(this);
        });
        $("[name=btn-shutdown]").on("click", this, function (e) {
            e.data.shutdown(this);
        });
        $("[name=btn-more] [name=btn-start]").on("click", this, function (e) {
            e.data.start(this);
        });
        $("[name=btn-more] [name=shutdown]").on("click", this, function (e) {
            e.data.shutdown(this);
        });
        $("[name=btn-more] [name=btn-reset]").on("click", this, function (e) {
            e.data.reset(this);
        });
    }
    start(on) {
        console.log("start", this.instanes);
        this.instanes.forEach(function (item, index, err) {
            $.post("instance/"+item, {action: 'start'}, function (data, status) {
                    console.log("success", status, data);
                })
                .fail(function (e) {
                    $("#errors").append(
                        `<div class="alert alert-danger alert-dismissible fade show" role="alert">
                            ${this.type} ${this.url}: ${e.responseText}
                            <button type="button" class="close" data-dismiss="alert" aria-label="Close">
                                <span aria-hidden="true">&times;</span>
                            </button>
                          </div>`);
                    console.log(this);
                })
        });
    }
    shutdown(on) {
        console.log("shutdown", this.instanes);
    }
    reset(on) {
        console.log("reset", this.instanes);
    }
    suspend(on) {
        console.log("suspend", this.instanes);
    }
    resume(on) {
        console.log("shutdown", this.shutdown());
    }
    create (on) {
        console.log("shutdown", this.instanes);
    }
}