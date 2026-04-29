import {Controller} from "./controller.js";
import {porttable} from "../widget/port/porttable.js";
import {checkbox} from "../widget/common/checkbox.js";
import {InterfaceApi} from "../api/interfaceapi.js";
import {confirmaction} from "../widget/common/confirmaction.js";


class CheckBoxCtl extends checkbox {
}


export class Port extends Controller {
    // {
    //   id: '#network #port',
    //   bridge: bridge of network,
    //   name: name of network,
    // }
    constructor(props) {
        super(props);
        this.name = props.name;
        this.uuid = props.uuid;
        this.bridge = props.bridge;
        this.confirm = props.confirm;

        this.checkbox = new CheckBoxCtl(props);
        this.uuids = this.checkbox.uuids;
        this.table = new porttable({
            id: this.child('#display-table'),
            uuid: this.uuid,
            name: this.name,
            bridge: this.bridge,
        });
        $(this.child('#remove')).on("click", (e) => {
            let uuids = this.uuids.store.slice();
            new confirmaction({
                id: this.confirm,
                action: "remove",
                name: uuids.join(", "),
                message: "remove",
            }).onsubmit(() => {
                uuids.forEach((item, index, err) => {
                    let values = item.split(',');
                    if (values.length === 2) {
                        new InterfaceApi({
                            inst: values[0],
                            uuids: values[1]
                        }).delete();
                    }
                });
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
}
