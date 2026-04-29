import {Controller} from "./controller.js";
import {GraphicsApi} from "../api/graphicsapi.js";
import {GraphicsTableWid} from "../widget/graphics/graphicstable.js";
import {CheckboxWid} from "../widget/common/checkbox.js";
import {ConfirmActionWid} from "../widget/common/confirmaction.js";


class CheckBoxCtl extends CheckboxWid {
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

        this.CheckboxWid = new CheckBoxCtl(props);
        this.uuids = this.CheckboxWid.uuids;
        this.table = new GraphicsTableWid({
            id: this.child('#display-table'),
            inst: this.inst,
        });

        // register button's click.
        $(this.child('#remove')).on("click", (e) => {
            let uuids = this.uuids.store.slice();
            new ConfirmActionWid({
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
                this.CheckboxWid.refresh();
            });
        });
        this.table.refresh((e) => {
            this.CheckboxWid.refresh();
        });
    }

    create(data) {
        new GraphicsApi({inst: this.inst, name: this.name}).create(data);
    }

    edit(data) {
        new GraphicsApi({inst: this.inst, name: this.name}).edit(data);
    }
}
