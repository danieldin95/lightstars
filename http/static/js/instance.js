import {InstanceApi} from './api/instance.js'

export class Instance {

    constructor() {
        this.disk = new Disk();
        this.interface = new Interface();

        // Register click handle.
        $("instance-start, instance-more-start").on("click", this, function (e) {
            new InstanceApi($(this).attr("data")).start();
        });
        $("instance-more-shutdown").on("click", this, function (e) {
            new InstanceApi($(this).attr("data")).shutdown();
        });
        $("instance-more-reset").on("click", this, function (e) {
            new InstanceApi($(this).attr("data")).reset();
        });
        $("instance-more-suspend").on("click", this, function (e) {
            new InstanceApi($(this).attr("data")).suspend();
        });
        $("instance-more-resume").on("click", this, function (e) {
            new InstanceApi($(this).attr("data")).resume();
        });
        $("instance-more-destroy").on("click", this, function (e) {
            new InstanceApi($(this).attr("data")).destroy();
        });
        $("instance-more-remove").on("click", this, function (e) {
            new InstanceApi($(this).attr("data")).remove();
        });
    }
}


export class Disk {

    constructor() {
        this.disks = [];

        let disabled = this.disable;
        let documentOne = $("disk-on-one input");
        let documentAll = $("disk-on-all input");

        for (let i = 0; i < documentOne.length; i++) {
            documentOne.eq(i).on("change", this.disks, function(e) {
                let uuid = $(this).attr("data");
                if ($(this).prop("checked")) {
                    e.data.push(uuid)
                } else {
                    e.data = e.data.filter(v => v != uuid);
                }
                disabled(e.data.length == 0);
            });
        }
        documentAll.on("change", this.disks, function(e) {
            if ($(this).prop("checked")) {
                documentOne.each(function (index, element) {
                    e.data.push($(this).attr("data"));
                    $(element).prop("checked", true);
                });
            } else {
                documentOne.each(function (index, element) {
                    e.data = [];
                    $(element).prop("checked", false);
                });
            }
            disabled(e.data.length == 0);
        });

        // Disabled firstly.
        disabled(this.disks.length === 0);
    }

    disable(is) {
        if (is) {
            $("disk-edit button").addClass('disabled');
            $("disk-remove button").addClass('disabled');
        } else {
            $("disk-edit button").removeClass('disabled');
            $("disk-remove button").removeClass('disabled');
        }
    }
}

export class Interface {

    constructor() {
        this.interfaces = [];

        let disabled = this.disable;
        let documentOne = $("interface-on-one input");
        let documentAll = $("interface-on-all input");

        for (let i = 0; i < documentOne.length; i++) {
            documentOne.eq(i).on("change", this.interfaces, function(e) {
                let uuid = $(this).attr("data");
                if ($(this).prop("checked")) {
                    e.data.push(uuid)
                } else {
                    e.data = e.data.filter(v => v != uuid);
                }
                disabled(e.data.length == 0);
            });
        }

        documentAll.on("change", this.interfaces, function(e) {
            if ($(this).prop("checked")) {
                documentOne.each(function (index, element) {
                    e.data.push($(this).attr("data"));
                    $(element).prop("checked", true);
                });
            } else {
                documentOne.each(function (index, element) {
                    e.data = [];
                    $(element).prop("checked", false);
                });
            }
            disabled(e.data.length == 0);
        });

        // Disabled firstly.
        disabled(this.interfaces.length === 0);
    }

    disable(is) {
        if (is) {
            $("interface-edit button").addClass('disabled');
            $("interface-remove button").addClass('disabled');
        } else {
            $("interface-edit button").removeClass('disabled');
            $("interface-remove button").removeClass('disabled');
        }
    }
}