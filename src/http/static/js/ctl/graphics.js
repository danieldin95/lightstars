import {Ctl} from "./ctl.js";
import {GraphicsApi} from "../api/graphics.js";
import {GraphicsTable} from "../widget/graphics/table.js";
import {CheckBox} from "../widget/common/checkbox.js";


class CheckBoxCtl extends CheckBox {
}


export class GraphicsCtl extends Ctl {
    // {
    //   id: '#instance #graphics',
    //   uuid: uuid of instance,
    //   name: name of instance,
    // }
    constructor(props) {
        super(props);
        this.name = props.name;
        this.inst = props.uuid;

        this.checkbox = new CheckBox(props);
        this.uuids = this.checkbox.uuids;
        this.table = new GraphicsTable({
            id: this.child('#display-table'),
            inst: this.inst,
        });

        // register button's click.
        $(this.child('#remove')).on("click", this, function (e) {
            new GraphicsApi({
                inst: e.data.inst,
                uuids: e.data.uuids.store,
                name: e.data.name}).delete();
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
        new GraphicsApi({inst: this.inst, name: this.name}).create(data);
    }

    edit(data) {
        new GraphicsApi({inst: this.inst, name: this.name}).edit(data);
    }
}
