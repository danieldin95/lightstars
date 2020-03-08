import {GraphicsApi} from "./api/graphics.js";
import {CheckBoxTop} from "./com/utils.js";
import {GraphicsTable} from "./widget/graphics/table.js";

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

        this.checkbox = new Checkbox(props);
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


export class Checkbox {
    // {
    //   id: '#instane #graphics'
    // }
    constructor(props) {
        this.id = props.id;
        this.uuids = {store: [], id: this.id};

        this.top = new CheckBoxTop({
            one: `${this.id} #on-one`,
            all: `${this.id} #on-all`,
            change: (e) => {
                this.change(this.uuids, e);
            },
        });

        // disabled firstly.
        this.change(this.uuids, this.uuids);
    }

    refresh() {
        this.top.refresh();
    }

    change(record, from) {
        record.store = from.store;

        if (from.store.length !== 1) {
            $(`${record.id} #edit`).addClass('disabled');
        } else {
            $(`${record.id} #edit`).removeClass('disabled');
        }
    }
}