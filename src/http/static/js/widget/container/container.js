import {Widget} from "../widget.js";


export class Container {
    // {
    //    parent: "#Container",
    //    default: "instances" // set default panel
    //    force: false, // force to apply default.
    // }
    constructor(props) {
        this.props = props;
        this.parent  = props.parent ? props.parent : '';
        this.current = props.current ? props.current : '';
        this.force = props.force;
        console.log([this.parent, this.current].join(" "));
        console.log('Base', props);
    }

    render() {
        $(this.parent).html(this.template());
    }

    loading() {
        console.log("Base", "implement me")
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
        $(document).attr("title", name + ' - LightStar');
    }
}
