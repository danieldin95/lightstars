import {InterfaceApi} from "../api/interface.js";
import {InterfaceTable} from "../widget/interface/table.js";
import {CheckBoxTab} from "../widget/checkbox/checkbox.js";


class CheckBox extends CheckBoxTab {
}


export class InterfaceCtl {
    // {
    //   id: '#instance #interface',
    //   uuid: uuid of instance,
    //   name: name of instance,
    // }
    constructor(props) {
        this.id = props.id;
        this.props = props;
        this.name = props.name;
        this.inst = props.uuid;

        this.checkbox = new CheckBox(props);
        this.uuids = this.checkbox.uuids;
        this.table = new InterfaceTable({
            id: `${this.id} #display-table`,
            inst: this.inst,
        });

        // register buttons's click
        $(`${this.id} #remove`).on("click", (e) => {
            new InterfaceApi({
                inst: this.inst,
                uuids: this.uuids.store
            }).delete();
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
        new InterfaceApi({inst: this.inst}).create(data);
    }

    edit(data) {
        new InterfaceApi({inst: this.inst}).edit(data);
    }
}
