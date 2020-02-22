import {InstanceApi} from "./api/instance.js";
import {ListenChangeAll} from "./com/utils.js";

export class Instances {

    constructor() {
        this.instanceOn = new InstanceOn();
        this.instances = this.instanceOn.uuids;

        // Register click handle.
        $("instance-console").on("click", this.instances, function (e) {
            new InstanceApi(e.data).console();
        });
        $("instance-start, instance-more-start").on("click", this.instances, function (e) {
            new InstanceApi(e.data).start();
        });
        $("instance-more-shutdown").on("click", this.instances, function (e) {
            new InstanceApi(e.data).shutdown();
        });
        $("instance-more-reset").on("click", this.instances, function (e) {
            new InstanceApi(e.data).reset();
        });
        $("instance-more-suspend").on("click", this.instances, function (e) {
            new InstanceApi(e.data).suspend();
        });
        $("instance-more-resume").on("click", this.instances, function (e) {
            new InstanceApi(e.data).resume();
        });
        $("instance-more-destroy").on("click", this.instances, function (e) {
            new InstanceApi(e.data).destroy();
        });
        $("instance-more-remove").on("click", this.instances, function (e) {
            new InstanceApi(e.data).remove();
        });
    }

    create(data) {
        new InstanceApi().create(data);
    }
}

export class InstanceOn {

    constructor() {
        this.uuids = [];

        let disabled = this.disable;
        ListenChangeAll(this.uuids, "instance-on-one input", "instance-on-all input", function (e) {
           disabled(e.data.length == 0);
        });

        // Disabled firstly.
        disabled(this.uuids.length === 0);
    }

    disable(is) {
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
    }
}