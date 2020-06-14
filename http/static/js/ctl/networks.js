import {Ctl} from "./ctl.js";
import {NetworkApi} from "../api/network.js";
import {NetworkTable} from "../widget/network/table.js";
import {CheckBox} from "../widget/checkbox/checkbox.js";


class CheckBoxCtl extends CheckBox {
}


export class NetworksCtl extends Ctl {
    // {
    //   id: "#networks",
    //   onthis: function (e) {},
    // }
    constructor(props) {
        super(props);
        this.checkbox = new CheckBoxCtl(props);
        this.uuids = this.checkbox.uuids;
        this.table = new NetworkTable({id: `${this.id} #display-table`});

        // register buttons's click.
        $(this.child('#delete')).on("click", (e) => {
            new NetworkApi({uuids: this.uuids.store}).delete();
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
