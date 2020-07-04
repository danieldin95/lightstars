import {Ctl} from './ctl.js';
import {InstanceApi} from "../api/instance.js";
import {InstanceTable} from "../widget/instance/table.js";
import {CheckBox} from "../widget/checkbox/checkbox.js";


class CheckboxCtl extends CheckBox {
    change(from) {
        super.change(from);
        if (from.store.length === 0) {
            $(this.child('#start')).addClass('disabled');
            $(this.child('#console')).addClass('disabled');
            $(this.child('#shutdown')).addClass('disabled');
            $(this.child('#more')).addClass('disabled');
        } else {
            $(this.child('#start')).removeClass('disabled');
            $(this.child('#console')).removeClass('disabled');
            $(this.child('#shutdown')).removeClass('disabled');
            $(this.child('#more')).removeClass('disabled');
        }
    }
}


export class InstanceCtl extends Ctl {
    // {
    //   id: '#instances'
    //   onthis: function (e) {},
    // }
    constructor(props) {
        super(props);
        this.checkbox = new CheckboxCtl(props);
        this.uuids = this.checkbox.uuids;
        this.table = new InstanceTable({id: `${this.id} #display-body`});

        // register buttons's click.
        $(this.child('#console')).on("click", this.uuids, function (e) {
            let props = {uuids: e.data.store, passwd: {}};
            e.data.store.forEach(function (v) {
                props.passwd[v] = $(`input[data=${v}]`).attr('passwd');
            });
            new InstanceApi(props).console();
        });
        $(this.child('#start')).on("click", this.uuids, function (e) {
            new InstanceApi({uuids: e.data.store}).start();
        });
        $(this.child('#more-start')).on("click", this.uuids, function (e) {
            new InstanceApi({uuids: e.data.store}).start();
        });
        $(this.child('#more-shutdown')).on("click", this.uuids, function (e) {
            new InstanceApi({uuids: e.data.store}).shutdown();
        });
        $(this.child('#more-reset')).on("click", this.uuids, function (e) {
            new InstanceApi({uuids: e.data.store}).reset();
        });
        $(this.child('#more-suspend')).on("click", this.uuids, function (e) {
            new InstanceApi({uuids: e.data.store}).suspend();
        });
        $(this.child('#more-resume')).on("click", this.uuids, function (e) {
            new InstanceApi({uuids: e.data.store}).resume();
        });
        $(this.child('#more-destroy')).on("click", this.uuids, function (e) {
            new InstanceApi({uuids: e.data.store}).destroy();
        });

        // refresh table and register refresh click.
        $(this.child('#refresh')).on("click", (e) => {
            this.refresh();
        });
        this.refresh();
    }

    refresh() {
        this.table.refresh((e) => {
            this.checkbox.refresh();

            // register click on this table row.
            let func = this.props.onthis;
            if (func) {
                $(this.child('#on-this')).on('click', function(e) {
                    func({uuid: $(this).attr('data')});
                });
            }
        });
    }

    create(data) {
        new InstanceApi().create(data);
    }
}
