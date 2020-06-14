import {NetworkApi} from "../api/network.js";
import {NetworkTable} from "../widget/network/table.js";
import {CheckBoxTab} from "../widget/checkbox/checkbox.js";


class CheckBox extends CheckBoxTab {
}


export class NetworksCtl {
    // {
    //   id: "#networks",
    //   onthis: function (e) {},
    // }
    constructor(props) {
        this.id = props.id;
        this.props = props;
        this.checkbox = new CheckBox(props);
        this.uuids = this.checkbox.uuids;
        this.table = new NetworkTable({id: `${this.id} #display-table`});

        // register buttons's click.
        $(`${this.id} #delete`).on("click", (e) => {
            new NetworkApi({uuids: this.uuids.store}).delete();
        });

        // refresh table and register refresh click.
        $(`${this.id} #refresh`).on("click", (e) => {
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
                $(`${this.id} #on-this`).on('click', function(e) {
                    func({uuid: $(this).attr('data')});
                });
            }
        });
    }
}
