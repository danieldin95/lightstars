import {Location} from "../lib/location.js";


export class Container {
    // {
    //    parent: "#container",
    //    default: "instances" // set default panel
    // }
    constructor(props) {
        this.props = props;
        this.parent  = props.parent ? props.parent : '';
        this.current = props.current ? props.current : '';
        console.log('Container', props, [this.parent, this.current].join(" "));
        this._alias = Container._alias;
    }

    render() {
        $(this.parent).html(this.template());
    }

    loading() {
        console.log("Container", "implement me")
    }

    compile(tmpl, data) {
        return template.compile(tmpl)(data);
    }

    template(v) {
        return `<div id="${this.current}">TODO ${this.current}</div>`;
    }

    id(id) {
        if (id) {
            return [this.parent, this.current, id].join(" ");
        }
        return [this.parent, this.current].join(" ")
    }

    title(name) {
        if (!this._alias) {
            this._alias = "LightStar"
        }
        $(document).attr("title", `${name} - ${this._alias}`);
    }

    url(page) {
        let query = Location.query();
        return page + "?" + query
    }

    static alias(value) {
        if (value !== undefined) {
            this._alias = value;
        }
        return this._alias;
    }
}
