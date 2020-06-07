import {NetworkApi} from "../api/network.js";
import {NetworkTable} from "../widget/network/table.js";
import {CheckBoxTab} from "../widget/checkbox/checkbox.js";


class CheckBox extends CheckBoxTab {
}


export class NetworkCtl {
    // {
    //   id: "#networks"
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
            this.table.refresh((e) => {
                this.checkbox.refresh();
            })
        });
        this.table.refresh((e) => {
            this.checkbox.refresh();
        });
    }

    create(data) {
        new NetworkApi().create(data);
    }
}
