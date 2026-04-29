import {Controller} from "./controller.js";
import {GraphicsApi} from "../api/graphics.js";
import {GraphicsTable} from "../widget/graphics/table.js";
import {CheckBox} from "../widget/common/checkbox.js";
import {ConfirmAction} from "../widget/common/confirm.js";


class CheckBoxCtl extends CheckBox {
}


export class GraphicsCtl extends Controller {
    // {
    //   id: '#instance #graphics',
    //   uuid: uuid of instance,
    //   name: name of instance,
    // }
    constructor(props) {
        super(props);
        this.name = props.name;
        this.inst = props.uuid;
        this.confirm = props.confirm;

        this.checkbox = new CheckBoxCtl(props);
        this.uuids = this.checkbox.uuids;
        this.table = new GraphicsTable({
            id: this.child('#display-table'),
            inst: this.inst,
        });

        // register button's click.
        $(this.child('#remove')).on("click", (e) => {
            let uuids = this.uuids.store.slice();
            new ConfirmAction({
                id: this.confirm,
                action: "remove",
                name: uuids.join(", "),
                message: "remove",
            }).onsubmit(() => {
                new GraphicsApi({
                    inst: this.inst,
                    uuids: uuids,
                    name: this.name}).delete();
            });
            $(this.confirm).modal("show");
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
