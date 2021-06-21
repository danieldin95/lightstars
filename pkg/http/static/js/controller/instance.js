import {Controller} from './controller.js';
import {InstanceApi} from "../api/instance.js";
import {InstanceTable} from "../widget/instance/table.js";
import {CheckBox} from "../widget/common/checkbox.js";
import {InstanceFooter} from "../widget/instance/footer.js";


class CheckboxCtl extends CheckBox {
    change(from) {
        super.change(from);
        if (from.store.length === 0) {
            $(this.child('#start')).attr("disabled","disabled");
            $(this.child('#console')).attr("disabled","disabled");
            $(this.child('#shutdown')).attr("disabled","disabled");
            $(this.child('#more')).attr("disabled","disabled");
        } else {
            $(this.child('#start')).removeAttr('disabled');
            $(this.child('#console')).removeAttr('disabled');
            $(this.child('#shutdown')).removeAttr('disabled');
            $(this.child('#more')).removeAttr('disabled');
        }
    }
}


export class InstanceCtl extends Controller {
    // {
    //   id: '#instances'
    //   onthis: function (e) {},
    // }
    constructor(props) {
        super(props);
        this.checkbox = new CheckboxCtl(props);
        this.uuids = this.checkbox.uuids;
        this.table = new InstanceTable({id: `${this.id} #display-body`});
        this.footer = new InstanceFooter({id: `${this.id} #footer`});

        // register buttons's click.
        $(this.child('#console')).on("click", this.uuids, function (e) {
            let props = {uuids: e.data.store, passwd: {}, name: {}};
            e.data.store.forEach(function (v) {
                let obj = $('input[data='+v+']');
                props.passwd[v] = obj.attr('passwd');
                props.name[v] = obj.attr('name');
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
        this.footer.refresh((e) => {
            //TODO
        })
    }

    create(data) {
        new InstanceApi().create(data);
    }
}
