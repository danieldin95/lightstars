import {Ctl} from './ctl.js'
import {InterfaceApi} from "../api/interface.js";
import {InterfaceTable} from "../widget/interface/table.js";
import {CheckBox} from "../widget/common/checkbox.js";


class CheckBoxCtl extends CheckBox {
}


export class InterfaceCtl extends Ctl {
    // {
    //   id: '#instance #interface',
    //   uuid: uuid of instance,
    //   name: name of instance,
    // }
    constructor(props) {
        super(props);
        this.name = props.name;
        this.inst = props.uuid;

        this.checkbox = new CheckBoxCtl(props);
        this.uuids = this.checkbox.uuids;
        this.table = new InterfaceTable({
            id: this.child('#display-table'),
            inst: this.inst,
        });

        // register buttons's click
        $(this.child('#remove')).on("click", (e) => {
            new InterfaceApi({
                inst: this.inst,
                uuids: this.uuids.store
            }).delete();
        });
        // refresh table and register refresh click.
        $(this.child('#refresh')).on("click", (e) => {
            this.table.refresh((e) => {
                this.checkbox.refresh();
            });
        });
        this.table.refresh((e) => {
            this.checkbox.refresh();
        });
    }

    create(data) {
        new InterfaceApi({inst: this.inst}).create(data);
    }

    edit(data) {
        new InterfaceApi({inst: this.inst}).edit(data);
    }
}
