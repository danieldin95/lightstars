import {GraphicsApi} from "../api/graphics.js";
import {GraphicsTable} from "../widget/graphics/table.js";
import {CheckBoxTab} from "../widget/checkbox/checkbox.js";


class CheckBox extends CheckBoxTab {
}


export class Graphics {
    // {
    //   id: '#instance #graphics',
    //   uuid: uuid of instance,
    //   name: name of instance,
    // }
    constructor(props) {
        this.id = props.id;
        this.name = props.name;
        this.inst = props.uuid;

        this.checkbox = new CheckBox(props);
        this.uuids = this.checkbox.uuids;
        this.table = new GraphicsTable({
            id: `${this.id} #display-table`,
            inst: this.inst,
        });

        // register button's click.
        $(`${this.id} #remove`).on("click", this, function (e) {
            new GraphicsApi({
                inst: e.data.inst,
                uuids: e.data.uuids.store,
                name: e.data.name}).delete();
        });

        // refresh table and register refresh click.
        $(`${this.id} #refresh`).on("click", (e) => {
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