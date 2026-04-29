import {Controller} from './controller.js'
import {InterfaceApi} from "../api/interfaceapi.js";
import {InterfaceTableWid} from "../widget/interface/interfacetable.js";
import {CheckboxWid} from "../widget/common/checkbox.js";
import {ConfirmActionWid} from "../widget/common/confirmaction.js";


class CheckBoxCtl extends CheckboxWid {
}


export class InterfaceCtl extends Controller {
    // {
    //   id: '#instance #interface',
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
        this.table = new InterfaceTableWid({
            id: this.child('#display-table'),
            inst: this.inst,
        });

        // register buttons's click
        $(this.child('#remove')).on("click", (e) => {
            let uuids = this.uuids.store.slice();
            new ConfirmActionWid({
                id: this.confirm,
                action: "remove",
                name: uuids.join(", "),
                message: "remove",
            }).onsubmit(() => {
                new InterfaceApi({
                    inst: this.inst,
                    uuids: uuids
                }).delete();
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
        new InterfaceApi({inst: this.inst}).create(data);
    }

    edit(data) {
        new InterfaceApi({inst: this.inst}).edit(data);
    }
}
