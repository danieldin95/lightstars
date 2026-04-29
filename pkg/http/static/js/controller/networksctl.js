import {Controller} from "./controller.js";
import {NetworkApi} from "../api/networkapi.js";
import {networktable} from "../widget/network/networktable.js";
import {checkbox} from "../widget/common/checkbox.js";
import {confirmaction} from "../widget/common/confirmaction.js";


class CheckBoxCtl extends checkbox {
}


export class Networks extends Controller {
    // {
    //   id: "#networks",
    //   onthis: function (e) {},
    // }
    constructor(props) {
        super(props);
        this.checkbox = new CheckBoxCtl(props);
        this.uuids = this.checkbox.uuids;
        this.table = new networktable({id: `${this.id} #display-table`});
        this.confirm = props.confirm;

        // register buttons's click.
        $(this.child('#delete')).on("click", (e) => {
            let uuids = this.uuids.store.slice();
            new confirmaction({
                id: this.confirm,
                action: "remove",
                name: uuids.join(", "),
                message: "remove",
            }).onsubmit(() => {
                new NetworkApi({uuids: uuids}).delete();
            });
            $(this.confirm).modal("show");
        });

        // refresh table and register refresh click.
        $(this.child('#refresh')).on("click", (e) => {
            this.refresh();
        });
        this.refresh();
    }

    create(data) {
        new NetworkApi().create(data);
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
}
