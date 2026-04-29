import {Controller} from "./controller.js";
import {DiskApi} from "../api/diskapi.js";
import {DiskTableWid} from "../widget/disk/disktable.js";
import {CheckboxWid} from "../widget/common/checkbox.js";


class CheckBoxCtl extends CheckboxWid {
}


export class DiskCtl extends Controller {
    // {
    //   id: '#instance #disk',
    //   uuid: uuid of instance,
    //   name: name of instance,
    //   onRemove: function {},
    // }
    constructor(props) {
        super(props);
        this.name = props.name;
        this.inst = props.uuid;

        this.CheckboxWid = new CheckBoxCtl(props);
        this.uuids = this.CheckboxWid.uuids;
        this.table = new DiskTableWid({
            id: this.child('#display-table'),
            inst: this.inst,
        });

        // register button's click.
        $(this.child('#remove')).on("click", (e) => {
            let data = {inst: this.inst, uuids: this.uuids.store};
            if (this.props.onRemove) {
                this.props.onRemove(data);
            } else {
                new DiskApi(data).delete();
            }
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
        new DiskApi({inst: this.inst}).create(data);
    }

    edit(data) {
        new DiskApi({inst: this.inst}).edit(data);
    }
}
