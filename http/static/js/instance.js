import {AlertDanger} from "./alert.js";
import {AlertWarn} from "./alert.js";
import {AlertSuccess} from "./alert.js";

export class Instance {

    constructor() {
        this.instanes = [];

        let disabled = function(is) {
            if (is) {
                $("instance-start button").addClass('disabled');
                $("instance-console button").addClass('disabled');
                $("instance-shutdown button").addClass('disabled');
                $("instance-more button").addClass('disabled');
            } else {
                $("instance-start button").removeClass('disabled');
                $("instance-console button").removeClass('disabled');
                $("instance-shutdown button").removeClass('disabled');
                $("instance-more button").removeClass('disabled');
            }
        };

        let instance_dom = $("instance-on-one input");
        for (let i = 0; i < instance_dom.length; i++) {
            instance_dom.eq(i).on("change", this, function(e) {
                let uuid = $(this).attr("data");
                if ($(this).prop("checked")) {
                    e.data.instanes.push(uuid)
                } else {
                    e.data.instanes = e.data.instanes.filter(v => v != uuid);
                }
                disabled(e.data.instanes.length == 0);
            });
        }
        $("instance-on-all input").on("change", this, function(e) {
            if ($(this).prop("checked")) {
                instance_dom.each(function (index, element) {
                    e.data.instanes.push($(this).attr("data"));
                    $(element).prop("checked", true);
                });
            } else {
                instance_dom.each(function (index, element) {
                    e.data.instanes = [];
                    $(element).prop("checked", false);
                });
            }
            disabled(e.data.instanes.length == 0);
        });

        // Disabled firstly.
        disabled(this.instanes.length == 0);

        // Register click handle.
        $("instance-console").on("click", this, function (e) {
            e.data.console(this);
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
        $("instance-more-suspend").on("click", this, function (e) {
            e.data.suspend(this);
        });
        $("instance-more-resume").on("click", this, function (e) {
            e.data.resume(this);
        });
        $("instance-more-destroy").on("click", this, function (e) {
            e.data.destroy(this);
        });
        $("instance-more-remove").on("click", this, function (e) {
            e.data.remove(this);
        });
    }

    start(on) {
        this.instanes.forEach(function (item, index, err) {
            let data = {action: 'start'};

            $.put("api/instance/"+item, JSON.stringify(data), function (data, status) {
                $("infos").append(AlertSuccess(`start instance '${item}' success`));
            }).fail(function (e) {
                $("errors").append(AlertDanger((`${this.type} ${this.url}: ${e.responseText}`)));
            });
        });
    }

    shutdown(on) {
        this.instanes.forEach(function (item, index, err) {
            let data = {action: 'shutdown'};

            $.put("api/instance/"+item, JSON.stringify(data), function (data, status) {
                $("infos").append(AlertSuccess(`shutdown instance '${item}' success`));
            }).fail(function (e) {
                $("errors").append(AlertWarn((`${this.type} ${this.url}: ${e.responseText}`)));
            });
        });
    }

    reset(on) {
        this.instanes.forEach(function (item, index, err) {
            let data = {action: 'reset'};

            $.put("api/instance/"+item, JSON.stringify(data), function (data, status) {
                $("infos").append(AlertSuccess(`reset instance '${item}' success`));
            }).fail(function (e) {
                $("errors").append(AlertDanger((`${this.type} ${this.url}: ${e.responseText}`)));
            });
        });
    }

    suspend(on) {
        this.instanes.forEach(function (item, index, err) {
            let data = {action: 'suspend'};

            $.put("api/instance/"+item, JSON.stringify(data), function (data, status) {
                $("infos").append(AlertSuccess(`suspend instance '${item}' success`));
            }).fail(function (e) {
                $("errors").append(AlertDanger((`${this.type} ${this.url}: ${e.responseText}`)));
            });
        });
    }

    resume(on) {
        this.instanes.forEach(function (item, index, err) {
            let data = {action: 'resume'};

            $.put("api/instance/"+item, JSON.stringify(data), function (data, status) {
                $("infos").append(AlertSuccess(`resume instance '${item}' success`));
            }).fail(function (e) {
                $("errors").append(AlertDanger((`${this.type} ${this.url}: ${e.responseText}`)));
            })
        });
    }

    destroy(on) {
        this.instanes.forEach(function (item, index, err) {
            let data = {action: 'destroy'};

            $.put("api/instance/"+item, JSON.stringify(data), function (data, status) {
                $("infos").append(AlertSuccess(`destroy instance '${item}' success`));
            }).fail(function (e) {
                $("errors").append(AlertDanger((`${this.type} ${this.url}: ${e.responseText}`)));
            })
        });
    }

    remove(on) {
        this.instanes.forEach(function (item, index, err) {
            let data = {action: 'remove'};

            $.put("api/instance/"+item, JSON.stringify(data), function (data, status) {
                $("infos").append(AlertSuccess(`remove instance '${item}' success`));
            }).fail(function (e) {
                $("errors").append(AlertDanger((`${this.type} ${this.url}: ${e.responseText}`)));
            })
        });
    }

    console(on) {
        this.instanes.forEach(function (item, index, err) {
            window.open("/static/console?instance="+item);
        });
    }

    create (data) {
        $.post("api/instance", JSON.stringify(data), function (data, status) {
            $("infos").append(AlertSuccess(`create instance '${data.name}' success`));
        }).fail(function (e) {
            $("errors").append(AlertDanger((`${this.type} ${this.url}: ${e.responseText}`)));
        })
    }
}