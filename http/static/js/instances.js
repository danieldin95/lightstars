import {InstanceApi} from "./api/instance.js";
import {InstanceTable} from "./widget/instance/table.js";
import {CheckBoxTab} from "./widget/checkbox/checkbox.js";


class Checkbox extends CheckBoxTab {
    change(from) {
        super.change(from);

        if (from.store.length === 0) {
            $(`${this.uuids.id} #start`).addClass('disabled');
            $(`${this.uuids.id} #console`).addClass('disabled');
            $(`${this.uuids.id} #shutdown`).addClass('disabled');
            $(`${this.uuids.id} #more`).addClass('disabled');
        } else {
            $(`${this.uuids.id} #start`).removeClass('disabled');
            $(`${this.uuids.id} #console`).removeClass('disabled');
            $(`${this.uuids.id} #shutdown`).removeClass('disabled');
            $(`${this.uuids.id} #more`).removeClass('disabled');
        }
    }
}


export class Instances {
    // {
    //   id: '#instances'
    // }
    constructor(props) {
        this.id = props.id;
        this.props = props;
        this.checkbox = new Checkbox(props);
        this.uuids = this.checkbox.uuids;
        this.table = new InstanceTable({id: `${this.id} #display-body`});

        // register buttons's click.
        $(`${this.id} #console`).on("click", this.uuids, function (e) {
            let props = {uuids: e.data.store, passwd: {}};
            e.data.store.forEach(function (v) {
                props.passwd[v] = $(`input[data=${v}]`).attr('passwd');
            });
            new InstanceApi(props).console();
        });
        $(`${this.id} #start, ${this.id} #more-start`).on("click", this.uuids, function (e) {
            new InstanceApi({uuids: e.data.store}).start();
        });
        $(`${this.id} #more-shutdown`).on("click", this.uuids, function (e) {
            new InstanceApi({uuids: e.data.store}).shutdown();
        });
        $(`${this.id} #more-reset`).on("click", this.uuids, function (e) {
            new InstanceApi({uuids: e.data.store}).reset();
        });
        $(`${this.id} #more-suspend`).on("click", this.uuids, function (e) {
            new InstanceApi({uuids: e.data.store}).suspend();
        });
        $(`${this.id} #more-resume`).on("click", this.uuids, function (e) {
            new InstanceApi({uuids: e.data.store}).resume();
        });
        $(`${this.id} #more-destroy`).on("click", this.uuids, function (e) {
            new InstanceApi({uuids: e.data.store}).destroy();
        });
        $(`${this.id} #more-remove`).on("click", this.uuids, function (e) {
            new InstanceApi({uuids: e.data.store}).remove();
        });

        // refresh table and register refresh click.
        $(`${this.id} #refresh`).on("click", (e) => {
            this.table.refresh((e) => {
                this.checkbox.refresh();
            });
        });
        this.table.refresh((e) => {
            this.checkbox.refresh();
        });
    }

    create(data) {
        new InstanceApi().create(data);
    }
}