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
        this.instance = props.uuid;

        this.checkbox = new Checkbox(props);
        this.graphics = this.checkbox.graphics;
        this.table = new GraphicsTable({
            id: `${this.id} #display-table`,
            instance: this.instance,
        });

        // register button's click.
        $(`${this.id} #remove`).on("click", this, function (e) {
            new GraphicsApi({
                instance: e.data.instance,
                uuids: e.data.graphics.store,
                name: e.data.name}).delete();
        });

        // refresh table and register refresh click.
        let refresh = function (your) {
            your.table.refresh(your.checkbox, function (e) {
                e.data.refresh();
            });
        };
        $(`${this.id} #refresh`).on("click", this, function (e) {
            refresh(e.data);
        });
        refresh(this);
    }

    create(data) {
        new GraphicsApi({instance: this.instance, name: this.name}).create(data);
    }

    edit(data) {
        new GraphicsApi({instance: this.instance, name: this.name}).edit(data);
    }
}


export class Checkbox {
    // {
    //   id: '#instane #graphics'
    // }
    constructor(props) {
        this.id = props.id;
        this.graphics = {store: [], id: this.id};

        let record = this.graphics;
        let change = this.change;

        this.top = new CheckBoxTop({
            one: `${this.id} #on-one`,
            all: `${this.id} #on-all`,
            change: function (e) {
                change(record, e);
            },
        });

        // disabled firstly.
        change(record, this.graphics);
    }

    refresh() {
        this.top.refresh();
    }

    change(record, from) {
        record.store = from.store;

        if (from.store.length != 1) {
            $(`${record.id} #edit`).addClass('disabled');
        } else {
            $(`${record.id} #edit`).removeClass('disabled');
        }
    }
}